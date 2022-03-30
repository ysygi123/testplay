package main

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"math/rand"
	"strconv"
	"testplay/model/rocketmq"
	"testplay/mq/rocketmq_i"
	_ "testplay/mq/rocketmq_service/producer"
	"time"
)

func main() {
	sendTransactionMessage()
	time.Sleep(1000000 * time.Second)
}

func sendDefaultMessage() {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 50; i++ {
		tagID := uint8(rand.Intn(1) + 1)
		reg := "这是一个奇怪的东西-" + strconv.Itoa(i) + " : tagID = " + strconv.Itoa(int(tagID))
		fmt.Println(reg)
		newRm := &rocketmq.RMyMessage{
			TagID:   tagID,
			Type:    1,
			UID:     12,
			ToUID:   13,
			OrderID: rand.Int63n(3),
			Message: reg,
		}
		b, _ := jsoniter.Marshal(newRm)
		res, err := rocketmq_i.GRocketProducerManager.SendMessage(context.Background(), "T1", [][]byte{b})
		//err := x.SendMyManyMessagesWithOneTagSync(context.Background(), "T1", vadsf, "TAG"+strconv.Itoa(tagID))
		fmt.Println(err)
		fmt.Println(res)
	}
}

func sendTransactionMessage() {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 50; i++ {
		tagID := uint8(rand.Intn(1) + 1)
		reg := "这是一个奇怪的东西-" + strconv.Itoa(i) + " : tagID = " + strconv.Itoa(int(tagID))
		fmt.Println(reg)
		newRm := &rocketmq.RMyMessage{
			TagID:   tagID,
			Type:    1,
			UID:     12,
			ToUID:   13,
			OrderID: rand.Int63n(3),
			Message: reg,
		}
		b, _ := jsoniter.Marshal(newRm)
		res, err := rocketmq_i.GRocketProducerManager.SendTransactionMessage(context.Background(), rocketmq_i.TOPICTRANSACTION, b)
		//err := x.SendMyManyMessagesWithOneTagSync(context.Background(), "T1", vadsf, "TAG"+strconv.Itoa(tagID))
		fmt.Println(err)
		fmt.Println(res)
	}
}
