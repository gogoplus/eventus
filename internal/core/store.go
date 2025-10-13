package core

import (
    "context"
    "encoding/json"
    "fmt"
    "path/filepath"
    "sync/atomic"
    "time"

    "github.com/dgraph-io/badger/v4"
)

const (
    metaSeqKey = "es:meta:seq"
    prefix     = "es:"
)

// Event is the domain payload
type Event struct {
    Type string          `json:"type"`
    Data json.RawMessage `json:"data"`
}

// Envelope wraps event with metadata
type Envelope struct {
    Seq       uint64 `json:"seq"`
    When      int64  `json:"when"`
    Committed bool   `json:"committed"`
    Ev        Event  `json:"event"`
}

// Store is a durable append-only store backed by BadgerDB
type Store struct {
    db  *badger.DB
    dir string
    seq uint64
}

// NewStore opens the DB under dir
func NewStore(dir string) (*Store, error) {
    opts := badger.DefaultOptions(filepath.Clean(dir)).WithLogger(nil)
    db, err := badger.Open(opts)
    if err != nil {
        return nil, err
    }
    s := &Store{db: db, dir: dir}
    if err := s.loadSeq(); err != nil {
        _ = db.Close()
        return nil, err
    }
    return s, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) loadSeq() error {
    return s.db.View(func(txn *badger.Txn) error {
        item, err := txn.Get([]byte(metaSeqKey))
        if err == badger.ErrKeyNotFound {
            atomic.StoreUint64(&s.seq, 0)
            return nil
        }
        if err != nil {
            return err
        }
        v, err := item.ValueCopy(nil)
        if err != nil {
            return err
        }
        var cur uint64
        if err := json.Unmarshal(v, &cur); err != nil {
            return err
        }
        atomic.StoreUint64(&s.seq, cur)
        return nil
    })
}

func (s *Store) AppendPending(ctx context.Context, ev Event) (uint64, error) {
    var seq uint64
    err := s.db.Update(func(txn *badger.Txn) error {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }
        cur := atomic.AddUint64(&s.seq, 1)
        seq = cur
        env := Envelope{Seq: seq, When: time.Now().UnixNano(), Committed: false, Ev: ev}
        raw, err := json.Marshal(env)
        if err != nil {
            return err
        }
        k := []byte(fmt.Sprintf("%s%020d", prefix, seq))
        if err := txn.Set(k, raw); err != nil {
            return err
        }
        metaRaw, _ := json.Marshal(seq)
        if err := txn.Set([]byte(metaSeqKey), metaRaw); err != nil {
            return err
        }
        return nil
    })
    if err != nil {
        return 0, err
    }
    return seq, nil
}

func (s *Store) MarkCommitted(ctx context.Context, seq uint64) error {
    key := []byte(fmt.Sprintf("%s%020d", prefix, seq))
    return s.db.Update(func(txn *badger.Txn) error {
        item, err := txn.Get(key)
        if err != nil {
            return err
        }
        v, err := item.ValueCopy(nil)
        if err != nil {
            return err
        }
        var e Envelope
        if err := json.Unmarshal(v, &e); err != nil {
            return err
        }
        if e.Committed {
            return nil
        }
        e.Committed = true
        raw, _ := json.Marshal(e)
        return txn.Set(key, raw)
    })
}

func (s *Store) ReadPending(ctx context.Context, limit int) ([]Envelope, error) {
    var out []Envelope
    err := s.db.View(func(txn *badger.Txn) error {
        it := txn.NewIterator(badger.DefaultIteratorOptions)
        defer it.Close()
        pfx := []byte(prefix)
        for it.Seek(pfx); it.ValidForPrefix(pfx) && (limit <= 0 || len(out) < limit); it.Next() {
            item := it.Item()
            v, err := item.ValueCopy(nil)
            if err != nil {
                return err
            }
            var e Envelope
            if err := json.Unmarshal(v, &e); err != nil {
                return err
            }
            if !e.Committed {
                out = append(out, e)
            }
        }
        return nil
    })
    return out, err
}

func (s *Store) ReadCommitted(ctx context.Context, since uint64, limit int) ([]Envelope, error) {
    var out []Envelope
    err := s.db.View(func(txn *badger.Txn) error {
        it := txn.NewIterator(badger.DefaultIteratorOptions)
        defer it.Close()
        start := []byte(fmt.Sprintf("%s%020d", prefix, since+1))
        pfx := []byte(prefix)
        for it.Seek(start); it.ValidForPrefix(pfx) && (limit <= 0 || len(out) < limit); it.Next() {
            item := it.Item()
            v, err := item.ValueCopy(nil)
            if err != nil {
                return err
            }
            var e Envelope
            if err := json.Unmarshal(v, &e); err != nil {
                return err
            }
            if e.Committed {
                out = append(out, e)
            }
        }
        return nil
    })
    return out, err
}

func (s *Store) LastSeq() uint64 { return atomic.LoadUint64(&s.seq) }
