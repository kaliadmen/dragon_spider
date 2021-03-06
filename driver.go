package dragonSpider

import (
	"database/sql"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

//OpenDb opens a database connection
func (ds *DragonSpider) OpenDb(dbType, dsn string) (*sql.DB, error) {
	if strings.ToLower(dbType) == "postgres" || strings.ToLower(dbType) == "postgresql" {
		dbType = "pgx"
	}

	if strings.ToLower(dbType) == "mysql" || strings.ToLower(dbType) == "mariadb" {
		dbType = "mysql"
	}

	if strings.ToLower(dbType) == "sqlite" {
		dbType = "sqlite3"
	}

	db, err := sql.Open(dbType, dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
