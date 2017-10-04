package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"unicode"

	"reflect"

	"bufio"

	"otremblay.com/jkl"
)

// def get_editor do
// 	[System.get_env("EDITOR"), "nano", "vim", "vi"]
// 	|> Enum.find(nil, fn (ed) -> System.find_executable(ed) != nil end)
//   end
var editors = []string{os.Getenv("EDITOR"), "nano", "vim", "vi"}

// GetEditor returns the path to an editor, taking $EDITOR in account
func GetEditor() string {
	for _, ed := range editors {
		if p, err := exec.LookPath(ed); err == nil {
			return p
		}
	}
	log.Fatal("No editor available; use flags.")
	return ""
}

func copyInitial(dst io.WriteSeeker, initial io.Reader) {
	io.Copy(dst, initial)
	dst.Seek(0, 0)
}

func GetIssueFromTmpFile(initial io.Reader, editMeta *jkl.EditMeta) (*jkl.JiraIssue, error) {
	f, err := ioutil.TempFile(os.TempDir(), "jkl")
	if err != nil {
		return nil, fmt.Errorf("Error opening tempfile: %v", err)
	}
	copyInitial(f, initial)
	f2, err := GetTextFromFile(f)
	if err != nil {
		return nil, fmt.Errorf("Error reading tempfile: %v", err)
	}
	return IssueFromReader(f2, editMeta), nil
}

func GetTextFromTmpFile(initial io.Reader) (io.Reader, error) {
	f, err := ioutil.TempFile(os.TempDir(), "jkl")
	if err != nil {
		return nil, err
	}
	copyInitial(f, initial)
	return GetTextFromFile(f)
}

func GetTextFromSpecifiedFile(filename string, initial io.Reader) (io.Reader, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	if fi, err := f.Stat(); err == nil && fi.Size() == 0 {
		copyInitial(f, initial)
	}
	return GetTextFromFile(f)
}

func GetTextFromFile(file *os.File) (io.Reader, error) {
	var err error
	if !*SilentMode {
		cmd := exec.Command(GetEditor(), file.Name())
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		_, err = file.Seek(0, 0)
	}
	return file, err
}

func GetIssueFromFile(filename string, initial io.Reader, editMeta *jkl.EditMeta) (*jkl.JiraIssue, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	if fi, err := f.Stat(); err == nil && fi.Size() == 0 {
		copyInitial(f, initial)
	}
	f2, err := GetTextFromFile(f)
	if err != nil {
		return nil, err
	}
	return IssueFromReader(f2, editMeta), nil
}

var spacex = regexp.MustCompile(`\s`)

func IssueFromReader(f io.Reader, editMeta *jkl.EditMeta) *jkl.JiraIssue {
	iss := &jkl.JiraIssue{Fields: &jkl.Fields{}}
	riss := reflect.ValueOf(iss).Elem()
	fieldsField := riss.FieldByName("Fields").Elem()
	currentField := reflect.Value{}
	currFieldName := ""
	brd := bufio.NewReader(f)
	for {
		b, _, err := brd.ReadLine()
		if err != nil {
			break
		}
		parts := strings.Split(string(b), ":")
		potentialField := spacex.ReplaceAllString(parts[0], "")

		// Is the current line a field in an issue directly?
		// Also special cases: Objects that have a deeper depth
		// have specific fields "flattened" for ease of use.
		// I think this loop could be made more general, to account
		// for deeper objects. Then again, there's not that many fields
		// I actually care about yet.
		// Custom fields are gonna be hell.

		if newfield := fieldsField.FieldByName(potentialField); newfield.IsValid() {
			parts = parts[1:len(parts)]
			if potentialField == "IssueType" {
				if len(parts) > 0 {
					iss.Fields.IssueType = &jkl.IssueType{}
					currentField = reflect.Value{}
					currFieldName = potentialField
					f2 := newfield.Elem()
					f3 := f2.FieldByName("Name")
					f3.SetString(strings.TrimSpace(strings.Join(parts, ":")))
				}
			} else if potentialField == "Project" {
				if len(parts) > 0 {
					iss.Fields.Project = &jkl.Project{}
					currentField = reflect.Value{}
					currFieldName = potentialField
					f2 := newfield.Elem()
					f3 := f2.FieldByName("Key")
					f3.SetString(strings.TrimSpace(strings.Join(parts, ":")))
				}
			} else if potentialField == "Parent" {
				if len(parts) > 0 {
					iss.Fields.Parent = &jkl.JiraIssue{}
					currentField = reflect.Value{}
					currFieldName = potentialField
					f2 := newfield.Elem()
					f3 := f2.FieldByName("Key")
					f3.SetString(strings.TrimSpace(strings.Join(parts, ":")))
				}
			} else {
				currFieldName = potentialField
				currentField = newfield
			}
		} else if editMeta != nil {
			// If it's not valid, throw it at the createmeta. It will probably end up in ExtraFields.

		}
		if currentField.IsValid() {
			newString := currentField.String() + "\n" + strings.Join(parts, ":")
			if currFieldName != "Description" {
				newString = strings.TrimSpace(newString)
			} else if currentField.String() == "" {
				newString = strings.TrimLeftFunc(newString, unicode.IsSpace)
			}
			currentField.SetString(newString)
		}
	}
	return iss
}

func IssueFromList(list []string, editMeta *jkl.EditMeta) *jkl.JiraIssue {
	return IssueFromReader(bytes.NewBufferString(strings.Join(list, "\n")), editMeta)
}
