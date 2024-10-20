package main

import (
	"flag"
	"fmt"
	"github.com/axlle-com/blog/pkg/common/config"
	"github.com/axlle-com/blog/pkg/common/db"
	mGallery "github.com/axlle-com/blog/pkg/gallery/db/migrate"
	dbPost "github.com/axlle-com/blog/pkg/post/db"
	mPost "github.com/axlle-com/blog/pkg/post/db/migrate"
	dbTemplate "github.com/axlle-com/blog/pkg/template/db"
	mTemplate "github.com/axlle-com/blog/pkg/template/db/migrate"
	dbUser "github.com/axlle-com/blog/pkg/user/db"
	mUser "github.com/axlle-com/blog/pkg/user/db/migrate"
	"os"
)

func main() {
	var command string
	db.Init(config.Config().DBUrl())
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
	"seed-teat": func() {
		seedTest()
	},
	"migrate": func() {
		migrate()
	},
	"refill": func() {
		db.Cache().ResetUsersSession()
		rollback()
		migrate()
		seedTest()
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
	mUser.Migrate()
	mPost.Migrate()
	mTemplate.Migrate()
	mGallery.Migrate()
}

func rollback() {
	mUser.Rollback()
	mPost.Rollback()
	mTemplate.Rollback()
	mGallery.Rollback()
}

func seedTest() {
	dbUser.SeedPermissions()
	dbUser.SeedRoles()
	dbUser.SeedUsers(5)
	dbUser.SeedUsersDefault()

	dbTemplate.SeedTemplate(5)
	dbPost.SeedPostCategory(10)
	dbPost.SeedPosts(100)
}
