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
- [Quickstart — Mock mode (no Kafka)](#quickstart---mock-mode-no-kafka)
  - [Linux / macOS (bash)](#linux--macos-bash)
  - [Windows (PowerShell)](#windows-powershell)
- [Quickstart — Real Kafka (Docker)](#quickstart---real-kafka-docker)
- [Example usage (programmatic API)](#example-usage-programmatic-api)
- [CLI examples (producer & consumer)](#cli-examples-producer--consumer)
- [API Reference (summary)](#api-reference-summary)
- [Testing & CI](#testing--ci)
- [Preparing the repo for GitHub (step-by-step)](#preparing-the-repo-for-github-step-by-step)
  - [Initialize, commit, push to empty repo](#initialize-commit-push-to-empty-repo)
  - [Create a GitHub Release (manual / `gh` / REST API)](#create-a-github-release-manual--gh--rest-api)
- [Contributing & CODE_OF_CONDUCT](#contributing--code_of_conduct)
- [Troubleshooting & FAQ](#troubleshooting--faq)
- [Roadmap & Governance](#roadmap--governance)
- [License](#license)

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


Eventus is a production-ready Go library implementing:
- Durable outbox pattern (BadgerDB)
- Transactional publish (Append -> Publish -> MarkCommitted)
- Kafka adapter (segmentio/kafka-go)
- Recovery worker that ensures eventual delivery

See /docs for full usage, architecture, and contribution guidelines.
