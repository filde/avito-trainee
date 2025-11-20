package models

type User struct {
	UserID   string `gorm:"primaryKey" json:"user_id"`
	Username string `gorm:"not null" json:"username"`
	TeamName string `gorm:"not null" json:"-"`
	IsActive bool   `gorm:"not null" json:"is_active"`
	Team     `gorm:"foreignKey:TeamName;references:TeamName" json:"-"`
}
