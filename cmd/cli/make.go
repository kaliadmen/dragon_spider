package main

import (
	"errors"
	"github.com/fatih/color"
	"strings"
)

func runMake(arg2, arg3 string) error {
	switch strings.ToLower(arg2) {
	case "auth":
		err := makeAuth()
		if err != nil {
			gracefulExit(err)
		}
		color.Yellow(" -users, tokens and remember_token migrations created and executed")
		color.Yellow(" -user and token created")
		color.Yellow(" -auth middleware created")
		color.Yellow("")
		color.Yellow("Please add user and token models to data/models.go!")
		color.Yellow("Don't forget to use the appropriate middleware in your routes!")

	case "handler":
		if arg3 == "" {
			gracefulExit(errors.New("handler must have a name"))
		}

		err := makeHandler(arg3)
		if err != nil {
			gracefulExit(err)
		}

	case "key":
		err := makeKey()
		if err != nil {
			gracefulExit(err)
		}

	case "migration":
		if arg3 == "" {
			gracefulExit(errors.New("migration must have a name"))
		}

		err := makeMigrations(arg3)
		if err != nil {
			gracefulExit(err)
		}

	case "model":
		if arg3 == "" {
			gracefulExit(errors.New("model must have a name"))
		}

		err := makeModel(arg3)
		if err != nil {
			gracefulExit(err)
		}

	case "session":
		err := makeSessionTable()
		if err != nil {
			gracefulExit(err)
		}

	}

	return nil
}
