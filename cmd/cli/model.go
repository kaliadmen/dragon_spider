package main

import (
	"errors"
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"strings"
)

func makeModel(name string) error {
	data, err := templateFs.ReadFile("templates/data/model.go.txt")
	if err != nil {
		return err
	}

	model := string(data)
	pluralClient := pluralize.NewClient()

	var modelName = name
	var tableName = name

	if pluralClient.IsPlural(name) {
		modelName = pluralClient.Singular(name)
		tableName = strings.ToLower(tableName)
	} else {
		tableName = strings.ToLower(pluralClient.Plural(name))
	}

	fileName := ds.RootPath + "/data/" + strings.ToLower(modelName) + ".go"
	if fileExists(fileName) {
		return errors.New(fileName + "already exist")
	}

	model = strings.ReplaceAll(model, "$MODELNAME$", strcase.ToCamel(modelName))
	model = strings.ReplaceAll(model, "$TABLENAME$", tableName)
	model = strings.ReplaceAll(model, "$?$", strings.ToLower(string(name[0])))

	err = copyDataToFile([]byte(model), fileName)
	if err != nil {
		return err
	}

	return nil
}
