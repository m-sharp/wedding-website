package migrations

import (
	"context"

	"github.com/m-sharp/wedding-website/lib"
)

const (
	createRSVPTable = `CREATE TABLE rsvp(
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		content JSON NOT NULL,
		ctime DATETIME
	);`
	destroyRSVPTable = `DROP TABLE rsvp;`
)

type Migration2 struct{}

func (m *Migration2) Upgrade(ctx context.Context, client *lib.DBClient) error {
	if _, err := client.Db.ExecContext(ctx, createRSVPTable); err != nil {
		return lib.NewDBError(createRSVPTable, err)
	}
	return nil
}

func (m *Migration2) Downgrade(ctx context.Context, client *lib.DBClient) error {
	if _, err := client.Db.ExecContext(ctx, destroyRSVPTable); err != nil {
		return lib.NewDBError(destroyRSVPTable, err)
	}
	return nil
}
