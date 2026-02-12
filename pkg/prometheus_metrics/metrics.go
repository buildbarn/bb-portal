package prometheusmetrics

import (
	"context"
	"log/slog"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/internal/database/dbauthservice"
	"github.com/buildbarn/bb-storage/pkg/util"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// UnauthenticatedUsersLabel The label for unauthenticated users in Invocations
	UnauthenticatedUsersLabel = "Unauthenticated"

	// AuthenticatedUsersLabel The label for authenticated users in Invocations
	AuthenticatedUsersLabel = "Authenticated"
)

var (
	// Invocations A gauge for invocations
	Invocations = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "buildbarn",
		Subsystem: "portal",
		Name:      "invocations",
		Help:      "Number of Invocations by authentication status",
	}, []string{"AuthStatus"})

	// AuthenticatedUsersCount A gauge for the number of authenticated users
	AuthenticatedUsersCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "buildbarn",
		Subsystem: "portal",
		Name:      "authenticated_user_count",
		Help:      "Number of Unique Authenticated Users",
	})
)

// RegisterMetrics registers the metrics with Prometheus
func init() {
	prometheus.MustRegister(Invocations)
	prometheus.MustRegister(AuthenticatedUsersCount)
}

// SyncMetrics synchronizes the Prometheus service with
// data from the database.
func SyncMetrics(db *ent.Client) {
	ctx := dbauthservice.NewContextWithDbAuthServiceBypass(context.Background())

	err := SyncInvocations(ctx, db)
	if err != nil {
		slog.Error("Failed to synchronize Prometheus metric Invocations", "err", err)
	}

	err = SyncAuthenticatedUsersCount(ctx, db)
	if err != nil {
		slog.Error("Failed to synchronize Prometheus metric AuthenticatedUsersCount", "err", err)
	}
}

// SyncInvocations synchronizes the Prometheus metric Invocations
// with invocations from the database.
func SyncInvocations(ctx context.Context, db *ent.Client) error {
	authenticatedInvocationsCount, err := db.BazelInvocation.Query().
		Where(bazelinvocation.HasAuthenticatedUser()).
		Count(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to count invocations with authenticated users")
	}
	Invocations.WithLabelValues(AuthenticatedUsersLabel).Set(float64(authenticatedInvocationsCount))

	unauthenticatedInvocationsCount, err := db.BazelInvocation.Query().
		Where(bazelinvocation.Not(bazelinvocation.HasAuthenticatedUser())).
		Count(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to count invocations with unauthenticated users")
	}
	Invocations.WithLabelValues(UnauthenticatedUsersLabel).Set(float64(unauthenticatedInvocationsCount))

	return nil
}

// SyncAuthenticatedUsersCount synchronizes the Prometheus metric
// AuthenticatedUsersCount with the users from the database.
func SyncAuthenticatedUsersCount(ctx context.Context, db *ent.Client) error {
	userCount, err := db.AuthenticatedUser.Query().Count(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to count authenticated users")
	}
	AuthenticatedUsersCount.Set(float64(userCount))
	return nil
}
