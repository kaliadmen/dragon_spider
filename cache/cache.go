package cache

import (
	"bytes"
	"encoding/gob"
)

type Cache interface {
	Has(string) (bool, error)
	Get(string) (interface{}, error)
	Set(string, interface{}, ...int) error
	Delete(string2 string) error
	DeleteByMatch(string) error
	DeleteAll() error
}

type Entry map[string]interface{}

func encode(item Entry) ([]byte, error) {
	bBuffer := bytes.Buffer{}
	e := gob.NewEncoder(&bBuffer)
	err := e.Encode(item)
	if err != nil {
		return nil, err
	}
	return bBuffer.Bytes(), nil
}

func decode(str string) (Entry, error) {
	item := Entry{}
	bBuffer := bytes.Buffer{}
	bBuffer.Write([]byte(str))
	d := gob.NewDecoder(&bBuffer)
	err := d.Decode(&item)
	if err != nil {
		return nil, err
	}
	return item, nil
}
