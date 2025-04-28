package structs

type RegisterRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=32"`
	Username string `validate:"required,min=2,max=16"`
}

type RegisterConfirmRequest struct {
	Email string `validate:"required,email"`
	OTP   string `validate:"required,min=6,max=6"`
}

type LoginRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=32"`
}

type IsTokenValidRequest struct {
	Token string `validate:"required"`
}

type LinkGithubRequest struct {
	GithubCode string `validate:"required"`
}

type UnlinkGithubRequest struct {
}

type GetDeveloperProfileRequest struct {
	DeveloperID int64 `validate:"required"`
}
