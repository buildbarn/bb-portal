package authmetadataextraction_test

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/buildbarn/bb-portal/pkg/authmetadataextraction"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-storage/pkg/auth"
	jmespath "github.com/buildbarn/bb-storage/pkg/proto/configuration/jmespath"
	"github.com/buildbarn/bb-storage/pkg/util"
	"github.com/stretchr/testify/require"
)

// TestAuthMetadataExtractorsFromConfiguration Tests creating extractors from a configuration.
func TestAuthMetadataExtractorsFromConfiguration(t *testing.T) {
	externalIDExtractor := exampleExternalIDExtractor()
	displayNameExtractor := exampleDisplayNameExtractor()
	userInfoExtractor := exampleUserInfoExtractor()

	t.Run("NilExtractorConfig", func(t *testing.T) {
		var config *bb_portal.AuthMetadataExtractorConfiguration = nil
		extractors, err := authmetadataextraction.AuthMetadataExtractorsFromConfiguration(config, nil)
		require.NoError(t, err)
		require.Nil(t, extractors)
	})

	t.Run("EmptyExtractorConfig", func(t *testing.T) {
		config := &bb_portal.AuthMetadataExtractorConfiguration{}
		extractors, err := authmetadataextraction.AuthMetadataExtractorsFromConfiguration(config, nil)
		require.Error(t, err)
		require.Nil(t, extractors)
	})

	t.Run("FullExtractorConfig", func(t *testing.T) {
		config := &bb_portal.AuthMetadataExtractorConfiguration{
			ExternalIdExtractionJmespathExpression:  &externalIDExtractor,
			DisplayNameExtractionJmespathExpression: &displayNameExtractor,
			UserInfoExtractionJmespathExpression:    &userInfoExtractor,
		}

		extractors, err := authmetadataextraction.AuthMetadataExtractorsFromConfiguration(config, nil)
		require.NoError(t, err)
		require.NotNil(t, extractors.ExternalIDJmespathExpression)
		require.NotNil(t, extractors.DisplayNameJmespathExpression)
		require.NotNil(t, extractors.UserInfoJmespathExpression)
	})

	t.Run("ExternalId", func(t *testing.T) {
		config := &bb_portal.AuthMetadataExtractorConfiguration{
			ExternalIdExtractionJmespathExpression:  &externalIDExtractor,
			DisplayNameExtractionJmespathExpression: nil,
			UserInfoExtractionJmespathExpression:    nil,
		}

		extractors, err := authmetadataextraction.AuthMetadataExtractorsFromConfiguration(config, nil)
		require.NoError(t, err)
		require.NotNil(t, extractors.ExternalIDJmespathExpression)
		require.Nil(t, extractors.DisplayNameJmespathExpression)
		require.Nil(t, extractors.UserInfoJmespathExpression)
	})

	t.Run("ExternalIdDisplayName", func(t *testing.T) {
		config := &bb_portal.AuthMetadataExtractorConfiguration{
			ExternalIdExtractionJmespathExpression:  &externalIDExtractor,
			DisplayNameExtractionJmespathExpression: &displayNameExtractor,
			UserInfoExtractionJmespathExpression:    nil,
		}

		extractors, err := authmetadataextraction.AuthMetadataExtractorsFromConfiguration(config, nil)
		require.NoError(t, err)
		require.NotNil(t, extractors.ExternalIDJmespathExpression)
		require.NotNil(t, extractors.DisplayNameJmespathExpression)
		require.Nil(t, extractors.UserInfoJmespathExpression)
	})

	t.Run("ExternalIdUserInfo", func(t *testing.T) {
		config := &bb_portal.AuthMetadataExtractorConfiguration{
			ExternalIdExtractionJmespathExpression:  &externalIDExtractor,
			DisplayNameExtractionJmespathExpression: nil,
			UserInfoExtractionJmespathExpression:    &userInfoExtractor,
		}

		extractors, err := authmetadataextraction.AuthMetadataExtractorsFromConfiguration(config, nil)
		require.NoError(t, err)
		require.NotNil(t, extractors.ExternalIDJmespathExpression)
		require.Nil(t, extractors.DisplayNameJmespathExpression)
		require.NotNil(t, extractors.UserInfoJmespathExpression)
	})
}

// AuthenticatedUserSummaryFromContext Tests user summary creation from authentication metadata in the context.
func TestAuthenticatedUserSummaryFromContext_ValidUserSummaryAttributes(t *testing.T) {
	extractors := exampleAuthMetadataExtractors()
	validExternalID := exampleExternalID()
	validDisplayName := exampleDisplayName()
	validUserInfo := exampleUserInfo()

	t.Run("NilExtractors", func(t *testing.T) {
		authMetadata, err := authMetadataFromFields(&validExternalID, &validDisplayName, validUserInfo)
		require.NoError(t, err)
		ctx := auth.NewContextWithAuthenticationMetadata(context.Background(), authMetadata)
		userSummary := authmetadataextraction.AuthenticatedUserSummaryFromContext(ctx, nil)

		require.Nil(t, userSummary)
	})

	t.Run("ValidUserSummary", func(t *testing.T) {
		authMetadata := util.Must(authMetadataFromFields(&validExternalID, &validDisplayName, validUserInfo))
		ctx := auth.NewContextWithAuthenticationMetadata(context.Background(), authMetadata)
		userSummary := authmetadataextraction.AuthenticatedUserSummaryFromContext(ctx, extractors)

		require.Equal(t, userSummary.ExternalID, validExternalID)
		require.Equal(t, *userSummary.DisplayName, validDisplayName)
		require.Equal(t, userSummary.UserInfo, validUserInfo)
	})
}

