package main

import (
	"bytes"
	"errors"
	"flag"
	"io"

	"otremblay.com/jkl"
)

type CommentCmd struct {
	args     []string
	file     string
	issueKey string
}

func NewCommentCmd(args []string) (*CommentCmd, error) {
	ccmd := &CommentCmd{}
	f := flag.NewFlagSet("comments", flag.ExitOnError)
	f.StringVar(&ccmd.file, "f", "", "File to get issue comment from")
	f.Parse(args)
	if len(f.Args()) < 1 {
		return nil, ErrNotEnoughArgs
	}
	ccmd.issueKey = f.Arg(0)
	return ccmd, nil
}

var ErrNotEnoughArgs = errors.New("Not enough arguments")

func (ccmd *CommentCmd) Comment() error {
	var b = bytes.NewBufferString("")
	var comment io.Reader
	var err error
	if ccmd.file != "" {
		comment, err = GetTextFromSpecifiedFile(ccmd.file, b)
		if err != nil {
			return err
		}
	} else {
		comment, err = GetTextFromTmpFile(b)
		if err != nil {
			return err
		}

	}

	io.Copy(b, comment)

	return jkl.AddComment(ccmd.issueKey, b.String())
}
