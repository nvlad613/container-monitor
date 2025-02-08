package infra

import (
	"backend/config"
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	_ "github.com/uptrace/bun/driver/pgdriver"
)

func InitDb(conf config.DatabaseConnection) (*bun.DB, error) {
	connector := pgdriver.NewConnector(
		pgdriver.WithAddr(fmt.Sprintf("%s:%d", conf.Hostname, conf.Port)),
		pgdriver.WithUser(conf.User.Login),
		pgdriver.WithPassword(conf.User.Password),
		pgdriver.WithDatabase(conf.Database),
		pgdriver.WithInsecure(!conf.TLSEnabled),
	)

	db := bun.NewDB(
		sql.OpenDB(connector),
		pgdialect.New(),
	)
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to open connection with db: %w", err)
	}

	return db, nil
}
