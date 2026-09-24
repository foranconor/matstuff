package domain

import "time"

type Material struct {
	ID          int
	Name        string
	Nickname    string
	Treatment   string
	Blurb       string
	Units       string
	Length      float64
	Width       float64
	Thickness   float64
	MaxSpan     float64
	Density     float64
	MaxOverhang float64
	Radius      float64
	Color       string
	Published   bool
	Archived    bool
	Created     time.Time
	Modified    time.Time
}

type Uses struct {
	ID          int
	MaterialID  int
	Exterior    bool
	Interior    bool
	EarlyAccess bool
	Stringers   bool
	Risers      bool
	Treads      bool
	Timber      bool
	Panel       bool
	Published   bool
	Archived    bool
	Handrail    bool
	Created     time.Time
	Modified    time.Time
}

type MaterialWithUses struct {
	Material Material
	Uses     Uses
}

type MaterialSupplier struct {
	ID           int
	Priority     int
	SupplierID   int
	MaterialID   int
	Price        float64
	LeadTime     int
	SupplierName string
	MaterialName string
	Created      time.Time
	Modified     time.Time
}
