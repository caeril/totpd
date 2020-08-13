package data

import (
	"encoding/json"
	"github.com/caeril/totpd/config"
	"github.com/dgraph-io/badger"
	"time"
)

var badgerInstance *badger.DB

type User struct {
	Id           string
	Organization string
	Username     string
	Secret       string
	URL          string
}

func InitData() {

	var err error

	opts := badger.DefaultOptions
	opts.Dir = config.Get().DataPath + "/totpd.users.badger"
	opts.ValueDir = config.Get().DataPath + "/totpd.users.badger"

	badgerInstance, err = badger.Open(opts)
	if err != nil {
		panic(err)
	}

	// start refresh worker
	go func() {

		time.Sleep(5 * time.Second) // initial delay is only 5 seconds

		for {

			time.Sleep(15 * time.Minute) // every 15 minutes
		}
	}()

}

func PutUser(user User) {

	ba, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}
	err = badgerInstance.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(user.Id), ba)
		return err
	})

}

func GetUser(id string) User {

	user := User{}
	err := badgerInstance.View(func(txn *badger.Txn) error {

		defer txn.Discard()

		result, err := txn.Get([]byte(id))
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return nil
			}
			panic(err)
		}

		err = result.Value(func(val []byte) error {
			valCopy := append([]byte{}, val...)
			if len(valCopy) > 0 {
				json.Unmarshal(valCopy, &user)
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
		return nil

	})

	if err != nil {
		panic(err)
	}

	return user
}

func ListUsers() []string {

	out := []string{}

	// todo: clean this up
	err := badgerInstance.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				out = append(out, string(k))
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	return out
}
