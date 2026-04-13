package domain

import "time"

type Bracket struct {
	ID           int
	Name         string
	Nickname     string
	SupplierID   *int
	Price        float64
	Published    bool
	Archived     bool
	SupplierName string
	Created      time.Time
}
