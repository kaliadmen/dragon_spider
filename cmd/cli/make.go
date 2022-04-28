package main

import (
	"errors"
	"strings"
)

func runMake(arg2, arg3 string) error {
	switch strings.ToLower(arg2) {
	case "auth":
		err := makeAuth()
		if err != nil {
			gracefulExit(err)
		}

	case "handler":
		if arg3 == "" {
			gracefulExit(errors.New("handler must have a name"))
		}

		err := makeHandler(arg3)
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
