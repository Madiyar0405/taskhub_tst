package grpc

import (
	authv1 "company/gen/auth"
	"context"
	"fmt"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"time"
)

type Client struct {
	api authv1.AuthClient
	log *slog.Logger
}

func New(ctx context.Context,
	log *slog.Logger,
	addr string,
	timeout time.Duration,
	retriesCount int) (*Client, error) {
	const op = "grpc.New"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	// TODO: защищенное соединение
	cc, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return &Client{
		api: authv1.NewAuthClient(cc),
	}, nil
}

func InterceptorLogger(l *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, lvl grpclog.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func (c *Client) IsTokenValid(ctx context.Context, token string) (bool, error) {
	const op = "grpc.IsTokenValid"

	resp, err := c.api.IsTokenValid(ctx, &authv1.IsTokenValidRequest{
		Token: token,
	})

	if err != nil {
		return false, fmt.Errorf("%s:%w", op, err)
	}

	return resp.IsTokenValid, nil
}

func (c *Client) IsGithubLinked(ctx context.Context, id int64) (bool, error) {
	const op = "grpc.IsGithubLinked"

	resp, err := c.api.IsGithubLinked(ctx, &authv1.IsGithubLinkedRequest{
		DevID: id,
	})

	if err != nil {
		return false, fmt.Errorf("%s:%w", op, err)
	}

	return resp.IsGithubLinked, nil
}
