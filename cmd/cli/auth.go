package main

import (
	"fmt"
	"time"
)

func makeAuth() error {
	//make migrations
	dbType := convertDbType(ds.Db.DatabaseType)
	filename := fmt.Sprintf("%d_create_auth_tables", time.Now().UnixNano())
	upFile := ds.RootPath + "/migrations/" + filename + ".up.sql"
	downFile := ds.RootPath + "/migrations/" + filename + ".down.sql"

	err := makeFileFromTemplate("templates/migrations/auth_tables."+dbType+".up.sql", upFile)
	if err != nil {
		gracefulExit(err)
	}

	err = makeFileFromTemplate("templates/migrations/auth_tables."+dbType+".down.sql", downFile)
	if err != nil {
		gracefulExit(err)
	}
	//run migrations
	err = runMigration("up", "")
	if err != nil {
		gracefulExit(err)
	}
	//copy files
	err = makeFileFromTemplate("templates/data/user.go.txt", ds.RootPath+"/data/user.go")
	if err != nil {
		gracefulExit(err)
	}

	err = makeFileFromTemplate("templates/data/token.go.txt", ds.RootPath+"/data/token.go")
	if err != nil {
		gracefulExit(err)
	}

	//copy middleware
	err = makeFileFromTemplate("templates/middleware/auth_web.go.txt", ds.RootPath+"/middleware/auth_web.go")
	if err != nil {
		gracefulExit(err)
	}

	err = makeFileFromTemplate("templates/middleware/auth_api.go.txt", ds.RootPath+"/middleware/auth_api.go")
	if err != nil {
		gracefulExit(err)
	}

	return nil
}
