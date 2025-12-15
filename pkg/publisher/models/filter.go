package models

import (
	"github.com/axlle-com/blog/app/errutil"
	"github.com/axlle-com/blog/app/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PublisherFilter struct {
	UUIDs []uuid.UUID
	Query string
	URL   *string

	models.Filter
}

func NewPublisherFilter() *PublisherFilter {
	return &PublisherFilter{
		UUIDs: []uuid.UUID{},
		Query: "",
	}
}

func (f *PublisherFilter) GetUUIDs() []uuid.UUID {
	return f.UUIDs
}

func (f *PublisherFilter) GetQuery() string {
	return f.Query
}

func (f *PublisherFilter) GetURL() *string {
	return f.URL
}

func (f *PublisherFilter) SetUUIDs(uuids []uuid.UUID) {
	f.UUIDs = uuids
}

func (f *PublisherFilter) SetQuery(query string) {
	f.Query = query
}

func (f *PublisherFilter) ValidateQuery(ctx *gin.Context) (*PublisherFilter, *errutil.Errors) {
	err := f.Filter.ValidateQuery(ctx, f)
	return f, err
}
