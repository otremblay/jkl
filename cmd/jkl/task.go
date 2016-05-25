package main

import (
	"errors"
	"fmt"

	"otremblay.com/jkl"
)

type TaskCmd struct{}

func (t *TaskCmd) Handle(args []string) error {
	c := len(args)
	if c == 1 {
		return t.Get(args[0])
	}
	if c == 2 {
		return t.Transition(args[0], args[1])
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

func (t *TaskCmd) Transition(taskKey, transition string) error {
	return nil
}
