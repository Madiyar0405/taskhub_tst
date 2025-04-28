package entities

import "time"

type DeveloperProfile struct {
	DeveloperID    int64
	GithubID       int64
	Username       string
	Bio            *string
	GithubUsername *string
	AvatarURL      *string
	CVURL          *string
	Email          string
	CreatedAt      time.Time

	IsGithubLinked bool
}
