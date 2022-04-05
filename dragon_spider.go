package dragonSpider

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/kaliadmen/dragon_spider/render"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"
)

const version = "1.0.0"

type DragonSpider struct {
	AppName  string
	Debug    bool
	Version  string
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	Render   *render.Render
	Routes   *chi.Mux
	RootPath string
	config   config
}

type config struct {
	port     string
	renderer string
}

//New creates directories structure for application and assign root path to application
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

	//create loggers
	infoLog, errorLog := ds.startLoggers()
	ds.InfoLog = infoLog
	ds.ErrorLog = errorLog

	//set application variables
	ds.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	ds.Version = version
	ds.RootPath = rp

	//set config variables
	ds.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"),
	}

	//set up render engine
	ds.createRenderer()

	//set routes
	ds.Routes = ds.routes().(*chi.Mux)

	return nil
}

//InitDir creates initial directories
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

	ds.InfoLog.Printf("Listening on port: %s", os.Getenv("PORT"))

	err := srv.ListenAndServe()
	if err != nil {
		ds.ErrorLog.Fatal(err)
	}
}

//createRenderer creates a render engine for template files
func (ds *DragonSpider) createRenderer() {
	engine := render.Render{
		Renderer: ds.config.renderer,
		RootPath: ds.RootPath,
		//Secure:     false,
		Port: ds.config.port,
		//ServerName: "",
	}

	ds.Render = &engine
}
