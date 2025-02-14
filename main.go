package main

import (
	"fmt"
	"os"
	"os/user"
	"yap/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Hello %s! Welcome to yappanese!\n", user.Username)
		fmt.Println("Try to write some code into the command-line")
		repl.Start(os.Stdin, os.Stdout)
	}
}
