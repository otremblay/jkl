package main

import (
	"errors"
	"flag"
)

type EditCommentCmd struct {
	args     []string
	file     string
	issueKey string
}

func NewEditCommentCmd(args []string) (*EditCommentCmd, error) {
	ccmd := &EditCommentCmd{}
	f := flag.NewFlagSet("comments", flag.ExitOnError)
	f.StringVar(&ccmd.file, "f", "", "File to get issue comment from")
	f.Parse(args)
	if len(f.Args()) < 1 {
		return nil, ErrNotEnoughArgs
	}
	ccmd.issueKey = f.Arg(0)
	return ccmd, nil
}

func (e *EditCommentCmd) Run() error {
	return errors.New("Not implemented")
}
