package models

type Team struct {
	TeamName string `gorm:"primaryKey" json:"team_name"`
	Members  []User `gorm:"foreignKey:TeamName;references:TeamName" json:"members"`
}

type TeamResponse struct {
	Team *Team `json:"team"`
}
