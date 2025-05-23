package model

import "time"

type Admin struct {
	Id        int       `gorm:"primaryKey;autoIncrement"` // Otomatik artan birincil anahtar
	Name      string    `gorm:"size:45;"`
	LastName  string    `gorm:"size:45;"`
	Email     string    `gorm:"size:45;"`       // Tekil ve boş olamaz
	Password  string    `gorm:"not null"`       // Boş olamaz
	CreatedAt time.Time `gorm:"autoCreateTime"` // Otomatik oluşturulma zamanı
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type AdminToken struct {
	Id        int       `gorm:"primaryKey;autoIncrement"`
	AdminId   int       `gorm:"not null;index"`                                  // Kullanıcı ID, index eklenmeli
	Token     string    `gorm:"type:varchar(255);not null"`                      // Token, boş olamaz
	Admin     Admin     `gorm:"foreignKey:AdminId;constraint:OnDelete:CASCADE;"` // Foreign Key ilişkisi
	CreatedAt time.Time `gorm:"autoCreateTime"`                                  // Otomatik oluşturulma zamanı
	UpdatedAt time.Time `gorm:"autoUpdateTime"`                                  // Otomatik güncellenme zamanı
}
