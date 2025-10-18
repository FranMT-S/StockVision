package filters

import (
	"api/config"
)

type Filters struct {
	Query    SanatizerString
	Page     int
	PageSize int
	Sort     Sort
}

func (f *Filters) Normalize() {
	if f.Page < 1 {
		f.Page = 1
	}

	if f.PageSize < 1 {
		f.PageSize = config.DefaultConstants().DefaultPageSize
	}

	if f.PageSize > config.DefaultConstants().MaxPageSize {
		f.PageSize = config.DefaultConstants().MaxPageSize
	}

	if !f.Sort.IsValid() {
		f.Sort = ASC
	}

	if f.Query == "" {
		f.Query = ""
	}
}

func (f Filters) Offset() int {
	return (f.Page - 1) * f.PageSize
}
