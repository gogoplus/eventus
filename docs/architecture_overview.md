# Architecture Overview

Eventus implements the durable outbox pattern:

- Write event Envelope to local store (Badger) with Committed=false.
- Publisher attempts to send to Kafka.
- On success, mark Envelope Committed=true.
- RetryWorker periodically scans for uncommitted envelopes and republishes.
