package service

import (
	"encoding/json"
	"testplay/utils"
)

type Ex1Msg struct {
	Uid   int    `json:"uid"`
	ToUid int    `json:"to_uid"`
	Msg   string `json:"msg"`
}

func Ex1Send() {
	em := &Ex1Msg{
		Uid:   1,
		ToUid: 1,
		Msg:   "adsf",
	}
	b, err := json.Marshal(em)
	if err != nil {
		return
	}
	utils.GloableProducerDBManager.SendMessageSync("testTopic", b)
}
