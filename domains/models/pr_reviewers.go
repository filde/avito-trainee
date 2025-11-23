package models

type PullRequestReviewers struct {
	PullReqID  string
	ReviewerID string
	PR         PullRequest `gorm:"foreignKey:PullReqID;references:PullRequestID" json:"-"`
	Reviewer   User        `gorm:"foreignKey:ReviewerID;references:UserID" json:"-"`
}
