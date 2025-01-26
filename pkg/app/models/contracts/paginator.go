package contracts

import (
	"html/template"
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
	AddQueryString(s string)
	GetQuery() template.URL
}
