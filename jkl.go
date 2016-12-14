package jkl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var defaultIssue = &JiraIssue{}

func bootHttpClient() {
	if httpClient == nil {
		httpClient = NewJiraClient("")
	}
}

func Create(issue *JiraIssue) error {
	bootHttpClient()
	payload, err := formatPayload(issue)
	if err != nil {
		return err
	}
	//	fmt.Println(issue)
	resp, err := httpClient.Post("api/2/issue", payload)
	if err != nil {
		fmt.Println(resp.StatusCode)
		return err
	}
	if resp.StatusCode >= 400 {
		io.Copy(os.Stderr, resp.Body)
	}
	return nil
}

func Edit(issue *JiraIssue) error {
	bootHttpClient()
	payload, err := formatPayload(issue)
	if err != nil {
		return err
	}
	resp, err := httpClient.Put("api/2/issue/"+issue.Key, payload)
	if err != nil {
		fmt.Println(resp.StatusCode)
		return err
	}
	if resp.StatusCode >= 400 {
		io.Copy(os.Stderr, resp.Body)
	}
	return nil
}

func List(jql string) ([]*JiraIssue, error) {
	bootHttpClient()
	path := "api/2/search?fields=*all&maxResults=1000"
	if jql != "" {
		path += "&jql=" + url.QueryEscape(jql)
	}
	resp, err := httpClient.Get(path)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	dec := json.NewDecoder(resp.Body)
	var issues = &Search{}
	err = dec.Decode(issues)
	if err != nil {
		b, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(b))
		return nil, err
	}
	return issues.Issues, nil
}

func GetIssue(taskKey string) (*JiraIssue, error) {
	bootHttpClient()
	path := "api/2/issue/" + taskKey
	resp, err := httpClient.Get(path)
	if err != nil {
		fmt.Println(resp.StatusCode)
		return nil, err
	}
	dec := json.NewDecoder(resp.Body)
	var issue = &JiraIssue{}
	err = dec.Decode(issue)
	if err != nil {
		b, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(b))
		return nil, err
	}
	return issue, nil
}

func AddComment(taskKey string, comment string) error {
	bootHttpClient()
	var b []byte
	payload := bytes.NewBuffer(b)
	enc := json.NewEncoder(payload)
	enc.Encode(map[string]string{"body": comment})
	resp, err := httpClient.Post("api/2/issue/"+taskKey+"/comment", payload)
	if err != nil {
		fmt.Println(resp.StatusCode)
		return err
	}
	if resp.StatusCode >= 400 {
		io.Copy(os.Stderr, resp.Body)
	}
	return nil
}

func formatPayload(issue *JiraIssue) (io.Reader, error) {
	if issue.Fields != nil &&
		issue.Fields.Project != nil &&
		issue.Fields.Project.Key == "" {
		issue.Fields.Project.Key = os.Getenv("JIRA_PROJECT")
	}
	var b []byte
	payload := bytes.NewBuffer(b)
	enc := json.NewEncoder(payload)
	err := enc.Encode(issue)
	fmt.Println(payload.String())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return payload, nil
}

func FindRCFile() string {
	dir, err := os.Getwd()
	if err != nil {

		log.Fatalln(err)
	}
	path := strings.Split(dir, "/")
	for i := len(path); i > 0; i-- {
		dotenvpath := strings.Join(path[0:i], "/") + "/.jklrc"
		err := godotenv.Load(dotenvpath)
		if err == nil {
			return dotenvpath
		}
	}
	log.Fatalln("No .jklrc found")
	return ""
}
