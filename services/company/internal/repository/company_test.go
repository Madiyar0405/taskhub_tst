package repository

import (
	"company/internal/storage"
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func testDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { db.Close() })

	return db, mock
}

func TestStorage_GetCompany(t *testing.T) {
	db, mock := testDB(t)
	repo := Storage{Db: db}
	ctx := context.Background()

	createdAt := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"company_id", "company_name", "description", "website", "logo", "created_at",
		}).AddRow(1, "Test Corp", "Test Description", "https://test.com", "logo.png", createdAt)

		mock.ExpectPrepare(`SELECT company_id, company_name, description, website, logo, created_at FROM companies WHERE company_id = \$1`).
			ExpectQuery().WithArgs(1).WillReturnRows(rows)

		company, err := repo.GetCompany(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), company.ID)
		assert.Equal(t, "Test Corp", company.CompanyName)
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectPrepare(`SELECT company_id, company_name, description, website, logo, created_at FROM companies WHERE company_id = \$1`).
			ExpectQuery().WithArgs(999).WillReturnError(sql.ErrNoRows)

		_, err := repo.GetCompany(ctx, 999)
		assert.ErrorIs(t, err, storage.ErrNoRecordFound)
	})

	t.Run("DatabaseError", func(t *testing.T) {
		mock.ExpectPrepare(`SELECT company_id, company_name, description, website, logo, created_at FROM companies WHERE company_id = \$1`).
			ExpectQuery().WithArgs(1).WillReturnError(errors.New("connection failed"))

		_, err := repo.GetCompany(ctx, 1)
		assert.ErrorContains(t, err, "connection failed")
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}
