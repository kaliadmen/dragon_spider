package main

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"strings"
)

func setup(arg1, arg2 string) {
	if arg1 != "new" && arg1 != "version" && arg1 != "help" {
		err := godotenv.Load()
		if err != nil {
			gracefulExit(err)
		}

		rootPath, err := os.Getwd()
		if err != nil {
			gracefulExit(err)
		}

		ds.RootPath = rootPath
		if os.Getenv("DATABASE_TYPE") != "" {
			ds.Db.DatabaseType = os.Getenv("DATABASE_TYPE")
		} else {
			ds.Db.DatabaseType = "sqlite"
		}
	}

}

func GetDSN() string {
	dbType := convertDbType(ds.Db.DatabaseType)

	switch dbType {

	case "postgres":
		var dsn string
		if os.Getenv("DATABASE_PASSWORD") != "" {
			dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_PASSWORD"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSL_MODE"),
			)
		} else {
			dsn = fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSL_MODE"),
			)
		}

		return dsn

	case "mysql":
		return "mysql://" + ds.CreateDSN()

	case "sqlite":
		return "sqlite3://" + ds.CreateDSN()

	default:
		return ""
	}
}

func showHelp() {
	color.Yellow(`Available commands:
	help                           -show help commands
	version	                       -print application version
	up	                           -take server out of maintenance mode
	down	                       -puts server in maintenance mode
	migrate(up)                    -runs all up migrations that have not been ran
	migrate down                   -reverses most recent migration
	migrate down all               -runs all down migrations
	migrate down n[int]            -runs n number of migrations, migrates down if negative number is passed
	migrate reset                  -runs all down migration in reverse, and runs all up migrations
	make migration <name> <format> -creates an up and down migration in migrations directory; format=sql/fizz (default fizz)
	make auth                      -creates migrations, models, and middleware for authentication, and runs migrations
	make database                  -create a sqlite database in tmp directory
	make handler <name>            -creates a bare handler in handlers directory
	make model <name>              -creates a bare model in data directory
	make popConfig                 -creates a config directory and a bare pop config file for gobuffalo migrations
	make session                   -creates a database table or cache entry for session store
	make mail <name>               -creates a html and plaintext email template in mail directory
`)
}

func convertDbType(dbType string) string {
	switch dbType {
	case "sqlite", "sqlite3":
		return "sqlite"

	case "mariadb", "mysql":
		return "mysql"

	case "postgres", "postgresql", "pgx":
		return "postgres"

	default:
		return ""
	}
}

func updateSourceFiles(path string, fi os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if fi.IsDir() {
		return nil
	}

	matched, err := filepath.Match("*.go", fi.Name())
	if err != nil {
		return err
	}

	if matched {
		read, err := os.ReadFile(path)
		if err != nil {
			gracefulExit(err)
		}
		contents := strings.Replace(string(read), "${APP_NAME}", appURL, -1)

		err = os.WriteFile(path, []byte(contents), 0)
		if err != nil {
			gracefulExit(err)
		}

	}

	return nil
}

func updateSource() {
	err := filepath.Walk(".", updateSourceFiles)
	if err != nil {
		gracefulExit(err)
	}
}

func validatePopConfig() {
	dbType := convertDbType(ds.Db.DatabaseType)

	if dbType == "" {
		gracefulExit(errors.New("no database connection supplied in .env file"))
	}

	if !fileExists(ds.RootPath + "/config/database.yml") {
		gracefulExit(errors.New(ds.RootPath + "/config/database.yml does not exist"))
	}

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
