package main

import (
	"fmt"
	"strconv"
	"time"
)

func makeMigrations(name string) error {
	dbType := convertDbType(ds.Db.DatabaseType)
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), name)

	upFile := ds.RootPath + "/migrations/" + filename + "." + dbType + ".up.sql"
	downFile := ds.RootPath + "/migrations/" + filename + "." + dbType + ".down.sql"

	err := makeFileFromTemplate("templates/migrations/migration."+dbType+".up.sql", upFile)
	if err != nil {
		gracefulExit(err)
	}

	err = makeFileFromTemplate("templates/migrations/migration."+dbType+".down.sql", downFile)
	if err != nil {
		gracefulExit(err)
	}

	return nil
}

func runMigration(migrationType, step string) error {
	dsn := GetDSN()

	switch migrationType {
	case "up":
		err := ds.MigrateUp(dsn)
		if err != nil {
			return err
		}

	case "down":
		if step == "" {
			step = "1"
		}

		if step == "all" {
			err := ds.MigrateDownAll(dsn)
			if err != nil {
				return err
			}
		}

		step, err := strconv.Atoi(step)
		if err != nil {
			gracefulExit(err)
		}

		if step > 1{
			err := ds.Steps(dsn, step)
			if err != nil {
				return err
			}
		} else {
			err := ds.Steps(dsn, -1)
			if err != nil {
				return err
			}
		}

	case "reset":
		err := ds.MigrateDownAll(dsn)
		if err != nil {
			return err
		}
		err = ds.MigrateUp(dsn)
		if err != nil {
			return err
		}

	default:
		showHelp()
	}

	return nil
}
