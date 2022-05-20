package main

import "fmt"

func makeSqliteDb() error {
	err := ds.CreateDirs(ds.RootPath + "/db/sqlite")
	if err != nil {
		return err
	}

	if !fileExists(ds.RootPath + "db/sqlite/app.db") {
		err := ds.CreateFile(fmt.Sprintf("%s/db/sqlite/app.db", ds.RootPath))
		if err != nil {
			return err
		}
	}
	return nil
}
