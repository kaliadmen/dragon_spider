package main

import "fmt"

func makeSqliteDb() error {
	err := ds.CreateDir(ds.RootPath + "/tmp/sqlite")
	if err != nil {
		return err
	}
	if !fileExists(ds.RootPath + "tmp/sqlite/app.db") {
		err := ds.CreateFile(fmt.Sprintf("%s/tmp/sqlite/app.db", ds.RootPath))
		if err != nil {
			return err
		}
	}
	return nil
}
