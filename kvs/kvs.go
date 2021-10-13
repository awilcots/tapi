package kvs

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dgraph-io/badger"
)

const dbLoc = "/tmp/badger"

type KVS struct {
	*badger.DB
}

func NewKVS() *KVS {
	db, err := badger.Open(badger.DefaultOptions(dbLoc).WithLogger(nil))
	if err != nil {
		panic(err)
	}

	return &KVS{db}
}

func (k KVS) Create(expiry time.Duration, data string) error {
	key := []byte(strconv.Itoa(int(time.Now().Unix())))

	return k.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry(key, []byte(data)).WithTTL(expiry)
		return txn.SetEntry(e)
	})
}

func (k KVS) ReadAll() error {
	k.View(func(txn *badger.Txn) error {
		iter := txn.NewIterator(badger.DefaultIteratorOptions)
		defer iter.Close()
		for iter.Rewind(); iter.Valid(); iter.Next() {
			item := iter.Item()

			expiresAt := time.Unix(int64(item.ExpiresAt()), 0).Format(time.RFC850)
			fmt.Printf("Expires At: %s\n\n", expiresAt)

			item.Value(func(val []byte) error {
				fmt.Println(string(val))
				return nil
			})

			fmt.Println("<--------------------------------->")
		}

		return nil
	})
	return nil
}
