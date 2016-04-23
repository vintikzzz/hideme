package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/vintikzzz/hideme/proxylist"
)

const baseURI = "http://hideme.ru/api/proxylist.txt?out=js&code="

func main() {
	proxylistCommand := flag.NewFlagSet("proxylist", flag.ExitOnError)
	codeFlag := proxylistCommand.String("code", "", "Your access code")
	if len(os.Args) == 1 {
		fmt.Println("usage: hideme <command> [<args>]")
		fmt.Println("The most commonly used commands are: ")
		fmt.Println(" proxylist   Load proxy list")
		return
	}
	switch os.Args[1] {
	case "proxylist":
		proxylistCommand.Parse(os.Args[2:])
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}
	if proxylistCommand.Parsed() {
		if *codeFlag == "" {
			fmt.Println("Please supply code using -code option.")
			return
		}
		prs, err := proxylist.Load(baseURI + *codeFlag)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		for _, pr := range prs {
			fmt.Println(pr.ToURL().String())
		}
	}

}
