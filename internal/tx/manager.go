package tx

import (
    "context"
    "fmt"

    "github.com/gogoplus/eventus/internal/core"
)

type Publisher interface {
    Publish(ctx context.Context, seq uint64, ev core.Event) error
}

type TxManager struct {
    store *core.Store
    pub   Publisher
}

func New(store *core.Store, pub Publisher) *TxManager { return &TxManager{store: store, pub: pub} }

func (m *TxManager) AppendAndPublish(ctx context.Context, ev core.Event) (uint64, error) {
    seq, err := m.store.AppendPending(ctx, ev)
    if err != nil {
        return 0, fmt.Errorf("append pending: %w", err)
    }
    if err := m.pub.Publish(ctx, seq, ev); err != nil {
        return seq, fmt.Errorf("publish failed: %w", err)
    }
    if err := m.store.MarkCommitted(ctx, seq); err != nil {
        return seq, fmt.Errorf("mark committed: %w", err)
    }
    return seq, nil
}
