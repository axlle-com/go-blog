package db

import (
	"encoding/json"
	"strings"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/pkg/settings/models"
	"github.com/axlle-com/blog/pkg/settings/service"
)

type settingsSeeder struct {
	config          contract.Config
	disk            contract.DiskService
	seedService     contract.SeedService
	settingsService *service.Service
}

type SettingsSeedData struct {
	Namespace string          `json:"namespace"`
	Scope     string          `json:"scope"`
	Key       string          `json:"key"`
	Type      string          `json:"type"`
	Value     json.RawMessage `json:"value"`
}

func NewSettingsSeeder(
	cfg contract.Config,
	disk contract.DiskService,
	seedService contract.SeedService,
	settingsService *service.Service,
) contract.Seeder {
	return &settingsSeeder{
		config:          cfg,
		disk:            disk,
		seedService:     seedService,
		settingsService: settingsService,
	}
}

func (s *settingsSeeder) Seed() error {
	return s.seedFromJSON((models.Setting{}).GetTable())
}

func (s *settingsSeeder) SeedTest(n int) error {
	return nil
}

func (s *settingsSeeder) seedFromJSON(moduleName string) error {
	files, err := s.seedService.GetFiles(s.config.Layout(), moduleName)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		logger.Infof("[settings][seeder][seedFromJSON] seed files not found for module: %s, skipping", moduleName)
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

		var list []SettingsSeedData
		if err := json.Unmarshal(data, &list); err != nil {
			return err
		}

		for _, item := range list {
			ns := strings.TrimSpace(item.Namespace)
			if ns == "" {
				ns = s.config.Layout()
			}

			scope := strings.TrimSpace(item.Scope)
			if scope == "" {
				scope = "global"
			}

			key := strings.TrimSpace(item.Key)
			if key == "" {
				logger.Infof("[settings][seeder][seedFromJSON] empty key for ns=%s scope=%s, skipping", ns, scope)
				continue
			}

			typ := strings.TrimSpace(strings.ToLower(item.Type))
			if typ == "" {
				logger.Infof("[settings][seeder][seedFromJSON] empty type for ns=%s scope=%s key=%s, skipping", ns, scope, key)
				continue
			}

			if err := s.settingsService.SaveTyped(ns, key, scope, typ, item.Value); err != nil {
				logger.Errorf("[settings][seeder][seedFromJSON] error saving setting ns=%s scope=%s key=%s: %v", ns, scope, key, err)
				continue
			}
		}

		s.settingsService.ResetCache()

		if err := s.seedService.MarkApplied(name); err != nil {
			return err
		}

		logger.Infof("[settings][seeder][seedFromJSON] seeded %d items from JSON (%s)", len(list), name)
	}

	return nil
}
