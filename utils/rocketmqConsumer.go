package utils

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"os"
	"time"
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
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})),
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
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName("testGroup"),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})),
	)
	err := c.Subscribe("testTopic", consumer.MessageSelector{}, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i := range msgs {
			fmt.Printf("subscribe callback: %v \n", msgs[i])
		}

		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	// Note: start after subscribe
	err = c.Start()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	time.Sleep(time.Hour)
	err = c.Shutdown()
	if err != nil {
		fmt.Printf("shutdown Consumer error: %s", err.Error())
	}
}
