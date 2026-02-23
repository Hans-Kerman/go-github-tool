package github

import (
	"net/url"
	"strings"
	"time"
)

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

func FormatRepoPath(target string) (string, error) {
	u, err := url.Parse(target)
	//尝试解析https风格url
	if err != nil && u.Scheme != "" {
		path := u.Path
		path = strings.Trim(path, "/") //去除首尾'/'
		path = strings.TrimSuffix(path, ".git")
		return path, err
	}
	//尝试解析scp风格url
	idx := strings.Index(target, ":") //返回第一个:的位置
	path := target[idx+1:]
	path = strings.Trim(path, "/")
	path = strings.TrimSuffix(path, ".git")
	return path, err
}
