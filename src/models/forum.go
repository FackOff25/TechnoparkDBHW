package models

type Forum struct {
	Title   string `json:"title,omitempty" gorm:"column:title"`
	User    string `json:"user,omitempty" gorm:"column:user"`
	Slug    string `json:"slug,omitempty" gorm:"column:slug"`
	Posts   uint64 `json:"posts,omitempty" gorm:"column:posts"`
	Threads uint64 `json:"threads,omitempty" gorm:"column:threads"`
}

type ForumUser struct {
	Forum string `gorm:"column:forum"`
	User  string `gorm:"column:user"`
}
