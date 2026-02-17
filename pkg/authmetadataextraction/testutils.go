package authmetadataextraction

import (
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-storage/pkg/auth"
	jmespath "github.com/buildbarn/bb-storage/pkg/proto/configuration/jmespath"
	"github.com/buildbarn/bb-storage/pkg/util"
)

// ExampleAuthMetadataExtractors returns example extractors used for testing.
func ExampleAuthMetadataExtractors() *AuthMetadataExtractors {
	externalIDExtractor := ExampleExternalIDExtractor()
	displayNameExtractor := ExampleDisplayNameExtractor()
	userInfoExtractor := ExampleUserInfoExtractor()

	config := bb_portal.AuthMetadataExtractorConfiguration{
		ExternalIdExtractionJmespathExpression:  &externalIDExtractor,
		DisplayNameExtractionJmespathExpression: &displayNameExtractor,
		UserInfoExtractionJmespathExpression:    &userInfoExtractor,
	}

	return util.Must(AuthMetadataExtractorsFromConfiguration(&config, nil))
}

// AuthMetadataFromFields returns authentication metadata
// from fields used in testing.
func AuthMetadataFromFields(externalID, displayName *string, userInfo map[string]any) *auth.AuthenticationMetadata {
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

	return util.Must(auth.NewAuthenticationMetadataFromRaw(authenticationMetadata))
}

// ExampleExternalIDExtractor returns an example external ID
// extractor used for testing.
func ExampleExternalIDExtractor() jmespath.Expression {
	return jmespath.Expression{Expression: "authenticationMetadata.private.external_id"}
}

// ExampleDisplayNameExtractor returns an example display name
// extractor used for testing.
func ExampleDisplayNameExtractor() jmespath.Expression {
	return jmespath.Expression{Expression: "authenticationMetadata.private.display_name"}
}

// ExampleUserInfoExtractor returns an example user info
// extractor used for testing.
func ExampleUserInfoExtractor() jmespath.Expression {
	return jmespath.Expression{Expression: "authenticationMetadata.public"}
}

// ExampleExternalID returns an example external ID
// used for testing.
func ExampleExternalID() string {
	return "12345678-1234-1234-1234-123456789abc"
}

// ExampleDisplayName returns an example display name
// used for testing.
func ExampleDisplayName() string {
	return "example_username"
}

// ExampleUserInfo returns an example user info
// map used for testing.
func ExampleUserInfo() map[string]any {
	return map[string]any{
		"full_name": "Example Name",
		"age":       float64(30),
		"contact_information": map[string]any{
			"email":        "user@example.com",
			"phone_number": "800-555-0199",
		},
	}
}
