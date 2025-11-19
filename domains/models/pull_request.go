package models

import "time"

type PullRequest struct {
	PullRequestID     string `gorm:"primaryKey"`
	PullRequestName   string `gorm:"not null"`
	AuthorID          string
	Status            string   `gorm:"not null"`
	AssignedReviewers []string `gorm:"not null"`
	CreatedAt         *time.Time
	MergedAt          *time.Time
}
