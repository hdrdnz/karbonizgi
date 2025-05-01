package model

import "time"

type Admin struct {
	Id        int       `gorm:"primaryKey;autoIncrement"` // Otomatik artan birincil anahtar
	Email     string    `gorm:"size:45;"`                 // Tekil ve boş olamaz
	Password  string    `gorm:"not null"`                 // Boş olamaz
	CreatedAt time.Time `gorm:"autoCreateTime"`           // Otomatik oluşturulma zamanı
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
