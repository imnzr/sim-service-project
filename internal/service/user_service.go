package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/imnzr/sim-service-project/config"
	"github.com/imnzr/sim-service-project/internal/repository"
	"github.com/imnzr/sim-service-project/models"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, req *models.RegisterUser) error
	Login(ctx context.Context, req *models.LoginRequest) (*models.TokenResponse, error)
	GetUserProfile(ctx context.Context, userId uint) (*models.UserProfileResponse, error)
	GenerateAccessToken(userId uint, email string) (string, error)
	GenerateRefreshToken(userId uint) (string, error)
	ValidateToken(tokenString string) (jwt.MapClaims, error)
}

type UserServiceImplementation struct {
	UserRepo             repository.UserRepository
	jwtSecretKey         []byte
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewUserService(userRepo repository.UserRepository, cfg *config.AppConfig) UserService {
	return &UserServiceImplementation{
		UserRepo:             userRepo,
		jwtSecretKey:         []byte(cfg.JWTSecretKey),
		accessTokenDuration:  cfg.AccessTokenDuration,
		refreshTokenDuration: cfg.RefreshTokenDuration,
	}
}

// ValidateToken implements UserService.
func (service *UserServiceImplementation) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return service.jwtSecretKey, nil
	})

	if err != nil {
		log.Printf("token validation: %v", err)
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token claims")
}

// GetUserProfile implements UserService.
func (service *UserServiceImplementation) GetUserProfile(ctx context.Context, userId uint) (*models.UserProfileResponse, error) {
	user, err := service.UserRepo.GetUserById(ctx, userId)
	if err != nil {
		log.Printf("user profile not found")
		return nil, fmt.Errorf("user profile not found")
	}
	response := models.UserProfileResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	return &response, nil
}

// GenerateAccessToken implements UserService.
func (service *UserServiceImplementation) GenerateAccessToken(userId uint, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"email":   email,
		"expired": time.Now().Add(service.accessTokenDuration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(service.jwtSecretKey)
}

// GenerateRefreshToken implements UserService.
func (service *UserServiceImplementation) GenerateRefreshToken(userId uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"expired": time.Now().Add(service.refreshTokenDuration).Unix(),
		"type":    "refresh",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(service.jwtSecretKey)
}

// Login implements UserService.
func (service *UserServiceImplementation) Login(ctx context.Context, req *models.LoginRequest) (*models.TokenResponse, error) {
	user, err := service.UserRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if user == nil {
		log.Printf("no user found for email: %s", req.Email)
		return nil, fmt.Errorf("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid credentiasl")
	}

	// generate access token
	accessToken, err := service.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// generate refresh token
	refreshToken, err := service.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Register implements UserService.
func (service *UserServiceImplementation) Register(ctx context.Context, req *models.RegisterUser) error {

	_, err := service.UserRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("User already existed")
		return errors.New("user already existed")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("error hashing password user: %v", err)
		return errors.New("error hashing password")
	}

	user := &models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return service.UserRepo.CreateUser(ctx, user)
}
