package main

import (
	"flag"
	"os"
	"strings"

	"text/template"

	"otremblay.com/jkl"
)

var listTemplateStr string
var listTemplate *template.Template

func init() {
	flag.StringVar(&listTemplateStr, "listTemplate", "{{.Key}}\t({{.Fields.IssueType.Name}}{{if .Fields.Parent}} of {{.Fields.Parent.Key}}{{end}})\t{{.Fields.Summary}}\n", "Go template used in list command")
	listTemplate = template.Must(template.New("listTemplate").Parse(listTemplateStr))
}

func List(args []string) error {
	if issues, err := jkl.List(strings.Join(args, " ")); err != nil {
		return err
	} else {
		for _, issue := range issues {
			listTemplate.Execute(os.Stdout, issue)
		}
	}
	return nil
}
