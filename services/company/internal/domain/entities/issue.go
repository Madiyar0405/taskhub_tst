package entities

import "database/sql"

type Issue struct {
	ID    int64
	Title string
	Body  string

	AssignmentStatus sql.NullString
	SolutionText     sql.NullString
	SolutionStatus   sql.NullString
}
