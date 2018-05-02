package main

import (
	"errors"
	"flag"

	"otremblay.com/jkl"
)

type LinkCmd struct {
	args []string
}

func NewLinkCmd(args []string) (*LinkCmd, error) {
	ccmd := &LinkCmd{}
	f := flag.NewFlagSet("Link", flag.ExitOnError)
	f.Parse(args)
	ccmd.args = f.Args()
	return ccmd, nil
}

var ErrLinkNotEnoughArgs = errors.New("Not enough arguments, need at least two issue keys and a reason")

func (ccmd *LinkCmd) Link() error {
	return jkl.LinkIssue(ccmd.args)
}

func (ccmd *LinkCmd) Run() error {
	return ccmd.Link()
}