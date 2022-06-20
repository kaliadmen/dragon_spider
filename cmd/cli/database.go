package main

import (
	"errors"
	"fmt"
)

func makeSqliteDb() error {
	err := ds.CreateDirs(ds.RootPath + "/db/sqlite")
	if err != nil {
		return err
	}

	path := ds.RootPath + "/db/sqlite/app.db"
	if !fileExists(path) {
		err := ds.CreateFile(path)
		if err != nil {
			return err
		}
	} else {
		return errors.New(fmt.Sprintf("%s already exists", path))
	}

	return nil
}
