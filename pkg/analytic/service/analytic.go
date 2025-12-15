package service

import (
	"github.com/axlle-com/blog/app/api"
	app "github.com/axlle-com/blog/app/service/struct"
	"github.com/axlle-com/blog/pkg/analytic/models"
	"github.com/axlle-com/blog/pkg/analytic/repository"
)

type AnalyticService struct {
	analyticRepo repository.AnalyticRepository
	api          *api.Api
}

func NewAnalyticService(
	analyticRepository repository.AnalyticRepository,
	api *api.Api,
) *AnalyticService {
	return &AnalyticService{
		analyticRepo: analyticRepository,
		api:          api,
	}
}

func (s *AnalyticService) GetByID(id uint) (*models.Analytic, error) {
	return s.analyticRepo.GetByID(id)
}

func (s *AnalyticService) Aggregate(analytic *models.Analytic) *models.Analytic {
	return analytic
}

func (s *AnalyticService) Create(analytic *models.Analytic) (*models.Analytic, error) {
	if err := s.analyticRepo.Create(analytic); err != nil {
		return nil, err
	}
	return analytic, nil
}

func (s *AnalyticService) Update(analytic *models.Analytic) (*models.Analytic, error) {
	if err := s.analyticRepo.Update(analytic); err != nil {
		return nil, err
	}

	return analytic, nil
}

func (s *AnalyticService) Delete(analytic *models.Analytic) (err error) {
	return s.analyticRepo.Delete(analytic)
}

func (s *AnalyticService) SaveFromRequest(form *models.AnalyticRequest, found *models.Analytic) (analytic *models.Analytic, err error) {
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
