package rocketmq_service

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"testplay/mq/rocketmq_i"
	"testplay/utils"
)

func init() {
	rocketmq_i.GrocketmqManager.Register(new(T1))
}

type T1 struct {
	T1hI *T1h
}

func (t *T1) GetConsumerNum() int {
	return 2
}

func (t *T1) GetTopic() string {
	return "T1"
}

func (t *T1) GetGroupName() string {
	return "T1_c"
}

func (t *T1) GetRocketMqHandleC() rocketmq_i.RocketMqHandleC {
	t.T1hI = new(T1h)
	return t.T1hI
}

func (t *T1) GetSelector() consumer.MessageSelector {
	return consumer.MessageSelector{
		Type:       consumer.TAG,
		Expression: "TAG1",
	}
}

type T1h struct {
	instanceName string
}

func (t *T1h) SetInstanceName(name string) {
	t.instanceName = name
}

func (t *T1h) Consumer(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	fmt.Println("我的instance是", t.instanceName)
	for _, v := range msgs {
		fmt.Println(utils.B2S(v.Body))
		fmt.Println("==============================")
	}
	return consumer.ConsumeSuccess, nil
}
