package auth

import (
	authv1 "auth/gen/auth"
	"auth/internal/config"
	"auth/internal/domain/entities"
	"auth/internal/grpc/structs"
	"auth/internal/services/auth"
	"auth/internal/storage"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

type Auth interface {
	Login(ctx context.Context, email string, password string) (token string, err error)
	RegisterNewUser(ctx context.Context, email string, username string, password string) (userID int64, err error)
	RegisterConfirm(ctx context.Context, email string, otp string) (success bool, message string, err error)

	LinkGithubAccount(ctx context.Context, userID int64, githubID int64, githubUsername string, avatarURL string) (success bool, message string, err error)
	UnlinkGithubAccount(ctx context.Context, userID int64) (success bool, err error)
	DeveloperProfile(ctx context.Context, devID int64) (*entities.DeveloperProfile, error)
	UserInfo(ctx context.Context, id int64) (*entities.User, error)
}

type serverAPI struct {
	authv1.UnimplementedAuthServer
	auth Auth
	v    *validator.Validate

	oauthConfig *oauth2.Config
}

var (
	cfg = config.MustLoad()
)

func Register(gRPC *grpc.Server, auth Auth) {
	authv1.RegisterAuthServer(gRPC, &serverAPI{
		auth: auth,
		v:    validator.New(),
		oauthConfig: &oauth2.Config{
			ClientID:     cfg.Github.ClientID,
			ClientSecret: cfg.Github.ClientSecret,
			Endpoint:     github.Endpoint,
			RedirectURL:  cfg.Github.RedirectURL,
			Scopes:       []string{"read:user", "user:email"},
		},
	})
}

func (s *serverAPI) authorizeUser(ctx context.Context) (jwt.MapClaims, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no metadata in request")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, status.Error(codes.Unauthenticated, "no authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader[0], "Bearer ")

	_, err := s.IsTokenValid(context.Background(), &authv1.IsTokenValidRequest{
		Token: tokenString,
	})
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "invalid token claims")
	}

	return claims, nil
}

//////////////////////////////////////////

func (s *serverAPI) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	registerRequest := structs.RegisterRequest{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	}

	err := s.v.Struct(registerRequest)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetUsername(), req.GetPassword())
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrUserExists):
			return nil, status.Error(codes.AlreadyExists, "email address is already exists")
		case errors.Is(err, auth.ErrUsernameExists):
			return nil, status.Errorf(codes.AlreadyExists, "username is already taken")
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	return &authv1.RegisterResponse{
		UserID: userID,
	}, nil
}

func (s *serverAPI) RegisterConfirm(ctx context.Context, req *authv1.RegisterConfirmRequest) (*authv1.RegisterConfirmResponse, error) {
	registerConfirmRequest := structs.RegisterConfirmRequest{
		Email: req.Email,
		OTP:   req.Otp,
	}

	err := s.v.Struct(registerConfirmRequest)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	success, message, err := s.auth.RegisterConfirm(ctx, req.GetEmail(), req.GetOtp())
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrInvalidOrExpiredOTP):
			return nil, status.Error(codes.InvalidArgument, "Invalid OTP code")
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	return &authv1.RegisterConfirmResponse{
		Message: message,
		Success: success,
	}, nil
}

