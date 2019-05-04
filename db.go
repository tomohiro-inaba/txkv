package txkv

import (
	"errors"
	"fmt"
	"sync"
)

type DB struct {
	storage *CSVStorage
	rwlock  sync.RWMutex
	txs     map[string]*Tx
}

func Open(path string) (*DB, error) {
	csv, err := NewCSVStorage(path)
	if err != nil {
		return nil, err
	}

	db := &DB{}
	db.storage = csv
	db.txs = make(map[string]*Tx)
	return db, nil
}

func (db *DB) Begin(writable bool) *Tx {
	if writable {
		db.rwlock.Lock()
	} else {
		db.rwlock.RLock()
	}
	tx := NewTx(db, writable)
	db.txs[tx.id] = tx
	return tx
}

func (db *DB) Close() {
	for _, tx := range db.txs {
		tx.Rollback()
	}
	db.storage.close()
}

func (db *DB) GetTx(txid string) (*Tx, error) {
	if tx, ok := db.txs[txid]; ok {
		return tx, nil
	}
	return nil, errors.New(fmt.Sprintf("Not found: txid=%s\n", txid))
}

func (db *DB) read(key Key) (Value, error) {
	return db.storage.Read(key)
}

func (db *DB) write(key Key, value Value) error {
	return db.storage.Write(key, value)
}

func (db *DB) gc() (bool, error) {
	return db.storage.GC()
}
