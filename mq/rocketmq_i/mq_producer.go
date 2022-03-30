package rocketmq_i

import (
	"context"
	"errors"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

var GRocketProducerManager *RocketProducerManger

func init() {
	GRocketProducerManager = new(RocketProducerManger)
	GRocketProducerManager.DB = make(map[string]*unionRocketMqHandleI)
}

type RocketProducerManger struct {
	DB map[string]*unionRocketMqHandleI
}

type unionRocketMqHandleI struct {
	Rmhi                RocketMqHandleI
	Producer            rocketmq.Producer
	TransactionProducer rocketmq.TransactionProducer
}

func (r *RocketProducerManger) Register(ri RocketMqHandleI) {
	tmpU := new(unionRocketMqHandleI)
	tmpU.Rmhi = ri
	r.DB[ri.GetTopics()] = tmpU
}

// GetMyProducer 获取生产者
func (r *RocketProducerManger) GetMyProducer(topic string) (x *unionRocketMqHandleI, err error) {

	x, ok := r.DB[topic]
	if !ok {
		err = errors.New("not found")
		return
	}
	if !x.Rmhi.IsTransaction() {
		err = r.getDefaultProducer(x, topic)
	} else {
		err = r.getMyTransactionProducer(x, topic)
	}

	return
}

// getCommonOption 公共方法
func (r *RocketProducerManger) getCommonOption(x *unionRocketMqHandleI) []producer.Option {
	option := []producer.Option{
		producer.WithNsResolver(primitive.NewPassthroughResolver(ServerName)),
		producer.WithRetry(x.Rmhi.GetRetryTimes()),
		producer.WithDefaultTopicQueueNums(x.Rmhi.GetQueueNum()),
	}

	if selector := x.Rmhi.GetMessageSelector(); selector != nil {
		option = append(option, producer.WithQueueSelector(selector))
	}
	return option
}

func (r *RocketProducerManger) getDefaultProducer(x *unionRocketMqHandleI, topic string) (err error) {

	if x.Producer != nil {
		return
	}

	option := r.getCommonOption(x)
	var rproducer rocketmq.Producer
	rproducer, err = rocketmq.NewProducer(
		option...,
	)
	if err != nil {
		return
	}
	err = rproducer.Start()
	x.Producer = rproducer
	return
}

func (r *RocketProducerManger) getMyTransactionProducer(x *unionRocketMqHandleI, topic string) (err error) {
	if x.TransactionProducer != nil {
		return
	}
	f := x.Rmhi
	option := r.getCommonOption(x)
	tproducer, err := rocketmq.NewTransactionProducer(
		f.GetTransactionListener(),
		option...,
	)
	if err != nil {
		return
	}
	x.TransactionProducer = tproducer
	err = x.TransactionProducer.Start()
	return
}

// SendMessage ...发送消息
func (r *RocketProducerManger) SendMessage(ctx context.Context, topic string, messages [][]byte) (res *primitive.SendResult, err error) {
	unI, err := r.GetMyProducer(topic)
	if err != nil {
		return
	}
	msgs := make([]*primitive.Message, len(messages))
	for i, body := range messages {
		pMsg := primitive.NewMessage(unI.Rmhi.GetTopics(), body)
		err = unI.Rmhi.HandleMessage(pMsg, body)
		if err != nil {
			return
		}
		msgs[i] = pMsg
	}
	res, err = unI.Rmhi.SendMessage(ctx, unI.Producer, msgs...)
	return
}

// SendTransactionMessage 发送事务消息
func (r *RocketProducerManger) SendTransactionMessage(ctx context.Context, topic string, messages []byte) (res *primitive.TransactionSendResult, err error) {
	unI, err := r.GetMyProducer(topic)
	if err != nil {
		return
	}
	pMsg := primitive.NewMessage(unI.Rmhi.GetTopics(), messages)
	err = unI.Rmhi.HandleMessage(pMsg, messages)
	if err != nil {
		return
	}
	res, err = unI.Rmhi.SendTransactionMessage(ctx, unI.TransactionProducer, pMsg)
	return
}

const (
	DefaultProducer     = 0
	TransactionProducer = 1
)

type RocketMqHandleI interface {
	GetTopics() string
	GetMessageSelector() producer.QueueSelector
	HandleMessage(message *primitive.Message, body []byte) (err error)
	SendMessage(ctx context.Context, producer rocketmq.Producer, msgs ...*primitive.Message) (sendResult *primitive.SendResult, err error)
	SendTransactionMessage(ctx context.Context, transactionProducer rocketmq.TransactionProducer, msg *primitive.Message) (sendResult *primitive.TransactionSendResult, err error)
	GetQueueNum() int
	GetRetryTimes() int
	IsTransaction() bool
	GetTransactionListener() primitive.TransactionListener
}

type FatherMqHandler struct {
}

func (f *FatherMqHandler) GetTopics() string {
	return ""
}

// GetMessageSelector ...
func (f *FatherMqHandler) GetMessageSelector() producer.QueueSelector {
	return nil
}

func (f *FatherMqHandler) GetQueueNum() int {
	return 6
}

func (f *FatherMqHandler) GetRetryTimes() int {
	return 6
}

func (f *FatherMqHandler) HandleMessage(message *primitive.Message, body []byte) (err error) {
	return
}

func (f *FatherMqHandler) SendMessage(ctx context.Context, producer rocketmq.Producer, messages ...*primitive.Message) (sendResult *primitive.SendResult, err error) {
	return
}
func (f *FatherMqHandler) SendTransactionMessage(ctx context.Context, transactionProducer rocketmq.TransactionProducer, msg *primitive.Message) (sendResult *primitive.TransactionSendResult, err error) {
	return
}
func (f *FatherMqHandler) IsTransaction() bool {
	return false
}

func (f *FatherMqHandler) GetTransactionListener() primitive.TransactionListener {
	return nil
}
