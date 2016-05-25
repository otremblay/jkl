package jkl

import (
	"bytes"
	"log"
	"text/template"
)

type Search struct {
	Issues []*Issue `json:"issues"`
}

type IssueType struct {
	Name string `json:"name"`
}
type Project struct {
	Key string `json:"key,omitempty"`
}

type Author struct {
	Name        string
	DisplayName string
}

type Comment struct {
	Author *Author
	Body   string
}

type CommentColl struct {
	Comments []Comment
}

type Status struct {
	Name string
}

type Fields struct {
	*IssueType  `json:"issuetype,omitempty"`
	Assignee    *Author      `json:",omitempty"`
	Project     *Project     `json:"project,omitempty"`
	Summary     string       `json:"summary,omitempty"`
	Description string       `json:"description,omitempty"`
	Comment     *CommentColl `json:"comment,omitempty"`
	Parent      *Issue       `json:",omitempty"`
	Status      *Status      `json:",omitempty"`
}
type Issue struct {
	Key    string  `json:"key,omitempty"`
	Fields *Fields `json:"fields"`
}

func (i *Issue) String() string {
	var b = bytes.NewBuffer(nil)
	err := issueTmpl.Execute(b, i)
	if err != nil {
		log.Fatalln(err)
	}

	return b.String()
}

var commentTemplate = `{{if .Fields.Comment }}{{range .Fields.Comment.Comments}}{{.Author.DisplayName}}:
-----------------
{{.Body}}
-----------------

{{end}}{{end}}`

var issueTmplTxt = "\x1b[1m{{.Key}}\x1b[0m\t{{if .Fields.IssueType}}[{{.Fields.IssueType.Name}}]{{end}}\t{{.Fields.Summary}}\n\n" +
	"\x1b[1mStatus\x1b[0m:\t {{.Fields.Status.Name}}\n" +
	"\x1b[1mAssignee:\x1b[0m\t{{.Fields.Assignee.Name}}\n\n" +
	"\x1b[1mDescription:\x1b[0m   {{.Fields.Description}} \n\n" +
	"\x1b[1mComments:\x1b[0m\n\n" + commentTemplate

var issueTmpl = template.Must(template.New("issueTmpl").Parse(issueTmplTxt))
