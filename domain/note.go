package domain

import "time"

type Note struct {
	ID       int
	Content  string
	Created  time.Time
	Modified time.Time
}
