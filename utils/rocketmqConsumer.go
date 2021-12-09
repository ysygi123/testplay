package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/consumer"
)

func A()  {
	s := 256+255
	uint16Num := uint16(s)
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.LittleEndian, uint16Num)
	fmt.Println(buf.Bytes())
	consumer.NewPullConsumer()
}