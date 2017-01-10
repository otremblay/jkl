package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"strings"

	"otremblay.com/jkl"
)

var verbose = flag.Bool("v", false, "Output debug information about jkl")
var help = flag.Bool("h", false, "Outputs usage information message")

func main() {
	jkl.FindRCFile()
	flag.Parse()
	jkl.Verbose = verbose
	if err := runcmd(flag.Args()); err != nil {
		log.Println(err)
	}
}

func runcmd(args []string) error {
	if len(args) == 0 {
		if *help {
			fmt.Fprintln(os.Stderr, usage)
			flag.PrintDefaults()
			return nil
		}
		args = append(args, "list")
	}
	if strings.Contains(args[0], "~") {
		args = append([]string{"edit-comment"}, args...)
	}
	cmd, err := getCmd(args, 0)
	if err != nil {
		return err
	}
	return cmd.Run()
}

func getCmd(args []string, depth int) (Runner, error) {
	switch args[0] {
	case "list":
		return NewListCmd(args[1:])
	case "create":
		return NewCreateCmd(args[1:])
	case "task":
		tcmd := &TaskCmd{args: args[1:]}
		return tcmd, nil
	case "edit":
		return NewEditCmd(args[1:])
	case "comment":
		return NewCommentCmd(args[1:])
	case "edit-comment":
		return NewEditCommentCmd(args[1:])
	default:
		// Think about this real hard.
		// I want `jkl JIRA-1234 done` to move it to done.
		// I want `jkl JIRA-1234` to print out info
		// I want `jkl JIRA-1234 edit` to run the edit command.
		// I want `jkl JIRA-1234 comment` to run the comment command.
		// I want `jkl JIRA-1234 attach <filename>` to run the attach command.

		if depth == 0 {
			// Assume args[0] is a task key
			args[0], args[1] = args[1], args[0]
			return getCmd(args, depth+1)
		} else {
			// Swapping the first two args didn't help;
			// this means it's a transition.

			// tcmd, err := NewTransitionCommand(args)
			// if err != nil {return nil, err}
			// return tcmd, nil
		}
	}
	return nil, ErrTaskSubCommandNotFound
}

const usage = `Usage:
jkl [options] <command> [args]

Available commands:

list
create
edit
comment

`
