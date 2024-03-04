package models

import "gorm.io/gorm"

type Photo struct {
	gorm.Model
	ID       uint   `json:"id" gorm:"primaryKey"`
	Title    string `json:"title" gorm:"not null"`
	Caption  string `json:"caption" gorm:"not null"`
	PhotoUrl string `json:"photo_url" gorm:"not null"`
	UserID   uint   `json:"user_id" gorm:"not null"`
}
