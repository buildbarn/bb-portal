package authmetadataextraction

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-storage/pkg/auth"
	"github.com/buildbarn/bb-storage/pkg/clock"
	"github.com/buildbarn/bb-storage/pkg/jmespath"
	"github.com/buildbarn/bb-storage/pkg/program"
	"github.com/buildbarn/bb-storage/pkg/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthMetadataExtractors stores the JMESPath extractors for the authentication metadata.
type AuthMetadataExtractors struct {
	DisplayNameJmespathExpression *jmespath.Expression
	ExternalIDJmespathExpression  *jmespath.Expression
	UserInfoJmespathExpression    *jmespath.Expression
}

// AuthenticatedUserSummary holds details about an authenticated user
type AuthenticatedUserSummary struct {
	ExternalID  string
	DisplayName *string
	UserInfo    map[string]any
}

// AuthMetadataExtractorsFromConfiguration sets the extractors from the configuration.
func AuthMetadataExtractorsFromConfiguration(authMetadataExtractorConfig *bb_portal.AuthMetadataExtractorConfiguration, group program.Group) (*AuthMetadataExtractors, error) {
	if authMetadataExtractorConfig == nil {
		slog.Info("Did not create authentication metadata extractors because authMetadataKeyConfiguration is not configured")
		return nil, nil
	}

	var externalIDJmespathExpression *jmespath.Expression
	var displayNameJmespathExpression *jmespath.Expression
	var userInfoJmespathExpression *jmespath.Expression
	var err error

	if authMetadataExtractorConfig.ExternalIdExtractionJmespathExpression == nil {
		return nil, status.Error(codes.InvalidArgument, "ExternalIdExtractionJmespathExpression is not configured")
	}

	externalIDJmespathExpression, err = jmespath.NewExpressionFromConfiguration(authMetadataExtractorConfig.ExternalIdExtractionJmespathExpression, group, clock.SystemClock)
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to create external ID extractor")
	}

	if authMetadataExtractorConfig.DisplayNameExtractionJmespathExpression == nil {
		slog.Warn("DisplayNameExtractionJmespathExpression is not configured, display names will be not be extracted")
		displayNameJmespathExpression = nil
	} else {
		displayNameJmespathExpression, err = jmespath.NewExpressionFromConfiguration(authMetadataExtractorConfig.DisplayNameExtractionJmespathExpression, group, clock.SystemClock)
		if err != nil {
			return nil, util.StatusWrap(err, "Failed to create display name extractor")
		}
	}

	if authMetadataExtractorConfig.UserInfoExtractionJmespathExpression == nil {
		slog.Warn("UserInfoExtractionJmespathExpression is not configured, user info will be not be extracted")
		userInfoJmespathExpression = nil
	} else {
		userInfoJmespathExpression, err = jmespath.NewExpressionFromConfiguration(authMetadataExtractorConfig.UserInfoExtractionJmespathExpression, group, clock.SystemClock)
		if err != nil {
			return nil, util.StatusWrap(err, "Failed to create user info extractor")
		}
	}

	return &AuthMetadataExtractors{
		DisplayNameJmespathExpression: displayNameJmespathExpression,
		ExternalIDJmespathExpression:  externalIDJmespathExpression,
		UserInfoJmespathExpression:    userInfoJmespathExpression,
	}, nil
}

// AuthenticatedUserSummaryFromContext creates an AuthenticatedUserSummary from a context.
func AuthenticatedUserSummaryFromContext(ctx context.Context, extractors *AuthMetadataExtractors) *AuthenticatedUserSummary {
	if extractors == nil {
		return nil
	}

	var err error
	authenticationMetadataRaw := map[string]any{"authenticationMetadata": auth.AuthenticationMetadataFromContext(ctx).GetRaw()}

	var externalID string
	if externalID, err = extractExternalID(extractors.ExternalIDJmespathExpression, authenticationMetadataRaw); err != nil {
		slog.Error("Failed to extract external ID", "err", err)
		return nil
	}

	var displayName *string
	if displayName, err = extractDisplayName(extractors.DisplayNameJmespathExpression, authenticationMetadataRaw); err != nil {
		slog.Error("Failed to extract display name", "err", err)
	}

	var userInfo map[string]any
	if userInfo, err = extractUserInfo(extractors.UserInfoJmespathExpression, authenticationMetadataRaw); err != nil {
		slog.Error("Failed to extract user information", "err", err)
	}

	return &AuthenticatedUserSummary{
		ExternalID:  externalID,
		DisplayName: displayName,
		UserInfo:    userInfo,
	}
}

func extractExternalID(externalIDExtractor *jmespath.Expression, authenticationMetadataRaw map[string]any) (string, error) {
	externalIDRaw, err := externalIDExtractor.Search(authenticationMetadataRaw)
	if err != nil {
		return "", err
	}

	switch t := externalIDRaw.(type) {
	default:
		return "", fmt.Errorf("external ID got an unexpected type: expected string, got %T", externalIDRaw)
	case string:
		return t, nil
	}
}

func extractDisplayName(displayNameExtractor *jmespath.Expression, authenticationMetadataRaw map[string]any) (*string, error) {
	if displayNameExtractor == nil {
		return nil, nil
	}

	displayNameRaw, err := displayNameExtractor.Search(authenticationMetadataRaw)
	if err != nil {
		slog.Error("Failed to extract display name", "err", err)
		return nil, err
	}

	switch t := displayNameRaw.(type) {
	case string:
		return &t, nil
	default:
		return nil, fmt.Errorf("display name got an unexpected type: expected string, got %T", displayNameRaw)
	}
}

func extractUserInfo(userInfoExtractor *jmespath.Expression, authenticationMetadataRaw map[string]any) (map[string]any, error) {
	if userInfoExtractor == nil {
		return nil, nil
	}

	userInfoRaw, err := userInfoExtractor.Search(authenticationMetadataRaw)
	if err != nil {
		return nil, err
	}

	switch t := userInfoRaw.(type) {
	case map[string]any:
		return t, nil
	default:
		return nil, fmt.Errorf("user info got an unexpected type: expected map[string]any, got %T", userInfoRaw)
	}
}
