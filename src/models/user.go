package models

import "strings"

type User struct {
	FullName string `json:"fullname,omitempty" validate:"required" gorm:"column:fullname"`
	NickName string `json:"nickname,omitempty" gorm:"column:nickname;primaryKey"`
	Email    string `json:"email,omitempty" validate:"required" gorm:"column:email"`
	About    string `json:"about,omitempty" gorm:"column:about"`
}

func UserComp(us1, us2 *User) bool {
	return strings.ToLower(us1.NickName) > strings.ToLower(us2.NickName)
}
