package models

import (
	"time"
)

type Function struct {
	Name      string
	Call      string
	Trigger   string
	Method    string
	CreatedAt time.Time
}
