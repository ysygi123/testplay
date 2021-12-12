package utils

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

var NameServer = []string{"192.168.3.185:9876"}

type ProducerDBManager struct {
	producer   rocketmq.Producer
	NameServer []string
}

var GloableProducerDBManager *ProducerDBManager

func init() {
	GloableProducerDBManager = &ProducerDBManager{
		NameServer: NameServer,
	}
	err := GloableProducerDBManager.newProducer()
	if err != nil {
		panic("我草 去 -----" + err.Error())
	}
}

func (p *ProducerDBManager) newProducer() (err error) {
	mqProducer, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver(p.NameServer)),
		producer.WithRetry(2),
	)
	if err != nil {
		return
	}
	p.producer = mqProducer
	err = p.producer.Start()
	return
}

func (p *ProducerDBManager) GetProducer() rocketmq.Producer {
	return p.producer
}

func (p *ProducerDBManager) makeMessage(topic string, body []byte) (msg *primitive.Message) {
	msg = primitive.NewMessage(topic, body)
	return
}

func (p *ProducerDBManager) SendMessageSync(topic string, body []byte) {
	res, err := p.GetProducer().SendSync(context.Background(), p.makeMessage(topic, body))
	fmt.Println("[][][][][][][][][][][][][][][][][][][][][][][][]")
	fmt.Println(res.String(), err)
	fmt.Println("[][][][][][][][][][][][][][][][][][][][][][][][]")
}

func (p *ProducerDBManager) SendMessageSyncWithTag(topic, tag string, body []byte) {
	res, err := p.GetProducer().SendSync(context.Background(), p.makeMessage(topic, body).WithTag(tag))
	fmt.Println("[][][][][][][][][][][][][][][][][][][][][][][][]")
	fmt.Println(res.String(), err)
	fmt.Println("[][][][][][][][][][][][][][][][][][][][][][][][]")
}
