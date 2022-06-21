package main

import (
	"errors"
	"github.com/iancoleman/strcase"
	"io/ioutil"
	"strings"
)

func makeHandler(name string) error {
	fileName := ds.RootPath + "/handlers/" + strings.ToLower(name) + ".go"
	if fileExists(fileName) {
		return errors.New(fileName + "already exist")
	}

	data, err := templateFs.ReadFile("templates/handlers/handler.go.txt")
	if err != nil {
		return err
	}

	handler := string(data)
	handler = strings.ReplaceAll(handler, "$HANDLERNAME$", strcase.ToCamel(strings.ToLower(name)))

	err = ioutil.WriteFile(fileName, []byte(handler), 0644)
	if err != nil {
		return err
	}

	return nil
}
