package storage

import "errors"

var (
	ErrUserExists          = errors.New("email already in use")
	ErrUsernameExists      = errors.New("username is already taken")
	ErrNoRecordFound       = errors.New("no record found")
	ErrAppNotFound         = errors.New("app not found")
	ErrInvalidOrExpiredOTP = errors.New("invalid confirmation code")
	ErrOtpAlreadySent      = errors.New("otp is sent already")
	ErrGithubLinked        = errors.New("this github account is already linked")
)
