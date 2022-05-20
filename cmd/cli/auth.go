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

	err = makeFileFromTemplate("templates/data/remember_me.go.txt", ds.RootPath+"/data/remember_me.go")
	if err != nil {
		gracefulExit(err)
	}

	err = makeFileFromTemplate("templates/handlers/auth_handlers.go.txt", ds.RootPath+"/handlers/auth_handlers.go")
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

	err = makeFileFromTemplate("templates/middleware/remember_me.go.txt", ds.RootPath+"/middleware/remember_me.go")
	if err != nil {
		gracefulExit(err)
	}

	err = makeFileFromTemplate("templates/mailer/password_reset.html.tmpl", ds.RootPath+"/mail/password_reset.html.tmpl")
	if err != nil {
		gracefulExit(err)
	}

	err = makeFileFromTemplate("templates/mailer/password_reset.txt.tmpl", ds.RootPath+"/mail/password_reset.txt.tmpl")
	if err != nil {
		gracefulExit(err)
	}

	err = makeFileFromTemplate("templates/mailer/password_reset.html.tmpl", ds.RootPath+"/mail/password_reset.html.tmpl")
	if err != nil {
		gracefulExit(err)
	}

	err = makeFileFromTemplate("templates/views/forgot.jet.txt", ds.RootPath+"/views/forgot.jet")
	if err != nil {
		gracefulExit(err)
	}

	err = makeFileFromTemplate("templates/views/reset-password.jet.txt", ds.RootPath+"/views/reset-password.jet")
	if err != nil {
		gracefulExit(err)
	}

	err = makeFileFromTemplate("templates/views/login.jet.txt", ds.RootPath+"/views/login.jet")
	if err != nil {
		gracefulExit(err)
	}

	return nil
}
