package main

import (
	"fmt"

	"github.com/ryanguinchard/gator/internal/config"

	"os"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	commandsMap map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {

	if len(cmd.args) < 1 {
		fmt.Println("Username required")
		os.Exit(1)

	}

	s.cfg.SetUser(cmd.args[0])

	fmt.Printf("User set to: %s\n", s.cfg.CurrentUserName)

	return nil
}

func (c *commands) run(s *state, cmd command) error {

	fn, ok := c.commandsMap[cmd.name]
	if !ok {
		return fmt.Errorf("Unknown command: %s", cmd.name)
	}

	return fn(s, cmd)

}

func (c *commands) register(name string, fn func(*state, command) error) {
	c.commandsMap[name] = fn
}
