package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/google/shlex"
	"google.golang.org/grpc/credentials"
)

type credentialHelperRequest struct {
	URI string `json:"uri"`
}

type credentialHelperResponse struct {
	Headers map[string][]string `json:"headers"`
}

type credentialsHelper struct {
	command string
}

func NewCredentialsHelper(command string) credentials.PerRPCCredentials {
	return &credentialsHelper{command: command}
}

func (h credentialsHelper) GetRequestMetadata(_ context.Context, uri ...string) (map[string]string, error) {
	resp, err := h.getCredentialsFromHelper(uri[0])
	if err != nil {
		return nil, err
	}
	headers := make(map[string]string, len(resp.Headers))
	for key, values := range resp.Headers {
		headers[key] = strings.Join(values, ",")
	}
	return headers, nil
}

func (h credentialsHelper) getCredentialsFromHelper(uri string) (*credentialHelperResponse, error) {
	reqBytes, err := json.Marshal(credentialHelperRequest{
		URI: uri,
	})
	if err != nil {
		return nil, fmt.Errorf("could not marshal credential helper request: %w", err)
	}
	var out bytes.Buffer
	parts, err := shlex.Split(h.command)
	if err != nil {
		return nil, fmt.Errorf("could not parse command: %w", err)
	}
	cmd := exec.Command(parts[0], parts[1:]...) //nolint:gosec // G204 - Trusted input from local user.
	cmd.Stdin = bytes.NewReader(reqBytes)
	cmd.Stdout = &out
	if err = cmd.Run(); err != nil {
		return nil, fmt.Errorf("error running credential helper: %w", err)
	}
	var resp credentialHelperResponse
	err = json.Unmarshal(out.Bytes(), &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal credential helper response: %w", err)
	}
	return &resp, nil
}

func (h credentialsHelper) RequireTransportSecurity() bool {
	return false
}
