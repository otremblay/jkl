package main

import (
	"flag"
	"fmt"
	"log"

	"os"

	"github.com/joho/godotenv"
	"strings"
)

func main() {
	findRCFile()
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Print(usage)
		return
	}
	if err := runcmd(flag.Args()); err != nil {
		log.Println(err)
	}
}

func findRCFile() {
	dir, err := os.Getwd()
	if err != nil {

		log.Fatalln(err)
	}
	path := strings.Split(dir, "/")
	for i := len(path) - 1; i > 0; i-- {
		err := godotenv.Load(strings.Join(path[0:i], "/") + "/.jklrc")
		if err == nil {
			return
		}
	}
	log.Fatalln("No .jklrc found")
}

func runcmd(args []string) error {
	switch args[0] {
	case "list":
		lcmd, err := NewListCmd(flag.Args()[1:])
		if err != nil {
			return err
		}
		return lcmd.List()
	case "create":
		ccmd, err := NewCreateCmd(flag.Args()[1:])
		if err != nil {
			return err
		}
		return ccmd.Create()
	case "task":
		tcmd := &TaskCmd{}
		return tcmd.Handle(flag.Args()[1:])
	case "edit":
		ecmd := &EditCmd{}
		return ecmd.Edit(flag.Arg(1))
	case "comment":
		ccmd, err := NewCommentCmd(flag.Args()[1:])
		if err != nil {
			return err
		}
		return ccmd.Comment()
	}
	fmt.Println(usage)
	return nil
}

const usage = `Usage:
jkl [options] <command> [args]

Available commands:

list
create
edit
`
