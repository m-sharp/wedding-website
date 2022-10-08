package migrations

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/m-sharp/wedding-website/lib"
)

const (
	increment = `INSERT INTO migration (success, ctime) VALUES (1, ?);`
	decrement = `DELETE FROM migration WHERE ID=(SELECT MAX(id) FROM migration);`

	countMigrations = `SELECT COUNT(*) FROM migration;`
	checkForTable   = `SELECT count(*)
		FROM information_schema.tables
		WHERE table_schema = 'wedding'
		  AND table_name = 'migration'
		LIMIT 1;`
)

type Migration interface {
	Upgrade(ctx context.Context, client *lib.DBClient) error
	Downgrade(ctx context.Context, client *lib.DBClient) error
}

func RunAll(ctx context.Context, client *lib.DBClient, log *zap.Logger) error {
	log = log.Named("Migrator")
	log.Info("Running DB migrations...")

	startCount, err := GetCurrentMigrationCount(ctx, client)
	if err != nil {
		return err
	}

	var ran []Migration
	for i, migration := range getAllMigrations() {
		if i <= startCount {
			continue
		}

		log.Debug("Running migration", zap.Int("Migration Number", i))
		if err := migration.Upgrade(ctx, client); err != nil {
			log.Error("Error running migration", zap.Int("Migration Number", i), zap.Error(err))
			if innerErr := rollback(ctx, client, log, ran...); innerErr != nil {
				log.Error("Failed to rollback migrations", zap.Int("Migration Number", i), zap.Error(err))
			}
			return err
		}
		ran = append(ran, migration)

		if err := incrementMigrationTable(ctx, client); err != nil {
			log.Error("Failed to increment migration table", zap.Int("Migration Number", i), zap.Error(err))
			return err
		}
	}
	log.Info("Finished running migrations", zap.Int("Run Count", len(ran)))
	return nil
}

func rollback(ctx context.Context, client *lib.DBClient, log *zap.Logger, toRollback ...Migration) error {
	for i, migration := range toRollback {
		log = log.With(zap.Int("Migration Number", i))
		log.Debug("Rolling back migration")
		if err := migration.Downgrade(ctx, client); err != nil {
			return errors.New(fmt.Sprintf("Failed to roll back migration #%v: %s", i, err))
		}

		if err := decrementMigrationTable(ctx, client); err != nil {
			log.Error("Failed to decrement migration table", zap.Error(err))
			return err
		}
	}
	return nil
}

func GetCurrentMigrationCount(ctx context.Context, client *lib.DBClient) (int, error) {
	var tableCheck int
	if err := client.Db.QueryRowContext(ctx, checkForTable).Scan(&tableCheck); err != nil {
		return 0, fmt.Errorf("error checking for migration table: %w", err)
	} else if tableCheck == 0 {
		return 0, nil
	}

	var result int
	if err := client.Db.QueryRowContext(ctx, countMigrations).Scan(&result); err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, fmt.Errorf("error getting current migration count: %w", err)
	}
	return result, nil
}

func incrementMigrationTable(ctx context.Context, client *lib.DBClient) error {
	if _, err := client.Db.ExecContext(ctx, increment, time.Now()); err != nil {
		return lib.NewDBError(increment, err)
	}
	return nil
}

func decrementMigrationTable(ctx context.Context, client *lib.DBClient) error {
	if _, err := client.Db.ExecContext(ctx, decrement); err != nil {
		return lib.NewDBError(decrement, err)
	}
	return nil
}

func getAllMigrations() map[int]Migration {
	return map[int]Migration{
		1: &Migration1{},
		2: &Migration2{},
	}
}
