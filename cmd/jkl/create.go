package main

import (
	"bytes"
	"errors"
	"flag"
	"os"

	"otremblay.com/jkl"
)

type CreateCmd struct {
	args    []string
	project string
	file    string
}

func NewCreateCmd(args []string) (*CreateCmd, error) {
	ccmd := &CreateCmd{project: os.Getenv("JIRA_PROJECT")}
	f := flag.NewFlagSet("x", flag.ExitOnError)
	f.StringVar(&ccmd.project, "p", "", "Jira project key")
	f.StringVar(&ccmd.file, "f", "", "File to get issue description from")
	f.Parse(args)
	return ccmd, nil
}

var ErrCcmdJiraProjectRequired = errors.New("Jira project needs to be set")

func (ccmd *CreateCmd) Create() error {
	var b = bytes.NewBufferString(CREATE_TEMPLATE)
	var iss *jkl.JiraIssue
	var err error
	if ccmd.file != "" {
		iss, err = GetIssueFromFile(ccmd.file, b)
		if err != nil {
			return err
		}
	} else {
		iss, err = GetIssueFromTmpFile(b)
		if err != nil {
			return err
		}

	}
	if iss.Fields != nil &&
		(iss.Fields.Project == nil || iss.Fields.Project.Key == "") {
		iss.Fields.Project = &jkl.Project{Key: ccmd.project}
	}
	return jkl.Create(iss)
}

const CREATE_TEMPLATE = `Issue Type:
Summary:
Description:`
