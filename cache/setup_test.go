package cache

import (
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/dgraph-io/badger"
	"github.com/gomodule/redigo/redis"
	"log"
	"os"
	"testing"
	"time"
)

var redisTestCache RedisCache
var badgerTestCache BadgerCache

func TestMain(m *testing.M) {
	//setup for redis test
	mini, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	defer mini.Close()

	//create a redis pool using miniredis
	pool := redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", mini.Addr())
		},
		MaxIdle:     50,
		MaxActive:   1000,
		IdleTimeout: 240 * time.Second,
	}

	redisTestCache.Connection = &pool
	redisTestCache.Prefix = "test_dragon"

	defer func(Connection *redis.Pool) {
		err := Connection.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(redisTestCache.Connection)

	//setup for badger db
	//remove old badger test database
	_ = os.RemoveAll("./testdata/tmp/badger")

	//create directories
	if _, err := os.Stat("./testdata/tmp/"); os.IsNotExist(err) {
		err := os.Mkdir("./testdata/tmp/", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = os.MkdirAll("./testdata/tmp/badger", 0755)
	if err != nil {
		log.Fatal(err)
	}

	//create new badger database and cache
	db, _ := badger.Open(badger.DefaultOptions("./testdata/tmp/badger"))
	badgerTestCache.Connection = db
	badgerTestCache.Prefix = "test_dragon"

	os.Exit(m.Run())
}
