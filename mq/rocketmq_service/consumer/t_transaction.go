package consumer

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"testplay/mq/rocketmq_i"
	"testplay/utils"
)

func init() {
	rocketmq_i.GrocketmqManager.Register(new(T2Transaction))
}

type T2Transaction struct {
	T2TransactionhI *T2Transactionh
}

func (t *T2Transaction) GetConsumerNum() int {
	return 6
}

func (t *T2Transaction) GetTopic() string {
	return rocketmq_i.TOPICTRANSACTION
}

func (t *T2Transaction) GetGroupName() string {
	return "T1_TransactionTest_c"
}

func (t *T2Transaction) GetRocketMqHandleC() rocketmq_i.RocketMqHandleC {
	t.T2TransactionhI = new(T2Transactionh)
	return t.T2TransactionhI
}

func (t *T2Transaction) GetIsOrder() bool {
	return true
}

func (t *T2Transaction) GetSelector() consumer.MessageSelector {
	return consumer.MessageSelector{
		//Type:       consumer.TAG,
		//Expression: "TAG1",
	}
}

type T2Transactionh struct {
	instanceName string
}

func (t *T2Transactionh) SetInstanceName(name string) {
	t.instanceName = name
}

func (t *T2Transactionh) Consumer(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	fmt.Println("我的instance是", t.instanceName)
	for _, v := range msgs {
		fmt.Println(utils.B2S(v.Body), "\n 查看一下queue_id是多少:", v.Queue.QueueId)
		fmt.Println("==============================")
	}
	return consumer.ConsumeSuccess, nil
}