func (s *serverAPI) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	loginRequest := structs.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	err := s.v.Struct(loginRequest)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidCredentials):
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		default:
			return nil, status.Error(codes.Internal, "failed to log in")
		}
	}

	return &authv1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) LinkGithubAccount(ctx context.Context, req *authv1.LinkGithubRequest) (*authv1.LinkGithubResponse, error) {
	linkGithubRequest := structs.LinkGithubRequest{
		GithubCode: req.GithubCode,
	}
	err := s.v.Struct(linkGithubRequest)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	githubToken, err := s.oauthConfig.Exchange(ctx, linkGithubRequest.GithubCode)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to exchange code for token")
	}

	client := s.oauthConfig.Client(ctx, githubToken)

	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to fetch user info from GitHub")
	}

	defer resp.Body.Close()

	var githubUser struct {
		ID        int64  `json:"id"`
		Login     string `json:"login"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		return nil, status.Error(codes.Internal, "failed to decode GitHub user info")
	}

	claims, err := s.authorizeUser(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	userID, _ := claims["uid"].(float64)

	//if !ok {
	//	return nil, status.Error(codes.PermissionDenied, "no access for this resource")
	//}

	currentUserID := int64(userID)

	success, message, err := s.auth.LinkGithubAccount(ctx, currentUserID, githubUser.ID, githubUser.Login, githubUser.AvatarURL)
	if err != nil {
		switch {
		default:
			return nil, status.Error(codes.Internal, "failed to link github account")
		}
	}

	return &authv1.LinkGithubResponse{
		Success: success,
		Message: message,
	}, nil

}

func (s *serverAPI) UnlinkGithubAccount(ctx context.Context, req *authv1.UnlinkGithubRequest) (*authv1.UnlinkGithubResponse, error) {
	claims, err := s.authorizeUser(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	userID, _ := claims["uid"].(float64)

	currentUserID := int64(userID)

	success, err := s.auth.UnlinkGithubAccount(ctx, currentUserID)
	if err != nil {
		switch {
		default:
			return nil, status.Error(codes.Internal, "failed to unlink github account")
		}
	}

	return &authv1.UnlinkGithubResponse{
		Success: success,
	}, nil
}

func (s *serverAPI) DeveloperProfile(ctx context.Context, req *authv1.GetDeveloperProfileRequest) (*authv1.GetDeveloperProfileResponse, error) {
	GetDeveloperProfileRequest := structs.GetDeveloperProfileRequest{
		DeveloperID: req.DevID,
	}

	err := s.v.Struct(GetDeveloperProfileRequest)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	developer, err := s.auth.DeveloperProfile(ctx, req.DevID)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrNoRecordFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to get developer profile")
		}
	}

	bio := ""
	if developer.Bio != nil {
		bio = *developer.Bio
	}

	avatarUrl := ""
	if developer.AvatarURL != nil {
		avatarUrl = *developer.AvatarURL
	}

	cvUrl := ""
	if developer.CVURL != nil {
		cvUrl = *developer.CVURL
	}

	githubUsername := ""
	if developer.GithubUsername != nil {
		githubUsername = *developer.GithubUsername
	}

	return &authv1.GetDeveloperProfileResponse{
		IsGithubLinked: developer.IsGithubLinked,
		Username:       developer.Username,
		GithubUsername: githubUsername,
		AvatarUrl:      avatarUrl,
		CvUrl:          cvUrl,
		Bio:            bio,
		Email:          developer.Email,
	}, nil
}

func (s *serverAPI) UserInfo(ctx context.Context, req *authv1.UserInfoRequest) (*authv1.UserInfoResponse, error) {
	claims, err := s.authorizeUser(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	userID, _ := claims["uid"].(float64)

	currentUserID := int64(userID)

	user, err := s.auth.UserInfo(ctx, currentUserID)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrNoRecordFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to get user info")
		}
	}

	return &authv1.UserInfoResponse{
		UserId:    user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CompanyId: user.CompanyID.Int64,
	}, nil
}

/////////////////////////////////////////////////

func (s *serverAPI) IsTokenValid(ctx context.Context, req *authv1.IsTokenValidRequest) (*authv1.IsTokenValidResponse, error) {
	isTokenValidRequest := structs.IsTokenValidRequest{
		Token: req.Token,
	}
	err := s.v.Struct(isTokenValidRequest)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := jwt.Parse(req.GetToken(), func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, err
		}
		return []byte(cfg.JwtSignedString), nil
	})

	if err != nil {
		return &authv1.IsTokenValidResponse{IsTokenValid: false}, nil
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &authv1.IsTokenValidResponse{IsTokenValid: true}, nil
	}

	return &authv1.IsTokenValidResponse{
		IsTokenValid: false,
	}, nil
}

func (s *serverAPI) IsGithubLinked(ctx context.Context, req *authv1.IsGithubLinkedRequest) (*authv1.IsGithubLinkedResponse, error) {
	developer, err := s.auth.DeveloperProfile(ctx, req.DevID)

	if err != nil {
		switch {
		case errors.Is(err, storage.ErrNoRecordFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to get developer")
		}
	}

	if !developer.IsGithubLinked {
		return &authv1.IsGithubLinkedResponse{IsGithubLinked: false}, nil
	}

	return &authv1.IsGithubLinkedResponse{
		IsGithubLinked: true,
	}, nil
}
