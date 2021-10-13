package kvs

import (
	"fmt"
	"time"
)

type KVS struct{}

func NewKVS() *KVS {
	return new(KVS)
}

func (k KVS) Create(expiry time.Time, data string) error {
	fmt.Println(expiry, ":", data)
	return nil
}
