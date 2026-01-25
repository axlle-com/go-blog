package db

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"time"

	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/info_block/models"
	"github.com/axlle-com/blog/pkg/info_block/service"
	"github.com/bxcodec/faker/v3"
)

type seeder struct {
	config           contract.Config
	disk             contract.DiskService
	seedService      contract.SeedService
	api              *api.Api
	infoBlockService *service.InfoBlockService
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
	cfg contract.Config,
	disk contract.DiskService,
	seedService contract.SeedService,
	api *api.Api,
	infoBlockService *service.InfoBlockService,
) contract.Seeder {
	return &seeder{
		config:           cfg,
		disk:             disk,
		seedService:      seedService,
		api:              api,
		infoBlockService: infoBlockService,
	}
}

func (s *seeder) Seed() error {
	return s.seedFromJSON((&models.InfoBlock{}).GetTable())
}

func (s *seeder) seedFromJSON(moduleName string) error {
	idsUser := s.api.User.GetAllIds()
	if len(idsUser) == 0 {
		return nil
	}

	files, err := s.seedService.GetFiles(s.config.Layout(), moduleName)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		logger.Infof("[info_block][seeder][seedFromJSON] seed files not found for module: %s, skipping", moduleName)
		return nil
	}

	for name, seedPath := range files {
		data, err := s.disk.ReadFile(seedPath)
		if err != nil {
			return err
		}

		ok, err := s.seedService.IsApplied(name)
		if err != nil {
			return err
		}
		if ok {
			continue
		}

		var infoBlocksData []InfoBlockSeedData
		if err := json.Unmarshal(data, &infoBlocksData); err != nil {
			return err
		}

		for _, blockData := range infoBlocksData {
			templateName := blockData.Template

			filter := models.NewInfoBlockFilter()
			filter.Title = &blockData.Title
			filter.TemplateName = &templateName
			foundByParams, _ := s.infoBlockService.FindByFilter(filter)
			if foundByParams != nil {
				logger.Infof("[info_block][seeder][seedFromJSON] info block with title='%s' and template='%s' already exists (ID=%d), skipping", blockData.Title, templateName, foundByParams.ID)
				continue
			}

			infoBlock := models.InfoBlock{
				TemplateName: templateName,
				Media:        blockData.Media,
				Title:        blockData.Title,
				Description:  blockData.Description,
				Image:        blockData.Image,
				CreatedAt:    db.TimePtr(time.Now()),
				UpdatedAt:    db.TimePtr(time.Now()),
			}

			if blockData.UserID != nil {
				foundUser, _ := s.api.User.GetByID(*blockData.UserID)
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

		if err := s.seedService.MarkApplied(name); err != nil {
			return err
		}

		logger.Infof("[info_block][seeder][seedFromJSON] seeded %d info blocks from JSON (%s)", len(infoBlocksData), name)
	}

	return nil
}

func (s *seeder) SeedTest(n int) error {
	return s.infoBlocks(n)
}

func (s *seeder) infoBlocks(n int) error {
	idsUser := s.api.User.GetAllIds()
	templates := s.api.Template.GetAll()

	for i := 1; i <= n; i++ {
		var templateName string
		if len(templates) > 0 {
			templateName = templates[rand.Intn(len(templates))].GetName()
		}
		randomUserID := idsUser[rand.Intn(len(idsUser))]

		infoBlock := models.InfoBlock{
			TemplateName: templateName,
			Media:        db.StrPtr(faker.Word()),
			Title:        "TitleInfoBlock #" + strconv.Itoa(i),
			Description:  db.StrPtr(faker.Paragraph()),
			Image:        db.StrPtr("/static/img/404.svg"),
			CreatedAt:    db.TimePtr(time.Now()),
			UpdatedAt:    db.TimePtr(time.Now()),
			DeletedAt:    nil,
			UserID:       &randomUserID,
		}

		_, err := s.infoBlockService.Create(&infoBlock, nil)
		if err != nil {
			return err
		}
	}
	logger.Info("Database seeded infoBlock successfully!")
	return nil
}
