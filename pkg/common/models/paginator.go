package models

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"math"
	"net/url"
	"strconv"
)

type Paginator interface {
	SetPage()
	GetPage() int
	SetPageSize()
	GetPageSize() int
	SetTotal(int)
	GetTotal() int
	PageNumbers() []interface{}
	HasPages() bool
	SetQuery()
	GetQuery() template.URL
}
type paginator struct {
	*gin.Context `json:"-"`
	Total        int `json:"total"`
	Page         int `json:"page"`
	PageSize     int `json:"pageSize"`
	queryString  template.URL
	Query        url.Values `json:"query"`
}

func NewPaginator(c *gin.Context) Paginator {
	p := &paginator{Context: c}
	p.SetQuery()
	p.SetPage()
	p.SetPageSize()
	return p
}

func (p *paginator) SetPage() {
	pageStr := p.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	p.Page = page
}

func (p *paginator) SetQuery() {
	p.Query = p.Request.URL.Query()
	p.queryString = template.URL(p.Query.Encode())
}

func (p *paginator) GetQuery() template.URL {
	return p.queryString
}

func (p *paginator) SetTotal(total int) {
	p.Total = total
}

func (p *paginator) GetTotal() int {
	return p.Total
}

func (p *paginator) GetPage() int {
	return p.Page
}

func (p *paginator) SetPageSize() {
	pageSizeStr := p.DefaultQuery("pageSize", "20")

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 20
	}
	p.PageSize = pageSize
}

func (p *paginator) GetPageSize() int {
	return p.PageSize
}

func (p *paginator) HasPages() bool {
	return p.Total > p.PageSize
}

func (p *paginator) PageNumbers() []interface{} {
	totalPages := int(math.Ceil(float64(p.Total) / float64(p.PageSize)))
	var pages []interface{}
	if totalPages <= 7 {
		for i := 1; i <= totalPages; i++ {
			pages = append(pages, i)
		}
	} else {
		if p.Page <= 4 {
			for i := 1; i <= 5; i++ {
				pages = append(pages, i)
			}
			pages = append(pages, "...")
			pages = append(pages, totalPages)
		} else if p.Page >= totalPages-3 {
			pages = append(pages, 1)
			pages = append(pages, "...")
			for i := totalPages - 4; i <= totalPages; i++ {
				pages = append(pages, i)
			}
		} else {
			pages = append(pages, 1)
			pages = append(pages, "...")
			for i := p.Page - 1; i <= p.Page+1; i++ {
				pages = append(pages, i)
			}
			pages = append(pages, "...")
			pages = append(pages, totalPages)
		}
	}
	return pages
}
