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

type listissue jkl.JiraIssue

func (l *listissue) URL() string {
	i := jkl.JiraIssue(*l)
	return (&i).URL()
}

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
func (l *listissue) EFByName(name string) string {
	return (*jkl.JiraIssue)(l).EFByName(name)
}

type ListCmd struct {
	args    []string
	tmplstr string
	tmpl    *template.Template
}

func NewListCmd(args []string) (*ListCmd, error) {
	ccmd := &ListCmd{}
	f := flag.NewFlagSet("x", flag.ExitOnError)
	if *verbose {
		fmt.Println(&ccmd.tmplstr)
	}
	f.StringVar(&ccmd.tmplstr, "listTemplate", "{{.Color}}{{.Key}}{{if .Color}}\x1b[39m{{end}}\t({{.Fields.IssueType.Name}}{{if .Fields.Parent}} of {{.Fields.Parent.Key}}{{end}})\t{{.Fields.Summary}}\t{{if .Fields.Assignee}}[{{.Fields.Assignee.Name}}]{{end}}\n", "Go template used in list command")
	f.Parse(args)
	ccmd.args = f.Args()
	if len(ccmd.args) == 0 {
		proj := os.Getenv("JIRA_PROJECT")
		if proj != "" {
			proj = fmt.Sprintf(" and project = '%s'", proj)
		}
		ccmd.args = []string{fmt.Sprintf("sprint in openSprints() %s order by rank", proj)}
		if *verbose {
			fmt.Println("No arguments, running default command")
		}
	}
	if *verbose {
		fmt.Println(ccmd.args)
	}
	ccmd.tmpl = template.Must(template.New("listTemplate").Parse(ccmd.tmplstr))
	return ccmd, nil
}

func (l *ListCmd) List() error {
	if issues, err := jkl.List(strings.Join(l.args, " ")); err != nil {
		return err
	} else {
		for issue := range issues {
			var li listissue
			li = listissue(*issue)
			err := l.tmpl.Execute(os.Stdout, &li)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
	return nil
}

func (l *ListCmd) Run() error {
	return l.List()
}
