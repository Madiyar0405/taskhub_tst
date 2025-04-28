package auth

import (
	"auth/internal/domain/entities"
	"auth/internal/lib/jwt"
	"auth/internal/lib/logger/sl"
	"auth/internal/storage"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, username string, passHash []byte) (userID int64, err error)
	ConfirmUser(ctx context.Context, email string, otp string) (success bool, message string, err error)

	SaveDeveloper(ctx context.Context, userID int64, githubID int64, githubUsername string, avatarUrl string) (success bool, message string, err error)
	DeleteDeveloper(ctx context.Context, userID int64) (success bool, err error)
}

type UserProvider interface {
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	GetUserByID(ctx context.Context, id int64) (*entities.User, error)

	GetDeveloper(ctx context.Context, devID int64) (*entities.DeveloperProfile, error)
}

//type AppProvider interface {
//	App(ctx context.Context, appID int) (entities.App, error)
//}

func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	tokenTTL time.Duration) *Auth {
	return &Auth{
		userSaver:    userSaver,
		userProvider: userProvider,
		log:          log,
		tokenTTL:     tokenTTL,
	}
}

var (
	ErrInvalidCredentials   = errors.New("invalid authentication credentials")
	ErrInvalidAppId         = errors.New("invalid app id")
	ErrUserExists           = errors.New("user already exists")
	ErrUsernameExists       = errors.New("username is already taken")
	ErrInvalidGithubAccount = errors.New("invalid github account")
)

func (a *Auth) RegisterNewUser(ctx context.Context, email string, username string, password string) (int64, error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))

		return 0, fmt.Errorf("%s:%w", op, err)

	}

	id, err := a.userSaver.SaveUser(ctx, email, username, passHash)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrUserExists):
			a.log.Warn("user already exists", slog.String("error", err.Error()))
			return 0, fmt.Errorf("%s:%w", op, ErrUserExists)

		case errors.Is(err, storage.ErrUsernameExists):
			a.log.Warn("username is already taken", slog.String("error", err.Error()))
			return 0, fmt.Errorf("%s:%w", op, ErrUsernameExists)

		default:
			a.log.Error("failed to save user", sl.Err(err))
			return 0, fmt.Errorf("%s:%w", op, err)
		}
	}

	return id, nil
}

func (a *Auth) RegisterConfirm(ctx context.Context, email string, otp string) (bool, string, error) {
	const op = "auth.RegisterConfirm"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("confirming email")

	success, message, err := a.userSaver.ConfirmUser(ctx, email, otp)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrInvalidOrExpiredOTP):
			a.log.Warn("invalid OTP", slog.String("error", err.Error()))
			return false, "Invalid OTP", err
		default:
			a.log.Error("failed to confirm email", sl.Err(err))
			return false, "", err
		}
	}

	return success, message, nil

}

func (a *Auth) Login(ctx context.Context, email string, password string) (string, error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email))

	log.Info("attempting to login user")

	user, err := a.userProvider.GetUserByEmail(ctx, email)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrNoRecordFound):
			a.log.Warn("user not found", slog.String("error", err.Error()))
			return "", fmt.Errorf("%s:%w", op, ErrInvalidCredentials)

		default:
			a.log.Error("failed to get user", sl.Err(err))
			return "", fmt.Errorf("%s:%w", op, err)
		}
	}

	err = entities.Matches(user, password)
	if err != nil {
		a.log.Info("invalid credentials", sl.Err(err))

		return "", fmt.Errorf("%s:%w", op, ErrInvalidCredentials)
	}

	log.Info("logged in successfully")

	token, err := jwt.NewToken(*user, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token", sl.Err(err))

		return "", fmt.Errorf("%s:%w", op, err)
	}

	return token, nil
}

func (a *Auth) LinkGithubAccount(ctx context.Context, userID int64, githubID int64, githubUsername string, avatarURL string) (bool, string, error) {
	const op = "auth.LinkGithubAccount"

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("userID", userID),
		slog.Int64("githubID", githubID),
	)

	log.Info("attempting to link github")

	success, message, err := a.userSaver.SaveDeveloper(ctx, userID, githubID, githubUsername, avatarURL)
	if err != nil {
		// TODO: ERrors
		switch {
		case errors.Is(err, storage.ErrGithubLinked):
			a.log.Warn("github account is already taken", slog.String("error", err.Error()))
			return false, "github account is already taken", fmt.Errorf("%s:%w", op, err)
		default:
			a.log.Error("failed to link github account", sl.Err(err))
			return false, "failed to link github account", fmt.Errorf("%s:%w", op, err)
		}
	}

	log.Info("linked successfully")

	return success, message, nil

}

func (a *Auth) UnlinkGithubAccount(ctx context.Context, userID int64) (bool, error) {
	const op = "auth.UnlinkGithubAccount"

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("userID", userID))

	log.Info("attempting to unlink github")

	success, err := a.userSaver.DeleteDeveloper(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrNoRecordFound):
			a.log.Warn("no github account found", slog.String("error", err.Error()))
			return false, fmt.Errorf("%s:%w", op, err)
		default:
			a.log.Warn("failed to unlink github account", slog.String("error", err.Error()))
			return false, fmt.Errorf("%s:%w", op, err)
		}
	}

	log.Info("unlinked successfully")

	return success, nil
}

func (a *Auth) DeveloperProfile(ctx context.Context, devID int64) (*entities.DeveloperProfile, error) {
	const op = "auth.DeveloperProfile"

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("devID", devID))

	log.Info("attempting to get profile")

	developer, err := a.userProvider.GetDeveloper(ctx, devID)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrNoRecordFound):
			a.log.Warn("no developer profile found", slog.String("error", err.Error()))
			return nil, fmt.Errorf("%s:%w", op, err)
		default:
			a.log.Warn("failed to get developer profile", slog.String("error", err.Error()))
			return nil, fmt.Errorf("%s:%w", op, err)
		}
	}

	log.Info("profile got successfully")

	return developer, nil
}

func (a *Auth) UserInfo(ctx context.Context, id int64) (*entities.User, error) {
	const op = "auth.UserInfo"

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("userID", id))

	log.Info("attempting to get user info")

	user, err := a.userProvider.GetUserByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrNoRecordFound):
			a.log.Warn("no user found", slog.String("error", err.Error()))
			return nil, fmt.Errorf("%s:%w", op, err)
		default:
			a.log.Warn("failed to get user", slog.String("error", err.Error()))
			return nil, fmt.Errorf("%s:%w", op, err)
		}
	}

	log.Info("user got successfully")

	return user, nil
}
