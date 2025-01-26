package db

import (
	. "github.com/axlle-com/blog/pkg/app/models"
	. "github.com/axlle-com/blog/pkg/app/models/contracts"
)

func NewCache() Cache {
	return NewRedisCache()
}
