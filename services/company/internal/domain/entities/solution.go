package entities

import (
	"time"
)

type Solution struct {
	ID           int64
	AssignmentID int64
	SolutionText string
	Status       string

	AssignedAt  time.Time
	CompletedAt time.Time
}
