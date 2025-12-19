package dbauthservice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/buildbarn/bb-portal/internal/database/dbauthservice"
	"github.com/buildbarn/bb-storage/pkg/auth"
	"github.com/buildbarn/bb-storage/pkg/clock"
	"github.com/buildbarn/bb-storage/pkg/jmespath"
	"github.com/buildbarn/bb-storage/pkg/util"
	"github.com/google/uuid"

	// Needed to avoid cyclic dependencies in ent (https://entgo.io/docs/privacy#privacy-policy-registration)
	_ "github.com/buildbarn/bb-portal/ent/gen/ent/runtime"
	"github.com/stretchr/testify/require"
)

func BenchmarkDbAuthService_GetAuthorizedInstanceNames(b *testing.B) {
	testCases := []int{
		1,
		10,
		100,
		1_000,
		10_000,
		32_766,
	}

	adminCtx := dbauthservice.NewContextWithDbAuthServiceBypass(context.Background())
	authCtx := auth.NewContextWithAuthenticationMetadata(context.Background(),
		util.Must(auth.NewAuthenticationMetadataFromRaw(map[string]any{
			"private": map[string]any{
				"roles": []string{"admin"},
			},
		})),
	)
	db := setupTestDB(b).Ent()
	authorizer := auth.NewJMESPathExpressionAuthorizer(
		jmespath.MustCompile("contains(authenticationMetadata.private.roles, 'admin')"),
	)

	for _, numInstanceNames := range testCases {
		b.Run(fmt.Sprintf("NumInstanceNames=%d", numInstanceNames), func(b *testing.B) {
			db.BazelInvocation.Delete().Exec(adminCtx)
			db.InstanceName.Delete().Exec(adminCtx)
			for i := 0; i < numInstanceNames; i++ {
				instanceName, err := db.InstanceName.Create().
					SetName(fmt.Sprintf("InstanceName-%d", i)).
					Save(adminCtx)
				require.NoError(b, err)
				err = db.BazelInvocation.Create().
					SetInstanceNameID(instanceName.ID).
					SetInvocationID(uuid.New()).
					Exec(adminCtx)
				require.NoError(b, err)
			}

			authService := dbauthservice.NewDbAuthService(db, clock.SystemClock, authorizer, 0)
			authServiceCtx := dbauthservice.NewContextWithDbAuthService(authCtx, authService)

			for b.Loop() {
				count, err := db.BazelInvocation.Query().Count(authServiceCtx)
				require.NoError(b, err)
				require.Equal(b, numInstanceNames, count)
			}
			db.BazelInvocation.Delete().Exec(adminCtx)
			db.InstanceName.Delete().Exec(adminCtx)
		})
	}
}
