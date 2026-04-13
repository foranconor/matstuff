package domain

import "time"

type Contact struct {
	ID       int
	Name     string
	Phone    string
	Email    string
	Created  time.Time
	Modified time.Time
}
