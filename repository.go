package main

import (
	"encoding/binary"
	"fmt"

	"github.com/ReSc/septum"
	"github.com/boltdb/bolt"
)

type repository struct {
	db      *bolt.DB
	buckets [][]byte
	events  []byte
	timers  []byte
	dbPath  string
}

func newRepository() *repository {
	r := &repository{
		events: []byte("Events"),
		timers: []byte("Timers"),
	}
	r.buckets = [][]byte{
		r.events,
		r.timers,
	}
	return r
}

func (r *repository) Open(dbPath string) error {
	r.dbPath = dbPath
	db, err := bolt.Open(dbPath, 0600, nil)

	if err != nil {
		return fmt.Errorf("repository.Open: %s", err)
	}

	// setup buckets for events.
	err = db.Update(func(tx *bolt.Tx) error {

		for i := range r.buckets {
			_, err := tx.CreateBucketIfNotExists(r.buckets[i])
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("repository.Open: %s", err)
	}

	r.db = db
	return nil
}

func (r *repository) SaveEvent(e *septum.Event) error {
	return r.save(r.events, e)
}

func (r *repository) SaveTimer(e *septum.Event) error {
	return r.save(r.timers, e)
}

func (r *repository) encode(e *septum.Event) (key []byte, value []byte, err error) {
	key = uint64ToBytes(e.Id)
	value, err = e.MarshalBinary()
	return key, value, err
}

func (r *repository) Close() error {
	return r.db.Close()
}

func uint64ToBytes(n uint64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, n)
	return bytes
}

func (r *repository) save(bucket []byte, e *septum.Event) error {
	key, value, err := r.encode(e)
	if err != nil {
		return err
	}

	tx, err := r.db.Begin(true)
	if err != nil {
		return err
	}
	bb := tx.Bucket(bucket)
	if err := bb.Put(key, value); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	} else {
		return tx.Commit()
	}

}
