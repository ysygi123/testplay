package rocketmq_i

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"sync"
	"testplay/utils"
)

type RocketMqHandleCManger interface {
	GetRocketMqHandleC() RocketMqHandleC
	GetConsumerNum() int
	GetTopic() string
	GetGroupName() string
	GetSelector() consumer.MessageSelector
}

type RocketMqHandleC interface {
	SetInstanceName(string)
	Consumer(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error)
}

var GrocketmqManager *RocketMqManger

func init() {
	GrocketmqManager = NewRocketMqManger()
}

type RocketMqManger struct {
	DB map[string]RocketMqHandleCManger
}

func NewRocketMqManger() *RocketMqManger {
	return &RocketMqManger{
		DB: make(map[string]RocketMqHandleCManger),
	}
}

func (r *RocketMqManger) Register(c RocketMqHandleCManger) {
	r.DB[c.GetTopic()] = c
}

func (r *RocketMqManger) Start() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	for _, c := range r.DB {
		for i := 0; i < c.GetConsumerNum(); i++ {
			instanceName := c.GetTopic() + c.GetGroupName() + utils.GetRandomString(6)
			fmt.Println("【循环DB】 DB为" + c.GetTopic() + " groupName" + c.GetGroupName() + "instanceName=" + instanceName)
			_ = r.GetStartConsumer(instanceName, c)
		}
	}
	wg.Wait()
}

func (r *RocketMqManger) GetStartConsumer(instanceName string, c RocketMqHandleCManger) (err error) {
	if len(c.GetTopic()) == 0 {
		err = fmt.Errorf("[---------] 没有")
	}
	//ops := make([]consumer.Option, 0)
	ops := []consumer.Option{
		consumer.WithGroupName(c.GetGroupName()),
		consumer.WithNameServer(ServerName),
		consumer.WithInstance(instanceName),
	}

	MqPushConsumerSuccess, err := rocketmq.NewPushConsumer(
		ops...,
	//consumer.WithConsumerModel(consumer.BroadCasting),
	//consumer.OffsetStore()
	//consumer.WithConsumeFromWhere(consumer.ConsumeFromLastOffset),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	selector := c.GetSelector()
	h := c.GetRocketMqHandleC()
	h.SetInstanceName(instanceName)
	fmt.Println("[---------]1")
	err = MqPushConsumerSuccess.Subscribe(c.GetTopic(), selector, h.Consumer)
	if err != nil {
		return
	}
	err = MqPushConsumerSuccess.Start()
	return
}
