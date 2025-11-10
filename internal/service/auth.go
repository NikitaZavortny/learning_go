package service

import (
	"auth-server/internal/model"
	"auth-server/internal/repository"
	"auth-server/internal/utils"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo    *repository.UserRepository
	MailService *MailService
	jwtSecret   string
}

func NewAuthService(userRepo *repository.UserRepository, MailService *MailService, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		MailService: MailService,
		jwtSecret:   jwtSecret,
	}
}

func (s *AuthService) Register(req *model.RegisterRequest) (*model.TokenResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	activationlink := uuid.New().String()
	user, err := s.userRepo.CreateUser(req.Email, string(hashedPassword), activationlink)
	if err != nil {
		return nil, err
	}

	// Generate tokens
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Email, s.jwtSecret)
	if err != nil {
		return nil, err
	}
	s.MailService.SendActivationMAil(activationlink, user.Email)

	return &model.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) Login(req *model.LoginRequest) (*model.TokenResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate tokens
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Email, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &model.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) RefreshTokens(refreshToken string) (*model.TokenResponse, error) {
	claims, err := utils.ValidateToken(refreshToken, s.jwtSecret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	user, err := s.userRepo.GetUserByID(claims.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Generate new tokens
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := utils.GenerateRefreshToken(user.ID, user.Email, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &model.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AuthService) Activate(link string) (*model.TokenResponse, error) {
	_, err := s.userRepo.GetUserByLink(link)
	if err != nil {
		return nil, err
	}
	_, e := s.userRepo.Activate(link)

	if e != nil {
		return nil, e
	}
	return nil, nil
}
