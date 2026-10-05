package main

import (
	"database/sql"
	"fmt"

	"github.com/ryanguinchard/gator/internal/config"

	"github.com/ryanguinchard/gator/internal/database"

	"os"

	_ "github.com/lib/pq"
)

func main() {

	// Read the configuration
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		return
	}

	// Connect to the database
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		return
	}
	defer db.Close()

	dbQueries := database.New(db)
	programState := &state{db: dbQueries, cfg: &cfg}

	commands := commands{commandsMap: make(map[string]func(*state, command) error)}
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)

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
