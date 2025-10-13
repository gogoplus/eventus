# Usage Guide

Quick mock demo:
1. go mod tidy
2. go run ./cmd/eventus-producer --mock
3. go run ./cmd/eventus-consumer --mock

Real Kafka:
1. cd examples/kafka/real && docker-compose up -d
2. go run ./cmd/eventus-producer --brokers=localhost:9092 --topic=eventus.example
