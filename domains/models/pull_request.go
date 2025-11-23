package models

import "time"

type PullRequest struct {
	PullRequestID     string     `gorm:"primaryKey" json:"pull_request_id"`
	PullRequestName   string     `gorm:"not null" json:"pull_request_name"`
	AuthorID          string     `gorm:"not null" json:"author_id"`
	Status            string     `gorm:"not null;type:pr_status_enum" json:"status"`
	CreatedAt         *time.Time `json:"-"`
	MergedAt          *time.Time `json:"mergedAt,omitempty"`
	Author            User       `gorm:"foreignKey:AuthorID;references:UserID" json:"-"`
	AssignedReviewers []string   `json:"assigned_reviewers" gorm:"-"`
}

type PullRequestResponse struct {
	PR *PullRequest `json:"pr"`
}

type NewPRReviewer struct {
	PullRequestID string `json:"pull_request_id"`
	OldReviewerID string `json:"old_reviewer_id"`
}

type PRShort struct {
	PullRequestID   string `son:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}
