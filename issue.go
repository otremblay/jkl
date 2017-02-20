package jkl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"text/template"
)

type Search struct {
	Issues []*JiraIssue `json:"issues"`
}

type IssueType struct {
	Name   string `json:"name"`
	IconURL string `json:",omitempty"`
	Fields map[string]*FieldSpec
}

func (it *IssueType) RangeFieldSpecs() string {
	output := bytes.NewBuffer(nil)
	keys := make([]string, 0, len(it.Fields))
	for k, _ := range it.Fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintln(output, fmt.Sprintf("%s: %v", it.Fields[k].Name, it.Fields[k].AllowedValues))
	}
	return output.String()
}

type AllowedValue struct {
	Id string
	Self string
	Value string
}

func (a *AllowedValue) String() string{
	return a.Value
}

type FieldSpec struct {
	Name     string
	Required bool
	Schema   struct {
		Type string
		Custom string
		CustomId int
		Items string
	}
	Operations []string
	AllowedValues []*AllowedValue
	
}


type CreateMeta struct {
	Projects []*Project
}

type Project struct {
	Key        string `json:"key,omitempty"`
	Name string
	IssueTypes []*IssueType
}

type Author struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

func (a *Author) String() string {
	return a.DisplayName
}

type Comment struct {
	Id     string `json:"id"`
	Author *Author `json:"author"`
	Body   string `json:"body"`
}

type CommentColl struct {
	Comments []Comment
}

type Status struct {
	Name string
	IconURL string `json:",omitempty"`
}

type TimeTracking struct {
	OriginalEstimateSeconds  int
	RemainingEstimateSeconds int
}

func (a *Attachment) String() string {
	return fmt.Sprintf("Filename: [%s]\tAuthor:[%s]\tURL:[%s]", a.Filename, a.Author.Name, a.Content)
}

type Attachment struct {
	Filename string
	Author *Author
	Content string
}

type LinkType struct {
	Id string
	Name string
	Inward string
	Outward string
}

func (i *IssueLink) String() string {
	if i.InwardIssue != nil {
		return fmt.Sprintf("%s --> %s: %s", i.LinkType.Inward, i.InwardIssue.Key, i.InwardIssue.Fields.Summary)
	}
	if i.OutwardIssue != nil {
		return fmt.Sprintf("%s --> %s: %s", i.LinkType.Outward, i.OutwardIssue.Key, i.OutwardIssue.Fields.Summary)
	}
	return "Issue Link Got Weird"
}

type IssueLink struct {
	LinkType *LinkType `json:"type"`
	InwardIssue *JiraIssue
	OutwardIssue *JiraIssue
}

func (w *Worklog) String() string {
	return fmt.Sprintf("%s worked %s on %s", w.Author.Name, w.TimeSpent, w.Started)
}

type Worklog struct {
	Author *Author
	Comment string
	TimeSpent string
	TimeSpentSeconds int
	Started string
}

type Worklogs struct {
	Worklogs []*Worklog `json:",omitempty"`
}

func (f *Fields) UnmarshalJSON(b []byte) error{
	err := json.Unmarshal(b, &f.rawFields)
	if err != nil {
		fmt.Println("splosion")
		return err
	}
	f.rawExtraFields = map[string]json.RawMessage{}
	vf := reflect.ValueOf(f).Elem()
	for key, mess := range f.rawFields {
		field := vf.FieldByNameFunc(func(s string)bool{return strings.ToLower(key) == strings.ToLower(s)})
		if field.IsValid() {
		    objType := field.Type()
    		obj := reflect.New(objType).Interface()
			err := json.Unmarshal(mess, &obj)
			if err != nil {
				fmt.Fprintln(os.Stderr, objType, obj, string(mess))
				fmt.Fprintln(os.Stderr, errors.New(fmt.Sprintf("%s [%s]: %s","Error allocating field",key, err)))
			}
			field.Set(reflect.ValueOf(obj).Elem())
		} else {
			f.rawExtraFields[key] = mess
		}
	}

	return nil
}

type Priority struct {
	Name string
}

type Fields struct {
	*IssueType   `json:"issuetype,omitempty"`
	Assignee     *Author       `json:",omitempty"`
	Project      *Project      `json:"project,omitempty"`
	Summary      string        `json:"summary,omitempty"`
	Description  string        `json:"description,omitempty"`
	Comment      *CommentColl  `json:"comment,omitempty"`
	Parent       *JiraIssue    `json:",omitempty"`
	Status       *Status       `json:",omitempty"`
	TimeTracking *TimeTracking `json:"timetracking,omitempty"`
	Attachment   []*Attachment   `json:"attachment,omitempty"`
	IssueLinks   []*IssueLink    `json:"issueLinks,omitempty"`
	Priority     *Priority `json:",omitempty"`
	Worklog *Worklogs `json:"worklog,omitempty"`
	rawFields map[string]json.RawMessage
	rawExtraFields map[string]json.RawMessage
	ExtraFields map[string]interface{} `json:"-"`
}

