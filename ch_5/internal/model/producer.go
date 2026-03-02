package model

import (
	"context"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Producer struct {
	Client *kgo.Client
	Topic  string
}

func (producer *Producer) Close() {
	producer.Client.Close()
}

func (producer *Producer) Produce(ctx context.Context, value []byte) error {
	record := &kgo.Record{
		Topic: producer.Topic,
		Value: value,
	}
	results := producer.Client.ProduceSync(ctx, record)
	return results.FirstErr()
}

func NewProducer(brokers []string, topic string) *Producer {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.AllowAutoTopicCreation(),
		kgo.RequestTimeoutOverhead(10*time.Second),
		kgo.ProduceRequestTimeout(30*time.Second),
	)
	if err != nil {
		panic(err)
	}

	producer := &Producer{
		Client: client,
		Topic:  topic,
	}
	return producer
}
