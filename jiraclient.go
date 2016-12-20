package jkl

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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
	if *Verbose {
		fmt.Println("Jira root:", j.jiraRoot)
	}
	return j
}

func (j *JiraClient) Do(req *http.Request) (*http.Response, error) {
	var err error
	req.SetBasicAuth(os.Getenv("JIRA_USER"), os.Getenv("JIRA_PASSWORD"))
	if *Verbose {
		fmt.Println("Jira User: ", os.Getenv("JIRA_USER"))
		fmt.Println("Jira Password: ", os.Getenv("JIRA_PASSWORD"))
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json, text/plain, text/html")
	req.URL, err = url.Parse(j.jiraRoot + "rest/" + req.URL.RequestURI())
	if err != nil {
		return nil, err
	}
	resp, err := j.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		fmt.Println("Status code:", resp.StatusCode)
		if *Verbose {
			fmt.Println("Headers:")
			fmt.Println(resp.Header)
		}
		fmt.Println("Response:")
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		fmt.Println(string(b))
		return nil, errors.New("Some http error happened.")
	}
	return resp, nil
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
