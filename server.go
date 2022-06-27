package dragonSpider

import (
	"database/sql"
	"fmt"
	"github.com/dgraph-io/badger"
	"github.com/gomodule/redigo/redis"
	"net/http"
	"os"
	"time"
)

func (ds *DragonSpider) ListenAndServe() error {
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
	go ds.listenRPC()

	return srv.ListenAndServe()

}