func (f *Fields) PrettyRemaining() string {
	return PrettySeconds(f.TimeTracking.RemainingEstimateSeconds)
}

func (f *Fields) PrettyOriginalEstimate() string {
	return PrettySeconds(f.TimeTracking.OriginalEstimateSeconds)
}

func PrettySeconds(seconds int) string {
	//This works because it's an integer division.
	days := seconds / 3600 / 8
	hours := seconds/3600 - (days * 8)
	minutes := (seconds - (hours * 3600) - (days * 8 * 3600)) / 60
	seconds = (seconds - (hours * 3600) - (minutes * 60) - (days * 8 * 3600))

	return fmt.Sprintf("%dd %2dh %2dm %2ds", days, hours, minutes, seconds)
}

type Transition struct{
	Id string `json:"id"`
	Name string `json:"name"`
}


type Schema struct {
	System string
	Custom string
	CustomId int
}


type EditMeta struct {
	Fields map[string]*FieldSpec
}



type JiraIssue struct {
	Key    string  `json:"key,omitempty"`
	Fields *Fields `json:"fields,omitempty"`
	Transitions []*Transition `json:"transitions,omitempty"`
	EditMeta *EditMeta `json:"editmeta,omitempty"`
}


var sprintRegexp =  regexp.MustCompile(`name=([^,]+),`)
func (i *JiraIssue) UnmarshalJSON(b []byte) error {
	tmp := map[string]json.RawMessage{}
	if len(b) == 0 {
		return nil
	}
	i.Fields = &Fields{}
	i.EditMeta = &EditMeta{Fields:map[string]*FieldSpec{}}
	err := json.Unmarshal(b, &tmp)
	if err != nil && *Verbose {
		fmt.Fprintln(os.Stderr, errors.New(fmt.Sprintf("%s: %s", "Error unpacking raw json", err)))
	}
	err = json.Unmarshal(tmp["fields"], &i.Fields)
	if err != nil && *Verbose {
		fmt.Fprintln(os.Stderr, errors.New(fmt.Sprintf("%s: %s", "Error unpacking fields", err)))
	}
	err = json.Unmarshal(tmp["transitions"], &i.Transitions)
	if err != nil && *Verbose {
		fmt.Fprintln(os.Stderr,  errors.New(fmt.Sprintf("%s: %s", "Error unpacking transitions", err)))
	}
	err = json.Unmarshal(tmp["editmeta"], &i.EditMeta)
	if err != nil && *Verbose{
		fmt.Fprintln(os.Stderr, errors.New(fmt.Sprintf("%s: %s", "Error unpacking EditMeta", err)))
	}
	err = json.Unmarshal(tmp["key"], &i.Key)
	if err != nil && *Verbose {
		fmt.Fprintln(os.Stderr, errors.New(fmt.Sprintf("%s: %s", "Error unpacking key", err)))
	}

	i.Fields.ExtraFields = map[string]interface{}{}
	for k, v := range i.Fields.rawExtraFields {
		if f, ok := i.EditMeta.Fields[k]; ok {
			if f.Schema.Custom == "com.pyxis.greenhopper.jira:gh-sprint" {
				results := sprintRegexp.FindStringSubmatch(string(v))
				if len(results) == 2 {
					i.Fields.ExtraFields[k] = results[1]
				}
			} else {switch f.Schema.Type {
			case "user":
				a := &Author{}
				json.Unmarshal(v, &a)
				i.Fields.ExtraFields[k] = a
		    case "option":
		    	val := &AllowedValue{}
		    	err = json.Unmarshal(v, &val)
		    	if err != nil {panic(err)}
		    	i.Fields.ExtraFields[k] = val
		    case "array":
		    	if f.Schema.Items == "option" {
		    		val := []*AllowedValue{}
		    		err = json.Unmarshal(v, &val)
		    		if err != nil {panic(err)}
		    		i.Fields.ExtraFields[k] = val
		    		continue
		    	}
		    	fallthrough
			default:
				if string(v) != "null" {
				i.Fields.ExtraFields[k] = string(v)
				}
			}
			}
		}
	}
	
	return nil
}

