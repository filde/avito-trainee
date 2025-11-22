package models

import "time"

type PullRequest struct {
	PullRequestID     string     `gorm:"primaryKey" json:"pull_request_id"`
	PullRequestName   string     `gorm:"not null" json:"pull_request_name"`
	AuthorID          string     `gorm:"not null" json:"author_id"`
	Status            string     `gorm:"not null;type:status_enum" json:"status"`
	CreatedAt         *time.Time `json:"-"`
	MergedAt          *time.Time `json:"mergedAt;omitempty"`
	User              `gorm:"foreignKey:AuthorID;references:UserID" json:"-"`
	Reviewers         []*User  `gorm:"many2many:pr_reviewers" json:"-"`
	AssignedReviewers []string `json:"assigned_reviewers" gorm:"-"`
}

type PullRequestResponse struct {
	PR *PullRequest `json:"pr"`
}

type NewPRReviewer struct {
	PullRequestID string `json:"pull_request_id"`
	OldReviewerID string `json:"old_reviewer_id"`
}
