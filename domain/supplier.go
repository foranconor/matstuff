package domain

import "time"

type Supplier struct {
	ID          int
	Name        string
	Website     string
	ContactID   int
	ContactName string
	Created     time.Time
	Modified    time.Time
}

type SupplierWithContact struct {
	Supplier Supplier
	Contact  Contact
}
