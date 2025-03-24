package grpcweb

import (
	"context"
	"testing"

	"github.com/buildbarn/bb-storage/pkg/auth"
	auth_pb "github.com/buildbarn/bb-storage/pkg/proto/auth"
	"github.com/jmespath/go-jmespath"
	"google.golang.org/protobuf/types/known/structpb"
)

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
