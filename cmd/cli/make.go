package main

import (
	"errors"
	"fmt"
	"github.com/iancoleman/strcase"
	"io/ioutil"
	"strings"
	"time"
)

func runMake(arg2, arg3 string) error {
	switch strings.ToLower(arg2) {
	case "auth":
		err := runAuth()
		if err != nil {
			gracefulExit(err)
		}

	case "handler":
		if arg3 == "" {
			gracefulExit(errors.New("handler must have a name"))
		}

		fileName := ds.RootPath + "/handlers/" + strings.ToLower(arg3) + ".go"
		if fileExists(fileName) {
			gracefulExit(errors.New(fileName + "already exist"))
		}

		data, err := templateFs.ReadFile("templates/handlers/handler.go.txt")
		if err != nil {
			gracefulExit(err)
		}

		handler := string(data)
		handler = strings.ReplaceAll(handler, "$HANDLERNAME$", strcase.ToCamel(strings.ToLower(arg3)))

		err = ioutil.WriteFile(fileName, []byte(handler), 0644)
		if err != nil {
			gracefulExit(err)
		}

	case "migration":
		dbType := ds.Db.DatabaseType
		if arg3 == "" {
			gracefulExit(errors.New("migration must have a name"))
		}

		filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), arg3)

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

	}

	return nil
}
