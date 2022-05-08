package cache

import "testing"

func TestBadgerCache_Has(t *testing.T) {
	err := badgerTestCache.Delete("foo")
	if err != nil {
		t.Error(err)
	}

	inCache, err := badgerTestCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("foo is in cache, when it should not be")
	}

	err = badgerTestCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	inCache, err = badgerTestCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if !inCache {
		t.Error("foo is not in cache, when it should be")
	}

}

func TestBadgerCache_Get(t *testing.T) {
	err := badgerTestCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	val, err := badgerTestCache.Get("foo")
	if err != nil {
		t.Error(err)
	}

	if val != "bar" {
		t.Error("incorrect value received from cache")
	}

}

func TestBadgerCache_Set(t *testing.T) {
	if err := badgerTestCache.Set("", ""); err == nil {
		t.Error("set allowed a blank entry")
	}

	err := badgerTestCache.Set("foo", "bar", 1000)
	if err != nil {
		t.Error(err)
	}

	inCache, err := badgerTestCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if !inCache {
		t.Error("foo is not cache, when it should be")
	}

}

func TestBadgerCache_Delete(t *testing.T) {
	err := badgerTestCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	err = badgerTestCache.Delete("foo")
	if err != nil {
		t.Error(err)
	}

	inCache, err := badgerTestCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("foo is in cache, when it should have been deleted")
	}
}

func TestBadgerCache_DeleteByMatch(t *testing.T) {
	err := badgerTestCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	err = badgerTestCache.Set("bar", "foo")
	if err != nil {
		t.Error(err)
	}

	err = badgerTestCache.DeleteByMatch("foo")
	if err != nil {
		t.Error(err)
	}

	inCache, err := badgerTestCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("foo is in cache, when it should have been deleted")
	}

	inCache, err = badgerTestCache.Has("bar")
	if err != nil {
		t.Error(err)
	}

	if !inCache {
		t.Error("bar is not in cache, when it should be")
	}

	err = badgerTestCache.DeleteAll()
	if err != nil {
		panic(err)
	}
}

func TestBadgerCache_DeleteAll(t *testing.T) {
	err := badgerTestCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	err = badgerTestCache.DeleteAll()
	if err != nil {
		t.Error(err)
	}

	inCache, err := badgerTestCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("foo is in cache, when it should have been deleted")
	}

}
