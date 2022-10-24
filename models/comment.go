package models

import (
	"time"
)

type Comment struct {
	GormModel
	UserId  uint
	PhotoId uint   `json:"photo_id" form:"photo_id"`
	Message string `gorm:"not null" json:"message" form:"message"`
	User    *User  `json:"user"`
	Photo   *Photo `json:"photo"`
}

type UserComment struct {
	Id       uint   `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type PhotoComment struct {
	Id       uint   `json:"id"`
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoUrl string `json:"photo_url"`
	UserId   uint   `json:"user_id"`
}

type GetComment struct {
	Id        uint         `json:"id"`
	Message   string       `json:"message"`
	PhotoId   uint         `json:"photo_id"`
	UserId    uint         `json:"user_id"`
	UpdatedAt time.Time    `json:"updated_at"`
	CreatedAt time.Time    `json:"created_at"`
	User      UserComment  `json:"user"`
	Photo     PhotoComment `json:"photo"`
}

// func (c *Comment) BeforeCreate(tx *gorm.DB) (err error) {
// 	_, errCreate := govalidator.ValidateStruct(c)

// 	if errCreate != nil {
// 		err = errCreate
// 		return
// 	}

// 	err = nil
// 	return
// }

// func (c *Comment) BeforeUpdate(tx *gorm.DB) (err error) {
// 	_, errCreate := govalidator.ValidateStruct(c)

// 	if errCreate != nil {
// 		err = errCreate
// 		return
// 	}

// 	err = nil
// 	return
// }
