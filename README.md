# Eventus ⚡ v1.1.0

[![Go Reference](https://pkg.go.dev/badge/github.com/gogoplus/eventus.svg)](https://pkg.go.dev/github.com/gogoplus/eventus)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)]

> Eventus — durable outbox, transactional publish, and reliable delivery for Go services.  
> Designed for production use with BadgerDB (local durable store) and Kafka (delivery).

---

## Table of contents

- [Why Eventus?](#why-eventus)
- [Key Features](#key-features)
- [Quick Installation](#quick-installation)
  - [Prerequisites](#prerequisites)
  - [Clone and build (fast)](#clone-and-build-fast)


---

## Why Eventus?

Microservices often need reliable, crash-safe message delivery from the service that produces an event to the eventual consumers. Eventus implements a **durable outbox** plus a **transactional publish** pattern:

1. Append event to a local durable store (BadgerDB) as an *Envelope* (uncommitted).
2. Publish the event to Kafka.
3. On successful publish, mark the envelope as *Committed*.
4. If publish fails, the envelope remains uncommitted and is retried by a `RetryWorker`.

This guarantees that no produced event is lost even if the producer crashes before publishing — while keeping publish/commit semantics simple and auditable.

---

## Key Features

- Durable append-only local storage (BadgerDB) for outbound events.
- Transactional flow: `AppendPending -> Publish -> MarkCommitted`.
- Kafka adapter (using `github.com/segmentio/kafka-go`) for low-latency delivery.
- Background retry worker for eventual delivery of failed publishes.
- Example CLI apps (producer & consumer) and Docker-compose example for Kafka.
- CI-ready: tests, recommended GitHub Actions workflow included.
- Community-friendly docs, contribution guidelines, and changelog.

---

## Quick Installation

### Prerequisites

- Go `>= 1.23`
- (Optional) Docker & Docker Compose to run Kafka locally
- Git (to clone and push to GitHub)
- (Optional) `gh` CLI if you will create releases automatically

### Clone and build (fast)

```bash
git clone https://github.com/gogoplus/eventus.git
cd eventus
go mod tidy
go build ./...


Quickstart — Mock mode (no Kafka)

Mock mode is useful for local demos and CI without a running Kafka instance.

Linux / macOS (bash)
cd eventus
go test ./internal/core -v           # unit tests for the core store
# run the producer in mock mode (writer does not connect to Kafka)
go run ./cmd/eventus-producer --mock
# run consumer (reads committed events from local store)
go run ./cmd/eventus-consumer --mock

Windows (PowerShell)
cd eventus
go test ./internal/core -v
go run ./cmd/eventus-producer --mock
go run ./cmd/eventus-consumer --mock


In --mock mode the producer writes and marks committed locally; nothing is published to Kafka.

Quickstart — Real Kafka (Docker Compose)

Start Kafka and Zookeeper using the supplied Docker Compose:

cd examples/kafka/real
docker-compose up -d
sleep 8    # give Kafka time to start
cd ../../..
go run ./cmd/eventus-producer --brokers=localhost:9092 --topic=eventus.example
go run ./cmd/eventus-consumer --brokers=localhost:9092 --topic=eventus.example


Stop services:

cd examples/kafka/real
docker-compose down

Example usage (programmatic API)

Use pkg/eventus to integrate Eventus into applications.

Minimal producer example
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/gogoplus/eventus/internal/core"
    "github.com/gogoplus/eventus/pkg/eventus"
)

func main() {
    // Open local badger-backed store (path: ./data)
    store, err := core.NewStore("./data")
    if err != nil {
        log.Fatalf("open store: %v", err)
    }
    defer store.Close()

    // Create manager with Kafka brokers and topic
    mgr, err := eventus.NewManager(store, "localhost:9092", "eventus.example")
    if err != nil {
        log.Fatalf("new manager: %v", err)
    }

    // Publish an event and commit
    ev := core.Event{Type: "UserCreated", Data: []byte(`{"id": 123, "name":"alice"}`)}
    seq, err := mgr.PublishAndCommit(context.Background(), ev)
    if err != nil {
        log.Fatalf("publish failed: %v", err)
    }
    fmt.Printf("Published seq=%d\n", seq)
}

Starting the retry worker
ctx := context.Background()
mgr.StartRetry(ctx)    // background worker starts scanning pending envelopes
// later...
mgr.StopRetry()

CLI examples (producer & consumer)
Producer (example options)
# mock mode (no Kafka)
go run ./cmd/eventus-producer --mock

# real Kafka brokers and topic
go run ./cmd/eventus-producer --brokers=localhost:9092 --topic=eventus.example

Consumer (reads committed events from store)
# consumer reading local committed events
go run ./cmd/eventus-consumer --topic=eventus.example --brokers=localhost:9092

API Reference (summary)

See pkg/eventus/manager.go for public API details.

Types

core.Event — the event payload

type Event struct {
    Type string          `json:"type"`
    Data json.RawMessage `json:"data"`
}


core.Envelope — persisted metadata wrapper

type Envelope struct {
    Seq       uint64
    When      int64
    Committed bool
    Ev        Event
}


Key functions/methods

core.NewStore(dir string) (*Store, error) — open/create DB

(*Store) AppendPending(ctx context.Context, ev Event) (uint64, error) — append uncommitted

(*Store) MarkCommitted(ctx context.Context, seq uint64) error — mark committed

tx.New(store, publisher) — create transaction manager

pub.NewPublisher(brokers, topic) — Kafka publisher

pkg/eventus.NewManager(store, brokers, topic) — convenience wrapper

(*Manager) PublishAndCommit(ctx, ev) — append/publish/commit (atomic in app-level ordering)

Testing & CI
Run locally
go mod tidy
go test ./... -v



