package repository

import (
	"company/internal/domain/entities"
	"company/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

//TODO s *Storage

type Storage struct {
	Db *sql.DB
}

func (s *Storage) AddGithubIntegration(ctx context.Context, installationID int64, companyName string, logoURL string) (int64, error) {
	const op = "repository.company.AddGithubIntegration"

	var id int64

	query := `INSERT INTO companies (installation_id, company_name, logo) VALUES ($1, $2, $3) RETURNING company_id`

	args := []any{installationID, companyName, logoURL}
	err := s.Db.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	return id, nil
}

func (s *Storage) GetGithubIntegration(ctx context.Context, id int64) (installationID int64, err error) {
	const op = "repository.company.CompanyIntegration"

	query, err := s.Db.Prepare("SELECT installation_id FROM companies WHERE company_id = $1")
	if err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	err = query.QueryRowContext(ctx, id).Scan(&installationID)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return 0, storage.ErrNoRecordFound
		default:
			return 0, fmt.Errorf("%s:%w", op, err)
		}
	}

	return installationID, nil

}

func (s *Storage) GetAllCompanies(ctx context.Context) ([]*entities.Company, error) {
	const op = "repository.company.GetAllCompanies"

	query := "SELECT company_id, company_name, logo FROM companies"

	rows, err := s.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	//totalRecords := 0
	var companies []*entities.Company

	for rows.Next() {
		company := &entities.Company{}

		err := rows.Scan(
			&company.ID,
			&company.CompanyName,
			&company.LogoURL)

		if err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}

		companies = append(companies, company)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	//if totalRecords == 0 {
	//	return nil, storage.ErrNoRecordFound
	//}

	return companies, nil
}

func (s *Storage) GetCompany(ctx context.Context, id int64) (*entities.Company, error) {
	const op = "repository.company.GetCompany"

	query, err := s.Db.Prepare("SELECT company_id, company_name, description, website, logo, created_at FROM companies WHERE company_id = $1")
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	var company entities.Company

	err = query.QueryRowContext(ctx, id).Scan(
		&company.ID,
		&company.CompanyName,
		&company.Description,
		&company.WebsiteURL,
		&company.LogoURL,
		&company.CreatedAt)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, storage.ErrNoRecordFound
		default:
			return nil, fmt.Errorf("%s:%w", op, err)
		}
	}

	return &company, nil
}

func (s *Storage) IsCompanyRepresentative(ctx context.Context, userID int64, companyID int64) (bool, error) {
	const op = "repository.user.IsCompanyRepresentative"

	var isRepresentative bool
	query := `
        SELECT EXISTS (
            SELECT 1 FROM users u
            JOIN company_users cu ON u.id = cu.user_id
            WHERE u.id = $1 AND cu.company_id = $2 AND u.role = 'company'
        )
    `

	err := s.Db.QueryRowContext(ctx, query, userID, companyID).Scan(&isRepresentative)
	if err != nil {
		return false, fmt.Errorf("%s:%w", op, err)
	}

	return isRepresentative, nil
}
