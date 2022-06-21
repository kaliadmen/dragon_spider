package main

import (
	"github.com/gobuffalo/pop"
	"strconv"
)

func makeMigrations(name, migrationFormat string) error {
	var upMigration, downMigration string

	validatePopConfig()

	if migrationFormat == "fizz" || migrationFormat == "" {
		migrationFormat = "fizz"
		upBytes, err := templateFs.ReadFile("templates/migrations/migration_up.fizz")
		if err != nil {
			return err
		}

		downBytes, err := templateFs.ReadFile("templates/migrations/migration_down.fizz")
		if err != nil {
			return err
		}

		upMigration = string(upBytes)
		downMigration = string(downBytes)
	} else {
		migrationFormat = "sql"
	}

	err := ds.CreatePopMigrations([]byte(upMigration), []byte(downMigration), name, migrationFormat)
	if err != nil {
		return err
	}

	return nil
}

func runMigration(migrationType, step string) error {
	validatePopConfig()
	tx, err := ds.ConnectToPop()
	if err != nil {
		return err
	}

	defer func(tx *pop.Connection) {
		err := tx.Close()
		if err != nil {
			gracefulExit(err)
		}
	}(tx)

	switch migrationType {
	case "up":
		//err := ds.MigrateUp(dsn)
		err := ds.PopMigrateUp(tx)
		if err != nil {
			return err
		}

	case "down":
		if step == "" {
			step = "1"
		}

		if step == "all" {
			err := ds.PopMigrateDown(tx, -1)
			if err != nil {
				return err
			}
		}

		step, err := strconv.Atoi(step)
		if err != nil {
			gracefulExit(err)
		}

		if step > 1 {
			err := ds.PopMigrateDown(tx, step)
			if err != nil {
				return err
			}
		} else {
			err := ds.PopMigrateDown(tx, 1)
			if err != nil {
				return err
			}
		}

	case "reset":
		err := ds.PopReset(tx)
		if err != nil {
			return err
		}

	default:
		showHelp()
	}

	return nil
}
