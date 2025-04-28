package app

import (
	grpcapp "company/internal/app/grpc"
	"company/internal/app/http"
	"company/internal/app/http/webhook"
	authgrpc "company/internal/clients/auth/grpc"
	"company/internal/repository"
	"company/internal/services/company"
	"company/internal/storage/postgres"
	"log/slog"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
	HTTPSrv *http.App
	Storage *repository.Storage
}

func New(log *slog.Logger, grpcPort int, httpPort int, storagePath string, tokenTTL time.Duration, webhookSecret string, authClient *authgrpc.Client) *App {

	storage, err := postgres.New(storagePath)

	if err != nil {
		panic(err)
	}

	companyService := company.New(log, storage, storage, storage, tokenTTL)
	grpcApp := grpcapp.New(log, companyService, companyService, grpcPort, authClient)

	webhookHandler := webhook.NewHandler(log, companyService, companyService, webhookSecret)

	httpApp := http.New(log, httpPort, grpcPort, webhookHandler)

	return &App{
		GRPCSrv: grpcApp,
		HTTPSrv: httpApp,
		Storage: storage,
	}
}
