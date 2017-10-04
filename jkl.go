package jkl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func init() {
	x := false
	Verbose = &x
}

var Verbose *bool
var defaultIssue = &JiraIssue{}

func bootHttpClient() {
	if httpClient == nil {
		httpClient = NewJiraClient("")
	}
}

func Create(issue *JiraIssue) (*JiraIssue, error) {
	bootHttpClient()
	payload, err := formatPayload(issue)
	if err != nil {
		return nil, err
	}
	//	fmt.Println(issue)
	resp, err := httpClient.Post("api/2/issue", payload)
	if err != nil {
		fmt.Println(resp.StatusCode)
		return nil, err
	}
	if resp.StatusCode >= 400 {
		io.Copy(os.Stderr, resp.Body)
		return nil, errors.New(fmt.Sprintf("HTTP error, %v", resp.StatusCode))
	}
	dec := json.NewDecoder(resp.Body)
	issue = &JiraIssue{}
	err = dec.Decode(issue)
	if err != nil {
		b, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(b))
		return nil, err
	}
	return issue, nil
}

func GetCreateMeta(projectKey, issueType string) (*CreateMeta, error) {
	bootHttpClient()
	path := fmt.Sprintf("api/2/issue/createmeta?expand=projects.issuetypes.fields&issuetypeNames=%s&projectKeys=%s", strings.Title(strings.ToLower(issueType)), projectKey)
	fmt.Println(path)
	resp, err := httpClient.Get(path)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		fmt.Println("Status code:", resp.StatusCode)
		fmt.Println("Response:")
		fmt.Println(string(b))
		return nil, errors.New("Some http error happened.")
	}
	fmt.Println(string(b))
	dec := json.NewDecoder(bytes.NewBuffer(b))
	var createmeta = &CreateMeta{}
	err = dec.Decode(createmeta)
	if err != nil {
		b, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(b))
		return nil, err
	}

	return createmeta, nil
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
	if resp.StatusCode >= 400 {
		fmt.Println("Status code:", resp.StatusCode)
		fmt.Println("Response:")
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		fmt.Println(string(b))
		return nil, errors.New("Some http error happened.")
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

	path := "api/2/issue/" + taskKey + "?expand=transitions,operations,editmeta"
	resp, err := httpClient.Get(path)
	if err != nil {
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

func GetComment(taskKey string, commentId string) (*Comment, error) {
	bootHttpClient()
	path := "api/2/issue/" + taskKey + "/comment/" + commentId
	resp, err := httpClient.Get(path)
	if err != nil {
		fmt.Println(resp.StatusCode)
		return nil, err
	}
	dec := json.NewDecoder(resp.Body)
	var comment = &Comment{}
	err = dec.Decode(comment)
	if err != nil {
		b, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(b))
		return nil, err
	}
	return comment, nil
}

func EditComment(taskKey string, commentId string, comment *Comment) error {
	bootHttpClient()
	payload, err := serializePayload(comment)
	if err != nil {
		return err
	}
	resp, err := httpClient.Put("api/2/issue/"+taskKey+"/comment/"+commentId, payload)
	if err != nil {
		fmt.Println(resp.StatusCode)
		return err
	}
	if resp.StatusCode >= 400 {
		io.Copy(os.Stderr, resp.Body)
	}
	return nil
}

func DoTransition(taskKey string, transitionName string) error {
	iss, err := GetIssue(taskKey)
	if err != nil {
		return err
	}
	var t *Transition
	fmt.Println(iss.Transitions)
	for _, transition := range iss.Transitions {
		if strings.ToLower(transition.Name) == strings.ToLower(transitionName) {
			t = transition
			break
		}
	}
	if t == nil {
		return errors.New("Transition not found")
	}
	payload, err := serializePayload(map[string]interface{}{"transition": t})
	resp, err := httpClient.Post("api/2/issue/"+taskKey+"/transitions/", payload)
	if err != nil {
		fmt.Println(resp.StatusCode)
		return err
	}
	if resp.StatusCode >= 400 {
		io.Copy(os.Stderr, resp.Body)
	}
	return nil
}

func LogWork(taskKey string, workAmount string) error {
	bootHttpClient()
	payload, err := serializePayload(map[string]interface{}{"timeSpent": workAmount})
	resp, err := httpClient.Post("api/2/issue/"+taskKey+"/worklog", payload)
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
	return serializePayload(issue)
}

func serializePayload(i interface{}) (io.Reader, error) {
	var b []byte
	payload := bytes.NewBuffer(b)
	enc := json.NewEncoder(payload)
	err := enc.Encode(i)
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
