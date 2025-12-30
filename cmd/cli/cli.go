package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/axlle-com/blog/app/config"
	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/di"
	"github.com/axlle-com/blog/app/logger"
)

var container *di.Container

func main() {
	var command string

	cfg := config.Config()

	newDB, err := db.SetupDB(cfg)
	if err != nil {
		panic("db not initialized")
	}

	container = di.NewContainer(cfg, newDB)

	flag.StringVar(&command, "command", "", "Command to execute")
	flag.Parse()

	if command != "" {
		handleCommand(command)
	} else {
		fmt.Println("No task provided. Use -command=name to specify a task.")
	}

	os.Exit(1)
}

var Commands = map[string]func(){
	"hello": func() {
		fmt.Println("Hello!")
	},
	"seed": func() {
		seed()
	},
	"migrate": func() {
		migrate()
	},
	"seed-test": func() {
		seedTest()
	},
	"refill": func() {
		rollback()
		migrate()
		//seed()
		//seedTest()
	},
}

func handleCommand(command string) {
	if cmdFunc, exists := Commands[command]; exists {
		cmdFunc()
	} else {
		fmt.Println("Unknown command:", command)
	}
}

func migrate() {
	err := container.Migrator.Migrate()
	if err != nil {
		logger.Errorf("[cli][migrate] Error: %v", err)
	}
}

func rollback() {
	container.Cache.ResetUsersSession()
	err := container.Migrator.Rollback()
	if err != nil {
		logger.Errorf("[cli][rollback] Error: %v", err)
	}
}

func seedTest() {
	seed()
	err := container.Seeder.SeedTest(100)
	if err != nil {
		logger.Errorf("[cli][seedTest] Error: %v", err)
	}
}

func seed() {
	err := container.Seeder.Seed()
	if err != nil {
		logger.Errorf("[cli][seed] Error: %v", err)
	}
}
