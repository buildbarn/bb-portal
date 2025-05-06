package common

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/buildbarn/bb-storage/pkg/auth"
	auth_pb "github.com/buildbarn/bb-storage/pkg/proto/auth"
	"github.com/jmespath/go-jmespath"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestExtractContextFromRequest(t *testing.T) {
	t.Run("EmptyHeaders", func(t *testing.T) {
		req := &http.Request{
			Header: http.Header{},
		}
		ctx := ExtractContextFromRequest(req)

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok || len(md) != 0 {
			t.Error("Expected empty metadata for request with no headers")
		}
	})

	t.Run("SingleHeader", func(t *testing.T) {
		req := &http.Request{
			Header: http.Header{
				"X-Test-Header": []string{"value1"},
			},
		}
		ctx := ExtractContextFromRequest(req)

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			t.Fatal("Expected metadata to be present in context")
		}
		expected := metadata.Pairs("X-Test-Header", "value1")
		if !reflect.DeepEqual(md, expected) {
			t.Errorf("Expected metadata %v, got %v", expected, md)
		}
	})

	t.Run("MultipleHeaders", func(t *testing.T) {
		req := &http.Request{
			Header: http.Header{
				"X-Test-Header":    []string{"value1", "value2"},
				"X-Another-Header": []string{"value3"},
			},
		}
		ctx := ExtractContextFromRequest(req)

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			t.Fatal("Expected metadata to be present in context")
		}
		expected := metadata.Pairs(
			"X-Test-Header", "value1",
			"X-Test-Header", "value2",
			"X-Another-Header", "value3",
		)
		if !reflect.DeepEqual(md, expected) {
			t.Errorf("Expected metadata %v, got %v", expected, md)
		}
	})
}

func TestIsInstanceNameAllowed(t *testing.T) {
	a := auth.NewJMESPathExpressionAuthorizer(
		jmespath.MustCompile("contains(authenticationMetadata.private.permittedInstanceNames, instanceName) || instanceName == ''"),
	)

	ctx := auth.NewContextWithAuthenticationMetadata(context.Background(), auth.MustNewAuthenticationMetadataFromProto(&auth_pb.AuthenticationMetadata{
		Private: structpb.NewStructValue(&structpb.Struct{
			Fields: map[string]*structpb.Value{
				"permittedInstanceNames": structpb.NewListValue(&structpb.ListValue{
					Values: []*structpb.Value{
						structpb.NewStringValue("allowed"),
						structpb.NewStringValue("alsoAllowed"),
					},
				}),
			},
		}),
	}))

	t.Run("ValidInstanceNames", func(t *testing.T) {
		if !IsInstanceNameAllowed(ctx, a, "") {
			t.Error("Expected empty instance name to be allowed")
		}
		if !IsInstanceNameAllowed(ctx, a, "allowed") {
			t.Error("Expected instance name to be allowed")
		}
		if !IsInstanceNameAllowed(ctx, a, "alsoAllowed") {
			t.Error("Expected instance name to be allowed")
		}
	})

	t.Run("InvalidInstanceNames", func(t *testing.T) {
		if IsInstanceNameAllowed(ctx, a, "forbidden") {
			t.Error("Expected instance name to be forbidden")
		}
		if IsInstanceNameAllowed(ctx, a, "allowed/") {
			t.Error("Expected instance name to be forbidden")
		}
	})
}
