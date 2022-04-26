package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func runMake(arg2, arg3 string) error {
	switch strings.ToLower(arg2) {
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

	case "auth":
		err := runAuth()
		if err != nil {
			gracefulExit(err)
		}

	}

	return nil
}
