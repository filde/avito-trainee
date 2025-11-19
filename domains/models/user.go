package models

type User struct {
	UserID   string `gorm:"primaryKey"`
	Username string `gorm:"not null"`
	TeamName string `gorm:"not null"`
	IsActive bool   `gorm:"not null;default:true"`
	Team     `gorm:"foreignKey:TeamName;references:TeamName" json:"-"`
}
