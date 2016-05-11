package main

import (
	"errors"
	"fmt"

	"otremblay.com/jkl"
)

type TaskCmd struct{}

func (t *TaskCmd) Handle(args []string) error {
	if len(args) == 1 {
		return t.Get(args[0])
	}
	return ErrTaskSubCommandNotFound
}

var ErrTaskSubCommandNotFound = errors.New("Subcommand not found.")

func (t *TaskCmd) Get(taskKey string) error {
	issue, err := jkl.GetIssue(taskKey)
	if err != nil {
		return err
	}
	fmt.Println(issue)
	return nil
}
