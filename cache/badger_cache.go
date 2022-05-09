package cache

import (
	"errors"
	"fmt"
	"github.com/dgraph-io/badger"
	"time"
)

type BadgerCache struct {
	Connection *badger.DB
	Prefix     string
}

func (bc *BadgerCache) Has(key string) (bool, error) {
	_, err := bc.Get(key)
	if err != nil {
		return false, nil
	}

	return true, nil
}

func (bc *BadgerCache) Get(key string) (interface{}, error) {
	key = fmt.Sprintf("%s:%s", bc.Prefix, key)
	var fromCache []byte

	err := bc.Connection.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			fromCache = append([]byte{}, val...)
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	decoded, err := decode(string(fromCache))
	if err != nil {
		return nil, err
	}

	item := decoded[key]

	return item, nil
}

func (bc *BadgerCache) Set(key string, value interface{}, expires ...int) error {
	if key == "" || value == "" {
		return errors.New("blank entries are not allowed")
	}
	key = fmt.Sprintf("%s:%s", bc.Prefix, key)

	entry := Entry{}

	entry[key] = value
	encoded, err := encode(entry)
	if err != nil {
		return err
	}

	if len(expires) > 0 {
		err = bc.Connection.Update(func(txn *badger.Txn) error {
			e := badger.NewEntry([]byte(key), encoded).WithTTL(time.Second * time.Duration(expires[0]))
			err = txn.SetEntry(e)
			return err
		})
	} else {
		err = bc.Connection.Update(func(txn *badger.Txn) error {
			e := badger.NewEntry([]byte(key), encoded)
			err = txn.SetEntry(e)
			return err
		})
	}

	return nil
}

func (bc *BadgerCache) Delete(key string) error {
	key = fmt.Sprintf("%s:%s", bc.Prefix, key)

	err := bc.Connection.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(key))
		return err
	})

	return err
}

func (bc *BadgerCache) DeleteByMatch(key string) error {
	return bc.deleteByMatch(key)
}

func (bc *BadgerCache) DeleteAll() error {

	return bc.deleteByMatch("")
}

func (bc *BadgerCache) deleteByMatch(key string) error {
	key = fmt.Sprintf("%s:%s", bc.Prefix, key)
	//search cache for keys to be deleted
	deletedKeys := func(keysToDelete [][]byte) error {
		if err := bc.Connection.Update(func(txn *badger.Txn) error {
			for _, k := range keysToDelete {
				if err := txn.Delete(k); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return err
		}
		return nil
	}

	//limit how much cache is used at one time
	collectionSize := 100000

	err := bc.Connection.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.AllVersions = false
		opts.PrefetchValues = false
		i := txn.NewIterator(opts)

		defer i.Close()

		keysToDelete := make([][]byte, 0, collectionSize)
		keysCollected := 0

		for i.Seek([]byte(key)); i.ValidForPrefix([]byte(key)); i.Next() {
			k := i.Item().KeyCopy(nil)
			keysToDelete = append(keysToDelete, k)
			keysCollected++
			if keysCollected == collectionSize {
				if err := deletedKeys(keysToDelete); err != nil {
					return err
				}
			}
		}
		if keysCollected > 0 {
			if err := deletedKeys(keysToDelete); err != nil {
				return err
			}
		}

		return nil
	})

	return err

}
