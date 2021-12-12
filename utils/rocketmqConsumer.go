package utils

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func A() {
	s := 256 + 255
	uint16Num := uint16(s)
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.LittleEndian, uint16Num)
	fmt.Println(buf.Bytes())
	consumer.NewPullConsumer()
}

func Pull() {
	c, err := consumer.NewPullConsumer(
		consumer.WithGroupName("test_group_1"),
		consumer.WithNsResolver(primitive.NewPassthroughResolver(NameServer)),
	)
	if err != nil {
		panic(err)
	}
	err = c.Start()
	if err != nil {
		panic(err)
	}

	selecter := consumer.MessageSelector{}
	ctx := context.Background()
	var resp *primitive.PullResult
	for {
		resp, err = c.Pull(ctx, "testTopic", selecter, 10)
		if err != nil {
			fmt.Println(err, "我草你妈呀")
			continue
		}
		fmt.Println("------------------------------------------")
		fmt.Println(resp)
		fmt.Println("------------------------------------------")

	}
}

func Push() {
	MqPushConsumerSuccess, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName("notify_consumer_success"),
		consumer.WithNameServer(NameServer),
		//consumer.WithConsumeFromWhere(consumer.ConsumeFromLastOffset),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	selector := consumer.MessageSelector{
		Type: consumer.TAG,
	}
	fmt.Println("养狗死全家1")
	err = MqPushConsumerSuccess.Subscribe("testTopic", selector, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i := range msgs {
			fmt.Println("看看这里是我收到的消息啊啊啊啊==========" + string(msgs[i].Body))
		}
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		panic(fmt.Sprintf("我草我看看是不是这个了consumer subscribe TopicNotifySucess err:%v", err))
		return
	}
	fmt.Println("养狗死全家2")
	err = MqPushConsumerSuccess.Start()
	if err != nil {
		panic(err)
	}
}

func GetStartConsumer(topic, groupName, instanceName string, f func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error)) (err error) {
	if len(topic) == 0 {
		err = fmt.Errorf("养狗死全家 没有")
	}
	MqPushConsumerSuccess, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName(groupName),
		consumer.WithNameServer(NameServer),
		consumer.WithInstance(instanceName),
		//consumer.WithConsumeFromWhere(consumer.ConsumeFromLastOffset),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	selector := consumer.MessageSelector{
		Type: consumer.TAG,
	}
	fmt.Println("养狗死全家1")
	err = MqPushConsumerSuccess.Subscribe(topic, selector, f)
	if err != nil {
		return
	}
	err = MqPushConsumerSuccess.Start()
	return
}

func F1(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for i := range msgs {
		fmt.Println("养狗死全家 消费者 【1】 ==========" + string(msgs[i].Body))
	}
	return consumer.ConsumeSuccess, nil
}

func F2(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for i := range msgs {
		fmt.Println("养狗死全家 消费者 【2】 ==========" + string(msgs[i].Body))
	}
	return consumer.ConsumeSuccess, nil
}
