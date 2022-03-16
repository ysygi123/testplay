package rocketmq_i

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

type RocketMqHandleI interface {
	GetTopics() []string
}

type FatherMqHandler struct {
	producer rocketmq.Producer
}

func (f *FatherMqHandler) GetTopics() []string {
	return []string{}
}

func (f *FatherMqHandler) GetMyProducer() (rocketmq.Producer, error) {
	if f.producer != nil {
		return f.producer, nil
	}
	var err error
	f.producer, err = rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver(ServerName)),
		producer.WithRetry(2),
		producer.WithDefaultTopicQueueNums(6),
	)
	if err != nil {
		return nil, err
	}
	err = f.producer.Start()
	return f.producer, err
}

func (f *FatherMqHandler) SendMyManyMessagesWithOneTagSync(ctx context.Context, topic string, messages [][]byte, tag string) (err error) {
	msgs := make([]*primitive.Message, len(messages))
	for i, body := range messages {
		msgs[i] = f.makeMessage(topic, body).WithTag(tag)
	}
	p, err := f.GetMyProducer()
	if err != nil {
		return
	}
	res, err := p.SendSync(ctx, msgs...)
	if err != nil {
		return
	}
	fmt.Println("======", res.String(), "======")
	return
}

func (f *FatherMqHandler) SendMyMessagesSync(ctx context.Context, topic string, messages [][]byte) (err error) {
	msgs := make([]*primitive.Message, len(messages))
	for i, body := range messages {
		msgs[i] = f.makeMessage(topic, body)
	}
	p, err := f.GetMyProducer()
	if err != nil {
		return
	}
	res, err := p.SendSync(ctx, msgs...)
	if err != nil {
		return
	}
	fmt.Println("======", res.String(), "======")
	return
}

func (f *FatherMqHandler) makeMessage(topic string, body []byte) (msg *primitive.Message) {
	msg = primitive.NewMessage(topic, body)
	return
}
