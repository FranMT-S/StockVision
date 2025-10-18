package scopes

import (
	"api/models/filters"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func TimeOrder(order filters.Sort) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order("time " + order.String())
	}
}

// OrderList orders the query by the given order map
// example: OrderList(map[string]filters.Sort{"time": filters.DESC, "company": filters.ASC})
func OrderList(order map[string]filters.Sort) func(db *gorm.DB) *gorm.DB {

	orders := []string{}
	return func(db *gorm.DB) *gorm.DB {
		for col, dir := range order {
			orders = append(orders, fmt.Sprintf("%s %s", col, dir))

		}
		return db.Order(strings.Join(orders, ", "))
	}
}
