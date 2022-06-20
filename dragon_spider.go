package dragonSpider

import (
	"database/sql"
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
	"github.com/dgraph-io/badger"
	"github.com/go-chi/chi/v5"
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
	"github.com/kaliadmen/dragon_spider/cache"
	"github.com/kaliadmen/dragon_spider/filesystems/miniofs"
	"github.com/kaliadmen/dragon_spider/filesystems/s3fs"
	"github.com/kaliadmen/dragon_spider/filesystems/sftpfs"
	"github.com/kaliadmen/dragon_spider/render"
	"github.com/kaliadmen/dragon_spider/session"
	"github.com/kaliadmen/mailer"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const version = "1.0.1"

var appRedisCache *cache.RedisCache
var redisPool *redis.Pool
var appBadgerCache *cache.BadgerCache
var badgerConnection *badger.DB
var logFile *os.File

//var mailer mailer.Mail

//DragonSpider is an overall type for the Dragon Spider package.
//Members exported in this type are available to any application that uses it
type DragonSpider struct {
	AppName       string
	Debug         bool
	Version       string
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	Render        *render.Render
	JetTemplate   *jet.Set
	Routes        *chi.Mux
	Session       *scs.SessionManager
	Db            Database
	RootPath      string
	config        config
	EncryptionKey string
	Cache         cache.Cache
	Scheduler     *cron.Cron
	Mail          mailer.Mail
	Server        Server
	FileSystems   map[string]any
	Minio         miniofs.Minio
	S3            s3fs.S3
	SFTP          sftpfs.SFTP
}

type config struct {
	port        string
	renderer    string
	sessionType string
	cookie      cookieConfig
	database    databaseConfig
	redis       redisConfig
	uploads     uploadConfig
}

type Server struct {
	ServerName string
	Port       string
	Secure     bool
	URL        string
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
			"mail",
			"views",
			"data",
			"migrations",
			"public",
			"tmp",
			"logs",
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

	//set application variables
	ds.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	ds.Version = version
	ds.RootPath = rp
	ds.EncryptionKey = os.Getenv("KEY")

	//set mailer
	if os.Getenv("SMTP_PORT") != "" {
		ds.Mail = ds.createMailer()
	}

	//create loggers
	infoLog, errorLog, file := ds.startLoggers()
	ds.InfoLog = infoLog
	ds.ErrorLog = errorLog
	logFile = file

	//connect to a database
	if os.Getenv("DATABASE_TYPE") != "" {
		err := ds.useDatabase(strings.ToLower(os.Getenv("DATABASE_TYPE")))
		if err != nil {
			ds.ErrorLog.Println(err)
		}

	}

	//initialize scheduler
	scheduler := cron.New()
	ds.Scheduler = scheduler

	//create a cache
	if strings.ToLower(os.Getenv("CACHE")) == "redis" || os.Getenv("SESSION_TYPE") == "redis" {
		appRedisCache = ds.createRedisClientCache()
		ds.Cache = appRedisCache
		redisPool = appRedisCache.Connection
	}

	if strings.ToLower(os.Getenv("CACHE")) == "badger" || os.Getenv("SESSION_TYPE") == "badger" {
		appBadgerCache = ds.createBadgerClientCache()
		ds.Cache = appBadgerCache
		badgerConnection = appBadgerCache.Connection

		//garbage collection
		_, err := ds.Scheduler.AddFunc("@daily", func() {
			_ = appBadgerCache.Connection.RunValueLogGC(0.7)
		})
		if err != nil {
			return err
		}

	}

	//file uploads
	types := strings.Split(os.Getenv("ALLOWED_FILETYPES"), ",")
	var mimeTypes []string

	for _, mimetype := range types {
		mimeTypes = append(mimeTypes, mimetype)
	}

	var maxUploadSize int64
	if maxSize, err := strconv.Atoi(os.Getenv("MAX_UPLOAD_SIZE")); err != nil {
		//15 MB
		maxUploadSize = 15 << 20
	} else {
		maxUploadSize = int64(maxSize)
	}

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
		redis: redisConfig{
			host:     os.Getenv("REDIS_HOST"),
			password: os.Getenv("REDIS_PASSWORD"),
			prefix:   os.Getenv("CACHE_PREFIX"),
		},
		uploads: uploadConfig{
			maxUploadSize:    maxUploadSize,
			allowedMimeTypes: mimeTypes,
		},
	}
	//set app server variables
	secure := true
	if strings.ToLower(os.Getenv("SECURE")) == "false" {
		secure = false
	}

	ds.Server = Server{
		ServerName: os.Getenv("SERVER_NAME"),
		Port:       os.Getenv("PORT"),
		Secure:     secure,
		URL:        os.Getenv("APP_URL"),
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

	switch ds.config.sessionType {
	case "redis":
		ses.RedisPool = appRedisCache.Connection

	case "badger":
		ses.BadgerConnection = appBadgerCache.Connection

	case "mysql", "sqlite", "postgres", "postgresql", "mariadb":
		ses.DbPool = ds.Db.Pool
	}

	//create session
	ds.Session = ses.InitSession()

	if ds.Debug {
		var jetViews = jet.NewSet(
			jet.NewOSFileSystemLoader(fmt.Sprintf("%s/views", rp)), jet.InDevelopmentMode())
		ds.JetTemplate = jetViews
	} else {
		var jetViews = jet.NewSet(
			jet.NewOSFileSystemLoader(fmt.Sprintf("%s/views", rp)))
		ds.JetTemplate = jetViews
	}

	//set up render engine
	ds.createRenderer()

	//set routes
	ds.Routes = ds.routes().(*chi.Mux)

	//listen for mail on mail channel
	if ds.Mail.Port != 0 {
		go ds.Mail.ListenForMail()
		ds.InfoLog.Println("Listening for mail on port", ds.Mail.Port)
	} else {
		ds.InfoLog.Println("No mailer port set in .env file")
	}

	ds.FileSystems = ds.CreateFileSystem()

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

func (ds *DragonSpider) useDatabase(dbType string) error {
	pool, err := ds.OpenDb(dbType, ds.CreateDSN())
	if err != nil {
		ds.ErrorLog.Println(err)
		if dbType == "sqlite" {
			ds.InfoLog.Println("Check your env file or try running ./dragon_spider make database")
		}
		os.Exit(1)
	}
	ds.Db = Database{
		DatabaseType: dbType,
		Pool:         pool,
	}

	return nil
}

//createRedisPool creates a redis database pool for cache
func (ds *DragonSpider) createRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   10000,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				ds.config.redis.host,
				redis.DialPassword(ds.config.redis.password))
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
	}
}

