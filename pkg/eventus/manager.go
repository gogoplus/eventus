package eventus

import (
    "context"

    "github.com/gogoplus/eventus/internal/core"
    "github.com/gogoplus/eventus/internal/publisher"
    "github.com/gogoplus/eventus/internal/recovery"
    "github.com/gogoplus/eventus/internal/tx"
)

type Manager struct {
    Store *core.Store
    Pub   *publisher.KafkaPublisher
    Tx    *tx.TxManager
    Rwk   *recovery.RetryWorker
}

func NewManager(store *core.Store, brokers, topic string) (*Manager, error) {
    pub, err := publisher.NewPublisher(brokers, topic)
    if err != nil {
        return nil, err
    }
    txMgr := tx.New(store, pub)
    rwk := recovery.NewRetryWorker(store, pub, recovery.Options{})
    return &Manager{Store: store, Pub: pub, Tx: txMgr, Rwk: rwk}, nil
}

func (m *Manager) PublishAndCommit(ctx context.Context, ev core.Event) (uint64, error) {
    return m.Tx.AppendAndPublish(ctx, ev)
}

func (m *Manager) StartRetry(ctx context.Context) { m.Rwk.Start(ctx) }
func (m *Manager) StopRetry()                    { m.Rwk.Stop() }
