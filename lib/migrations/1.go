package migrations

import (
	"context"

	"github.com/m-sharp/wedding-website/lib"
)

const (
	createMigrationTable = `CREATE TABLE migration(
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		success BIT,
		ctime DATETIME
	);`
	destroyMigrationTable = `DROP TABLE migration;`
)

type Migration1 struct{}

func (m *Migration1) Upgrade(ctx context.Context, client *lib.DBClient) error {
	if _, err := client.Db.ExecContext(ctx, createMigrationTable); err != nil {
		return lib.NewDBError(createMigrationTable, err)
	}
	return nil
}

func (m *Migration1) Downgrade(ctx context.Context, client *lib.DBClient) error {
	if _, err := client.Db.ExecContext(ctx, destroyMigrationTable); err != nil {
		return lib.NewDBError(destroyMigrationTable, err)
	}
	return nil
}
