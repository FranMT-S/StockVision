package scopes

import (
	"api/models/filters"

	"gorm.io/gorm"
)

func SortCompany(sort filters.Sort) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order("company " + sort.String())
	}
}
