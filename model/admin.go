package model

import "time"

type Admin struct {
	Id        int       `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"size:45;"`
	LastName  string    `gorm:"size:45;"`
	Email     string    `gorm:"size:45;"`
	Password  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type AdminToken struct {
	Id        int       `gorm:"primaryKey;autoIncrement"`
	AdminId   int       `gorm:"not null;index"`
	Token     string    `gorm:"type:varchar(255);not null"`
	Admin     Admin     `gorm:"foreignKey:AdminId;constraint:OnDelete:CASCADE;"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
