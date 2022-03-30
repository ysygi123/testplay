package rocketmq

type RMyMessage struct {
	TagID   uint8  `json:"tag_id"`
	Type    uint8  `json:"type"`
	UID     int    `json:"uid"`
	ToUID   int    `json:"to_uid"`
	OrderID int64  `json:"order_id"`
	Message string `json:"message"`
}
