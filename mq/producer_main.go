package main

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"testplay/mq/rocketmq_i"
	"testplay/utils"
	"time"
)

func main() {
	x := new(rocketmq_i.FatherMqHandler)
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 10; i++ {
		tagID := rand.Intn(1) + 1
		reg := "这是一个奇怪的东西-" + strconv.Itoa(i) + " : tagID = " + strconv.Itoa(tagID)
		fmt.Println(reg)
		vadsf := [][]byte{
			utils.S2B(reg),
		}
		//err := x.SendMyMessagesSync(context.Background(), "T1", vadsf)
		err := x.SendMyManyMessagesWithOneTagSync(context.Background(), "T1", vadsf, "TAG"+strconv.Itoa(tagID))
		fmt.Println(err)
	}
}
