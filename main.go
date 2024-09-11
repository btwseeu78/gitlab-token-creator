package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

type Tokenrespone struct {
	Id               string `json:"id"`
	Token            string `json:"token"`
	Token_Expires_At string `json:"token_expires_at"`
}

func main() {
	token := flag.String("token", "", "Private Technical Token For runner creation")
	gitlabUrl := flag.String("url", "", "gitlab url for project creation")
	flag.Parse()

	if *token == "" || *gitlabUrl == "" {
		fmt.Println("Please provide token and gitlab url")
	}

	//url := "https://gitlab.com/api/v4/user/runners"
	url := fmt.Sprintf("%s/api/v4/user/runners", *gitlabUrl)
	method := "POST"
	groupId := "56239669"
	tag_list := []string{"helm", "guild", "runner"}
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("runner_type", "group_type")
	_ = writer.WriteField("group_id", groupId)
	_ = writer.WriteField("description", "test runner")
	_ = writer.WriteField("tag_list", strings.Join(tag_list, ","))
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("PRIVATE-TOKEN", *token)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
	var structToken Tokenrespone
	if err := json.Unmarshal(body, &structToken); err != nil {
		fmt.Println("unable to parse token")
	}
	fmt.Printf("Token for group: %s is %v", groupId, structToken.Token)

}
