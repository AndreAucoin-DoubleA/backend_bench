package redpanda

import (
	"backend_bench/internal/model"

	"github.com/twmb/franz-go/pkg/kgo"
)

func NewConsumerWithConfig(brokers []string, topic, group string) (*model.Consumer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup(group),
		kgo.ConsumeTopics(topic),
	)

	if err != nil {
		return nil, err
	}

	return &model.Consumer{
		Client: client,
		Topic:  topic,
		Group:  group,
	}, nil
}
