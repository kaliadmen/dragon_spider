package dragonSpider

import "database/sql"

type initPaths struct {
	//root of application
	rootPath string
	//directories available to application
	dirNames []string
}

//cookieConfig holds data for cookie configuration
type cookieConfig struct {
	name       string
	lifetime   string
	persistent string
	Secure     string
	domain     string
}

type databaseConfig struct {
	dsn      string
	database string
}

type Database struct {
	DatabaseType string
	Pool         *sql.DB
}

type redisConfig struct {
	host     string
	password string
	prefix   string
}
