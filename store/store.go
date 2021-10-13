package store

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/awilcots/tapi/kvs"
)

type StoreStrategy struct {
	Pipe bool
	File bool
	Args bool
}

type Store struct {
	Storer
}

type Storer interface {
	Create(time.Time, string) error
}

func (ss StoreStrategy) Store(ttl string, args []string) error {
	store := newKVSStore()
	expiry := getExpiryFromTTL(ttl)

	switch true {
	case ss.Pipe:
		return store.pipeStrat(expiry)
	case ss.File:
		return store.fileStrat(expiry, args[0])
	case ss.Args:
		return store.argStrat(expiry, args[0])
	default:
		return errors.New("no strategy discerned, this error shouldn't happen, scream and cry plz")
	}
}

func newKVSStore() *Store {
	return &Store{kvs.NewKVS()}
}

func (s Store) pipeStrat(expiry time.Time) error {
	var data string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		data += "\n" + scanner.Text()
	}

	return s.Create(expiry, data)
}

func (s Store) fileStrat(expiry time.Time, file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("unable to read file %q: %v", file, err)
	}

	return s.Create(expiry, string(data))
}

func (s Store) argStrat(expiry time.Time, data string) error {
	return s.Create(expiry, data)
}

func getExpiryFromTTL(ttl string) time.Time {
	numS := ttl[:len(ttl)-1]
	unit := ttl[len(ttl)-1]

	if unit == 'd' {
		num, _ := strconv.Atoi(numS)
		num *= 24

		ttl = strconv.Itoa(num) + "h"
	}

	dur, err := time.ParseDuration(ttl)
	if err != nil {
		panic(err)
	}

	return time.Now().Add(dur)
}
