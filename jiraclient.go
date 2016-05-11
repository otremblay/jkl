package jkl

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

var j, _ = cookiejar.New(nil)

var httpClient *JiraClient

type JiraClient struct {
	*http.Client
	jiraRoot string
}

func NewJiraClient(jiraRoot string) *JiraClient {
	j := &JiraClient{
		&http.Client{
			Jar: j,
		},
		jiraRoot,
	}
	if j.jiraRoot == "" {
		j.jiraRoot = os.Getenv("JIRA_ROOT")
	}
	return j
}

func (j *JiraClient) Do(req *http.Request) (*http.Response, error) {
	var err error
	req.SetBasicAuth(os.Getenv("JIRA_USER"), os.Getenv("JIRA_PASSWORD"))
	req.Header.Add("Content-Type", "application/json")
	req.URL, err = url.Parse(j.jiraRoot + "rest/" + req.URL.RequestURI())
	if err != nil {
		return nil, err
	}
	return j.Client.Do(req)
}

func (j *JiraClient) Put(path string, payload io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("PUT", path, payload)
	if err != nil {
		return nil, err
	}
	return j.Do(req)
}

func (j *JiraClient) Post(path string, payload io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", path, payload)
	if err != nil {
		return nil, err
	}
	return j.Do(req)
}

func (j *JiraClient) Get(path string) (*http.Response, error) {
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	return j.Do(req)
}
