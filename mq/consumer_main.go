package main

import (
	"fmt"
	"testplay/mq/rocketmq_i"
	_ "testplay/mq/rocketmq_service/consumer"
)

func main() {
	fmt.Println("############################ \n#                          # \n#       start              # \n#                          # \n############################ ")
	rocketmq_i.GrocketmqManager.Start()
}