func TestAuthenticatedUserSummaryFromContext_NilUserSummaryAttributes(t *testing.T) {
	extractors := exampleAuthMetadataExtractors()
	validExternalID := exampleExternalID()
	validDisplayName := exampleDisplayName()
	validUserInfo := exampleUserInfo()

	type setup struct {
		externalID  *string
		displayName *string
		userInfo    map[string]any
	}

	tests := []struct {
		name          string
		setup         setup
		desiredOutput string
	}{
		{
			name: "NilExternalId",
			setup: setup{
				externalID:  nil,
				displayName: &validDisplayName,
				userInfo:    validUserInfo,
			},
			desiredOutput: "external ID got an unexpected type: expected string, got <nil>",
		},
		{
			name: "NilDisplayName",
			setup: setup{
				externalID:  &validExternalID,
				displayName: nil,
				userInfo:    validUserInfo,
			},
			desiredOutput: "display name got an unexpected type: expected string, got <nil>",
		},
		{
			name: "NilUserInfo",
			setup: setup{
				externalID:  &validExternalID,
				displayName: &validDisplayName,
				userInfo:    nil,
			},
			desiredOutput: "user info got an unexpected type: expected map[string]any, got <nil>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture error messages
			originalLogger := slog.Default()
			var buf bytes.Buffer
			testLogger := slog.New(slog.NewTextHandler(&buf, nil))
			slog.SetDefault(testLogger)
			defer slog.SetDefault(originalLogger)

			authMetadata, err := authMetadataFromFields(tt.setup.externalID, tt.setup.displayName, tt.setup.userInfo)
			require.NoError(t, err)
			ctx := auth.NewContextWithAuthenticationMetadata(context.Background(), authMetadata)
			userSummary := authmetadataextraction.AuthenticatedUserSummaryFromContext(ctx, extractors)

			require.Contains(t, buf.String(), tt.desiredOutput)
			if tt.setup.externalID == nil {
				require.Nil(t, userSummary)
			} else {
				require.Equal(t, userSummary.ExternalID, *tt.setup.externalID)

				if tt.setup.displayName == nil {
					require.Nil(t, userSummary.DisplayName)
				} else {
					require.Equal(t, *userSummary.DisplayName, *tt.setup.displayName)
				}

				if tt.setup.userInfo == nil {
					require.Nil(t, userSummary.UserInfo)
				} else {
					require.Equal(t, userSummary.UserInfo, tt.setup.userInfo)
				}
			}
		})
	}
}

func exampleExternalIDExtractor() jmespath.Expression {
	return jmespath.Expression{Expression: "authenticationMetadata.private.external_id"}
}

func exampleDisplayNameExtractor() jmespath.Expression {
	return jmespath.Expression{Expression: "authenticationMetadata.private.display_name"}
}

func exampleUserInfoExtractor() jmespath.Expression {
	return jmespath.Expression{Expression: "authenticationMetadata.public"}
}

func exampleAuthMetadataExtractors() *authmetadataextraction.AuthMetadataExtractors {
	externalIDExtractor := exampleExternalIDExtractor()
	displayNameExtractor := exampleDisplayNameExtractor()
	userInfoExtractor := exampleUserInfoExtractor()

	config := bb_portal.AuthMetadataExtractorConfiguration{
		ExternalIdExtractionJmespathExpression:  &externalIDExtractor,
		DisplayNameExtractionJmespathExpression: &displayNameExtractor,
		UserInfoExtractionJmespathExpression:    &userInfoExtractor,
	}

	return util.Must(authmetadataextraction.AuthMetadataExtractorsFromConfiguration(&config, nil))
}

func exampleExternalID() string {
	return "12345678-1234-1234-1234-123456789abc"
}

func exampleDisplayName() string {
	return "example_username"
}

func exampleUserInfo() map[string]any {
	return map[string]any{
		"full_name": "Example Name",
		"age":       float64(30),
		"contact_information": map[string]any{
			"email":        "user@example.com",
			"phone_number": "800-555-0199",
		},
	}
}

func authMetadataFromFields(externalID, displayName *string, userInfo map[string]any) (*auth.AuthenticationMetadata, error) {
	var private map[string]any
	var public map[string]any
	if externalID == nil && displayName == nil {
		private = nil
	} else {
		private = map[string]any{}
		if externalID != nil {
			private["external_id"] = *externalID
		}
		if displayName != nil {
			private["display_name"] = *displayName
		}
	}

	if userInfo != nil {
		public = userInfo
	}

	authenticationMetadata := map[string]any{
		"private": private,
		"public":  public,
	}

	return auth.NewAuthenticationMetadataFromRaw(authenticationMetadata)
}
