package github

import "time"

// 用于接受的Issue,不能用于构造发送的请求
type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	User      *User
	CreatedAt time.Time `json:"created_at"`
	Body      string    //md格式
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}
