package main

import (
	"flag"

	"otremblay.com/jkl"
)

type AttachCmd struct {
	args    []string
	file    string
	taskKey string
}

func NewAttachCmd(args []string) (*AttachCmd, error) {
	ccmd := &AttachCmd{}
	f := flag.NewFlagSet("x", flag.ExitOnError)
	f.Parse(args)
	if len(f.Args()) < 2 {
		return nil, ErrNotEnoughArgs
	}
	ccmd.taskKey = f.Arg(0)
	ccmd.file = f.Arg(1)
	return ccmd, nil
}

func (ecmd *AttachCmd) Attach() error {
	return jkl.Attach(ecmd.taskKey, ecmd.file)
}

func (ecmd *AttachCmd) Run() error {
	return ecmd.Attach()
}
