package jkl

import (
	"os"
	"testing"
	"text/template"
)

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