func (i *JiraIssue) URL() string {
	return os.Getenv("JIRA_ROOT") + "browse/" + i.Key
}

func (i *JiraIssue) String() string {
	var b = bytes.NewBuffer(nil)
	var tmpl *template.Template = issueTmpl
	if os.Getenv("JKLNOCOLOR") == "true" {
		tmpl = issueTmplNoColor
	}
	err := tmpl.Execute(b, i)
	if err != nil {
		log.Fatalln(err)
	}

	return b.String()
}

func (i *JiraIssue) PrintExtraFields() string{
	sorter := map[string]string{}
	b := bytes.NewBuffer(nil)
	for k, v := range i.Fields.ExtraFields {
		if f, ok := i.EditMeta.Fields[k]; ok && v != nil {
			sorter[f.Name] = fmt.Sprintf("%s: %s", f.Name, v)
		}
	}
	keys := make([]string, 0, len(sorter))
	for k, _ := range sorter {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, v := range keys {
		fmt.Fprintln(b, sorter[v])
	}
	return b.String()
}

var commentTemplate = `{{if .Fields.Comment }}{{$k := .Key}}{{range .Fields.Comment.Comments}}{{.Author.DisplayName}} [~{{.Author.Name}}] ({{$k}}`+CommentIdSeparator+`{{.Id}}):
-----------------
{{.Body}}
-----------------

{{end}}{{end}}`

var issueTmplTxt = "\x1b[1m{{.Key}}\x1b[0m\t{{if .Fields.IssueType}}[{{.Fields.IssueType.Name}}]{{end}}\t{{.Fields.Summary}}\n\n" +
    "\x1b[1mURL\x1b[0m: {{.URL}}\n\n" +
	"{{if .Fields.Status}}\x1b[1mStatus\x1b[0m:\t {{.Fields.Status.Name}}\n{{end}}" +
	"{{if .Fields.Priority}}\x1b[1mStatus\x1b[0m:\t {{.Fields.Priority.Name}}\n{{end}}" +
	"\x1b[1mTransitions\x1b[0m: {{range .Transitions}}[{{.Name}}] {{end}}\n"+
	"{{if .Fields.Assignee}}\x1b[1mAssignee:\x1b[0m\t{{.Fields.Assignee.Name}}\n{{end}}\n" +
	"\x1b[1mTime Remaining/Original Estimate:\x1b[0m\t{{.Fields.PrettyRemaining}} / {{.Fields.PrettyOriginalEstimate}}\n\n" +
	"{{$i := .}}{{range $k, $v := .Fields.ExtraFields}}{{with index $i.EditMeta.Fields $k}}\x1b[1m{{.Name}}\x1b[0m{{end}}: {{$v}}\n{{end}}\n\n"+
	"\x1b[1mDescription:\x1b[0m   {{.Fields.Description}} \n\n" +
	"\x1b[1mIssue Links\x1b[0m: \n{{range .Fields.IssueLinks}}\t{{.}}\n{{end}}\n\n" + 
	"\x1b[1mComments:\x1b[0m\n\n" + commentTemplate +
	"Worklog:\n{{range .Fields.Worklog.Worklogs}}\t{{.}}\n{{end}}"

var issueTmplNoColorTxt = "{{.Key}}\t{{if .Fields.IssueType}}[{{.Fields.IssueType.Name}}]{{end}}\t{{.Fields.Summary}}\n\n" +
    "URL: {{.URL}}\n\n" +
	"{{if .Fields.Status}}Status:\t {{.Fields.Status.Name}}\n{{end}}" +
	"Transitions: {{range .Transitions}}[{{.Name}}] {{end}}\n"+
	"{{if .Fields.Assignee}}Assignee:\t{{.Fields.Assignee.Name}}\n{{end}}\n" +
	"Time Remaining/Original Estimate:\t{{.Fields.PrettyRemaining}} / {{.Fields.PrettyOriginalEstimate}}\n\n" +
	"{{.PrintExtraFields}}\n\n"+
	"Description:   {{.Fields.Description}} \n\n" +
	"Issue Links: \n{{range .Fields.IssueLinks}}\t{{.}}\n{{end}}\n\n" + 
	"Comments:\n\n" + commentTemplate+
	"Worklog:\n{{range .Fields.Worklog.Worklogs}}\t{{.}}\n{{end}}"

var CommentIdSeparator = "~"
var issueTmpl = template.Must(template.New("issueTmpl").Parse(issueTmplTxt))
var issueTmplNoColor = template.Must(template.New("issueTmplNoColor").Parse(issueTmplNoColorTxt))