package main

import (
	"fmt"

	"github.com/ryanguinchard/gator/internal/config"

	"os"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		return
	}

	programState := &state{cfg: &cfg}

	commands := commands{commandsMap: make(map[string]func(*state, command) error)}
	commands.register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("No command provided")
		os.Exit(1)
	}

	cmd := command{name: os.Args[1], args: os.Args[2:]}

	err = commands.run(programState, cmd)
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
	}
}
