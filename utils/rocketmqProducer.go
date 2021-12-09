package utils

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

type ProducerDBManager struct {
	producer rocketmq.Producer
	NameServer []string
}

var GloableProducerDBManager *ProducerDBManager

func init()  {
	GloableProducerDBManager = &ProducerDBManager{
		NameServer: []string{"127.0.0.1:9876"},
	}
	err := GloableProducerDBManager.newProducer()
	if err != nil {
		panic("我草 去 -----" + err.Error())
	}
}

func (p *ProducerDBManager)newProducer() (err error) {
	mqProducer, err := rocketmq.NewProducer(
		producer.WithGroupName("test_group_1"),
		producer.WithNameServer(p.NameServer),
		producer.WithRetry(3),
		)
	if err != nil {
		return
	}
	p.producer = mqProducer
	err = p.producer.Start()
	return
}

func (p *ProducerDBManager)GetProducer() rocketmq.Producer {
	return p.producer
}

func (p *ProducerDBManager)makeMessage(topic string, body []byte) (msg *primitive.Message) {
	msg = primitive.NewMessage(topic, body)
	return
}

func (p *ProducerDBManager)SendMessageSync(topic string, body []byte)  {
	res, err := p.GetProducer().SendSync(context.Background(), p.makeMessage(topic, body))
	fmt.Println("[][][][][][][][][][][][][][][][][][][][][][][][]")
	fmt.Println(res.String(), err)
	fmt.Println("[][][][][][][][][][][][][][][][][][][][][][][][]")
}