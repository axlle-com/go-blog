package main

import (
	"flag"
	"fmt"
	"github.com/axlle-com/blog/app"
	"github.com/axlle-com/blog/app/config"
	db2 "github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/pkg/message/db"
	"os"

	mGallery "github.com/axlle-com/blog/pkg/gallery/db/migrate"
	sInfoBlock "github.com/axlle-com/blog/pkg/info_block/db"
	mInfoBlock "github.com/axlle-com/blog/pkg/info_block/db/migrate"
	mMessage "github.com/axlle-com/blog/pkg/message/db/migrate"
	postDB "github.com/axlle-com/blog/pkg/post/db"
	postMigrate "github.com/axlle-com/blog/pkg/post/db/migrate"
	dbTemplate "github.com/axlle-com/blog/pkg/template/db"
	mTemplate "github.com/axlle-com/blog/pkg/template/db/migrate"
	dbUser "github.com/axlle-com/blog/pkg/user/db"
	mUser "github.com/axlle-com/blog/pkg/user/db/migrate"
	userRepository "github.com/axlle-com/blog/pkg/user/repository"
)

func main() {
	var command string
	db2.Init(config.Config().DBUrl())
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
	"seed-test": func() {
		seedTest()
	},
	"migrate": func() {
		migrate()
	},
	"refill": func() {
		db2.NewCache().ResetUsersSession()
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
	mUser.NewMigrator().Migrate()
	postMigrate.NewMigrator().Migrate()
	mTemplate.NewMigrator().Migrate()
	mGallery.NewMigrator().Migrate()
	mInfoBlock.NewMigrator().Migrate()
	mMessage.NewMigrator().Migrate()
}

func rollback() {
	mUser.NewMigrator().Rollback()
	postMigrate.NewMigrator().Rollback()
	mTemplate.NewMigrator().Rollback()
	mGallery.NewMigrator().Rollback()
	mInfoBlock.NewMigrator().Rollback()
	mMessage.NewMigrator().Rollback()
}

func seedTest() {
	container := app.New(nil)

	userSeeder := dbUser.NewSeeder(
		container.UserRepo,
		userRepository.NewRoleRepo(),
		userRepository.NewPermissionRepo(),
	)
	userSeeder.SeedTest(5)
	userSeeder.Seed()

	dbTemplate.NewSeeder(
		container.TemplateRepo,
	).SeedTest(10)

	postDB.NewSeeder(
		container.PostRepo,
		container.PostService,
		container.CategoryRepo,
		container.UserProvider,
		container.TemplateProvider,
	).SeedTest(100)

	sInfoBlock.NewSeeder(
		container.InfoBlockService,
		container.TemplateProvider,
		container.UserProvider,
	).SeedTest(20)

	db.NewMessageSeeder(
		container.MessageService,
		container.UserProvider,
	).SeedTest(200)
}
