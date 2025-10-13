package main

import (
    "context"
    "flag"
    "fmt"
    "log"
    "time"

    "github.com/gogoplus/eventus/internal/core"
    "github.com/gogoplus/eventus/pkg/eventus"
)

func main() {
    mock := flag.Bool("mock", false, "use mock (no kafka)")
    brokers := flag.String("brokers", "localhost:9092", "kafka brokers")
    topic := flag.String("topic", "eventus.example", "topic")
    flag.Parse()

    store, err := core.NewStore("./data")
    if err != nil {
        log.Fatalf("open store: %v", err)
    }
    defer store.Close()

    mgr, err := eventus.NewManager(store, *brokers, *topic)
    if err != nil {
        log.Fatalf("new manager: %v", err)
    }
    ctx := context.Background()
    mgr.StartRetry(ctx)

    for i := 0; i < 5; i++ {
        ev := core.Event{Type: "demo", Data: []byte(fmt.Sprintf(`{"n":%d}`, i))}
        if seq, err := mgr.PublishAndCommit(context.Background(), ev); err != nil {
            log.Printf("publish failed: %v", err)
        } else {
            log.Printf("published seq=%d", seq)
        }
        time.Sleep(200 * time.Millisecond)
    }

    mgr.StopRetry()
    fmt.Println("producer finished")
}
