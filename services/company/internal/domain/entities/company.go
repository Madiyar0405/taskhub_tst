package entities

import (
	"database/sql"
	"time"
)

type Company struct {
	ID          int64
	CompanyName string
	Description sql.NullString
	WebsiteURL  sql.NullString
	LogoURL     string
	CreatedAt   time.Time
}
