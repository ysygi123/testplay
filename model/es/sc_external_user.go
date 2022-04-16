package es

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/olivere/elastic/v7"
	"testplay/utils"
	"time"
)

// ExternalProfile 外部联系人的自定义展示信息
type ExternalProfile struct {
	ExternalCorpName string          `json:"external_corp_name"`
	ExternalAttr     []*ExternalAttr `json:"external_attr"`
}

// ExternalAttr 外部联系人属性
type ExternalAttr struct {
	Type        int              `json:"type"` // 0 Text, 1 Web 2 Miniprogram
	Name        string           `json:"name"`
	Text        *TextAttr        `json:"text,omiempty"`
	Web         *WebAttr         `json:"web,omiempty"`
	MiniProgram *MiniprogramAttr `json:"miniprogram,omiempty"`
}

// TextAttr ...
type TextAttr struct {
	Value string `json:"value"`
}

// WebAttr ...
type WebAttr struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

// MiniprogramAttr ...
type MiniprogramAttr struct {
	AppID    string `json:"appid"`
	PagePath string `json:"pagepath"`
	Title    string `json:"title"`
}

// FollowTag 跟进人的备注
type FollowTag struct {
	GroupName string `json:"group_name"`
	TagName   string `json:"tag_name"`
	TagID     string `json:"tag_id"`
	Type      int    `json:"type"`
}

// FollowUserInfo 外部联系人跟进用户信息
type FollowUserInfo struct {
	CorpUserID     int          `json:"corp_user_id"`
	UserID         string       `json:"userid"`
	Remark         string       `json:"remark"`
	Name           string       `json:"name"`
	Description    string       `json:"description"`
	CreateTime     int64        `json:"createtime"`
	LostTime       int64        `json:"lost_time"`
	Tags           []*FollowTag `json:"tags"`
	Avatar         string       `json:"avatar"`
	TagID          []string     `json:"tag_id"`
	State          string       `json:"state"`
	StateText      string       `json:"state_text"`
	Status         int          `json:"status"`
	RemarkCorpName string       `json:"remark_corp_name"`
	RemarkMobiles  []string     `json:"remark_mobiles"`
	OpenUserID     string       `json:"open_user_id"`
	OperUserID     string       `json:"oper_user_id"`
	AddWay         int          `json:"add_way"`
	From           int          `json:"from"`
	FromText       string       `json:"from_text"`
}

// FollowTag 跟进人的备注
//type FollowTag struct {
//	GroupName string `json:"group_name"`
//	TagName   string `json:"tag_name"`
//	TagID     string `json:"tag_id"`
//	Type      int    `json:"type"`
//}

// ScExternalUser 外部联系人表
type ScExternalUser struct {
	ID                int64             `json:"id" xorm:"pk autoincr BIGINT(20)"`
	CID               int               `json:"cid" xorm:"not null default 0 comment('系统内部公司ID') INT(11)"`
	Avatar            string            `json:"avatar"  xorm:"not null default '' comment('外部联系人ID') VARCHAR(256)"`
	ExternalUserid    string            `json:"external_userid" xorm:"not null default '' comment('外部联系人ID') VARCHAR(64)"`
	Name              string            `json:"name" xorm:"not null default '' comment('外部联系人名称') VARCHAR(32)"`
	Mobile            string            `json:"mobile" xorm:"not null default '' comment('手机号') VARCHAR(16)"`
	Position          string            `json:"position" xorm:"not null default '' comment('职位') VARCHAR(32)"`
	CorpName          string            `json:"corp_name" xorm:"not null default '' comment('公司名称') VARCHAR(32)"`
	CorpFullName      string            `json:"corp_full_name" xorm:"not null default '' comment('公司全称') VARCHAR(32)"`
	LTV               string            `json:"ltv" xorm:"not null default 0.00 comment('ltv冗余') DECIMAL(10,2)"`
	Status            int               `json:"status" xorm:"not null comment('状态') TINYINT(1)"`
	Type              int               `json:"type" xorm:"not null default 0 comment('类型 1个人微信 2企业微信') TINYINT(4)"`
	Gender            int               `json:"gender" xorm:"not null default 0 comment('性别 0-未知 1-男性 2-女性') TINYINT(4)"`
	Age               int               `json:"age" xorm:"not null default 0 comment('年龄') TINYINT(4)"`
	Birthday          string            `json:"birthday" xorm:"not null default '' comment('生日') VARCHAR(16)"`
	Province          string            `json:"province" xorm:"not null default '' comment('省份') VARCHAR(16)"`
	City              string            `json:"city" xorm:"not null default '' comment('城市') VARCHAR(16)"`
	DyAccount         string            `json:"dy_account" xorm:"not null default '' comment('抖音号') VARCHAR(32)"`
	Unionid           string            `json:"unionid" xorm:"not null default '' comment('微信unionid') VARCHAR(128)"`
	ExternalProfile   *ExternalProfile  `json:"-" xorm:"comment('自定义展示信息') JSON"`
	FollowUser        []*FollowUserInfo `json:"follow_user" xorm:"comment('跟进客服信息') JSON"`
	AddTime           int64             `json:"add_time" xorm:"not null default 0 comment('最早添加时间') INT(11)"`
	CreatedAt         time.Time         `json:"created_at" xorm:"DATETIME"`
	UpdatedAt         time.Time         `json:"updated_at" xorm:"DATETIME"`
	Tags              []*FollowTag      `xorm:"-" json:"tags"`   //虚拟字段
	Remark            string            `json:"remark" xorm:"-"` //虚拟字段
	StateText         string            `json:"state_text" xorm:"-"`
	IsJoinGroup       int               `json:"is_join_group" xrom:"not null comment('是否加入群聊 1：已加入 2：未加入') TINYINT(4)"`
	CurrentSaveIsLost bool              `json:"-" xorm:"-"` //虚拟字段，是否流失，调用save方法时会为其赋值，可在外部进行判断使用
}

type ExternalOption struct {
	QueryBase
	Name     string
	CorpName string
}

func GetExternalWithLimit(op *ExternalOption, ctx context.Context) (list []*ScExternalUser, err error) {
	list = make([]*ScExternalUser, 0)
	es := utils.NewClient().SetIndex("sc-external-user")
	newS := es.EsClient.Search(es.Index)

	if op.Page != 0 && op.Size != 0 {
		newS.From((op.Page - 1) * op.Size).Size(op.Size)
	}

	if len(op.Sort) != 0 {
		for _, v := range op.Sort {
			newS.Sort(v.SortField, v.Ascending)
		}
	}
	q := elastic.NewBoolQuery()
	if len(op.Name) != 0 {
		q.Must(elastic.NewMatchQuery("name", op.Name))
	}
	if len(op.CorpName) != 0 {
		q.Must(elastic.NewMatchQuery("corp_name", op.CorpName))
	}
	newS.Query(q)
	searchES, err := newS.Do(ctx)
	if err != nil {
		return
	}

	err = jsoniter.Unmarshal(GetDefaultEsJsonByte(searchES.Hits.Hits), &list)
	if err != nil {
		return
	}

	fmt.Println(list)
	return
}
