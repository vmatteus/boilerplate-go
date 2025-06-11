package domain

import (
	"time"
)

// User represents a user entity
type User struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Name      string     `json:"name" gorm:"size:255;not null"`
	Email     string     `json:"email" gorm:"size:255;uniqueIndex;not null"`
	CreatedAt *time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName returns the table name for the User entity
func (User) TableName() string {
	return "users"
}
