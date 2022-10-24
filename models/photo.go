package models

import (
	"time"
)

type Photo struct {
	GormModel
	Title    string    `gorm:"not null" json:"title" form:"title"`
	Caption  string    `json:"caption,omitempty" form:"caption"`
	PhotoUrl string    `gorm:"not null" json:"photo_url" form:"photo_url"`
	Comments []Comment `gorm:"constraint:OnUpdate:CASCADE, OnDelete:SET NULL;" json:"comments"`
	User     *User     `json:"user"`
	UserId   uint
}

type UserPhoto struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

type GetPhotos struct {
	Id       uint      `json:"id"`
	Title    string    `json:"title"`
	Caption  string    `json:"caption"`
	PhotoUrl string    `json:"photo_url"`
	UserId   uint      `json:"user_id"`
	CreateAt time.Time `json:"created_at"`
	UpdateAt time.Time `json:"updated_at"`
	User     UserPhoto `json:"user"`
}

// func (p *Photo) BeforeCreate(tx *gorm.DB) (err error) {
// 	_, errCreate := govalidator.ValidateStruct(p)

// 	if errCreate != nil {
// 		err = errCreate
// 		return
// 	}

// 	err = nil
// 	return
// }

// func (p *Photo) BeforeUpdate(tx *gorm.DB) (err error) {
// 	_, errCreate := govalidator.ValidateStruct(p)

// 	if errCreate != nil {
// 		err = errCreate
// 		return
// 	}

// 	err = nil
// 	return
// }
