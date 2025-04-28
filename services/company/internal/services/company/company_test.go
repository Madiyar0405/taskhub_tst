package company

import (
	"company/internal/domain/entities"
	"company/internal/storage"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"testing"
)

type MockCompanyProvider struct {
	mock.Mock
}

func (m *MockCompanyProvider) AddGithubIntegration(ctx context.Context, installationID int64, companyName string, logoURL string) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockCompanyProvider) GetGithubIntegration(ctx context.Context, id int64) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockCompanyProvider) GetAllCompanies(ctx context.Context) ([]*entities.Company, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockCompanyProvider) GetCompany(ctx context.Context, id int64) (*entities.Company, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.Company), args.Error(1)
}

func TestCompanyService_CompanyInfo(t *testing.T) {
	mockProvider := new(MockCompanyProvider)
	service := Company{
		log:          slog.Default(),
		compProvider: mockProvider,
	}

	expectedCompany := &entities.Company{
		ID:          1,
		CompanyName: "Test Corp",
	}

	t.Run("Success", func(t *testing.T) {
		mockProvider.On("GetCompany", context.Background(), int64(1)).
			Return(expectedCompany, nil).Once()

		result, err := service.CompanyInfo(context.Background(), 1)
		assert.NoError(t, err)
		assert.Equal(t, expectedCompany, result)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockProvider.On("GetCompany", context.Background(), int64(2)).
			Return((*entities.Company)(nil), storage.ErrNoRecordFound).Once()

		_, err := service.CompanyInfo(context.Background(), 2)
		assert.ErrorIs(t, err, storage.ErrNoRecordFound)
	})

	t.Run("DatabaseError", func(t *testing.T) {
		mockProvider.On("GetCompany", context.Background(), int64(3)).
			Return((*entities.Company)(nil), errors.New("db error")).Once()

		_, err := service.CompanyInfo(context.Background(), 3)
		assert.ErrorContains(t, err, "db error")
	})

	mockProvider.AssertExpectations(t)
}
