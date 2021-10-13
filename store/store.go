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

type Storer interface {
	Create(time.Duration, string) error
	ReadAll() error
}

type StoreStrategy struct {
	Pipe bool
	File bool
	Args bool
}

type Store struct {
	Storer
	Strategy StoreStrategy
}

func (s Store) Save(ttl string, args []string) error {
	expiry := getExpiryFromTTL(ttl)

	switch true {
	case s.Strategy.Pipe:
		return s.pipeStrat(expiry)
	case s.Strategy.File:
		return s.fileStrat(expiry, args[0])
	case s.Strategy.Args:
		return s.argStrat(expiry, args[0])
	default:
		return errors.New("no strategy discerned, this error shouldn't happen, scream and cry plz")
	}
}

func (s Store) Show() {
	if err := s.ReadAll(); err != nil {
		panic(err)
	}
}

func NewKVSStore() *Store {
	return &Store{kvs.NewKVS(), StoreStrategy{}}
}

func (s Store) pipeStrat(expiry time.Duration) error {
	var data string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		data += "\n" + scanner.Text()
	}

	return s.Create(expiry, data)
}

func (s Store) fileStrat(expiry time.Duration, file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("unable to read file %q: %v", file, err)
	}

	return s.Create(expiry, string(data))
}

func (s Store) argStrat(expiry time.Duration, data string) error {
	return s.Create(expiry, data)
}

func getExpiryFromTTL(ttl string) time.Duration {
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

	return dur
}
