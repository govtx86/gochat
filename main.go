package main

import (
	"fmt"
	"os"

	"github.com/govtx86/gochat/internal/app"
)

func main() {
	if len(os.Args) == 1 || (len(os.Args) > 1 && os.Args[1] == "client") {
		app.RunClient()
	} else if len(os.Args) > 1 && os.Args[1] == "server" {
		app.RunServer()
	} else {
		fmt.Println("Usage: gochat [client|server]")
	}
}
