package jkl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"
)

var j, _ = cookiejar.New(nil)

var httpClient *JiraClient

type JiraClient struct {
	*http.Client
	jiraRoot string
}

func init() {
	x := false
	Verbose = &x
}

func NewJiraClient(jiraRoot string) *JiraClient {
	jc := &JiraClient{
		&http.Client{
			Jar: j,
		},
		jiraRoot,
	}
	if jc.jiraRoot == "" {
		jc.jiraRoot = os.Getenv("JIRA_ROOT")
	}
	if cookiefile := os.Getenv("JIRA_COOKIEFILE"); cookiefile != "" {
		makeNewFile := false
		f, err := os.Open(cookiefile)
		server := jc.jiraRoot + "rest/gadget/1.0/login"
		u, _ := url.Parse(server)
		if err != nil {
			makeNewFile = true
		} else {
			if stat, err := f.Stat(); err == nil {
				if time.Now().Sub(stat.ModTime()).Minutes() > 60 {
					makeNewFile = true
				} else {
					var cookies []*http.Cookie
					dec := json.NewDecoder(f)
					dec.Decode(&cookies)
					u, _ = url.Parse(jc.jiraRoot)
					jc.Jar.SetCookies(u, cookies)
				}
			}
			f.Close()
		}
		if makeNewFile {
			f, err = os.Create(cookiefile)
			if err != nil {
				panic(err)
			}

			http.DefaultClient.Jar = j
			form := url.Values{}
			form.Add("os_username", os.Getenv("JIRA_USER"))
			form.Add("os_password", os.Getenv("JIRA_PASSWORD"))
			req, _ := http.NewRequest("POST", server, strings.NewReader(form.Encode()))
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			resp, err := http.DefaultClient.Do(req)
			if err != nil || resp.StatusCode >= 400 {
				fmt.Println(resp.Header)
				fmt.Println(resp.Status)
				fmt.Println(err)
			}
			b := bytes.NewBuffer(nil)
			enc := json.NewEncoder(b)
			enc.Encode(j.Cookies(u))
			io.Copy(f, b)
			f.Close()
		}
	}

	if *Verbose {
		fmt.Println("Jira root:", jc.jiraRoot)
	}
	return jc
}

func (j *JiraClient) DoLess(req *http.Request) (*http.Response, error) {
	var err error
	if os.Getenv("JIRA_COOKIEFILE") == "" {
		req.SetBasicAuth(os.Getenv("JIRA_USER"), os.Getenv("JIRA_PASSWORD"))
	}
	if *Verbose {
		fmt.Println("Jira User: ", os.Getenv("JIRA_USER"))
		fmt.Println("Jira Password: ", os.Getenv("JIRA_PASSWORD"))
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json, text/plain, text/html")
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

func (j *JiraClient) Do(req *http.Request) (*http.Response, error) {
	var err error
	req.URL, err = url.Parse(j.jiraRoot + "rest/" + req.URL.RequestURI())
	if err != nil {
		return nil, err
	}
	return j.DoLess(req)
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
