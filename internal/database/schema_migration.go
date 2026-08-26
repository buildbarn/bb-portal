package database

import (
	"context"
	"database/sql"
)

// RenameLegacyActionsTable renames the table used by the former Action entity
// before Ent migrates the renamed ActionExecution entity. Ent's automatic
// migrator does not infer table renames and would otherwise create an empty
// action_executions table alongside the existing actions table.
func RenameLegacyActionsTable(ctx context.Context, connection *sql.DB) error {
	_, err := connection.ExecContext(ctx, `
DO $$
BEGIN
    IF to_regclass('action_executions') IS NULL
       AND to_regclass('actions') IS NOT NULL THEN
        ALTER TABLE actions RENAME TO action_executions;
    END IF;
END
$$;
`)
	return err
}
