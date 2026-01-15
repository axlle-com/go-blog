package migrate

import (
	"errors"
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"gorm.io/gorm/clause"
)

type SeedService struct {
	db   contract.DB
	disk contract.DiskService
}

func NewSeedService(db contract.DB, disk contract.DiskService) contract.SeedService {
	inst := &SeedService{
		db:   db,
		disk: disk,
	}

	err := inst.db.PostgreSQL().AutoMigrate(
		Seed{},
	)

	if err != nil {
		logger.Errorf("[app][service][migrate][NewSeedService] Error:%+v", err)
	}

	return inst
}

func (s *SeedService) GetFiles(layout, moduleName string) (map[string]string, error) {
	seedDir := filepath.ToSlash(filepath.Join("services", "db", layout, "seed", moduleName))

	entries, err := s.disk.ReadDir(seedDir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			logger.Infof("seed dir not found: %s, skip", seedDir)
			return map[string]string{}, nil
		}
		return nil, err
	}

	out := make(map[string]string, 0)

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		name := e.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".json") {
			continue
		}

		// полный путь для DiskService
		full := seedDir + "/" + name
		out[name] = full
	}

	return out, nil
}

func (s *SeedService) MarkApplied(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("seed name is empty")
	}

	now := time.Now()
	rec := Seed{
		Name:      name,
		CreatedAt: &now,
	}

	return s.db.PostgreSQL().
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoNothing: true,
		}).
		Create(&rec).Error
}

func (s *SeedService) IsApplied(name string) (bool, error) {
	if strings.TrimSpace(name) == "" {
		return false, errors.New("seed name is empty")
	}

	var count int64
	err := s.db.PostgreSQL().
		Model(&Seed{}).
		Where("name = ?", name).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
