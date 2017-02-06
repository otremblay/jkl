package main

import (
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"strings"

	"otremblay.com/jkl"
)

type EditCommentCmd struct {
	args      []string
	file      string
	issueKey  string
	commentId string
}

func NewEditCommentCmd(args []string) (*EditCommentCmd, error) {
	ccmd := &EditCommentCmd{}
	f := flag.NewFlagSet("editcomments", flag.ExitOnError)
	f.StringVar(&ccmd.file, "f", "", "File to get issue comment from")
	f.Parse(args)
	if len(f.Args()) < 1 {
		return nil, ErrNotEnoughArgs
	}
	ids := strings.Split(f.Arg(0), jkl.CommentIdSeparator)
	ccmd.issueKey = ids[0]
	if len(ids) < 2 {
		if len(f.Args()) == 2 {
			ccmd.commentId = f.Args()[1]
		} else {
			return nil, ErrNotEnoughArgs
		}
	} else {
		ccmd.commentId = ids[1]
	}
	return ccmd, nil
}

func (e *EditCommentCmd) Run() error {


	// Get Comment
	comm, err := jkl.GetComment(e.issueKey, e.commentId)
	if err != nil {
		return err
	}
		b := bytes.NewBufferString(comm.Body)
	var rdr io.Reader
	if e.file != "" {
		rdr, err = GetTextFromSpecifiedFile(e.file, b)

		if err != nil {
			return err
		}
	} else {
		rdr, err = GetTextFromTmpFile(b)
		if err != nil {
			return err
		}
	}
	btxt, err := ioutil.ReadAll(rdr)
	txt := strings.TrimSpace(string(btxt))
	if txt == "" || txt == strings.TrimSpace(comm.Body) {
		// Nothing to do
		return nil
	}
	comm.Body = txt
	return jkl.EditComment(e.issueKey, e.commentId, comm)
}