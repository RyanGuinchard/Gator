package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/ryanguinchard/gator/internal/config"
	"github.com/ryanguinchard/gator/internal/database"

	"os"
)

type state struct {
	db  *database.Queries
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

	userName := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), userName)
	if err != nil {
		fmt.Printf("user %s not found: %v\n", userName, err)
		os.Exit(1)
	}

	err = s.cfg.SetUser(userName)
	if err != nil {
		fmt.Printf("Error setting user in config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("User set to: %s\n", s.cfg.CurrentUserName)

	return nil
}

func handlerRegister(s *state, cmd command) error {

	if len(cmd.args) < 1 {
		fmt.Println("A name is required to register a user")
		os.Exit(1)
	}

	userName := cmd.args[0]

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      userName,
	})

	if err != nil {
		fmt.Printf("Error creating user: %v\n", err)
		os.Exit(1)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		fmt.Printf("Error setting user in config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("User created: %s\n", user.Name)
	fmt.Printf("Logged User Data: %+v\n", user)

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