//createBadgerConnection creates a new connection to badger db
func (ds *DragonSpider) createBadgerConnection() (*badger.DB, error) {
	if ds.Debug {
		db, err := badger.Open(badger.DefaultOptions(ds.RootPath + "/db/badger"))
		if err != nil {
			return nil, err
		}

		return db, nil
	}

	db, err := badger.Open(badger.DefaultOptions(ds.RootPath + "/db/badger").WithLogger(nil))
	if err != nil {
		return nil, err
	}

	return db, nil
}

//createRedisClientCache creates a redis cache
func (ds *DragonSpider) createRedisClientCache() *cache.RedisCache {
	cacheClient := cache.RedisCache{
		Connection: ds.createRedisPool(),
		Prefix:     ds.config.redis.prefix,
	}

	return &cacheClient
}

//createBadgerClientCache creates a badger db cache
func (ds *DragonSpider) createBadgerClientCache() *cache.BadgerCache {
	connection, err := ds.createBadgerConnection()
	if err != nil {
		ds.ErrorLog.Println("Couldn't create badger connection", err)
		return nil
	}

	cacheClient := cache.BadgerCache{
		Connection: connection,
		Prefix:     os.Getenv("CACHE_PREFIX"),
	}

	return &cacheClient
}

func (ds *DragonSpider) createMailer() mailer.Mail {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	m := mailer.Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Templates:   ds.RootPath + "/mail",
		Host:        os.Getenv("SMTP_HOST"),
		Port:        port,
		Username:    os.Getenv("SMTP_USERNAME"),
		Password:    os.Getenv("SMTP_PASSWORD"),
		Encryption:  os.Getenv("SMTP_ENCRYPTION"),
		FromAddress: os.Getenv("FROM_ADDRESS"),
		FromName:    os.Getenv("FROM_NAME"),
		Jobs:        make(chan mailer.Message, 20),
		Results:     make(chan mailer.Result, 20),
		API:         os.Getenv("MAILER_API"),
		APIKey:      os.Getenv("MAILER_KEY"),
		APIUrl:      os.Getenv("MAILER_URL"),
	}

	return m
}

