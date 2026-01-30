package service

import (
	"github.com/axlle-com/blog/app/api"
	app "github.com/axlle-com/blog/app/service/struct"
	"github.com/axlle-com/blog/pkg/analytic/models"
	"github.com/axlle-com/blog/pkg/analytic/repository"
)

type Service struct {
	api  *api.Api
	repo repository.AnalyticRepository
}

func NewService(
	api *api.Api,
	analyticRepository repository.AnalyticRepository,
) *Service {
	return &Service{
		api:  api,
		repo: analyticRepository,
	}
}

func (s *Service) GetByID(id uint) (*models.Analytic, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Aggregate(analytic *models.Analytic) *models.Analytic {
	return analytic
}

func (s *Service) Create(analytic *models.Analytic) (*models.Analytic, error) {
	if err := s.repo.Create(analytic); err != nil {
		return nil, err
	}
	return analytic, nil
}

func (s *Service) Update(analytic *models.Analytic) (*models.Analytic, error) {
	if err := s.repo.Update(analytic); err != nil {
		return nil, err
	}

	return analytic, nil
}

func (s *Service) Delete(analytic *models.Analytic) (err error) {
	return s.repo.Delete(analytic)
}

func (s *Service) SaveFromRequest(form *models.AnalyticRequest, found *models.Analytic) (analytic *models.Analytic, err error) {
	analyticForm := app.LoadStruct(&models.Analytic{}, form).(*models.Analytic)

	if found == nil {
		analytic, err = s.Create(analyticForm)
	} else {
		analyticForm.ID = found.ID
		analytic, err = s.Update(analyticForm)
	}

	if err != nil {
		return
	}

	return
}
