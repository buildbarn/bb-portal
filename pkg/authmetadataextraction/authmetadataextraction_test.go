package authmetadataextraction_test

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/buildbarn/bb-portal/pkg/authmetadataextraction"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-storage/pkg/auth"
	"github.com/stretchr/testify/require"
)

// TestAuthMetadataExtractorsFromConfiguration Tests creating extractors from a configuration.
func TestAuthMetadataExtractorsFromConfiguration(t *testing.T) {
	externalIDExtractor := authmetadataextraction.ExampleExternalIDExtractor()
	displayNameExtractor := authmetadataextraction.ExampleDisplayNameExtractor()
	userInfoExtractor := authmetadataextraction.ExampleUserInfoExtractor()

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
	extractors := authmetadataextraction.ExampleAuthMetadataExtractors()
	validExternalID := authmetadataextraction.ExampleExternalID()
	validDisplayName := authmetadataextraction.ExampleDisplayName()
	validUserInfo := authmetadataextraction.ExampleUserInfo()
	authMetadata := authmetadataextraction.AuthMetadataFromFields(&validExternalID, &validDisplayName, validUserInfo)

	t.Run("NilExtractors", func(t *testing.T) {
		ctx := auth.NewContextWithAuthenticationMetadata(context.Background(), authMetadata)
		userSummary := authmetadataextraction.AuthenticatedUserSummaryFromContext(ctx, nil)

		require.Nil(t, userSummary)
	})

	t.Run("ValidUserSummary", func(t *testing.T) {
		ctx := auth.NewContextWithAuthenticationMetadata(context.Background(), authMetadata)
		userSummary := authmetadataextraction.AuthenticatedUserSummaryFromContext(ctx, extractors)

		require.Equal(t, userSummary.ExternalID, validExternalID)
		require.Equal(t, *userSummary.DisplayName, validDisplayName)
		require.Equal(t, userSummary.UserInfo, validUserInfo)
	})
}

func TestAuthenticatedUserSummaryFromContext_NilUserSummaryAttributes(t *testing.T) {
	extractors := authmetadataextraction.ExampleAuthMetadataExtractors()
	validExternalID := authmetadataextraction.ExampleExternalID()
	validDisplayName := authmetadataextraction.ExampleDisplayName()
	validUserInfo := authmetadataextraction.ExampleUserInfo()

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

			authMetadata := authmetadataextraction.AuthMetadataFromFields(tt.setup.externalID, tt.setup.displayName, tt.setup.userInfo)
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
