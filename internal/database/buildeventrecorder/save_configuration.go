package buildeventrecorder

import (
	"context"
	"database/sql"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/internal/database"
	"github.com/buildbarn/bb-storage/pkg/util"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// saveConfigurationBatch is an efficient implementation of save configuration for
// a batch of configuration events.
func (r *BuildEventRecorder) saveConfigurationBatch(ctx context.Context, batch []BuildEventWithInfo) error {
	if len(batch) == 0 {
		return nil
	}

	ctx, span := r.tracer.Start(
		ctx,
		"BuildEventRecorder.saveConfigurationBatch",
		trace.WithAttributes(
			attribute.Int("batch_size", len(batch)),
			attribute.String("invocation.id", r.InvocationID),
			attribute.String("invocation.instance_name", r.InstanceName),
		),
	)
	defer span.End()

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return util.StatusWrap(err, "Failed to start transaction")
	}
	defer func() { tx.Rollback() }()

	row, err := tx.Sqlc().LockBazelInvocationCompletion(ctx, r.InvocationDbID)
	if err != nil {
		return util.StatusWrap(err, "Failed to lock bep completed for invocation")
	}
	if row.BepCompleted {
		return status.Error(codes.InvalidArgument, "Attempted to configure targets for an invocation but the invocation was already completed.")
	}

	if err = createConfigurationsBulk(ctx, r.InvocationDbID, tx, batch); err != nil {
		return util.StatusWrap(err, "Failed to bulk insert configurations")
	}

	if err := r.saveHandledEventsForBatch(ctx, batch, tx); err != nil {
		return util.StatusWrap(err, "Failed to bulk insert event metadata")
	}

	if err = tx.Commit(); err != nil {
		return util.StatusWrap(err, "Failed to commit batch of target configured events")
	}

	return nil
}

func createConfigurationsBulk(ctx context.Context, invocationDbID int64, tx database.Handle, batch []BuildEventWithInfo) error {
	configBuilders := make([]*ent.ConfigurationCreate, 0, len(batch))

	for _, info := range batch {
		be := info.Event

		configID := be.GetId().GetConfiguration().GetId()
		if configID == "" {
			return status.Error(codes.InvalidArgument, "Received configuration event with empty configuration ID")
		}
		config := be.GetConfiguration()
		if config == nil {
			return status.Error(codes.InvalidArgument, "Received non configuration event to batch configuration method")
		}

		create := tx.Ent().Configuration.Create().
			SetConfigurationID(configID).
			SetBazelInvocationID(invocationDbID).
			SetIsTool(config.IsTool)

		if config.Mnemonic != "" {
			create.SetMnemonic(config.Mnemonic)
		}
		if config.PlatformName != "" {
			create.SetPlatformName(config.PlatformName)
		}
		if config.Cpu != "" {
			create.SetCPU(config.Cpu)
		}
		if len(config.MakeVariable) > 0 {
			create.SetMakeVariables(config.MakeVariable)
		}

		configBuilders = append(configBuilders, create)
	}

	err := tx.Ent().Configuration.
		CreateBulk(configBuilders...).
		Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "Failed to bulk insert configurations to database")
	}
	return nil
}
