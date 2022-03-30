package producer

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"sync"
	"sync/atomic"
	"testplay/mq/rocketmq_i"
	"time"
)

func init() {
	rocketmq_i.GRocketProducerManager.Register(new(TransactionTestProducer))
}

type DemoListener struct {
	localTrans       *sync.Map
	transactionIndex int32
}

func NewDemoListener() *DemoListener {
	return &DemoListener{
		localTrans: new(sync.Map),
	}
}

func (dl *DemoListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	nextIndex := atomic.AddInt32(&dl.transactionIndex, 1)
	fmt.Printf("nextIndex: %v for transactionID: %v\n", nextIndex, msg.TransactionId)
	status := nextIndex % 3
	dl.localTrans.Store(msg.TransactionId, primitive.LocalTransactionState(status+1))

	fmt.Printf("dl")
	return primitive.CommitMessageState
}

func (dl *DemoListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	fmt.Printf("%v msg transactionID : %v\n", time.Now(), msg.TransactionId)
	v, existed := dl.localTrans.Load(msg.TransactionId)
	if !existed {
		fmt.Printf("unknow msg: %v, return Commit", msg)
		return primitive.CommitMessageState
	}
	state := v.(primitive.LocalTransactionState)
	switch state {
	case 1:
		fmt.Printf("checkLocalTransaction COMMIT_MESSAGE: %v\n", msg)
		return primitive.CommitMessageState
	case 2:
		fmt.Printf("checkLocalTransaction ROLLBACK_MESSAGE: %v\n", msg)
		return primitive.RollbackMessageState
	case 3:
		fmt.Printf("checkLocalTransaction unknow: %v\n", msg)
		return primitive.UnknowState
	default:
		fmt.Printf("checkLocalTransaction default COMMIT_MESSAGE: %v\n", msg)
		return primitive.CommitMessageState
	}
}

type TransactionTestProducer struct {
	rocketmq_i.FatherMqHandler
}

func (t *TransactionTestProducer) IsTransaction() bool {
	return true
}

func (t *TransactionTestProducer) GetTopics() string {
	return rocketmq_i.TOPICTRANSACTION
}

func (t *TransactionTestProducer) GetTransactionListener() primitive.TransactionListener {
	return NewDemoListener()
}

func (t *TransactionTestProducer) SendTransactionMessage(ctx context.Context, transactionProducer rocketmq.TransactionProducer, msg *primitive.Message) (sendResult *primitive.TransactionSendResult, err error) {
	sendResult, err = transactionProducer.SendMessageInTransaction(ctx, msg)

	return
}
