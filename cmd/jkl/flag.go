package main

import (
	"errors"
	"flag"

	"otremblay.com/jkl"
)

type FlagCmd struct {
	args []string
	flg  bool
}

func NewFlagCmd(args []string, flg bool) (*FlagCmd, error) {
	ccmd := &FlagCmd{flg: flg}
	f := flag.NewFlagSet("flag", flag.ExitOnError)
	f.Parse(args)
	if len(f.Args()) < 1 {
		return nil, ErrFlagNotEnoughArgs
	}
	ccmd.args = f.Args()
	return ccmd, nil
}

var ErrFlagNotEnoughArgs = errors.New("Not enough arguments, need at least one issue key")

func (ccmd *FlagCmd) Flag() error {
	return jkl.FlagIssue(ccmd.args, ccmd.flg)
}

func (ccmd *FlagCmd) Run() error {
	return ccmd.Flag()
}
