package main

import (
	"errors"
	"fmt"
	"strings"

	"otremblay.com/jkl"
)

type TaskCmd struct {
	args []string
}

// TODO: split in individual commands.
func (t *TaskCmd) Handle() error {
	c := len(t.args)
	if c == 1 {
		return t.Get(t.args[0])
	}
	if c == 2 {
		// fmt.Println(t.args)
		err := t.Transition(t.args[0], t.args[1])
		if err != nil {
			//fmt.Println(err)
			return t.Log(t.args[0], strings.Join(t.args[1:], " "))
		}
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
	return jkl.DoTransition(taskKey, transition)
}

func (t *TaskCmd) Log(taskKey, time string) error {
	return jkl.LogWork(taskKey, time)
}

func (t *TaskCmd) Run() error {
	return t.Handle()
}
