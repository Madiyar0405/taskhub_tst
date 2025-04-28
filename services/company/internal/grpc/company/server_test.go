package company

import (
	companyv1 "company/gen/company"
	"company/internal/domain/entities"
	"company/internal/storage"
	"context"
	"database/sql"
	"errors"
	"github.com/go-playground/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

type MockCompany struct {
	mock.Mock
}

func (m *MockCompany) AddCompany(ctx context.Context, installationID int64, companyName string, companyLogo string) (companyID int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockCompany) CompanyGithubIntegration(ctx context.Context, id int64) (installationID int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockCompany) AllCompaniesInfo(ctx context.Context) ([]*entities.Company, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockCompany) CompanyInfo(ctx context.Context, id int64) (*entities.Company, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Company), args.Error(1)
}

func TestServerAPI_Company(t *testing.T) {
	mockService := new(MockCompany)
	v := validator.New()
	server := &serverAPI{
		company: mockService,
		v:       v,
	}

	createdAt := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	successCompany := &entities.Company{
		ID:          1,
		CompanyName: "Tech Corp",
		Description: sql.NullString{String: "Tech company", Valid: true},
		WebsiteURL:  sql.NullString{String: "https://tech.com", Valid: true},
		LogoURL:     "logo.png",
		CreatedAt:   createdAt,
	}

	tests := []struct {
		name        string
		request     *companyv1.GetCompanyRequest
		mockSetup   func()
		expected    *companyv1.GetCompanyResponse
		expectedErr codes.Code
	}{
		{
			name: "success",
			request: &companyv1.GetCompanyRequest{
				CompanyId: 1,
			},
			mockSetup: func() {
				mockService.On("CompanyInfo", context.Background(), int64(1)).
					Return(successCompany, nil).Once()
			},
			expected: &companyv1.GetCompanyResponse{
				CompanyId:   1,
				CompanyName: "Tech Corp",
				Description: "Tech company",
				WebsiteUrl:  "https://tech.com",
				Logo:        "logo.png",
				CreatedAt:   timestamppb.New(createdAt),
			},
		},
		{
			name: "not found",
			request: &companyv1.GetCompanyRequest{
				CompanyId: 2,
			},
			mockSetup: func() {
				mockService.On("CompanyInfo", context.Background(), int64(2)).
					Return(nil, storage.ErrNoRecordFound).Once()
			},
			expectedErr: codes.NotFound,
		},
		{
			name: "internal error",
			request: &companyv1.GetCompanyRequest{
				CompanyId: 3,
			},
			mockSetup: func() {
				mockService.On("CompanyInfo", context.Background(), int64(3)).
					Return(nil, errors.New("database error")).Once()
			},
			expectedErr: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			resp, err := server.Company(context.Background(), tt.request)

			if tt.expectedErr != codes.OK {
				assert.Nil(t, resp)
				assert.Equal(t, tt.expectedErr, status.Code(err))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, resp)
			}
		})
	}

	mockService.AssertExpectations(t)
}
