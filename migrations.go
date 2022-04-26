package dragonSpider

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func (ds *DragonSpider) MigrateUp(dsn string) error {
	m, err := migrate.New("file://"+ds.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	_, _ = m.Close()

	if err := m.Up(); err != nil {
		log.Println("Error running migration:", err)
		return err
	}

	return nil
}

func (ds *DragonSpider) MigrateDownAll(dsn string) error {
	m, err := migrate.New("file://"+ds.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	_, _ = m.Close()

	if err := m.Down(); err != nil {
		log.Println("Error running migration:", err)
		return err
	}

	return nil

}

func (ds *DragonSpider) Steps(dsn string, n int) error {
	m, err := migrate.New("file://"+ds.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	_, _ = m.Close()

	if err := m.Steps(n); err != nil {
		log.Println("Error running migration:", err)
		return err
	}

	return nil
}

func (ds *DragonSpider) ForceMigrate(dsn string) error {
	m, err := migrate.New("file://"+ds.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	_, _ = m.Close()

	if err := m.Force(-1); err != nil {
		log.Println("Error running migration:", err)
		return err
	}

	return nil

}
