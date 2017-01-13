package jkl

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"text/template"
)

type Search struct {
	Issues []*JiraIssue `json:"issues"`
}

type IssueType struct {
	Name   string `json:"name"`
	Fields map[string]FieldSpec
}

type FieldSpec struct {
	Name     string
	Required bool
	Schema   struct {
		Type string
	}
}

type Project struct {
	Key        string `json:"key,omitempty"`
	IssueTypes []IssueType
}

type Author struct {
	Name        string
	DisplayName string
}

type Comment struct {
	Id     string
	Author *Author
	Body   string
}

type CommentColl struct {
	Comments []Comment
}

type Status struct {
	Name string
}

type TimeTracking struct {
	OriginalEstimateSeconds  int
	RemainingEstimateSeconds int
}

type Fields struct {
	*IssueType   `json:"issuetype,omitempty"`
	Assignee     *Author       `json:",omitempty"`
	Project      *Project      `json:"project,omitempty"`
	Summary      string        `json:"summary,omitempty"`
	Description  string        `json:"description,omitempty"`
	Comment      *CommentColl  `json:"comment,omitempty"`
	Parent       *JiraIssue    `json:",omitempty"`
	Status       *Status       `json:",omitempty"`
	TimeTracking *TimeTracking `json:"timetracking,omitempty"`
}

func (f *Fields) PrettyRemaining() string {
	return PrettySeconds(f.TimeTracking.RemainingEstimateSeconds)
}

func (f *Fields) PrettyOriginalEstimate() string {
	return PrettySeconds(f.TimeTracking.OriginalEstimateSeconds)
}

func PrettySeconds(seconds int) string {
	//This works because it's an integer division.
	days := seconds / 3600 / 8
	hours := seconds/3600 - (days * 8)
	minutes := (seconds - (hours * 3600) - (days * 8 * 3600)) / 60
	seconds = (seconds - (hours * 3600) - (minutes * 60) - (days * 8 * 3600))

	return fmt.Sprintf("%dd %2dh %2dm %2ds", days, hours, minutes, seconds)
}

type JiraIssue struct {
	Key    string  `json:"key,omitempty"`
	Fields *Fields `json:"fields"`
}

func (i *JiraIssue) URL() string {
	return os.Getenv("JIRA_ROOT") + "browse/" + i.Key
}

func (i *JiraIssue) String() string {
	var b = bytes.NewBuffer(nil)
	var tmpl *template.Template = issueTmpl
	if os.Getenv("JKLNOCOLOR") == "true" {
		tmpl = issueTmplNoColor
	}
	err := tmpl.Execute(b, i)
	if err != nil {
		log.Fatalln(err)
	}

	return b.String()
}

var commentTemplate = `{{if .Fields.Comment }}{{$k := .Key}}{{range .Fields.Comment.Comments}}{{.Author.DisplayName}} ({{$k}}#{{.Id}}):
-----------------
{{.Body}}
-----------------

{{end}}{{end}}`

var issueTmplTxt = "\x1b[1m{{.Key}}\x1b[0m\t{{if .Fields.IssueType}}[{{.Fields.IssueType.Name}}]{{end}}\t{{.Fields.Summary}}\n\n" +
    "\x1b[1mURL\x1b[0m: {{.URL}}\n\n" +
	"{{if .Fields.Status}}\x1b[1mStatus\x1b[0m:\t {{.Fields.Status.Name}}\n{{end}}" +
	"{{if .Fields.Assignee}}\x1b[1mAssignee:\x1b[0m\t{{.Fields.Assignee.Name}}\n{{end}}\n" +
	"\x1b[1mTime Remaining/Original Estimate:\x1b[0m\t{{.Fields.PrettyRemaining}} / {{.Fields.PrettyOriginalEstimate}}\n\n" +
	"\x1b[1mDescription:\x1b[0m   {{.Fields.Description}} \n\n" +
	"\x1b[1mComments:\x1b[0m\n\n" + commentTemplate

var issueTmplNoColorTxt = "{{.Key}}\t{{if .Fields.IssueType}}[{{.Fields.IssueType.Name}}]{{end}}\t{{.Fields.Summary}}\n\n" +
    "URL: {{.URL}}\n\n" +
	"{{if .Fields.Status}}Status:\t {{.Fields.Status.Name}}\n{{end}}" +
	"{{if .Fields.Assignee}}Assignee:\t{{.Fields.Assignee.Name}}\n{{end}}\n" +
	"Time Remaining/Original Estimate:\t{{.Fields.PrettyRemaining}} / {{.Fields.PrettyOriginalEstimate}}\n\n" +
	"Description:   {{.Fields.Description}} \n\n" +
	"Comments:\n\n" + commentTemplate

var issueTmpl = template.Must(template.New("issueTmpl").Parse(issueTmplTxt))
var issueTmplNoColor = template.Must(template.New("issueTmplNoColor").Parse(issueTmplNoColorTxt))