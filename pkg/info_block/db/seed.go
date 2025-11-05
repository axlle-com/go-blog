package db

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/axlle-com/blog/pkg/info_block/service"
	template "github.com/axlle-com/blog/pkg/template/provider"
	user "github.com/axlle-com/blog/pkg/user/provider"
	"github.com/bxcodec/faker/v3"
)

type seeder struct {
	infoBlockService *service.InfoBlockService
	templateProvider template.TemplateProvider
	userProvider     user.UserProvider
	config           contract.Config
}

type InfoBlockSeedData struct {
	ID          *uint   `json:"id"`
	UserID      *uint   `json:"user_id"`
	Template    string  `json:"template"`
	Media       *string `json:"media"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Image       *string `json:"image"`
}

func NewSeeder(
	infoBlockService *service.InfoBlockService,
	templateProvider template.TemplateProvider,
	user user.UserProvider,
	cfg contract.Config,
) contract.Seeder {
	return &seeder{
		infoBlockService: infoBlockService,
		templateProvider: templateProvider,
		userProvider:     user,
		config:           cfg,
	}
}

func (s *seeder) Seed() error {
	return s.seedFromJSON("info_blocks")
}

func (s *seeder) seedFromJSON(moduleName string) error {
	layout := s.config.Layout()
	seedPath := s.config.SrcFolderBuilder("db", layout, "seed", fmt.Sprintf("%s.json", moduleName))

	// Проверяем существование файла
	if _, err := os.Stat(seedPath); os.IsNotExist(err) {
		logger.Infof("[info_block][seeder][seedFromJSON] seed file not found: %s, skipping", seedPath)
		return nil
	}

	// Читаем JSON файл
	data, err := os.ReadFile(seedPath)
	if err != nil {
		return err
	}

	var infoBlocksData []InfoBlockSeedData
	if err := json.Unmarshal(data, &infoBlocksData); err != nil {
		return err
	}

	idsUser := s.userProvider.GetAllIds()
	if len(idsUser) == 0 {
		return nil
	}

	for _, blockData := range infoBlocksData {
		resourceName := "info_blocks"
		var templateID *uint
		if blockData.Template != "" {
			tpl, err := s.templateProvider.GetByNameAndResource(blockData.Template, resourceName)
			if err != nil {
				logger.Errorf("[info_block][seeder][seedFromJSON] template not found: name=%s, resource=%s, error=%v", blockData.Template, resourceName, err)
				continue
			}
			id := tpl.GetID()
			templateID = &id
		}

		filter := models.NewInfoBlockFilter()
		filter.Title = &blockData.Title
		filter.TemplateID = templateID
		foundByParams, _ := s.infoBlockService.FindByFilter(filter)
		if foundByParams != nil {
			logger.Infof("[info_block][seeder][seedFromJSON] info block with title='%s' and template_id=%v already exists (ID=%d), skipping", blockData.Title, templateID, foundByParams.ID)
			continue
		}

		infoBlock := models.InfoBlock{
			TemplateID:  templateID,
			Media:       blockData.Media,
			Title:       blockData.Title,
			Description: blockData.Description,
			Image:       blockData.Image,
			CreatedAt:   db.TimePtr(time.Now()),
			UpdatedAt:   db.TimePtr(time.Now()),
		}

		if blockData.UserID != nil {
			foundUser, _ := s.userProvider.GetByID(*blockData.UserID)
			if foundUser != nil {
				userID := foundUser.GetID()
				infoBlock.UserID = &userID
			}
		}

		_, err := s.infoBlockService.Create(&infoBlock, nil)
		if err != nil {
			logger.Errorf("[info_block][seeder][seedFromJSON] error creating info block: %v", err)
			continue
		}
	}

	logger.Infof("[info_block][seeder][seedFromJSON] seeded %d info blocks from JSON", len(infoBlocksData))
	return nil
}

func (s *seeder) SeedTest(n int) error {
	return s.infoBlocks(n)
}

func (s *seeder) infoBlocks(n int) error {
	idsUser := s.userProvider.GetAllIds()
	ids := s.templateProvider.GetAllIds()

	for i := 1; i <= n; i++ {
		randomID := ids[rand.Intn(len(ids))]
		randomUserID := idsUser[rand.Intn(len(idsUser))]

		infoBlock := models.InfoBlock{
			TemplateID:  &randomID,
			Media:       db.StrPtr(faker.Word()),
			Title:       "TitleInfoBlock #" + strconv.Itoa(i),
			Description: db.StrPtr(faker.Paragraph()),
			Image:       db.StrPtr("/public/img/404.svg"),
			CreatedAt:   db.TimePtr(time.Now()),
			UpdatedAt:   db.TimePtr(time.Now()),
			DeletedAt:   nil,
			UserID:      &randomUserID,
		}

		_, err := s.infoBlockService.Create(&infoBlock, nil)
		if err != nil {
			return err
		}
	}
	logger.Info("Database seeded infoBlock successfully!")
	return nil
}
