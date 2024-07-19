package main

import (
	"bufio"
	"flag"
	"fmt"
	DB "github.com/axlle-com/blog/pkg/common/db"
	"github.com/axlle-com/blog/pkg/common/models"
	post "github.com/axlle-com/blog/pkg/post/db"
	user "github.com/axlle-com/blog/pkg/user/db"
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
		user.SeedUsers(100)
		post.SeedPosts(100)
	},
	"migrate": func() {
		db := DB.GetDB()
		err := db.AutoMigrate(
			&models.Post{},
			&models.User{},
			&models.PostCategory{},
			&models.Template{},
		)
		if err != nil {
			log.Fatalln(err)
		}
	},
	"refill": func() {
		db := DB.GetDB()
		err := db.Migrator().DropTable(
			&models.Post{},
			&models.User{},
			&models.PostCategory{},
			&models.Template{},
		)
		if err != nil {
			log.Fatalln(err)
		}
		err = db.AutoMigrate(
			&models.Post{},
			&models.User{},
			&models.PostCategory{},
			&models.Template{},
		)
		if err != nil {
			log.Fatalln(err)
		}
		user.SeedUsers(100)
		post.SeedPosts(100)
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
