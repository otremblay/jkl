package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"text/template"

	"otremblay.com/jkl"
)

type EditCmd struct {
	args    []string
	project string
	file    string
	taskKey string
}

func NewEditCmd(args []string) (*EditCmd, error) {
	ccmd := &EditCmd{project: os.Getenv("JIRA_PROJECT")}
	f := flag.NewFlagSet("x", flag.ExitOnError)
	f.StringVar(&ccmd.project, "p", "", "Jira project key")
	f.StringVar(&ccmd.file, "f", "", "File to get issue description from")
	f.Parse(args)
	ccmd.taskKey = f.Arg(0)
	return ccmd, nil
}

func (ecmd *EditCmd) Edit() error {
	b := bytes.NewBuffer(nil)
	iss, err := jkl.GetIssue(ecmd.taskKey)
	if err != nil {
		return fmt.Errorf("Edit failed: %v", err)
	}
	err = editTmpl.Execute(b, iss)
	if err != nil {
		return err
	}

	if ecmd.file != "" {
		iss, err = GetIssueFromFile(ecmd.file, b, iss.EditMeta)

		if err != nil {
			return fmt.Errorf("Error getting issue from file: %v", err)
		}
	} else {
		iss, err = GetIssueFromTmpFile(b, iss.EditMeta)
		if err != nil {
			return fmt.Errorf("Error getting issue from temp file: %v", err)
		}

	}
	iss.Key = ecmd.taskKey
	return jkl.Edit(iss)
}

func (ecmd *EditCmd) Run() error {
	return ecmd.Edit()
}

const EDIT_TEMPLATE = `Summary: {{.Fields.Summary}}
Description: {{.Fields.Description}}
`

var editTmpl = template.Must(template.New("editTmpl").Parse(EDIT_TEMPLATE))
