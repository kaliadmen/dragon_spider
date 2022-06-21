package main

import (
	"github.com/gobuffalo/pop"
	"os"
	"strings"
)

func makeAuth() error {
	//make migrations
	dbType := convertDbType(ds.Db.DatabaseType)

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

	upBytes, err := templateFs.ReadFile("templates/migrations/auth_tables." + dbType + ".up.sql")
	if err != nil {
		return err
	}

	downBytes, err := templateFs.ReadFile("templates/migrations/auth_tables." + dbType + ".up.sql")
	if err != nil {
		return err
	}

	err = ds.CreatePopMigrations(upBytes, downBytes, "auth", "sql")
	if err != nil {
		return err
	}

	//run migrations
	err = ds.PopMigrateUp(tx)
	if err != nil {
		return err
	}

	//copy files
	err = makeFileFromTemplate("templates/data/user.go.txt", ds.RootPath+"/data/user.go")
	if err != nil {
		return err
	}

	err = makeFileFromTemplate("templates/data/token.go.txt", ds.RootPath+"/data/token.go")
	if err != nil {
		return err
	}

	err = makeFileFromTemplate("templates/data/remember_me.go.txt", ds.RootPath+"/data/remember_me.go")
	if err != nil {
		return err
	}

	data, err := templateFs.ReadFile("templates/handlers/auth_handlers.go.txt")
	if err != nil {
		return err
	}

	file := string(data)

	if appURL == "" && os.Getenv("APP_GITHUB_URL") != "" {
		appURL = os.Getenv("APP_GITHUB_URL")
	}

	file = strings.ReplaceAll(file, "${APP_NAME}", appURL)

	err = copyDataToFile([]byte(file), ds.RootPath+"/handlers/auth_handlers.go")
	if err != nil {
		return err
	}

	//copy middleware
	err = makeFileFromTemplate("templates/middleware/auth_web.go.txt", ds.RootPath+"/middleware/auth_web.go")
	if err != nil {
		return err
	}

	err = makeFileFromTemplate("templates/middleware/auth_api.go.txt", ds.RootPath+"/middleware/auth_api.go")
	if err != nil {
		return err
	}

	data, err = templateFs.ReadFile("templates/middleware/remember_me.go.txt")
	if err != nil {
		return err
	}

	file = string(data)

	if appURL == "" && os.Getenv("APP_GITHUB_URL") != "" {
		appURL = os.Getenv("APP_GITHUB_URL")
	}

	file = strings.ReplaceAll(file, "${APP_NAME}", appURL)

	err = copyDataToFile([]byte(file), ds.RootPath+"/middleware/remember_me.go")
	if err != nil {
		return err
	}

	err = makeFileFromTemplate("templates/mailer/password_reset.html.tmpl", ds.RootPath+"/mail/password_reset.html.tmpl")
	if err != nil {
		return err
	}

	err = makeFileFromTemplate("templates/mailer/password_reset.txt.tmpl", ds.RootPath+"/mail/password_reset.txt.tmpl")
	if err != nil {
		return err
	}

	err = makeFileFromTemplate("templates/views/forgot.jet.txt", ds.RootPath+"/views/forgot.jet")
	if err != nil {
		return err
	}

	err = makeFileFromTemplate("templates/views/reset-password.jet.txt", ds.RootPath+"/views/reset-password.jet")
	if err != nil {
		return err
	}

	err = makeFileFromTemplate("templates/views/login.jet.txt", ds.RootPath+"/views/login.jet")
	if err != nil {
		return err
	}

	return nil
}
