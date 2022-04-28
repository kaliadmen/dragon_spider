package main

import (
	"errors"
	"fmt"
	"time"
)

func makeSessionTable() error {
	dbType := ds.Db.DatabaseType

	if dbType == "" {
		gracefulExit(errors.New("you need to set a session type in .env file first"))
	}

	if dbType == "mariadb" {
		dbType = "mysql"
	}

	if dbType == "postgresql" {
		dbType = "postgres"
	}

	fileName := fmt.Sprintf("%d_create_session_table", time.Now().UnixNano())
	upFile := ds.RootPath + "/migrations/" + fileName + "." + dbType + ".up.sql"
	downFile := ds.RootPath + "/migrations/" + fileName + "." + dbType + ".down.sql"

	err := makeFileFromTemplate("templates/migrations/"+"session_table."+dbType+".up.sql", upFile)
	if err != nil {
		gracefulExit(err)
	}

	err = makeFileFromTemplate("templates/migrations/"+"session_table."+dbType+".down.sql", downFile)
	if err != nil {
		gracefulExit(err)
	}

	err = runMigration("up", "")
	if err != nil {
		gracefulExit(err)
	}

	return nil
}
