package entities

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID             int64
	Email          string
	Username       string
	Role           string
	HashedPassword []byte
	CreatedAt      time.Time
	UpdatedAt      time.Time

	CompanyID sql.NullInt64
}

func Matches(user *User, password string) error {
	err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))

	if err != nil {
		return err
	}

	return nil
}

type OTP struct {
	Code string
}
