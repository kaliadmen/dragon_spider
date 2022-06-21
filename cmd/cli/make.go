package main

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"strings"
)

func runMake(arg2, arg3, arg4 string) error {
	switch strings.ToLower(arg2) {
	case "auth":
		err := makeAuth()
		if err != nil {
			return err
		}
		color.Yellow(" -users, tokens and remember_token migrations created and executed")
		color.Yellow(" -user, token, and remember_me models created")
		color.Yellow(" -auth and remember_me middleware created")
		color.Yellow(" -password reset email templates created")
		color.Yellow(" -user login and password reset views created")
		color.Yellow("")
		color.Yellow("Please add user and token models to data/models.go!")
		color.Yellow("Don't forget to use the appropriate middleware in your routes!")

	case "handler":
		if arg3 == "" {
			return errors.New("handler must have a name")
		}

		err := makeHandler(arg3)
		if err != nil {
			return err
		}

	case "key":
		err := makeKey()
		if err != nil {
			return err
		}

	case "database":
		err := makeSqliteDb()
		if err != nil {
			return err
		}

	case "mail":
		if arg3 == "" {
			return errors.New("mail template must have a name")
		}

		err := makeMail(arg3)
		if err != nil {
			return err
		}

	case "migration":
		if arg3 == "" {
			return errors.New("migration must have a name")
		}

		err := makeMigrations(arg3, arg4)
		if err != nil {
			return err
		}

	case "model":
		if arg3 == "" {
			return errors.New("model must have a name")
		}

		err := makeModel(arg3)
		if err != nil {
			return err
		}
	case "popconfig":
		err := makePopConfig()
		if err != nil {
			return err
		}

	case "session":
		err := makeSessionTable()
		if err != nil {
			return err
		}
	default:
		return errors.New(fmt.Sprintf("unknown make command: %s", arg2))

	}

	return nil
}
