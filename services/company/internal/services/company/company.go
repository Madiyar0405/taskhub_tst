package company

import (
	"company/internal/domain/entities"
	"company/internal/storage"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

type Company struct {
	log           *slog.Logger
	compProvider  CompanyProvider
	issueProvider IssueProvider
	asSolProvider AssignmentSolutionProvider
	tokenTTL      time.Duration
}

type CompanyProvider interface {
	AddGithubIntegration(ctx context.Context, installationID int64, companyName string, logoURL string) (int64, error)
	GetGithubIntegration(ctx context.Context, id int64) (int64, error)
	GetCompany(ctx context.Context, id int64) (*entities.Company, error)
	GetAllCompanies(ctx context.Context) ([]*entities.Company, error)

	IsCompanyRepresentative(ctx context.Context, userID int64, companyID int64) (bool, error)
}

func New(
	log *slog.Logger,
	compProvider CompanyProvider,
	issueProvider IssueProvider,
	asSolProvider AssignmentSolutionProvider,
	tokenTTL time.Duration) *Company {
	return &Company{
		compProvider:  compProvider,
		issueProvider: issueProvider,
		asSolProvider: asSolProvider,
		log:           log,
		tokenTTL:      tokenTTL,
	}
}

func (c Company) CompanyGithubIntegration(ctx context.Context, id int64) (int64, error) {
	const op = "company.CompanyIntegration"

	log := c.log.With(
		slog.String("op", op),
		slog.Int64("id", id),
	)

	log.Info("getting information about github integration")

	installationID, err := c.compProvider.GetGithubIntegration(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrNoRecordFound):
			c.log.Warn("no company found", slog.String("error", err.Error()))
			return 0, fmt.Errorf("%s:%w", op, err)
		default:
			c.log.Warn("failed to get company integrations", slog.String("error", err.Error()))
			return 0, fmt.Errorf("%s:%w", op, err)
		}
	}

	log.Info("integration ID got successfully")

	return installationID, nil

}

func (c Company) AllCompaniesInfo(ctx context.Context) ([]*entities.Company, error) {
	const op = "company.AllCompaniesInfo"

	log := c.log.With(
		slog.String("op", op),
	)

	log.Info("get all companies")

	companies, err := c.compProvider.GetAllCompanies(ctx)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrNoRecordFound):
			c.log.Warn("no companies found", slog.String("error", err.Error()))
			return nil, fmt.Errorf("%s:%w", op, err)
		default:
			c.log.Warn("failed to get companies", slog.String("error", err.Error()))
			return nil, fmt.Errorf("%s:%w", op, err)
		}
	}

	log.Info("companies got successfully")

	return companies, nil
}

func (c Company) CompanyInfo(ctx context.Context, id int64) (*entities.Company, error) {
	const op = "company.CompanyInfo"

	log := c.log.With(
		slog.String("op", op),
	)

	log.Info("get information about company")

	company, err := c.compProvider.GetCompany(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrNoRecordFound):
			c.log.Warn("no company found", slog.String("error", err.Error()))
			return nil, fmt.Errorf("%s:%w", op, err)
		default:
			c.log.Warn("failed to get company information", slog.String("error", err.Error()))
			return nil, fmt.Errorf("%s:%w", op, err)
		}
	}

	log.Info("company information got successfully")

	return company, nil

}

func (c Company) AddCompany(ctx context.Context, installationID int64, companyName string, companyLogo string) (int64, error) {
	const op = "company.AddCompany"

	log := c.log.With(
		slog.String("op", op),
	)

	log.Info("catching information via webhook")

	companyID, err := c.compProvider.AddGithubIntegration(ctx, installationID, companyName, companyLogo)
	if err != nil {
		c.log.Error("company has not been added")
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	log.Info("company integrated successfully via github")

	return companyID, nil

}

func (c Company) VerifyCompanyRepresentative(ctx context.Context, userID int64, issueID int64) (bool, error) {
	const op = "company.VerifyCompanyRepresentative"

	log := c.log.With(
		slog.String("op", op),
		slog.Int64("userID", userID),
		slog.Int64("issueID", issueID),
	)

	log.Info("verifying company representative access")

	companyID, err := c.issueProvider.GetCompanyIDByIssueID(ctx, issueID)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrNoRecordFound):
			log.Warn("no company found for issue", slog.String("error", err.Error()))
			return false, fmt.Errorf("%s:%w", op, err)
		default:
			log.Warn("failed to get company for issue", slog.String("error", err.Error()))
			return false, fmt.Errorf("%s:%w", op, err)
		}
	}

	isRepresentative, err := c.compProvider.IsCompanyRepresentative(ctx, userID, companyID)
	if err != nil {
		log.Warn("failed to verify representative status", slog.String("error", err.Error()))
		return false, fmt.Errorf("%s:%w", op, err)
	}

	if !isRepresentative {
		log.Info("user is not a company representative",
			slog.Int64("userID", userID),
			slog.Int64("companyID", companyID))
	} else {
		log.Info("user verified as company representative")
	}

	return isRepresentative, nil
}
