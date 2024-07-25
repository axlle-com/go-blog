package models

import (
	"github.com/gin-gonic/gin"
	"math"
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
}
type paginator struct {
	*gin.Context
	total    int
	page     int
	pageSize int
}

func NewPaginator(c *gin.Context) Paginator {
	p := &paginator{Context: c}
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
	p.page = page
}

func (p *paginator) SetTotal(total int) {
	p.total = total
}

func (p *paginator) GetTotal() int {
	return p.total
}

func (p *paginator) GetPage() int {
	return p.page
}

func (p *paginator) SetPageSize() {
	pageSizeStr := p.DefaultQuery("pageSize", "20")

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 20
	}
	p.pageSize = pageSize
}

func (p *paginator) GetPageSize() int {
	return p.pageSize
}

func (p *paginator) HasPages() bool {
	return p.total > p.pageSize
}

func (p *paginator) PageNumbers() []interface{} {
	totalPages := int(math.Ceil(float64(p.total) / float64(p.pageSize)))
	var pages []interface{}
	if totalPages <= 7 {
		for i := 1; i <= totalPages; i++ {
			pages = append(pages, i)
		}
	} else {
		if p.page <= 4 {
			for i := 1; i <= 5; i++ {
				pages = append(pages, i)
			}
			pages = append(pages, "...")
			pages = append(pages, totalPages)
		} else if p.page >= totalPages-3 {
			pages = append(pages, 1)
			pages = append(pages, "...")
			for i := totalPages - 4; i <= totalPages; i++ {
				pages = append(pages, i)
			}
		} else {
			pages = append(pages, 1)
			pages = append(pages, "...")
			for i := p.page - 1; i <= p.page+1; i++ {
				pages = append(pages, i)
			}
			pages = append(pages, "...")
			pages = append(pages, totalPages)
		}
	}
	return pages
}
