package main

import (
	"fmt"
	"os"
	"os/user"
	"github.com/assimad8/go-interpreter/internal/repl"
)
func main(){
	user,err := user.Current()
	if err!=nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the EMAD programming language!\n",user.Username)
	fmt.Print("Feel free to type in Commands\n")
	repl.Start(os.Stdin,os.Stdout)
}