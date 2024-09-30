package uuidgql

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
)

// errExpectedString
var errExpectedString = errors.New("expected string")

// MarshalUUID gets the id
func MarshalUUID(u uuid.UUID) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, strconv.Quote(u.String()))
	})
}

// UnmarshalUUID unmarsalls the id
func UnmarshalUUID(v interface{}) (u uuid.UUID, err error) {
	s, ok := v.(string)
	if !ok {
		return u, fmt.Errorf("invalid type %T: %w", v, errExpectedString)
	}
	return uuid.Parse(s)
}
