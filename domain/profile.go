package domain

import "time"

type Profile struct {
	ID       int
	Name     string
	Created  time.Time
	Modified time.Time
}

type MaterialProfile struct {
	ID           int
	MaterialID   int
	ProfileID    int
	Price        float64
	MaterialName string
	ProfileName  string
	Created      time.Time
	Modified     time.Time
}
