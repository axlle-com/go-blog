package service

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	app "github.com/axlle-com/blog/app/service"
	"github.com/axlle-com/blog/pkg/template/http/request"
	"github.com/axlle-com/blog/pkg/template/models"
	templateRepository "github.com/axlle-com/blog/pkg/template/repository"
	userProvider "github.com/axlle-com/blog/pkg/user/provider"
)

type TemplateService struct {
	templateRepo templateRepository.TemplateRepository
	userProvider userProvider.UserProvider
}

func NewTemplateService(
	templateRepo templateRepository.TemplateRepository,
	userProvider userProvider.UserProvider,
) *TemplateService {
	return &TemplateService{
		templateRepo: templateRepo,
		userProvider: userProvider,
	}
}

func (s *TemplateService) GetByID(id uint) (*models.Template, error) {
	return s.templateRepo.GetByID(id)
}

func (s *TemplateService) Aggregate(template *models.Template) *models.Template {
	if template.UserID != nil && *template.UserID != 0 {
		var err error
		template.User, err = s.userProvider.GetByID(*template.UserID)
		if err != nil {
			logger.Error(err)
		}
	}

	return template
}

func (s *TemplateService) GetByIDs(ids []uint) ([]*models.Template, error) {
	return s.templateRepo.GetByIDs(ids)
}

func (s *TemplateService) Create(template *models.Template, user contract.User) (*models.Template, error) {
	if user != nil {
		id := user.GetID()
		template.UserID = &id
	}
	if err := s.templateRepo.Create(template); err != nil {
		return nil, err
	}
	return template, nil
}

func (s *TemplateService) Update(template *models.Template) (*models.Template, error) {
	if err := s.templateRepo.Update(template); err != nil {
		return nil, err
	}

	return template, nil
}

func (s *TemplateService) Delete(template *models.Template) (err error) {
	return s.templateRepo.Delete(template)
}

func (s *TemplateService) SaveFromRequest(form *request.TemplateRequest, found *models.Template, user contract.User) (template *models.Template, err error) {
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
