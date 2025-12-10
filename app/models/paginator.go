package models

import (
	"html/template"
	"net/url"
	"strconv"

	"github.com/axlle-com/blog/app/models/contract"
)

const DefaultPageSize = 20

type paginator struct {
	Total       int          `json:"total"`
	Page        int          `json:"page"`
	PageSize    int          `json:"pageSize"`
	Query       url.Values   `json:"query"`
	QueryString template.URL `json:"queryString"`
	URL         template.URL `json:"url"`
}

func FromQuery(query url.Values) contract.Paginator {
	p := &paginator{
		Query: query,
	}
	p.SetPage()
	p.SetPageSize()
	p.setQueryString()

	return p
}

func (p *paginator) Clone() contract.Paginator {
	newPaginator := &paginator{
		Total:       p.Total,
		Page:        p.Page,
		PageSize:    p.PageSize,
		Query:       cloneValues(p.Query),
		QueryString: p.QueryString,
		URL:         p.URL,
	}

	return newPaginator
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

func (p *paginator) SetURL(url string) {
	p.URL = template.URL(url)
}

func (p *paginator) GetURL() template.URL {
	return p.URL
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
		pageSize = DefaultPageSize
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

	query := "page=" + strconv.Itoa(p.Page) + size

	if qs := p.GetQuery(); qs != "" {
		query += "&" + string(qs)
	}

	return template.URL(query)
}

func (p *paginator) HasPages() bool {
	return p.PageSize > 0 && p.Total > p.PageSize
}

func (p *paginator) PageNumbers() []interface{} {
	if p.PageSize <= 0 {
		return nil
	}

	totalPages := p.Total / p.PageSize
	if p.Total%p.PageSize != 0 {
		totalPages++
	}
	if totalPages <= 0 {
		return nil
	}

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

func (p *paginator) setQueryString() {
	if p.Query == nil {
		return
	}

	query := make(url.Values)
	for key, values := range p.Query {
		if key == "page" || key == "pageSize" || key == "_csrf" {
			continue
		}
		query[key] = values
	}

	temp := query.Encode()
	if temp == "" {
		return
	}

	p.QueryString = template.URL(temp)
}

func (p *paginator) defaultQuery(key, defaultValue string) string {
	values, ok := p.Query[key]
	if ok && len(values) > 0 {
		return values[0]
	}
	return defaultValue
}

func cloneValues(v url.Values) url.Values {
	if v == nil {
		return nil
	}

	m := make(url.Values, len(v))
	for k, vv := range v {
		cp := make([]string, len(vv))
		copy(cp, vv)
		m[k] = cp
	}

	return m
}
