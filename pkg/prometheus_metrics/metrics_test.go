package prometheusmetrics_test

import (
	"context"
	"os"
	"testing"

	"github.com/buildbarn/bb-portal/internal/database/buildeventrecorder"
	"github.com/buildbarn/bb-portal/internal/database/embedded"
	"github.com/buildbarn/bb-portal/pkg/authmetadataextraction"
	prometheusmetrics "github.com/buildbarn/bb-portal/pkg/prometheus_metrics"
	"github.com/buildbarn/bb-portal/test/testutils"
	"github.com/buildbarn/bb-storage/pkg/auth"
	"github.com/buildbarn/bb-storage/pkg/util"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/require"
)

var dbProvider *embedded.DatabaseProvider

func TestMetrics(t *testing.T) {
	ctx := context.Background()
	dbProvider := util.Must(embedded.NewDatabaseProvider(os.Stderr))
	externalID := authmetadataextraction.ExampleExternalID()
	instanceName := "exampleInstanceName"
	displayName := authmetadataextraction.ExampleDisplayName()
	userInfo := authmetadataextraction.ExampleUserInfo()
	authMetadata := authmetadataextraction.AuthMetadataFromFields(&externalID, &displayName, userInfo)
	authMetadataCtx := auth.NewContextWithAuthenticationMetadata(
		ctx,
		authMetadata,
	)

	t.Run("NoAuthMetadata", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		newAuthenticatedUsersGauge := newAuthenticatedUsersGauge()
		_, err := buildeventrecorder.FindOrCreateAuthenticatedUser(
			ctx,
			db,
			nil,
			newAuthenticatedUsersGauge,
		)
		require.NoError(t, err)

		authenticatedUsersCount := int(testutil.ToFloat64(newAuthenticatedUsersGauge))
		require.Equal(t, authenticatedUsersCount, 0)
	})

	t.Run("AuthenticatedUser", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		newAuthenticatedUsersGauge := newAuthenticatedUsersGauge()
		_, err := buildeventrecorder.FindOrCreateAuthenticatedUser(
			authMetadataCtx,
			db,
			authmetadataextraction.ExampleAuthMetadataExtractors(),
			newAuthenticatedUsersGauge,
		)
		require.NoError(t, err)

		authenticatedUsersCount := int(testutil.ToFloat64(newAuthenticatedUsersGauge))
		require.Equal(t, authenticatedUsersCount, 1)
	})

	t.Run("UnathenticatedInvocation", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		invocationsGaugeVec := newInvocationsGaugeVec()
		_, _, err := buildeventrecorder.FindOrCreateInvocation(
			ctx,
			db,
			uuid.NewString(),
			instanceName,
			nil,
			nil,
			invocationsGaugeVec,
		)
		require.NoError(t, err)

		unauthenticatedInvocationsCount := getInvocationCount(invocationsGaugeVec, prometheusmetrics.UnauthenticatedUsersLabel)
		authenticatedInvocationsCount := getInvocationCount(invocationsGaugeVec, prometheusmetrics.AuthenticatedUsersLabel)
		require.Equal(t, unauthenticatedInvocationsCount, 1)
		require.Equal(t, authenticatedInvocationsCount, 0)
	})

	t.Run("AuthenticatedInvocation", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		invocationsGaugeVec := newInvocationsGaugeVec()
		authenticatedUsersGauge := newAuthenticatedUsersGauge()
		userDbID, err := buildeventrecorder.FindOrCreateAuthenticatedUser(
			authMetadataCtx,
			db,
			authmetadataextraction.ExampleAuthMetadataExtractors(),
			authenticatedUsersGauge,
		)
		require.NoError(t, err)

		_, _, err = buildeventrecorder.FindOrCreateInvocation(
			authMetadataCtx,
			db,
			uuid.NewString(),
			instanceName,
			nil,
			userDbID,
			invocationsGaugeVec,
		)
		require.NoError(t, err)

		unauthenticatedInvocationsCount := getInvocationCount(invocationsGaugeVec, prometheusmetrics.UnauthenticatedUsersLabel)
		authenticatedInvocationsCount := getInvocationCount(invocationsGaugeVec, prometheusmetrics.AuthenticatedUsersLabel)
		require.Equal(t, unauthenticatedInvocationsCount, 0)
		require.Equal(t, authenticatedInvocationsCount, 1)
	})

	t.Run("AuthenticatedAndUnauthenticatedInvocations", func(t *testing.T) {
		db := testutils.SetupTestDB(t, dbProvider)
		invocationsGaugeVec := newInvocationsGaugeVec()
		authenticatedUsersGauge := newAuthenticatedUsersGauge()
		userDbID, err := buildeventrecorder.FindOrCreateAuthenticatedUser(
			authMetadataCtx,
			db,
			authmetadataextraction.ExampleAuthMetadataExtractors(),
			authenticatedUsersGauge,
		)
		require.NoError(t, err)

		// Authenticated invocation
		_, _, err = buildeventrecorder.FindOrCreateInvocation(
			authMetadataCtx,
			db,
			uuid.NewString(),
			instanceName,
			nil,
			userDbID,
			invocationsGaugeVec,
		)
		require.NoError(t, err)

		// Unauthenticated invocation
		_, _, err = buildeventrecorder.FindOrCreateInvocation(
			ctx,
			db,
			uuid.NewString(),
			instanceName,
			nil,
			nil,
			invocationsGaugeVec,
		)
		require.NoError(t, err)

		unauthenticatedInvocationsCount := getInvocationCount(invocationsGaugeVec, prometheusmetrics.UnauthenticatedUsersLabel)
		authenticatedInvocationsCount := getInvocationCount(invocationsGaugeVec, prometheusmetrics.AuthenticatedUsersLabel)
		require.Equal(t, unauthenticatedInvocationsCount, 1)
		require.Equal(t, authenticatedInvocationsCount, 1)
	})
}

func newAuthenticatedUsersGauge() prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{})
}

func newInvocationsGaugeVec() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{}, []string{"AuthStatus"})
}

func getInvocationCount(gaugeVec *prometheus.GaugeVec, label string) int {
	return int(testutil.ToFloat64(gaugeVec.WithLabelValues(label)))
}
