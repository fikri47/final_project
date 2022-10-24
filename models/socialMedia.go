package models

import (
	"time"
)

type SocialMedia struct {
	GormModel
	Name           string `gorm:"not null" json:"name" form:"name"`
	SocialMediaUrl string `gorm:"not null" json:"social_media_url" form:"social_media_url"`
	User           *User  `json:"user"`
	UserId         uint   `json:"user_id"`
}

type UserSocialMedia struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
}

type GetSocialMedias struct {
	Id             uint            `json:"id"`
	Name           string          `json:"name"`
	SocialMediaUrl string          `json:"social_media_url"`
	UserId         uint            `json:"user_id"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdateAt       time.Time       `json:"updated_at"`
	User           UserSocialMedia `json:"user"`
}

// func (sm *SocialMedia) BeforeCreate(tx *gorm.DB) (err error) {
// 	_, errCreate := govalidator.ValidateStruct(sm)

// 	if errCreate != nil {
// 		err = errCreate
// 		return
// 	}

// 	err = nil
// 	return
// }

// func (sm *SocialMedia) BeforeUpdate(tx *gorm.DB) (err error) {
// 	_, errCreate := govalidator.ValidateStruct(sm)

// 	if errCreate != nil {
// 		err = errCreate
// 		return
// 	}

// 	err = nil
// 	return
// }
