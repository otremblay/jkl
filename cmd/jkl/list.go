package main

import (
	"flag"
	"os"
	"strings"

	"text/template"

	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"otremblay.com/jkl"
)

var listTemplateStr string
var listTemplate *template.Template

func init() {
	flag.StringVar(&listTemplateStr, "listTemplate", "{{.Color}}{{.Key}}{{if .Color}}\x1b[39m{{end}}\t({{.Fields.IssueType.Name}}{{if .Fields.Parent}} of {{.Fields.Parent.Key}}{{end}})\t{{.Fields.Summary}}\t[{{.Fields.Assignee.Name}}]\n", "Go template used in list command")
	listTemplate = template.Must(template.New("listTemplate").Parse(listTemplateStr))
}

type listissue jkl.Issue

func (l *listissue) Color() string {
	if os.Getenv("JKLNOCOLOR") == "true" || !terminal.IsTerminal(int(os.Stdout.Fd())) {
		return ""
	}
	if strings.Contains(os.Getenv("RED_ISSUE_STATUSES"), l.Fields.Status.Name) {
		return "\x1b[31m"
	}

	if strings.Contains(os.Getenv("GREEN_ISSUE_STATUSES"), l.Fields.Status.Name) {
		return "\x1b[32m"
	}
	if strings.Contains(os.Getenv("BLUE_ISSUE_STATUSES"), l.Fields.Status.Name) {
		return "\x1b[34m"
	}
	if strings.Contains(os.Getenv("YELLOW_ISSUE_STATUSES"), l.Fields.Status.Name) || os.Getenv("YELLOW_ISSUE_STATUSES") == "default" {
		return "\x1b[33m"
	}
	return ""
}

func List(args []string) error {
	if issues, err := jkl.List(strings.Join(args, " ")); err != nil {
		return err
	} else {
		for _, issue := range issues {
			var li listissue
			li = listissue(*issue)
			err := listTemplate.Execute(os.Stdout, &li)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
	return nil
}
