package message

import "time"

type Message struct {
	MsgId    int64     `json:"msg_id" xorm:"msg_id"`       // 消息id
	TypeId   int64     `json:"type_id" xorm:"type_id"`     // 消息类型,对应message_type
	Title    string    `json:"title" xorm:"title"`         // 消息标题
	Content  string    `json:"content" xorm:"content"`     // 消息内容
	CreateAt time.Time `json:"create_at" xorm:"create_at"` // 写入时间
}

func (*Message) TableName() string {
	return "message"
}

func (*Message) DatabaseName() string {
	return "student"
}
