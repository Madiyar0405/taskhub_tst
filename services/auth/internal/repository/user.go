package repository

import (
	"auth/internal/domain/entities"
	"auth/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
)

type Storage struct {
	Db *sql.DB
}

func (s *Storage) DeleteInactiveUsers(ctx context.Context) error {
	const op = "repository.user.DeleteInactiveUsers"
	query := `DELETE FROM users WHERE activated = FALSE`
	_, err := s.Db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

//
//func (s *Storage) App(ctx context.Context, appID int) (entities.App, error) {
//	//TODO implement me
//	panic("implement me")
//}

func (s *Storage) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	const op = "repository.user.GetUser"

	query, err := s.Db.Prepare("SELECT id, email, username, password_hash FROM users WHERE email = $1 AND activated = true")
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	var user entities.User

	err = query.QueryRowContext(ctx, email).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.HashedPassword)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, storage.ErrNoRecordFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, username string, passHash []byte) (int64, error) {
	const op = "repository.user.SaveUser"

	var id int64
	queryUsers := `INSERT INTO users (email, username, password_hash, activated) VALUES ($1, $2, $3, $4) RETURNING id`

	argsUsers := []any{email, username, passHash, false}
	err := s.Db.QueryRowContext(ctx, queryUsers, argsUsers...).Scan(&id)
	if err != nil {
		var pqErr *pq.Error
		switch {
		case errors.As(err, &pqErr) && pqErr.Code == "23505":
			if pqErr.Constraint == "users_email_key" {
				return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
			}
			if pqErr.Constraint == "users_username_key" {
				return 0, fmt.Errorf("%s: %w", op, storage.ErrUsernameExists)
			}
		default:
			return 0, fmt.Errorf("%s: %w", op, err)
		}
	}

	code, err := GenerateRandomCode()
	if err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}
	fmt.Println(code)

	err = SendEmailWithCode(email, code)
	if err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	queryOtps := `INSERT INTO otps (email, otp_code) VALUES ($1, $2)`

	argsOtps := []any{email, code}
	_, err = s.Db.ExecContext(ctx, queryOtps, argsOtps...)

	if err != nil {
		var pqErr *pq.Error
		switch {
		case errors.As(err, &pqErr) && pqErr.Code == "23505":
			return 0, fmt.Errorf("%s:%w", op, storage.ErrOtpAlreadySent)
		default:
			return 0, fmt.Errorf("%s: %w", op, err)
		}
	}

	return id, nil
}

func (s *Storage) ConfirmUser(ctx context.Context, email string, otp string) (bool, string, error) {
	const op = "repository.user.ConfirmUser"

	var isValid bool

	query := `
        SELECT EXISTS (
            SELECT 1 FROM otps
            WHERE email = $1 AND otp_code = $2 AND expires_at > NOW()
        )
    `

	args := []any{email, otp}
	err := s.Db.QueryRowContext(ctx, query, args...).Scan(&isValid)
	if err != nil {
		return false, "", fmt.Errorf("%s:%w", op, err)
	}

	if !isValid {
		return false, "Invalid or expired OTP, try again later", storage.ErrInvalidOrExpiredOTP
	}

	queryUsers := `UPDATE users SET activated = TRUE WHERE email= $1`
	_, err = s.Db.ExecContext(ctx, queryUsers, email)
	if err != nil {
		return false, "", fmt.Errorf("%s:%w", op, err)
	}

	queryOtps := `DELETE FROM otps WHERE email = $1`
	_, err = s.Db.ExecContext(ctx, queryOtps, email)
	if err != nil {
		return false, "", fmt.Errorf("%s:%w", op, err)
	}

	return true, "Email confirmed successfully", nil
}

///////////////////////// DEVELOPER /////////////////////////////

func (s *Storage) SaveDeveloper(ctx context.Context, userID, githubID int64, githubUsername string, avatarUrl string) (bool, string, error) {
	const op = "repository.user.SaveDeveloper"

	query := `INSERT INTO developers (developer_id, github_id, github_username, avatar_url) VALUES ($1, $2, $3, $4)`
	args := []any{userID, githubID, githubUsername, avatarUrl}

	_, err := s.Db.ExecContext(ctx, query, args...)
	if err != nil {
		var pqErr *pq.Error
		switch {
		case errors.As(err, &pqErr) && pqErr.Code == "23505":
			return false, "", fmt.Errorf("%s: %w", op, storage.ErrGithubLinked)
		default:
			return false, "", fmt.Errorf("%s: %w", op, err)
		}
	}

	return true, "Github account successfully linked", nil

}

func (s *Storage) DeleteDeveloper(ctx context.Context, userID int64) (bool, error) {
	const op = "repository.user.DeleteDeveloper"
	query := `DELETE FROM developers WHERE developer_id = $1`

	result, err := s.Db.ExecContext(ctx, query, userID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if rowsAffected == 0 {
		return false, storage.ErrNoRecordFound
	}

	return true, nil

}

func (s *Storage) GetDeveloper(ctx context.Context, devID int64) (*entities.DeveloperProfile, error) {
	const op = "repository.user.GetDeveloper"

	if devID < 1 {
		return nil, storage.ErrNoRecordFound
	}

	query, err := s.Db.Prepare(`
	SELECT 
		COALESCE(d.github_id IS NOT NULL, false) AS is_github_linked,
		d.bio,
		d.github_username,
		d.avatar_url,
		d.cv_url,
		u.email,
		u.username
-- 		u.created_at
	FROM 
		users u
	LEFT JOIN 
		developers d ON u.id = d.developer_id
	WHERE 
		u.id = $1
`)

	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	var developer entities.DeveloperProfile

	err = query.QueryRowContext(ctx, devID).Scan(
		&developer.IsGithubLinked,
		&developer.Bio,
		&developer.GithubUsername,
		&developer.AvatarURL,
		&developer.CVURL,
		&developer.Email,
		&developer.Username,
		//&developer.CreatedAt
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, storage.ErrNoRecordFound
		default:
			return nil, err
		}
	}

	return &developer, nil

}

func (s *Storage) GetUserByID(ctx context.Context, id int64) (*entities.User, error) {
	const op = "repository.user.GetUserByID"

	query, err := s.Db.Prepare(
		`SELECT           
    		u.id, 
            u.username, 
            u.email, 
            u.role,
            cu.company_id
        FROM users u
        LEFT JOIN company_users cu ON u.id = cu.user_id
        WHERE u.id = $1`)

	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	var user entities.User

	err = query.QueryRowContext(ctx, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Role,
		&user.CompanyID)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, storage.ErrNoRecordFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
