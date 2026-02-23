package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/charmbracelet/huh"
)

const IssuesURL = "https://api.github.com/repos/"

type issueRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// SendRequest请求发生错误的时候返回的json结构(用404的仓库测试得到)
type GitHubError struct {
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
	Status           string `json:"status"`
}

var TestRequest = issueRequest{
	Title: "test",
	Body:  "test string",
}

func MakeIssue() (op *issueRequest, err error) {
	op = new(issueRequest)
	err = huh.NewInput().
		Title("输入构造issue的标题").
		Prompt(">>?").
		Value(&op.Title).
		Run()
	if err != nil {
		err = fmt.Errorf("输入标题出现错误:%w", err)
		return
	}
	editor := os.Getenv("EDITOR")
	if editor == "" {
		err = huh.NewText().
			Title("输入构造的issue的主体内容").
			Value(&op.Body).
			ExternalEditor(false).
			Run()
		if err != nil {
			err = fmt.Errorf("输入主体出现错误:%w", err)
			return
		}
	} else {
		err = huh.NewText().
			Title("输入构造的issue的主体内容").
			Value(&op.Body).
			Run()
		if err != nil {
			err = fmt.Errorf("输入主体出现错误:%w", err)
			return
		}
	}
	return
}

func SendRequest(bodyData *issueRequest, target string) error {
	oauthKey := os.Getenv("ghPAT")

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(bodyData); err != nil {
		return fmt.Errorf("error while encode json: %w", err)
	}
	apiURL := IssuesURL + target + "/issues"
	req, err := http.NewRequest("POST", apiURL, &buf)
	if err != nil {
		return fmt.Errorf("error while makeHTTPRequest: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+oauthKey)

	fmt.Printf("向%s发送请求\n", apiURL)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error while sent request: %w", err)
	}
	defer res.Body.Close()
	resCode := res.StatusCode
	if resCode != 200 && resCode != 201 {
		bodyBytes, _ := io.ReadAll(res.Body)
		var errorBody GitHubError
		if err := json.Unmarshal(bodyBytes, &errorBody); err != nil {
			return fmt.Errorf("error when unmarshal error json: %w", err)
		} else {
			return fmt.Errorf("error in sent request: %s", errorBody.Message)
		}
	} else {
		//fmt.Println("request has been sent successfully.")
		return nil
	}
}
