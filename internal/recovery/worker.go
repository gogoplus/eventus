package recovery

import (
    "context"
    "fmt"
    "sync"
    "time"

    "github.com/gogoplus/eventus/internal/core"
    "github.com/gogoplus/eventus/internal/tx"
)

type Options struct {
    Interval time.Duration
    Batch    int
    Workers  int
}

type RetryWorker struct {
    store *core.Store
    pub   tx.Publisher
    opts  Options
    wg    sync.WaitGroup
    stop  chan struct{}
}

func NewRetryWorker(store *core.Store, pub tx.Publisher, opts Options) *RetryWorker {
    if opts.Interval == 0 { opts.Interval = 1 * time.Second }
    if opts.Batch == 0 { opts.Batch = 100 }
    if opts.Workers == 0 { opts.Workers = 4 }
    return &RetryWorker{store: store, pub: pub, opts: opts, stop: make(chan struct{})}
}

func (r *RetryWorker) Start(ctx context.Context) {
    r.wg.Add(1)
    go func() {
        defer r.wg.Done()
        ticker := time.NewTicker(r.opts.Interval)
        defer ticker.Stop()
        sem := make(chan struct{}, r.opts.Workers)
        for {
            select {
            case <-ctx.Done():
                return
            case <-r.stop:
                return
            case <-ticker.C:
                pending, err := r.store.ReadPending(ctx, r.opts.Batch)
                if err != nil {
                    fmt.Printf("retry read pending: %v\n", err)
                    continue
                }
                for _, env := range pending {
                    env := env
                    sem <- struct{}{}
                    r.wg.Add(1)
                    go func() {
                        defer r.wg.Done()
                        defer func() { <-sem }()
                        r.handle(ctx, env)
                    }()
                }
            }
        }
    }()
}

func (r *RetryWorker) Stop() {
    close(r.stop)
    r.wg.Wait()
}

func (r *RetryWorker) handle(ctx context.Context, env core.Envelope) {
    seq := env.Seq
    if err := r.pub.Publish(ctx, seq, env.Ev); err != nil {
        fmt.Printf("retry publish failed seq=%d err=%v\n", seq, err)
        return
    }
    if err := r.store.MarkCommitted(ctx, seq); err != nil {
        fmt.Printf("retry mark committed failed seq=%d err=%v\n", seq, err)
        return
    }
    fmt.Printf("retry published and committed seq=%d\n", seq)
}
