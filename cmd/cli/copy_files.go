package main

import (
	"embed"
	"errors"
	"io/ioutil"
	"os"
)

//go:embed templates
var templateFs embed.FS

func makeFileFromTemplate(templatePath, targetFile string) error {
	//check if file doesn't exist
	if fileExists(targetFile) {
		return errors.New(targetFile + " already exist")
	}

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

func fileExists(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}
