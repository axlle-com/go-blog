package db

import (
	"github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contracts"
)

func NewCache() contracts.Cache {
	return models.NewRedisCache()
}
