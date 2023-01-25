package models

type Vote struct {
	ThreadId uint64 `gorm:"column:thread_id"`
	NickName string `json:"nickname,omitempty" gorm:"column:nickname"`
	Voice    int64  `json:"voice,omitempty" gorm:"column:voice"`
}
