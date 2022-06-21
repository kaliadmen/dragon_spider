package main

import (
	"errors"
	"github.com/gobuffalo/pop"
	"os"
)

func makeSessionTable() error {
	dbType := convertDbType(ds.Db.DatabaseType)

	if dbType == "" || os.Getenv("SESSION_TYPE") == "" {
		return errors.New("you need to set a session type and/or a database in .env file first")
	}

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

	upBytes, err := templateFs.ReadFile("templates/migrations/session_table." + dbType + ".up.sql")
	if err != nil {
		return err
	}

	downBytes, err := templateFs.ReadFile("templates/migrations/session_table." + dbType + ".down.sql")
	if err != nil {
		return err
	}

	err = ds.CreatePopMigrations(upBytes, downBytes, "create_session_table", "sql")
	if err != nil {
		return err
	}

	//run migrations
	err = ds.PopMigrateUp(tx)
	if err != nil {
		return err
	}

	return nil
}
