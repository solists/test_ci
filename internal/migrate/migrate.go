package migrate

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

func Migrate(db *sqlx.DB, driver, migrationsPath string, migrateDown bool) (err error) {
	if err = goose.SetDialect(driver); err != nil {
		return fmt.Errorf("goose can't change dialect to %v", driver)
	}

	if migrateDown {
		if err = goose.Down(db.DB, migrationsPath); err != nil {
			return fmt.Errorf("can't up %v migrations: %v", driver, err)
		}

		return
	}

	if err = goose.Up(db.DB, "./migrations"); err != nil {
		return fmt.Errorf("failed to apply database migrations: %v", err)
	}

	return
}
