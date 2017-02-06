package main

import (
	"bytes"
	"errors"
	"flag"
	"io"
	"os"
	"fmt"
	"text/template"

	"otremblay.com/jkl"
)

type CreateCmd struct {
	args    []string
	project string
	file    string
	issuetype string
}

func NewCreateCmd(args []string) (*CreateCmd, error) {
	ccmd := &CreateCmd{project: os.Getenv("JIRA_PROJECT")}
	f := flag.NewFlagSet("x", flag.ExitOnError)
	f.StringVar(&ccmd.project, "p", "", "Jira project key")
	f.StringVar(&ccmd.file, "f", "", "File to get issue description from")
	f.Parse(args)
	ccmd.args = f.Args()
	return ccmd, nil
}

var ErrCcmdJiraProjectRequired = errors.New("Jira project needs to be set")

func (ccmd *CreateCmd) Create() error {
	var b = bytes.NewBuffer([]byte{})
	var readfile bool
	if fp := os.Getenv("JIRA_ISSUE_TEMPLATE"); fp != "" {
		if f, err := os.Open(fp); err == nil {
			_, err := io.Copy(b, f)
			if err == nil {
				readfile = true
			}

		}
	}
	
	if ccmd.project == "" {
		return ErrCcmdJiraProjectRequired
	}
	isstype := ""
	if len(ccmd.args) > 0 {
		isstype = ccmd.args[0]
	}
	cm, err := jkl.GetCreateMeta(ccmd.project, isstype)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Error getting the CreateMeta for project [%s] and issue types [%s]", ccmd.project, isstype), err)
	}
	
	if !readfile {
		createTemplate.Execute(b, cm)
	}
	var iss *jkl.JiraIssue
	// TODO: Evil badbad don't do this.
	em := &jkl.EditMeta{Fields: cm.Projects[0].IssueTypes[0].Fields}
	if ccmd.file != "" {
		iss, err = GetIssueFromFile(ccmd.file, b, em)
		if err != nil {
			return err
		}
	} else {
		iss, err = GetIssueFromTmpFile(b, em)
		if err != nil {
			return err
		}

	}
	if iss.Fields != nil &&
		(iss.Fields.Project == nil || iss.Fields.Project.Key == "") {
		iss.Fields.Project = &jkl.Project{Key: ccmd.project}
	}
	iss, err = jkl.Create(iss)
	if err != nil {
		return err
	}
	fmt.Println(iss.Key)
	return nil
}

func (ccmd *CreateCmd) Run() error {
	return ccmd.Create()
}

var createTemplate = template.Must(template.New("createissue").Parse(`{{range .Projects -}}
Project: {{.Key}}
{{range .IssueTypes -}}
Issue Type: {{.Name}}
Summary:
Description:
{{.RangeFieldSpecs}}
{{end}}
{{end}}`))
