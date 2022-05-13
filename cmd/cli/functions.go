package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"os"
)

func setup() {
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
	help                   -show help commands
	version	               -print application version
	migrate(up)            -runs all up migrations that have not been ran
	migrate down           -reverses most recent migration
	migrate down all       -runs all down migrations
	migrate n[int]         -runs n number of migrations, migrates down if negative number is passed
	migrate reset          -runs all down migration in reverse, and runs all up migrations
	make migrations <name> -creates an up and down migration in migrations directory
	make auth              -creates migrations, models, and middleware for authentication, and runs migrations
	make database          -create a sqlite database in tmp directory
	make handler <name>    -creates a bare handler in handlers directory
	make model <name>      -creates a bare model in data directory
	make session           -creates a database table or cache entry for session store
	make mail <name>       -creates a html and plaintext email template in mail directory
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
