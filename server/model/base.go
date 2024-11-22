package model

import (
	"time"
)

type Base struct {
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (base Base) GetKey() string {
	return base.Code
}
