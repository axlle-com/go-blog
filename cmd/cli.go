package main

import (
	"bufio"
	"flag"
	"fmt"
	db2 "github.com/axlle-com/blog/pkg/common/db"
	"github.com/axlle-com/blog/pkg/common/models"
	. "github.com/axlle-com/blog/pkg/user/db"
	"log"
	"os"
	"strings"
)

func main() {
	var command string
	flag.StringVar(&command, "command", "", "Command to execute")
	flag.Parse()

	if command != "" {
		handleCommand(command)
	} else {
		handleCommands()
	}
}

var Commands = map[string]func(){
	"hello": func() {
		fmt.Println("Hello!")
	},
	"seed": func() {
		SeedUsers(100)
	},
	"migrate": func() {
		db := db2.GetDB()
		if err := db.AutoMigrate(&models.Post{}, &models.User{}); err != nil {
			log.Fatalln(err)
		}
	},
}

func handleCommand(command string) {
	if cmdFunc, exists := Commands[command]; exists {
		cmdFunc()
	} else {
		fmt.Println("Unknown command:", command)
	}
}

func handleCommands() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter command: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading command:", err)
			continue
		}
		input = strings.TrimSpace(input)
		handleCommand(input)
	}
}
