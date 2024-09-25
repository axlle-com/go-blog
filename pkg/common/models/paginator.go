package models

import (
	"github.com/axlle-com/blog/pkg/common/models/contracts"
	"html/template"
	"math"
	"net/url"
	"strconv"
)

const DefaultPageSize = 20

type paginator struct {
	Total       int          `json:"total"`
	Page        int          `json:"page"`
	PageSize    int          `json:"pageSize"`
	Query       url.Values   `json:"query"`
	QueryString template.URL `json:"queryString"`
}

func Paginator(q url.Values) contracts.Paginator {
	p := &paginator{
		Query: q,
	}
	p.SetPage()
	p.SetPageSize()
	return p
}

func (p *paginator) defaultQuery(key, defaultValue string) string {
	values, ok := p.Query[key]
	if ok {
		return values[0]
	}
	return defaultValue
}

func (p *paginator) SetPage() {
	pageStr := p.defaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	p.Page = page
}

func (p *paginator) AddQueryString(s string) {
	if s == "" {
		return
	}
	p.QueryString += template.URL("&" + s)
}

func (p *paginator) GetQuery() template.URL {
	return p.QueryString
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
	pageSizeStr := p.defaultQuery("pageSize", strconv.Itoa(DefaultPageSize))

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 20
	}
	p.PageSize = pageSize
}

func (p *paginator) GetPageSize() int {
	return p.PageSize
}

func (p *paginator) PrintFullQuery() template.URL {
	size := ""
	if p.PageSize != DefaultPageSize {
		size = "&pageSize=" + strconv.Itoa(p.PageSize)
	}
	pageQuery := template.URL("page=" + strconv.Itoa(p.Page) + size)
	return pageQuery + p.GetQuery()
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
