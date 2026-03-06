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

	if configuration.GetFrontendSource().GetSource() == nil {
		return fmt.Errorf("No frontend source configured")
	}
	if configuration.GetFrontendConfig() == nil {
		return fmt.Errorf("No frontend configuration found")
	}

	switch s := configuration.FrontendSource.Source.(type) {
	case *bb_portal.FrontendService_FrontendSource_Proxy:
		return frontendProxy(configuration.FrontendConfig, router, s.Proxy)
	case *bb_portal.FrontendService_FrontendSource_Embedded:
		panic("Embedded not implemented yet")
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

func frontendProxy(configuration *frontend.PortalFrontendConfiguration, router *mux.Router, proxyUrl string) error {
	remote, err := url.Parse(proxyUrl)
	if err != nil {
		return util.StatusWrap(err, "Could not parse proxy url")
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)

	if err := validateFrontendConfig(configuration); err != nil {
		return util.StatusWrap(err, "Error validating frontend config")
	}

	marshalOptions := protojson.MarshalOptions{
		EmitUnpopulated: true,
	}
	protoConfig, err := marshalOptions.Marshal(configuration)

	s := fmt.Sprintf("<script>window.__env__ = %s</script>", string(protoConfig))

	// Intercept and modify the response
	proxy.ModifyResponse = func(r *http.Response) error {
		if r.Header.Get("Content-Type") != "text/html" {
			return nil
		}

		// Read the original body
		oldBody, err := io.ReadAll(r.Body)
		if err != nil {
			return err
		}
		defer r.Body.Close()

		// Perform the string replacement
		newBody := bytes.ReplaceAll(oldBody, []byte("<!-- BB_PORTAL_CONFIGURATION_PLACEHOLDER -->"), []byte(s))

		// Update the body and the Content-Length header
		r.Body = io.NopCloser(bytes.NewReader(newBody))
		r.Header.Set("Content-Length", strconv.Itoa(len(newBody)))
		return nil
	}

	router.PathPrefix("/").Handler(proxy)
	return nil
}
