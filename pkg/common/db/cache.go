package db

import (
	"github.com/axlle-com/blog/pkg/common/models"
	"github.com/axlle-com/blog/pkg/common/models/contracts"
)

func Cache() contracts.Cache {
	return models.Redis()
}
