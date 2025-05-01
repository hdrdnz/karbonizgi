package model

import (
	"time"
)

type User struct {
	Id          int       `gorm:"primaryKey;autoIncrement"` // Otomatik artan birincil anahtar
	Email       string    `gorm:"size:45;"`                 // Tekil ve boş olamaz
	Firstname   string    `gorm:"size:100;"`                // Maksimum 100 karakter
	Lastname    string    `gorm:"size:100;not null"`
	Password    string    `gorm:"not null"`                // Boş olamaz
	Username    string    `gorm:"unique;size:50;not null"` // Tekil, max 50 karakter
	UserType    string    `gorm:"size:100;not null"`       // ENUM Kullanımı
	CompanyName string    `gorm:"size:100;not null"`       // ENUM Kullanımı
	CreatedAt   time.Time `gorm:"autoCreateTime"`          // Otomatik oluşturulma zamanı
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

type UserToken struct {
	Id        int       `gorm:"primaryKey;autoIncrement"`
	UserId    int       `gorm:"not null;index"`                                 // Kullanıcı ID, index eklenmeli
	Token     string    `gorm:"type:varchar(255);not null"`                     // Token, boş olamaz
	User      User      `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"` // Foreign Key ilişkisi
	CreatedAt time.Time `gorm:"autoCreateTime"`                                 // Otomatik oluşturulma zamanı
	UpdatedAt time.Time `gorm:"autoUpdateTime"`                                 // Otomatik güncellenme zamanı
}

type UserAction struct {
	Id        int    `gorm:"primaryKey;autoIncrement"`
	UserId    int    `gorm:"not null;index"`
	Action    string `gorm:"type:varchar(255);not null"`
	Status    string `gorm:"type:varchar(255);not null"`
	User      User
	CreatedAt time.Time
	UpdatedAt time.Time
}

func Migrate() {
	Db.AutoMigrate(
		&User{},
		&UserToken{},
		&UserAction{},
	)
}
