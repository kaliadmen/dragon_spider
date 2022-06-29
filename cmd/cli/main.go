package main

import (
	"errors"
	"github.com/fatih/color"
	dragonSpider "github.com/kaliadmen/dragon_spider/v2"
	"os"
	"strings"
)

const version = "2.0.0"

var ds dragonSpider.DragonSpider

func main() {
	var message string

	arg1, arg2, arg3, arg4, err := validateArguments()

	if err != nil {
		gracefulExit(err)
	}

	setup(arg1, arg2)

	switch strings.ToLower(arg1) {
	case "help":
		showHelp()

	case "up":
		rpcClient(false)
	case "down":
		rpcClient(true)
	case "new":
		if arg2 == "" {
			gracefulExit(errors.New("new requires a name for the application"))
		}

		err := createApp(arg2, arg3, arg4)
		if err != nil {
			gracefulExit(err)
		}

	case "version":
		color.Yellow("Application version: " + version)

	case "make":
		if arg2 == "" {
			gracefulExit(errors.New("make requires subcommands: (migration|handler|model)"))
		}
		err = runMake(arg2, arg3, arg4)
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

func validateArguments() (string, string, string, string, error) {
	var arg1, arg2, arg3, arg4 string

	if len(os.Args) > 1 {
		arg1 = os.Args[1]

		if len(os.Args) >= 3 {
			arg2 = os.Args[2]
		}

		if len(os.Args) >= 4 {
			arg3 = os.Args[3]
		}
		if len(os.Args) >= 5 {
			arg4 = os.Args[4]
		}
	} else {
		color.Red("Error: commands required")
		showHelp()
		return "", "", "", "", errors.New("commands required")
	}

	return arg1, arg2, arg3, arg4, nil
}
