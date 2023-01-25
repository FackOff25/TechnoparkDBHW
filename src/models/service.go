package models

type ServiceStatus struct {
	UserCount   int64 `json:"user"`
	ForumCount  int64 `json:"forum"`
	ThreadCount int64 `json:"thread"`
	PostCount   int64 `json:"post"`
}
