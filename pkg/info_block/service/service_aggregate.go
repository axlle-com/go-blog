package service

import (
	"sync"

	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/service"
	"github.com/axlle-com/blog/pkg/info_block/models"
)

type AggregateService struct {
	api *api.Api
}

func NewAggregateService(
	api *api.Api,
) *AggregateService {
	return &AggregateService{
		api: api,
	}
}

func (s *AggregateService) Aggregate(infoBlock *models.InfoBlock) *models.InfoBlock {
	var wg sync.WaitGroup

	service.SafeGo(&wg, func() {
		if infoBlock.UserID != nil && *infoBlock.UserID != 0 {
			var err error
			infoBlock.User, err = s.api.User.GetByID(*infoBlock.UserID)
			if err != nil {
				logger.Error(err)
			}
		}
	})

	service.SafeGo(&wg, func() {
		if infoBlock.TemplateName != "" {
			tpl, err := s.api.Template.GetByName(infoBlock.TemplateName)
			if err != nil {
				logger.Error(err)
				return
			}
			infoBlock.Template = tpl
		}
	})

	service.SafeGo(&wg, func() {
		infoBlock.Galleries = s.api.Gallery.GetForResourceUUID(infoBlock.UUID.String())
	})

	wg.Wait()

	return infoBlock
}
