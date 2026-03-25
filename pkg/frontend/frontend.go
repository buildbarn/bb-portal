package frontend

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	"github.com/buildbarn/bb-portal/pkg/proto/configuration/bb_portal"
	"github.com/buildbarn/bb-portal/pkg/proto/configuration/frontend"
	"github.com/buildbarn/bb-storage/pkg/util"
	"github.com/gorilla/mux"
	"google.golang.org/protobuf/encoding/protojson"
)

// ServeFrontend serves the frontend.
//
// NOTE: This needs to be the last handler registered on the router, as it has
// a catch-all route.
func ServeFrontend(configuration *bb_portal.FrontendService, router *mux.Router) error {
	// Return 404 for all API requests not already handled.
	router.PathPrefix("/api/").Handler(router.NotFoundHandler)

	if configuration == nil {
		return fmt.Errorf("No frontend service configuration found")
	}
	if configuration.FrontendSource == nil {
		return fmt.Errorf("No frontend source configured")
	}

	switch s := configuration.FrontendSource.Source.(type) {
	case *bb_portal.FrontendService_FrontendSource_Proxy:
		return setupProxyHandler(router, s, configuration.FrontendConfig)
	case *bb_portal.FrontendService_FrontendSource_Embedded:
		return setupEmbeddedHandler(router, configuration.FrontendConfig)
	default:
		return fmt.Errorf("Unknown frontend source type: %T", s)
	}
}

func validateFrontendConfig(frontendConfig *frontend.PortalFrontendConfiguration) error {
	if frontendConfig == nil {
		return fmt.Errorf("No frontend configuration found")
	}
	if frontendConfig.FeatureFlags == nil {
		return fmt.Errorf("No feature flags found")
	}
	return nil
}

func injectFrontendConfigScript(html []byte, frontendConfig *frontend.PortalFrontendConfiguration) ([]byte, error) {
	if frontendConfig == nil {
		return nil, fmt.Errorf("No frontend configuration found")
	}
	marshalOptions := protojson.MarshalOptions{
		EmitUnpopulated: true,
	}
	configProto, err := marshalOptions.Marshal(frontendConfig)
	if err != nil {
		return nil, util.StatusWrap(err, "Failed to marshal frontend configuration")
	}
	configScript := fmt.Appendf(nil, "<script>window.__env__ = %s</script>", string(configProto))
	return bytes.ReplaceAll(html, []byte("<!-- BB_PORTAL_CONFIGURATION_PLACEHOLDER -->"), configScript), nil
}

func setupProxyHandler(router *mux.Router, sourceConfig *bb_portal.FrontendService_FrontendSource_Proxy, frontendConfig *frontend.PortalFrontendConfiguration) error {
	if sourceConfig == nil {
		return fmt.Errorf("Frontend Proxy configuration is empty")
	}
	if sourceConfig.Proxy == "" {
		return fmt.Errorf("Frontend Proxy URL is empty")
	}
	if err := validateFrontendConfig(frontendConfig); err != nil {
		return util.StatusWrap(err, "Error validating frontend config")
	}

	remote, err := url.Parse(sourceConfig.Proxy)
	if err != nil {
		return util.StatusWrap(err, "Could not parse proxy url")
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ModifyResponse = func(r *http.Response) error {
		if r.Header.Get("Content-Type") != "text/html" {
			return nil
		}

		oldIndexContent, err := io.ReadAll(r.Body)
		if err != nil {
			return err
		}
		defer r.Body.Close()

		newIndexContent, err := injectFrontendConfigScript(oldIndexContent, frontendConfig)
		if err != nil {
			return err
		}

		r.Body = io.NopCloser(bytes.NewReader(newIndexContent))
		r.Header.Set("Content-Length", strconv.Itoa(len(newIndexContent)))
		return nil
	}

	router.PathPrefix("/").Handler(proxy)
	return nil
}
