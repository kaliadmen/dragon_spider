package dragonSpider

import (
	"database/sql"
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/kaliadmen/dragon_spider/render"
	"github.com/kaliadmen/dragon_spider/session"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const version = "1.0.0"

//DragonSpider is an overall type for the Dragon Spider package.
//Members exported in this type are available to any application that uses it
type DragonSpider struct {
	AppName     string
	Debug       bool
	Version     string
	ErrorLog    *log.Logger
	InfoLog     *log.Logger
	Render      *render.Render
	JetTemplate *jet.Set
	Routes      *chi.Mux
	Session     *scs.SessionManager
	Db          Database
	RootPath    string
	config      config
}

type config struct {
	port        string
	renderer    string
	sessionType string
	cookie      cookieConfig
	database    databaseConfig
}

//New creates application config, reads the .env file, populate Dragon Spider type bases on .env values
//and creates the necessary directories and files if they don't exist
func (ds *DragonSpider) New(rp string) error {
	//create directories
	pathConfig := initPaths{
		rootPath: rp,
		dirNames: []string{
			"handlers",
			"middleware",
			"views",
			"data",
			"migrations",
			"public",
			"tmp",
			"logs",
			"sqlite",
		},
	}

	//create directory structure
	err := ds.InitDir(pathConfig)
	if err != nil {
		return err
	}

	//create empty .env file if it doesn't exist in root path
	err = ds.createDotEnv(rp)
	if err != nil {
		return err
	}

	//read .env file using godotenv package
	err = godotenv.Load(rp + "/.env")
	if err != nil {
		return err
	}

	//create sqlite database if it doesn't exist
	err = ds.createSqliteDb(rp)
	if err != nil {
		return err
	}

	//create loggers
	infoLog, errorLog := ds.startLoggers()
	ds.InfoLog = infoLog
	ds.ErrorLog = errorLog

	//connect to a database
	if os.Getenv("DATABASE_TYPE") != "" {
		pool, err := ds.OpenDb(os.Getenv("DATABASE_TYPE"), ds.CreateDSN())
		if err != nil {
			errorLog.Println(err)
			os.Exit(1)
		}
		ds.Db = Database{
			DatabaseType: os.Getenv("DATABASE_TYPE"),
			Pool:         pool,
		}
	} else {
		err := os.Setenv("DATABASE_TYPE", "sqlite")
		if err != nil {
			return err
		}
		pool, err := ds.OpenDb(os.Getenv("DATABASE_TYPE"), ds.CreateDSN())
		if err != nil {
			errorLog.Println(err)
			os.Exit(1)
		}
		ds.Db = Database{
			DatabaseType: "sqlite",
			Pool:         pool,
		}
	}

	//set application variables
	ds.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	ds.Version = version
	ds.RootPath = rp

	//set config variables
	ds.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"),
		cookie: cookieConfig{
			name:       os.Getenv("COOKIE_NAME"),
			lifetime:   os.Getenv("COOKIE_LIFETIME"),
			persistent: os.Getenv("COOKIE_PERSIST"),
			Secure:     os.Getenv("COOKIE_SECURE"),
			domain:     os.Getenv("COOKIE_DOMAIN"),
		},
		sessionType: os.Getenv("SESSION_TYPE"),
		database: databaseConfig{
			dsn:      ds.CreateDSN(),
			database: os.Getenv("DATABASE_TYPE"),
		},
	}

	//set session configuration
	ses := session.Session{
		CookieName:     ds.config.cookie.name,
		CookieLifetime: ds.config.cookie.lifetime,
		CookiePersist:  ds.config.cookie.persistent,
		CookieSecure:   ds.config.cookie.Secure,
		CookieDomain:   ds.config.cookie.domain,
		SessionType:    ds.config.sessionType,
	}

	//create session
	ds.Session = ses.InitSession()

	var jetViews = jet.NewSet(jet.NewOSFileSystemLoader(fmt.Sprintf("%s/views", rp)), jet.InDevelopmentMode())
	ds.JetTemplate = jetViews

	//set up render engine
	ds.createRenderer()

	//set routes
	ds.Routes = ds.routes().(*chi.Mux)

	return nil
}

//InitDir creates initial directories for Dragon Spider application
func (ds *DragonSpider) InitDir(p initPaths) error {
	root := p.rootPath

	for _, path := range p.dirNames {
		//create directory if it doesn't exist'
		err := ds.CreateDir(fmt.Sprintf("%s/%s", root, path))
		if err != nil {
			return err
		}
	}

	return nil
}

//createDotEnv creates a .env file
func (ds *DragonSpider) createDotEnv(path string) error {
	err := ds.CreateFile(fmt.Sprintf("%s/.env", path))
	if err != nil {
		return err
	}

	return nil
}

//createSqliteDb creates a sqlite database file
func (ds *DragonSpider) createSqliteDb(path string) error {
	if !FileExists(path + "/sqlite/app.db") {
		err := ds.CreateFile(fmt.Sprintf("%s/sqlite/app.db", path))
		if err != nil {
			return err
		}
	}
	return nil
}

//startLoggers creates and returns info logger and error logger
func (ds *DragonSpider) startLoggers() (*log.Logger, *log.Logger) {
	//TODO log info and error to files
	var infoLog *log.Logger
	var errorLog *log.Logger

	currentOS := runtime.GOOS

	//add color to log output if running on linux system
	if currentOS == "linux" {
		infoLog = log.New(os.Stdout, "\033[33mINFO\033[0m:\t", log.Ldate|log.Ltime)
		errorLog = log.New(os.Stdout, "\033[31mERROR\033[0m:\t", log.Ldate|log.Ltime|log.Lshortfile)

		return infoLog, errorLog
	}

	infoLog = log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog

}

//ListenAndServe starts the web server
func (ds *DragonSpider) ListenAndServe() {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		ErrorLog:     ds.ErrorLog,
		Handler:      ds.Routes,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	defer func(Pool *sql.DB) {
		err := Pool.Close()
		if err != nil {

		}
	}(ds.Db.Pool)

	ds.InfoLog.Printf("Listening on port: %s", os.Getenv("PORT"))

	err := srv.ListenAndServe()
	if err != nil {
		ds.ErrorLog.Fatal(err)
	}
}

//createRenderer creates a render engine for template files
func (ds *DragonSpider) createRenderer() {
	engine := render.Render{
		Renderer:    ds.config.renderer,
		JetTemplate: ds.JetTemplate,
		RootPath:    ds.RootPath,
		//Secure:     false,
		Port: ds.config.port,
		//ServerName: "",
		Session: ds.Session,
	}

	ds.Render = &engine
}

func (ds *DragonSpider) CreateDSN() string {
	var dsn string
	timeout := 5

	switch strings.ToLower(os.Getenv("DATABASE_TYPE")) {
	case "postgres", "postgresql":
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s timezone=UTCtimeout=%d",
			os.Getenv("DATABASE_HOST"),
			os.Getenv("DATABASE_PORT"),
			os.Getenv("DATABASE_USER"),
			os.Getenv("DATABASE_NAME"),
			os.Getenv("DATABASE_SSL_MODE"),
			timeout)
		if os.Getenv("DATABASE_PASSWORD") != "" {
			dsn = fmt.Sprintf("%s password=%s", dsn, os.Getenv("DATABASE_PASSWORD"))
		}

	default:
		if os.Getenv("SQLITE_FILE") != "" {
			dsn = os.Getenv("SQLITE_FILE")
		} else {
			dsn = "sqlite/app.db"
		}
	}

	return dsn
}
