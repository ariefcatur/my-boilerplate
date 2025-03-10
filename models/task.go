package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Title       string `json:"title" gorm:"not null"`
	Description string `json:"description"`
	Status      string `json:"status" gorm:"default:'pending'"`  // pending, completed
	Priority    string `json:"priority" gorm:"default:'medium'"` // low, medium, high
	UserID      uint   `json:"user_id"`
	User        User   `json:"-" gorm:"foreignkey:UserID"`
}
