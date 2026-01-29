package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/settings/models"
	"github.com/axlle-com/blog/pkg/settings/repository"
)

type Service struct {
	repo  repository.Repository
	cache contract.Cache
	ttl   time.Duration
}

func NewService(cache contract.Cache, repo repository.Repository, ttl time.Duration) *Service {
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}

	return &Service{
		repo:  repo,
		cache: cache,
		ttl:   ttl,
	}
}

func (s *Service) cacheKey(ns, key, scope string) string {
	return "settings|" + ns + "|" + key + "|" + scope
}

func (s *Service) Invalidate(ns, key, scope string) {
	if s.cache == nil {
		return
	}

	s.cache.DeleteCache(s.cacheKey(ns, key, scope))
}

func (s *Service) ResetCache() {
	if s.cache == nil {
		return
	}

	s.cache.DeleteByPrefix("settings|")
}

func (s *Service) GetSetting(ns, key, scope string) (models.Setting, bool) {
	return s.get(ns, key, scope)
}

func (s *Service) get(ns, key, scope string) (models.Setting, bool) {
	ck := s.cacheKey(ns, key, scope)

	if s.cache != nil {
		if raw, ok := s.cache.GetCache(ck); ok && raw != "" {
			var st models.Setting
			if err := json.Unmarshal([]byte(raw), &st); err == nil {
				return st, true
			}
			s.cache.DeleteCache(ck)
		}
	}

	st, err := s.repo.Get(ns, key, scope)
	if err != nil || st == nil {
		return models.Setting{}, false
	}

	if s.cache != nil {
		if b, err := json.Marshal(st); err == nil {
			s.cache.AddCacheTTL(ck, string(b), s.ttl)
		}
	}

	return *st, true
}

func (s *Service) GetString(ns, key, scope string) (string, bool) {
	st, ok := s.get(ns, key, scope)
	if !ok {
		return "", false
	}

	var value string
	if err := json.Unmarshal(st.Value, &value); err != nil {
		return "", false
	}

	return value, true
}

func (s *Service) GetBool(ns, key, scope string) (bool, bool) {
	st, ok := s.get(ns, key, scope)
	if !ok {
		return false, false
	}

	var value bool
	if err := json.Unmarshal(st.Value, &value); err != nil {
		return false, false
	}

	return value, true
}

func (s *Service) GetJSON(ns, key, scope string, out any) bool {
	st, ok := s.get(ns, key, scope)
	if !ok {
		return false
	}

	return json.Unmarshal(st.Value, out) == nil
}

func (s *Service) SaveString(ns, key, scope string, value string) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	st := &models.Setting{
		Namespace: ns,
		Key:       key,
		Scope:     scope,
		Type:      models.SettingTypeString,
		Value:     b,
	}

	if err := s.repo.Upsert(st); err != nil {
		return err
	}

	s.Invalidate(ns, key, scope)
	return nil
}

func (s *Service) SaveBool(ns, key, scope string, value bool) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	st := &models.Setting{
		Namespace: ns,
		Key:       key,
		Scope:     scope,
		Type:      models.SettingTypeBool,
		Value:     b,
	}

	if err := s.repo.Upsert(st); err != nil {
		return err
	}

	s.Invalidate(ns, key, scope)
	return nil
}

func (s *Service) SaveJSON(ns, key, scope string, value any) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	st := &models.Setting{
		Namespace: ns,
		Key:       key,
		Scope:     scope,
		Type:      models.SettingTypeJSON,
		Value:     b,
	}

	if err := s.repo.Upsert(st); err != nil {
		return err
	}

	s.Invalidate(ns, key, scope)

	return nil
}

func (s *Service) SaveTyped(ns, key, scope, typ string, value any) error {
	typ = strings.TrimSpace(strings.ToLower(typ))
	if typ == "" {
		return fmt.Errorf("setting type is empty")
	}

	var (
		body []byte
		err  error
	)

	switch typ {
	case models.SettingTypeString:
		switch v := value.(type) {
		case string:
			body, err = json.Marshal(v)
		case json.RawMessage:
			var out string
			if e := json.Unmarshal(v, &out); e != nil {
				return fmt.Errorf("setting value is not a string")
			}
			body, err = json.Marshal(out)
		default:
			return fmt.Errorf("setting value is not a string")
		}
	case models.SettingTypeBool:
		switch v := value.(type) {
		case bool:
			body, err = json.Marshal(v)
		case json.RawMessage:
			var out bool
			if e := json.Unmarshal(v, &out); e != nil {
				return fmt.Errorf("setting value is not a bool")
			}
			body, err = json.Marshal(out)
		default:
			return fmt.Errorf("setting value is not a bool")
		}
	case models.SettingTypeJSON:
		if raw, ok := value.(json.RawMessage); ok {
			if len(raw) == 0 {
				return fmt.Errorf("setting value is empty json")
			}
			body = raw
		} else {
			body, err = json.Marshal(value)
		}
	default:
		return fmt.Errorf("unsupported setting type: %s", typ)
	}

	if err != nil {
		return err
	}

	st := &models.Setting{
		Namespace: ns,
		Key:       key,
		Scope:     scope,
		Type:      typ,
		Value:     body,
	}

	if err := s.repo.Upsert(st); err != nil {
		return err
	}

	s.Invalidate(ns, key, scope)

	return nil
}
