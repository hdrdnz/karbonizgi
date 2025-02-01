package model

import "time"

type User struct {
	Email      string    `gorm:"email"`
	Password   string    `gorm:"password"`
	Username   string    `gorm:"user_name"`
	UserType   string    `gorm:"user_type"`
	Created_at time.Time `gorm:"user_type"`
	Updated_at time.Time `gorm:"updated_at"`
}
