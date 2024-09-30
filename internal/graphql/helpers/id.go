package helpers

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

// graphQLIDFromString graphQLIDFromString converts a string ID to a base64 encoded string for use as a GraphQL ID.
func graphQLIDFromString(input string) string {
	return base64.URLEncoding.EncodeToString([]byte(input))
}

// GraphQLIDFromTypeAndID Takes an id and returns the b64enc string.
func GraphQLIDFromTypeAndID(objType string, id int) string {
	return graphQLIDFromString(fmt.Sprintf("%s:%d", objType, id))
}

// GraphQLTypeAndIntIDFromID ID Decoder helper
func GraphQLTypeAndIntIDFromID(id string) (string, int, error) {
	bytes, err := base64.URLEncoding.DecodeString(id)
	if err != nil {
		return "", 0, fmt.Errorf("could not decode ID (id: %s): %w", id, err)
	}
	s := string(bytes)
	parts := strings.SplitN(s, ":", 2)
	intID, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, fmt.Errorf("could extract int ID (id: %s): %w", id, err)
	}
	return parts[0], intID, nil
}
