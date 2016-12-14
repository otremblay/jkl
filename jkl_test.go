package jkl

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"text/template"
)

func TestUnmarshalProjects(t *testing.T) {
	f, err := os.Open("projects.json")
	if err != nil {
		t.Error(err)
	}
	dec := json.NewDecoder(f)
	x := struct{ Projects []Project }{}

	err = dec.Decode(&x)
	if err != nil {
		t.Error(err)
	}
	for _, p := range x.Projects {
		for _, it := range p.IssueTypes {
			for sn, f := range it.Fields {
				fmt.Println(it.Name, sn, f.Name, f.Required, f.Schema.Type)
			}
		}
	}
}

type TestType struct {
	Field string
}

func (t *TestType) String() string {
	return t.Field
}

func TestStringerInTemplate(t *testing.T) {
	x := template.Must(template.New("stuff").Parse("{{.}}"))
	x.Execute(os.Stdout, &TestType{"This works"})
}
