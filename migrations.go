package dragonSpider

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gobuffalo/pop"

	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

func (ds *DragonSpider) ConnectToPop() (*pop.Connection, error) {
	popConnectionType := os.Getenv("POP_CONNECTION_TYPE")

	if popConnectionType == "" {
		popConnectionType = "development"
	}

	tx, err := pop.Connect(popConnectionType)
	if err != nil {
		return nil, err
	}

	return tx, nil

}

func (ds *DragonSpider) CreatePopMigrations(upMigration, downMigration []byte, migrationName, migrationType string) error {
	var migrationPath = ds.RootPath + "/migrations"

	err := pop.MigrationCreate(migrationPath, migrationName, migrationType, upMigration, downMigration)
	if err != nil {
		return err
	}

	return nil
}

func (ds *DragonSpider) PopMigrateUp(tx *pop.Connection) error {
	var migrationPath = ds.RootPath + "/migrations"

	fileMigrator, err := pop.NewFileMigrator(migrationPath, tx)
	if err != nil {
		return err
	}

	err = fileMigrator.Up()
	if err != nil {
		return err
	}

	return nil
}

func (ds *DragonSpider) PopMigrateDown(tx *pop.Connection, steps ...int) error {
	var migrationPath = ds.RootPath + "/migrations"

	n := 1
	if len(steps) > 0 {
		n = steps[0]
	}

	fileMigrator, err := pop.NewFileMigrator(migrationPath, tx)
	if err != nil {
		return err
	}

	err = fileMigrator.Down(n)
	if err != nil {
		return err
	}

	return nil
}

func (ds *DragonSpider) PopReset(tx *pop.Connection) error {
	var migrationPath = ds.RootPath + "/migrations"

	fileMigrator, err := pop.NewFileMigrator(migrationPath, tx)
	if err != nil {
		return err
	}

	err = fileMigrator.Reset()
	if err != nil {
		return err
	}

	return nil
}
