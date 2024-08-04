package main

import (
	"bufio"
	"flag"
	"fmt"
	DB "github.com/axlle-com/blog/pkg/common/db"
	common "github.com/axlle-com/blog/pkg/common/models"
	gallery "github.com/axlle-com/blog/pkg/gallery/models"
	post "github.com/axlle-com/blog/pkg/post/db"
	"github.com/axlle-com/blog/pkg/post/models"
	rights "github.com/axlle-com/blog/pkg/rights/db"
	template "github.com/axlle-com/blog/pkg/template/db"
	user "github.com/axlle-com/blog/pkg/user/db"
	"gorm.io/gorm"
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
		template.SeedTemplate(100)
		rights.SeedPermissions()
		rights.SeedRoles()
		user.SeedUsers(100)
		post.SeedPostCategory(100)
		post.SeedPosts(100)
	},
	"migrate": func() {
		db := DB.GetDB()
		err := db.AutoMigrate(
			&models.Post{},
			&common.User{},
			&models.PostCategory{},
			&common.Template{},
			&common.Role{},
			&common.Permission{},
			&gallery.Gallery{},
			&gallery.GalleryImage{},
			&gallery.GalleryHasResource{},
		)
		if err != nil {
			log.Fatalln(err)
		}
	},
	"refill": func() {
		db := DB.GetDB()
		dropIntermediateTables(db)
		err := db.Migrator().DropTable(
			&models.Post{},
			&common.User{},
			&models.PostCategory{},
			&common.Template{},
			&common.Role{},
			&common.Permission{},
			&gallery.Gallery{},
			&gallery.GalleryImage{},
			&gallery.GalleryHasResource{},
		)
		if err != nil {
			log.Fatalln(err)
		}
		err = db.AutoMigrate(
			&models.Post{},
			&common.User{},
			&models.PostCategory{},
			&common.Template{},
			&common.Role{},
			&common.Permission{},
			&gallery.Gallery{},
			&gallery.GalleryImage{},
			&gallery.GalleryHasResource{},
		)
		if err != nil {
			log.Fatalln(err)
		}
		template.SeedTemplate(100)
		rights.SeedPermissions()
		rights.SeedRoles()
		user.SeedUsers(100)
		post.SeedPostCategory(100)
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

func dropIntermediateTables(db *gorm.DB) {
	migrator := db.Migrator()
	intermediateTables := []string{
		"user_has_role",
		"user_has_permission",
		"role_has_permission",
	}
	for _, table := range intermediateTables {
		if err := migrator.DropTable(table); err != nil {
			fmt.Println("Error dropping table:", table, err)
			return
		}
		fmt.Println("Dropped intermediate table:", table)
	}
}
