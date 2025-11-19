package models

type User struct {
	UserID   string `gorm:"primaryKey"`
	Username string `gorm:"not null"`
	TeamName string
	IsActive bool `gorm:"not null;default:true"`
}
