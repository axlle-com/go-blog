package service

import (
	"encoding/json"
	"time"

	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/settings/models"
	"github.com/axlle-com/blog/pkg/settings/repository"
	"golang.org/x/sync/singleflight"
)

type CompanyInfoService struct {
	repo  repository.Repository
	cache contract.Cache
	ttl   time.Duration
	sf    singleflight.Group

	settings *Service
}

func NewCompanyInfoService(
	cache contract.Cache,
	repo repository.Repository,
	settings *Service,
	ttl time.Duration,
) *CompanyInfoService {
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}

	return &CompanyInfoService{
		repo:     repo,
		cache:    cache,
		ttl:      ttl,
		settings: settings,
	}
}

func (s *CompanyInfoService) GetCompanyInfo(ns, scope string) (contract.CompanyInfo, bool) {
	ck := s.cacheKeyCompanyInfo(ns, scope)

	if s.cache != nil {
		if raw, ok := s.cache.GetCache(ck); ok && raw != "" {
			var out models.CompanyInfo
			if err := json.Unmarshal([]byte(raw), &out); err == nil {
				return &out, true
			}

			s.cache.DeleteCache(ck)
		}
	}

	value, err, _ := s.sf.Do(ck, func() (any, error) {
		if s.cache != nil {
			if raw, ok := s.cache.GetCache(ck); ok && raw != "" {
				var out models.CompanyInfo
				if e := json.Unmarshal([]byte(raw), &out); e == nil {
					return &out, nil
				}

				s.cache.DeleteCache(ck)
			}
		}

		out, ok := s.getCompanyInfoFromDB(ns, scope)
		if !ok {
			return nil, nil
		}

		if s.cache != nil {
			if b, e := json.Marshal(out); e == nil {
				s.cache.AddCacheTTL(ck, string(b), s.ttl)
			}
		}

		return out, nil
	})

	if err != nil || value == nil {
		return &models.CompanyInfo{}, false
	}

	out, ok := value.(models.CompanyInfo)
	if !ok {
		return &models.CompanyInfo{}, false
	}

	return &out, true
}

func (s *CompanyInfoService) SaveCompanyInfo(ns, scope string, companyInfo models.CompanyInfo) error {
	current, ok := s.getCompanyInfoFromDB(ns, scope)

	if ok {
		if companyInfoEquals(current, companyInfo) {
			return nil
		}
	}

	if err := s.settings.SaveString(ns, models.CompanyEmailKey, scope, companyInfo.Email); err != nil {
		return err
	}
	if err := s.settings.SaveString(ns, models.CompanyNameKey, scope, companyInfo.Name); err != nil {
		return err
	}
	if err := s.settings.SaveString(ns, models.CompanyPhoneKey, scope, companyInfo.Phone); err != nil {
		return err
	}
	if err := s.settings.SaveString(ns, models.CompanyAddressKey, scope, companyInfo.Address); err != nil {
		return err
	}
	if err := s.settings.SaveString(ns, models.PolicyKey, scope, companyInfo.Policy); err != nil {
		return err
	}

	s.InvalidateCompanyInfo(ns, scope)

	return nil
}

func (s *CompanyInfoService) InvalidateCompanyInfo(ns, scope string) {
	if s.cache == nil {
		return
	}

	s.cache.DeleteCache(s.cacheKeyCompanyInfo(ns, scope))
}

func (s *CompanyInfoService) cacheKeyCompanyInfo(ns, scope string) string {
	return "settings|company_info|" + ns + "|" + scope
}

func (s *CompanyInfoService) getCompanyInfoFromDB(ns, scope string) (models.CompanyInfo, bool) {
	keys := []string{
		models.CompanyEmailKey,
		models.CompanyNameKey,
		models.CompanyPhoneKey,
		models.CompanyAddressKey,
		models.PolicyKey,
	}

	rows, err := s.repo.GetMany(ns, scope, keys)
	if err != nil || len(rows) == 0 {
		return models.CompanyInfo{}, false
	}

	var out models.CompanyInfo
	okAny := false

	for _, st := range rows {
		var value string
		if err := json.Unmarshal(st.Value, &value); err != nil {
			continue
		}

		switch st.Key {
		case models.CompanyEmailKey:
			out.Email = value
			okAny = true
		case models.CompanyNameKey:
			out.Name = value
			okAny = true
		case models.CompanyPhoneKey:
			out.Phone = value
			okAny = true
		case models.PolicyKey:
			out.Policy = value
			okAny = true
		case models.CompanyAddressKey:
			out.Address = value
			okAny = true
		}
	}

	if !okAny {
		return models.CompanyInfo{}, false
	}

	return out, true
}

func companyInfoEquals(a, b models.CompanyInfo) bool {
	return a.Email == b.Email &&
		a.Name == b.Name &&
		a.Phone == b.Phone &&
		a.Policy == b.Policy &&
		a.Address == b.Address
}
