package models

import (
	"github.com/go-openapi/strfmt"
)

type Thread struct {
	Id      uint64          `json:"id,omitempty" gorm:"column:id"`
	Title   string          `json:"title,omitempty" gorm:"column:title"`
	Author  string          `json:"author,omitempty" gorm:"column:author"`
	Forum   string          `json:"forum,omitempty" gorm:"column:forum"`
	Message string          `json:"message,omitempty" gorm:"column:message"`
	Votes   int64           `json:"votes,omitempty" gorm:"column:votes"`
	Slug    string          `json:"slug,omitempty" gorm:"column:slug"`
	Created strfmt.DateTime `json:"created,omitempty" gorm:"column:created"`
}
