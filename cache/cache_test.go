package cache

import (
	"testing"
)

func TestRedisCache_Has(t *testing.T) {
	err := redisTestCache.Delete("foo")
	if err != nil {
		t.Error(err)
	}

	inCache, err := redisTestCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("foo is in cache, when it should not be")
	}

	err = redisTestCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	inCache, err = redisTestCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if !inCache {
		t.Error("foo is not in cache, when it should be")
	}

}

func TestRedisCache_Get(t *testing.T) {
	err := redisTestCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	val, err := redisTestCache.Get("foo")
	if err != nil {
		t.Error(err)
	}

	if val != "bar" {
		t.Error("incorrect value received from cache")
	}

}

func TestRedisCache_Set(t *testing.T) {
	err := redisTestCache.Set("foo", "bar", 1000)
	if err != nil {
		t.Error(err)
	}

	inCache, err := redisTestCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if !inCache {
		t.Error("foo is not cache, when it should be")
	}

}

func TestRedisCache_Delete(t *testing.T) {
	err := redisTestCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	err = redisTestCache.Delete("foo")
	if err != nil {
		t.Error(err)
	}

	inCache, err := redisTestCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("foo is in cache, when it should have been deleted")
	}
}

func TestRedisCache_DeleteByMatch(t *testing.T) {
	err := redisTestCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	err = redisTestCache.Set("bar", "foo")
	if err != nil {
		t.Error(err)
	}

	err = redisTestCache.DeleteByMatch("foo")
	if err != nil {
		t.Error(err)
	}

	inCache, err := redisTestCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("foo is in cache, when it should have been deleted")
	}

	inCache, err = redisTestCache.Has("bar")
	if err != nil {
		t.Error(err)
	}

	if !inCache {
		t.Error("bar is not in cache, when it should be")
	}

	err = redisTestCache.DeleteAll()
	if err != nil {
		panic(err)
	}
}

func TestRedisCache_DeleteAll(t *testing.T) {
	err := redisTestCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	err = redisTestCache.DeleteAll()
	if err != nil {
		t.Error(err)
	}

	inCache, err := redisTestCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("foo is in cache, when it should have been deleted")
	}

}
