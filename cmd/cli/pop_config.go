package main

import (
	"errors"
	"fmt"
)

func makePopConfig() error {
	err := ds.CreateDir(ds.RootPath + "/config")
	if err != nil {
		return err
	}

	path := ds.RootPath + "/config/database.yml"
	if !fileExists(ds.RootPath + "/config/database.yml") {
		err := makeFileFromTemplate("templates/pop_config.txt", path)
		if err != nil {
			return err
		}
	} else {
		return errors.New(fmt.Sprintf("%s already exists", path))
	}
	return nil
}
