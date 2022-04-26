package main

import (
	"errors"
	"github.com/fatih/color"
	dragonSpider "github.com/kaliadmen/dragon_spider"
	"os"
	"strings"
)

const version = "1.0.0"

var ds dragonSpider.DragonSpider

func main() {
	var message string

	arg1, arg2, arg3, err := validateArguments()

	if err != nil {
		gracefulExit(err)
	}

	setup()

	switch strings.ToLower(arg1) {
	case "help":
		showHelp()

	case "version":
		color.Yellow("Application version: " + version)

	case "make":
		if arg2 == "" {
			gracefulExit(errors.New("make requires subcommands: (migration|handler|model)"))
		}
		err = runMake(arg2, arg3)
		if err != nil {
			gracefulExit(err)
		}

	case "migrate":
		if arg2 == "" {
			arg2 = "up"
		}
		err = runMigration(arg2, arg3)
		if err != nil {
			gracefulExit(err)
		}
		message = "Migration completed successfully!"

	default:
		showHelp()
	}

	gracefulExit(nil, message)
}

func validateArguments() (string, string, string, error) {
	var arg1, arg2, arg3 string

	if len(os.Args) > 1 {
		arg1 = os.Args[1]

		if len(os.Args) >= 3 {
			arg2 = os.Args[2]
		}

		if len(os.Args) >= 4 {
			arg3 = os.Args[3]
		}
	} else {
		color.Red("Error: commands required")
		showHelp()
		return "", "", "", errors.New("commands required")
	}

	return arg1, arg2, arg3, nil
}

func gracefulExit(err error, msg ...string) {
	message := ""
	if len(msg) > 0 {
		message = msg[0]
	}

	if err != nil {
		color.Red("Error: %v\n", err)
	}

	if len(message) > 0 {
		color.Yellow(message)
	} else {
		color.Green("Done!")
	}

	os.Exit(0)
}
