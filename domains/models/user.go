package models

type User struct {
	UserID   string `gorm:"primaryKey" json:"user_id"`
	Username string `gorm:"not null" json:"username"`
	TeamName string `gorm:"not null" json:"-"`
	IsActive bool   `gorm:"not null" json:"is_active"`
	Team     `gorm:"foreignKey:TeamName;references:TeamName" json:"-"`
}

type UserResponse struct {
	User *UserFull `json:"user"`
}

type UserFull struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive string `json:"is_active"`
}

type UsersPR struct {
	UserID       string     `json:"user_id"`
	PullRequests []*PRShort `json:"pull_requests" gorm:"-"`
}
