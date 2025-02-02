package model

import (
	"time"
)

type User struct {
	Id        uint      `gorm:"primaryKey;autoIncrement"` // Otomatik artan birincil anahtar
	Email     string    `gorm:"unique;not null"`          // Tekil ve boş olamaz
	Firstname string    `gorm:"size:100;not null"`        // Maksimum 100 karakter
	Lastname  string    `gorm:"size:100;not null"`
	Password  string    `gorm:"not null"`                // Boş olamaz
	Username  string    `gorm:"unique;size:50;not null"` // Tekil, max 50 karakter
	UserType  string    `gorm:"size:100;not null"`       // ENUM Kullanımı
	CreatedAt time.Time `gorm:"autoCreateTime"`          // Otomatik oluşturulma zamanı
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type UserToken struct {
	Id        uint      `gorm:"primaryKey;autoIncrement"`
	UserId    uint      `gorm:"not null;index"`                                 // Kullanıcı ID, index eklenmeli
	Token     string    `gorm:"type:varchar(255);not null"`                     // Token, boş olamaz
	User      User      `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"` // Foreign Key ilişkisi
	CreatedAt time.Time `gorm:"autoCreateTime"`                                 // Otomatik oluşturulma zamanı
	UpdatedAt time.Time `gorm:"autoUpdateTime"`                                 // Otomatik güncellenme zamanı
}

func Migrate() {
	Db.AutoMigrate(&User{}, &UserToken{})
}
