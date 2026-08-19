package frontend

import (
	"strings"
	"testing"

	"github.com/buildbarn/bb-portal/pkg/proto/configuration/frontend"
	"google.golang.org/protobuf/types/known/emptypb"
)

// The web frontend reads window.__env__ directly as a plain JS object,
// without running it through a protobuf JSON decoder (see
// frontend/src/utils/env.ts). It relies on protojson's documented
// behavior of emitting unset message-type fields as JSON `null` rather
// than omitting the key. Frontend feature-flag checks therefore must
// treat `null` the same as "unset", not just an absent key.
func TestInjectFrontendConfigScriptOmitsInstanceNameColumnsByDefault(t *testing.T) {
	html, err := injectFrontendConfigScript(
		[]byte("<!-- BB_PORTAL_CONFIGURATION_PLACEHOLDER -->"),
		&frontend.PortalFrontendConfiguration{
			FeatureFlags: &frontend.PortalFrontendConfiguration_FeatureFlags{
				Bes: &frontend.PortalFrontendConfiguration_FeatureFlags_BesFeatureFlags{},
			},
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, key := range []string{"columnInstanceNameBuilds", "columnInstanceNameInvocations"} {
		if !strings.Contains(string(html), `"`+key+`":null`) {
			t.Errorf("expected %q to be serialized as null when unset, got: %s", key, html)
		}
	}
}

func TestInjectFrontendConfigScriptEnablesInstanceNameColumns(t *testing.T) {
	html, err := injectFrontendConfigScript(
		[]byte("<!-- BB_PORTAL_CONFIGURATION_PLACEHOLDER -->"),
		&frontend.PortalFrontendConfiguration{
			FeatureFlags: &frontend.PortalFrontendConfiguration_FeatureFlags{
				Bes: &frontend.PortalFrontendConfiguration_FeatureFlags_BesFeatureFlags{
					ColumnInstanceNameBuilds:      &emptypb.Empty{},
					ColumnInstanceNameInvocations: &emptypb.Empty{},
				},
			},
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, key := range []string{"columnInstanceNameBuilds", "columnInstanceNameInvocations"} {
		if !strings.Contains(string(html), `"`+key+`":{}`) {
			t.Errorf("expected %q to be serialized as {} when enabled, got: %s", key, html)
		}
	}
}
