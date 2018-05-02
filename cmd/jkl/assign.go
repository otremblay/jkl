package main

import (
	"errors"
	"flag"

	"otremblay.com/jkl"
)

type AssignCmd struct {
	args     []string
	assignee string
	issueKey string
}

func NewAssignCmd(args []string) (*AssignCmd, error) {
	ccmd := &AssignCmd{}
	f := flag.NewFlagSet("assign", flag.ExitOnError)
	f.Parse(args)
	if len(f.Args()) < 2 {
		return nil, ErrAssignNotEnoughArgs
	}
	ccmd.issueKey = f.Arg(0)
	ccmd.assignee = f.Arg(1)
	return ccmd, nil
}

var ErrAssignNotEnoughArgs = errors.New("Not enough arguments, need issue key + assignee")

func (ccmd *AssignCmd) Assign() error {
	return jkl.Assign(ccmd.issueKey, ccmd.assignee)
}

func (ccmd *AssignCmd) Run() error {
	return ccmd.Assign()
}
