package cache

import (
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/gomodule/redigo/redis"
	"os"
	"testing"
	"time"
)

var redisTestCache RedisCache

func TestMain(m *testing.M) {
	mini, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	defer mini.Close()

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

	os.Exit(m.Run())
}
