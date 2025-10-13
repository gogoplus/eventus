package main

import (
    "context"
    "flag"
    "fmt"
    "log"
    "time"

    "github.com/gogoplus/eventus/internal/core"
)

func main() {
    brokers := flag.String("brokers", "localhost:9092", "kafka brokers")
    topic := flag.String("topic", "eventus.example", "topic")
    flag.Parse()

    _ = brokers
    _ = topic
    store, err := core.NewStore("./data")
    if err != nil {
        log.Fatalf("open store: %v", err)
    }
    defer store.Close()

    ctx := context.Background()
    for i := 0; i < 10; i++ {
        committed, err := store.ReadCommitted(ctx, 0, 100)
        if err != nil {
            log.Printf("read committed err: %v", err)
            time.Sleep(time.Second)
            continue
        }
        for _, e := range committed {
            fmt.Printf("consumed seq=%d type=%s data=%s\n", e.Seq, e.Ev.Type, string(e.Ev.Data))
        }
        time.Sleep(2 * time.Second)
    }
}
