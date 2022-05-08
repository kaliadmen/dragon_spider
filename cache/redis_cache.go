package cache

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type RedisCache struct {
	Connection *redis.Pool
	Prefix     string
}

func (rc *RedisCache) Has(str string) (bool, error) {
	key := fmt.Sprintf("%s:%s", rc.Prefix, str)
	conn := rc.Connection.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(conn)

	ok, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (rc *RedisCache) Get(str string) (interface{}, error) {
	key := fmt.Sprintf("%s:%s", rc.Prefix, str)
	conn := rc.Connection.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(conn)

	cacheEntry, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	decoded, err := decode(string(cacheEntry))
	if err != nil {
		return nil, err
	}

	item := decoded[key]

	return item, nil
}

func (rc *RedisCache) Set(str string, value interface{}, expires ...int) error {
	if str == "" || value == "" {
		return errors.New("blank entries are not allowed")
	}

	key := fmt.Sprintf("%s:%s", rc.Prefix, str)
	conn := rc.Connection.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(conn)

	entry := Entry{}
	entry[key] = value

	encoded, err := encode(entry)
	if err != nil {
		return err
	}

	if len(expires) > 0 {
		_, err := conn.Do("SETEX", key, expires[0], string(encoded))
		if err != nil {
			return err
		}
	} else {
		_, err := conn.Do("SET", key, string(encoded))
		if err != nil {
			return err
		}
	}

	return nil
}

func (rc *RedisCache) Delete(str string) error {
	key := fmt.Sprintf("%s:%s", rc.Prefix, str)
	conn := rc.Connection.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(conn)

	_, err := conn.Do("DEL", key)
	if err != nil {
		return err
	}

	return nil
}

func (rc *RedisCache) DeleteByMatch(str string) error {
	key := fmt.Sprintf("%s:%s", rc.Prefix, str)
	conn := rc.Connection.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(conn)

	keys, err := rc.getKeys(key)
	if err != nil {
		return err
	}

	for _, x := range keys {
		_, err := conn.Do("DEL", x)
		if err != nil {
			return err
		}
	}

	return nil
}

func (rc *RedisCache) DeleteAll() error {
	key := fmt.Sprintf("%s:", rc.Prefix)
	conn := rc.Connection.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(conn)

	keys, err := rc.getKeys(key)
	if err != nil {
		return err
	}

	for _, x := range keys {
		_, err := conn.Do("DEL", x)
		if err != nil {
			return err
		}
	}

	return nil
}

func (rc *RedisCache) getKeys(pattern string) ([]string, error) {
	conn := rc.Connection.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(conn)

	i := 0
	var keys []string

	for {
		arr, err := redis.Values(conn.Do("SCAN", i, "MATCH", fmt.Sprintf("%s*", pattern)))
		if err != nil {
			return keys, err
		}

		i, _ = redis.Int(arr[0], nil)
		k, _ := redis.Strings(arr[1], nil)
		keys = append(keys, k...)

		if i == 0 {
			break
		}
	}

	return keys, nil
}
