package publisher

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    kafka "github.com/segmentio/kafka-go"
    "github.com/gogoplus/eventus/internal/core"
)

type KafkaPublisher struct{ w *kafka.Writer }

func NewPublisher(brokers string, topic string) (*KafkaPublisher, error) {
    w := &kafka.Writer{Addr: kafka.TCP(brokers), Topic: topic, RequiredAcks: kafka.RequireAll, Async: false}
    return &KafkaPublisher{w: w}, nil
}

func (k *KafkaPublisher) Publish(ctx context.Context, seq uint64, ev core.Event) error {
    b, err := json.Marshal(ev)
    if err != nil {
        return err
    }
    msg := kafka.Message{Key: []byte(fmt.Sprintf("%d", seq)), Value: b, Time: time.Now()}
    ctx2, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    return k.w.WriteMessages(ctx2, msg)
}

func (k *KafkaPublisher) Close() error { return k.w.Close() }
