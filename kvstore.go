// Package kvstore provides support for local persistent key-value store.
//
// This implementation is a simple wrapper around boltdb.
//
// See https://github.com/boltdb/bolt
package kvstore

import (
	"encoding/binary"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

// ID is unique identifier of any object in a bucket.
type ID uint64

// Config is store manager configuration.
type Config struct {
	Filename string
}

// Manager is store manager.
type Manager struct {
	cfg Config

	db *bolt.DB
}

// NewManager creates new Manager instance.
func NewManager(cfg *Config) *Manager {
	return &Manager{cfg: *cfg}
}

// Open opens connection to store.
func (m *Manager) Open() error {
	if err := os.MkdirAll(filepath.Dir(m.cfg.Filename), 0700); err != nil {
		return errors.Wrap(err, "failed to create store directory")
	}

	db, err := bolt.Open(m.cfg.Filename, 0600, &bolt.Options{
		Timeout: 10 * time.Second,
	})

	if err != nil {
		return errors.Wrap(err, "failed to open store")
	}

	m.db = db

	return nil
}

// Close closes connection to store if one is open.
func (m *Manager) Close() error {
	if m.db != nil {
		return m.db.Close()
	}

	return nil
}

// Update invokes update transaction against store.
func (m *Manager) Update(updater func(trx *Trx) error) error {
	return m.db.Update(func(tx *bolt.Tx) error {
		return updater(&Trx{tx})
	})
}

// View invokes view transaction against store.
func (m *Manager) View(viewer func(trx *Trx) error) error {
	return m.db.View(func(tx *bolt.Tx) error {
		return viewer(&Trx{tx})
	})
}

// Trx is a transaction against store.
type Trx struct {
	tx *bolt.Tx
}

// Create generates ID for specified bucket, invokes idReceiver to setup an object and stores this object to store.
func (t *Trx) Create(bucket string, idReceiver func(id ID) interface{}) (interface{}, error) {
	b := t.tx.Bucket([]byte(bucket))
	id, err := b.NextSequence()

	if err != nil {
		return nil, errors.Wrap(err, "failed to generate object ID")
	}

	if err != nil {
		return nil, err
	}

	object := idReceiver(ID(id))
	bytes, err := json.Marshal(object)

	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal object to JSON")
	}

	if err := b.Put(itob(id), bytes); err != nil {
		return nil, errors.Wrap(err, "failed to put object to bucket")
	}

	return object, nil
}

// Put puts object specified by ID and stores it into receiver.
func (t *Trx) Put(bucket string, id ID, object interface{}) error {
	b := t.tx.Bucket([]byte(bucket))
	bytes, err := json.Marshal(object)

	if err != nil {
		return errors.Wrap(err, "failed to marshal object to JSON")
	}

	if err := b.Put(idtob(id), bytes); err != nil {
		return errors.Wrap(err, "failed to put object to bucket")
	}

	return nil
}

// Fetch fetches object specified by ID.
func (t *Trx) Fetch(bucket string, id ID, receiver interface{}) error {
	b := t.tx.Bucket([]byte(bucket))
	bytes := b.Get(idtob(id))

	if bytes == nil {
		return nil
	}

	return json.Unmarshal(bytes, &receiver)
}

// Delete deletes object specified by ID.
func (t *Trx) Delete(bucket string, id ID) error {
	b := t.tx.Bucket([]byte(bucket))

	return b.Delete(idtob(id))
}

// InitializeBucket creates bucket if one does not exist.
func (t *Trx) InitializeBucket(bucket string) error {
	_, err := t.tx.CreateBucketIfNotExists([]byte(bucket))

	return err
}

func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)

	return b
}

func idtob(v ID) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))

	return b
}
