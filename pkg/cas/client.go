package cas

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"github.com/bazelbuild/remote-apis-sdks/go/pkg/client"
	"github.com/bazelbuild/remote-apis-sdks/go/pkg/digest"

	"github.com/buildbarn/bb-portal/pkg/auth"
)

type ConnectionManager struct {
	params ManagerParams
}

type ManagerParams struct {
	// TLSCACertFile is the PEM file that contains TLS root certificates.
	TLSCACertFile string

	// CredentialsHelperCommand is a command to use as a Bazel credentials helper.
	CredentialsHelperCommand string
}

func NewConnectionManager(params ManagerParams) *ConnectionManager {
	return &ConnectionManager{
		params: params,
	}
}

type Client struct {
	client *client.Client
}

func (manager *ConnectionManager) GetClientForURI(ctx context.Context, uri *url.URL) (*Client, error) {
	instanceName := strings.Split(strings.TrimPrefix(uri.Path, "/"), "/")[0]
	dialParms := client.DialParams{
		Service:       uri.Hostname() + ":443",
		TLSCACertFile: manager.params.TLSCACertFile,
	}
	if manager.params.CredentialsHelperCommand != "" {
		credentialsHelper := auth.NewCredentialsHelper(manager.params.CredentialsHelperCommand)
		dialParms.UseExternalAuthToken = true
		dialParms.ExternalPerRPCCreds = &client.PerRPCCreds{Creds: credentialsHelper}
	}
	remoteAPIsClient, err := client.NewClient(ctx, instanceName, dialParms)
	if err != nil {
		return nil, fmt.Errorf("could not open a remote APIs SDK client: %w", err)
	}
	return &Client{
		client: remoteAPIsClient,
	}, nil
}

func (c *Client) ReadBlobToFile(ctx context.Context, uri *url.URL, fpath string) error {
	pathParts := strings.Split(uri.Path, "/")
	digestPath := strings.Join(pathParts[len(pathParts)-2:], "/")
	d, err := digest.NewFromString(digestPath)
	if err != nil {
		return fmt.Errorf("could not create digest from path %s: %w", uri.Path, err)
	}

	slog.InfoContext(ctx, "reading from CAS", "digest", d.String(), "file", fpath)

	_, err = c.client.ReadBlobToFile(ctx, d, fpath)
	if err != nil {
		return fmt.Errorf("could not read blob at %s to file %s: %w", uri.String(), fpath, err)
	}

	return nil
}

func (c *Client) Close() error {
	return c.client.Close()
}