//startLoggers creates and returns info logger and error logger
func (ds *DragonSpider) startLoggers() (*log.Logger, *log.Logger, *os.File) {

	var infoLog *log.Logger
	var errorLog *log.Logger
	var loggerFile *os.File

	currentOS := runtime.GOOS

	switch ds.Debug {
	case false:
		file, err := os.OpenFile(ds.RootPath+"/logs/appLog.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			_ = fmt.Sprintf("error opening file: %v", err)
		}

		infoLog = log.New(file, "INFO:\t", log.Ldate|log.Ltime)
		errorLog = log.New(file, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
		logFile = file

		return infoLog, errorLog, loggerFile

	default:
		//add color to log output if running on linux system
		if currentOS == "linux" {
			infoLog = log.New(os.Stdout, "\033[33mINFO\033[0m:\t", log.Ldate|log.Ltime)
			errorLog = log.New(os.Stdout, "\033[31mERROR\033[0m:\t", log.Ldate|log.Ltime|log.Lshortfile)

			return infoLog, errorLog, nil
		}

		infoLog = log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
		errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

		return infoLog, errorLog, nil

	}

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

	//close database and cache connections
	if ds.Db.Pool != nil {
		defer func(Pool *sql.DB) {
			err := Pool.Close()
			if err != nil {
				ds.ErrorLog.Println("Could not close database connection", err)
			}
		}(ds.Db.Pool)
	}

	if redisPool != nil {
		defer func(redisPool *redis.Pool) {
			err := redisPool.Close()
			if err != nil {
				ds.ErrorLog.Println("Could not close redis connection", err)
			}
		}(redisPool)
	}

	if badgerConnection != nil {
		defer func(badgerConnection *badger.DB) {
			err := badgerConnection.Close()
			if err != nil {
				ds.ErrorLog.Println("Could not close badger db connection", err)
			}
		}(badgerConnection)
	}

	if logFile != nil {
		defer func(logFile *os.File) {
			err := logFile.Close()
			if err != nil {
				ds.ErrorLog.Println("Could not close log file", err)
			}
		}(logFile)
	}

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
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DATABASE_USER")
	dbPass := os.Getenv("DATABASE_PASSWORD")
	dbName := os.Getenv("DATABASE_NAME")
	sslMode := os.Getenv("DATABASE_SSL_MODE")
	timeout := 5

	switch strings.ToLower(os.Getenv("DATABASE_TYPE")) {
	case "postgres", "postgresql":
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s timezone=UTCtimeout=%d",
			dbHost,
			dbPort,
			dbUser,
			dbName,
			sslMode,
			timeout)
		if dbPass != "" {
			dsn = fmt.Sprintf("%s password=%s", dsn, dbPass)
		}

	case "mysql", "mariadb":
		//username:password@protocol(address)/dbname?param=value
		dsn = fmt.Sprintf("%s@tcp(%s:%s)/%s?tls=%s&timeout=%d",
			dbUser,
			dbHost,
			dbPort,
			dbName,
			sslMode,
			timeout)
		if dbPass != "" {
			dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?tls=%s&timeout=%ds",
				dbUser,
				dbPass,
				dbHost,
				dbPort,
				dbName,
				sslMode,
				timeout)
		}

	case "sqlite":
		if os.Getenv("SQLITE_FILE") != "" {
			dsn = os.Getenv("SQLITE_FILE")
		} else {
			dsn = ds.RootPath + "/db/sqlite/app.db"
		}

	default:
		return dsn
	}

	return dsn
}

func (ds *DragonSpider) CreateFileSystem() map[string]any {
	fileSystems := make(map[string]any)

	if os.Getenv("MINIO_SECRET") != "" {
		useSSL := false
		if strings.ToLower(os.Getenv("MINIO_USESSL")) == "true" {
			useSSL = true
		}

		minio := miniofs.Minio{
			Endpoint: os.Getenv("MINIO_ENDPOINT"),
			Key:      os.Getenv("MINIO_KEY"),
			Secret:   os.Getenv("MINIO_SECRET"),
			UseSSL:   useSSL,
			Region:   os.Getenv("MINIO_REGION"),
			Bucket:   os.Getenv("MINIO_BUCKET"),
		}

		fileSystems["MINIO"] = minio
		ds.Minio = minio
	}

	if os.Getenv("SFTP_HOST") != "" {
		sftp := sftpfs.SFTP{
			Host:     os.Getenv("SFTP_HOST"),
			User:     os.Getenv("SFTP_USER"),
			Password: os.Getenv("SFTP_PASSWORD"),
			Port:     os.Getenv("SFTP_PORT"),
		}

		fileSystems["SFTP"] = sftp
		ds.SFTP = sftp
	}

	if os.Getenv("S3_KEY") != "" {
		s3 := s3fs.S3{
			Key:      os.Getenv("S3_KEY"),
			Secret:   os.Getenv("S3_SECRET"),
			Region:   os.Getenv("S3_REGION"),
			Endpoint: os.Getenv("S3_ENDPOINT"),
			Bucket:   os.Getenv("S3_BUCKET"),
		}

		fileSystems["S3"] = s3
		ds.S3 = s3
	}

	return fileSystems
}
