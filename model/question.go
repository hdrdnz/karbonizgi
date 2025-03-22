package model

import "time"

type QuestionTypes struct {
	Id          int       `gorm:"primaryKey;autoIncrement"`
	UserType    string    `gorm:"size:45;"`
	QuestionKey string    `gorm:"size:45;"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

type QuestionSubhead struct {
	Id              int    `gorm:"primaryKey;autoIncrement"`
	QuestionTypesId int    `gorm:"not null;index"`
	QuestionKey     string `gorm:"size:100;"`
	QuestionTypes   QuestionTypes
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}

type UserDetailScore struct {
	Id              int     `gorm:"primaryKey;autoIncrement"`
	UserId          int     `gorm:"not null;index"`
	QuestionTypesId int     `gorm:"not null;index"`
	TotalScore      float64 `gorm:"not null;index"`
	QuestionTypes   QuestionTypes
	User            User
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}
type UserScore struct {
	Id        int     `gorm:"primaryKey;autoIncrement"`
	UserId    int     `gorm:"not null;index"`
	Score     float64 `gorm:"not null;index"`
	User      User
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type SubScore struct {
	Id                int `gorm:"primaryKey;autoIncrement"`
	UserScoreId       int `gorm:"not null;index"`
	QuestionSubheadId int `gorm:"not null;index"`
	Score             int `gorm:"not null;index"`
	QuestionSubhead   QuestionSubhead
	UserScore         UserDetailScore
	CreatedAt         time.Time `gorm:"autoCreateTime"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime"`
}
