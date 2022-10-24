package models

import "time"

type GormModel struct {
	Id        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at omitempty"`
	UpdatedAt time.Time `json:"updated_at omitempty"`
}
