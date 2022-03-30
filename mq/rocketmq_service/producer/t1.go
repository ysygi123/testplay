package producer

import (
	"context"
	"errors"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	jsoniter "github.com/json-iterator/go"
	"strconv"
	rocketmq2 "testplay/model/rocketmq"
	"testplay/mq/rocketmq_i"
)

func init() {
	rocketmq_i.GRocketProducerManager.Register(new(T1Producer))
}

type T1Producer struct {
	rocketmq_i.FatherMqHandler
}

func (t *T1Producer) GetTopics() string {
	return "T1"
}

// GetMessageSelector ...
func (t *T1Producer) GetMessageSelector() producer.QueueSelector {
	return producer.NewHashQueueSelector()
}

func (t *T1Producer) HandleMessage(message *primitive.Message, body []byte) (err error) {
	RmsgT1 := new(rocketmq2.RMyMessage)
	err = jsoniter.Unmarshal(body, RmsgT1)
	if err != nil {
		return
	}
	message.WithShardingKey(strconv.FormatInt(RmsgT1.OrderID, 10))
	return
}

func (t *T1Producer) SendMessage(ctx context.Context, producer rocketmq.Producer, messages ...*primitive.Message) (res *primitive.SendResult, err error) {
	np, ok := producer.(rocketmq.Producer)
	if !ok {
		err = errors.New("不是 np类型啊")
	}
	res, err = np.SendSync(ctx, messages...)
	return
}
