package main

import (
	"os"
	"sync"
)

type DB struct {
	rwlock sync.RWMutex
	*dal
}

func Open(path string, options *Options) (*DB, error) {
	var err error

	options.pageSize = os.Getpagesize()
	dal, err := newDal(path, Options)
	if err != nil {
		return nil, err
	}

	db := &DB{
		sync.RWMutex,
		dal,
	}

	return db, nil
}

func (db *DB) Close() error {
	return db.close()
}

func (db *DB) ReadTx() *tx {
	db.rwlock.RLock()
	return newTx(db, False)
}

func (db *DB) WriteTx() *tx {
	db.rwlock.Lock()
	return newTx(db, True)
}
