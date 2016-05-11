package main

import (
	"strings"
	"testing"
)

func TestIssueFromList(t *testing.T) {
	iss := IssueFromList(strings.Split(`Description: Cowboys
 from
hell
Issue Type: Sometype
this is ignored
Summary: Dookienator
also ignored`, "\n"))
	AssertEqual(t, `Cowboys
 from
hell`, iss.Fields.Description)
	AssertEqual(t, "Sometype", iss.Fields.IssueType.Name)
}

func TestSpacex(t *testing.T) {
	AssertEqual(t, "Something", spacex.ReplaceAllString("Some thing", ""))
}

func AssertEqual(t *testing.T, expected interface{}, actual interface{}) {
	if expected != actual {
		t.Errorf(`Assertation failed!
Asserted: %v
Actual: %v`, expected, actual)
	}
}

func Assert(t *testing.T, fn func() bool) {

}
