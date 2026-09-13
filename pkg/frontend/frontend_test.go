package frontend

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	frontend_configuration "github.com/buildbarn/bb-portal/pkg/proto/configuration/frontend"
	"github.com/gorilla/mux"
)

func TestProxyUpdatesContentLengthAfterInjectingConfiguration(t *testing.T) {
	vite := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = io.WriteString(w, "<html><!-- BB_PORTAL_CONFIGURATION_PLACEHOLDER --></html>")
	}))
	defer vite.Close()

	router := mux.NewRouter()
	err := setupProxyHandler(
		router,
		&bb_portal.FrontendService_FrontendSource_Proxy{Proxy: vite.URL},
		&frontend_configuration.PortalFrontendConfiguration{
			FeatureFlags: &frontend_configuration.PortalFrontendConfiguration_FeatureFlags{},
		},
	)
	if err != nil {
		t.Fatalf("Failed to configure frontend proxy: %v", err)
	}
	portal := httptest.NewServer(router)
	defer portal.Close()

	response, err := portal.Client().Get(portal.URL)
	if err != nil {
		t.Fatalf("Failed to fetch proxied frontend: %v", err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Failed to read proxied frontend: %v", err)
	}
	if response.ContentLength != int64(len(body)) {
		t.Fatalf("Content-Length is %d, while body size is %d", response.ContentLength, len(body))
	}
	if !strings.Contains(string(body), "window.__env__") {
		t.Fatalf("Frontend configuration was not injected: %s", body)
	}
}
