package usecase

import (
	"errors"
	"socialnetwork/internal/domain"
	"socialnetwork/internal/middleware"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	userRepo    domain.UserRepository
	jwtSecret   string
}

func NewAuthUseCase(userRepo domain.UserRepository, jwtSecret string) *AuthUseCase {
	return &AuthUseCase{
		userRepo:    userRepo,
		jwtSecret:   jwtSecret,
	}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token     string       `json:"token"`
	User      domain.User `json:"user"`
}

func (a *AuthUseCase) Register(req *RegisterRequest) (*AuthResponse, error) {
	// Check if email already exists
	existingUser, err := a.userRepo.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &domain.User{
		ID:        uuid.New(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		FullName:  req.FullName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := a.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Generate token
	token, err := middleware.GenerateToken(user.ID.String(), a.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (a *AuthUseCase) Login(req *LoginRequest) (*AuthResponse, error) {
	user, err := a.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate token
	token, err := middleware.GenerateToken(user.ID.String(), a.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}
