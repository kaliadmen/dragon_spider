package main

import (
	"embed"
	"io/ioutil"
)

//go:embed templates
var templateFs embed.FS

func makeFileFromTemplate(templatePath, targetFile string) error {
	//check if file doesn't exist

	data, err := templateFs.ReadFile(templatePath)
	if err != nil {
		gracefulExit(err)
	}

	err = copyDataToFile(data, targetFile)
	if err != nil {
		gracefulExit(err)
	}

	return nil
}

func copyDataToFile(data []byte, writeTo string) error {
	err := ioutil.WriteFile(writeTo, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
