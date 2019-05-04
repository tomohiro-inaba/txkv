package txkv

import (
	"errors"

	"github.com/google/uuid"
)

type Tx struct {
	id       string
	db       *DB
	writable bool
	entries  map[Key]Value
	dirty    map[Key]bool
}

func NewTx(db *DB, writable bool) *Tx {
	id := uuid.New().String()
	entries := make(map[Key]Value)
	dirty := make(map[Key]bool)
	return &Tx{id, db, writable, entries, dirty}
}

func (tx *Tx) ID() string {
	return tx.id
}

func (tx *Tx) Read(key Key) (Value, error) {
	if value, ok := tx.entries[key]; ok {
		return value, nil
	}
	value, err := tx.db.read(key)
	if err != nil {
		return "", err
	}
	tx.entries[key] = value
	return value, nil
}

func (tx *Tx) Write(key Key, value Value) error {
	if !tx.writable {
		return errors.New("This tx isn't writable")
	}
	tx.entries[key] = value
	tx.dirty[key] = true
	return nil
}

func (tx *Tx) Commit() error {
	for k, v := range tx.entries {
		if tx.dirty[k] {
			tx.db.write(k, v)
		}
	}
	return tx.closeTx()
}

func (tx *Tx) Rollback() error {
	return tx.closeTx()
}

func (tx *Tx) GC() (bool, error) {
	return tx.db.gc()
}

func (tx *Tx) closeTx() error {
	delete(tx.db.txs, tx.id)
	if tx.writable {
		tx.db.rwlock.Unlock()
	} else {
		tx.db.rwlock.RUnlock()
	}
	return nil
}
