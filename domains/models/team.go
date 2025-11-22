package models

type Team struct {
	TeamName string `gorm:"primaryKey" json:"team_name"`
	Members  []User `gorm:"-" json:"members"`
}

type TeamResponse struct {
	Team *Team `json:"team"`
}
