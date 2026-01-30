package service

import (
	"github.com/axlle-com/blog/app/api"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	app "github.com/axlle-com/blog/app/service/struct"
	"github.com/axlle-com/blog/pkg/template/http/request"
	"github.com/axlle-com/blog/pkg/template/models"
	templateRepository "github.com/axlle-com/blog/pkg/template/repository"
)

type Service struct {
	api          *api.Api
	templateRepo templateRepository.TemplateRepository
}

func NewService(
	api *api.Api,
	templateRepo templateRepository.TemplateRepository,
) *Service {
	return &Service{
		api:          api,
		templateRepo: templateRepo,
	}
}

func (s *Service) GetByID(id uint) (*models.Template, error) {
	return s.templateRepo.GetByID(id)
}

func (s *Service) Aggregate(template *models.Template) *models.Template {
	if template.UserID != nil && *template.UserID != 0 {
		var err error
		template.User, err = s.api.User.GetByID(*template.UserID)
		if err != nil {
			logger.Error(err)
		}
	}

	return template
}

func (s *Service) GetByIDs(ids []uint) ([]*models.Template, error) {
	return s.templateRepo.GetByIDs(ids)
}

func (s *Service) Create(template *models.Template, user contract.User) (*models.Template, error) {
	if user != nil {
		id := user.GetID()
		template.UserID = &id
	}
	if err := s.templateRepo.Create(template); err != nil {
		return nil, err
	}
	return template, nil
}

func (s *Service) Update(template *models.Template) (*models.Template, error) {
	if err := s.templateRepo.Update(template); err != nil {
		return nil, err
	}

	return template, nil
}

func (s *Service) Delete(template *models.Template) (err error) {
	return s.templateRepo.Delete(template) // @todo обновить связи
}

func (s *Service) SaveFromRequest(form *request.TemplateRequest, found *models.Template, user contract.User) (template *models.Template, err error) {
	templateForm := app.LoadStruct(&models.Template{}, form).(*models.Template)

	if found == nil {
		template, err = s.Create(templateForm, user)
	} else {
		templateForm.ID = found.ID
		template, err = s.Update(templateForm)
	}

	if err != nil {
		return
	}

	return
}
