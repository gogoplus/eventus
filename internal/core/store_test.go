package core

import (
    "context"
    "os"
    "testing"
)

func TestAppendAndMarkCommit(t *testing.T) {
    dir := "./testdata_store"
    _ = os.RemoveAll(dir)
    s, err := NewStore(dir)
    if err != nil {
        t.Fatalf("open store: %v", err)
    }
    defer func() { s.Close(); os.RemoveAll(dir) }()

    seq, err := s.AppendPending(context.Background(), Event{Type: "t", Data: []byte("x")})
    if err != nil {
        t.Fatalf("append pending: %v", err)
    }
    if seq == 0 {
        t.Fatalf("invalid seq")
    }
    if err := s.MarkCommitted(context.Background(), seq); err != nil {
        t.Fatalf("mark committed: %v", err)
    }
}
