package service

import (
	"errors"
	"go-ticket/internal/domain"
	"go-ticket/pkg/hash"
	"go-ticket/pkg/jwt"

	"github.com/rs/zerolog/log"
)

type authService struct {
	userRepo  domain.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo domain.UserRepository, jwtSecret string) domain.AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *authService) Register(name, email, password string) (*domain.User, error) {
	existing, _ := s.userRepo.GetByEmail(email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := hash.HashPassword(password)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password during registration")
		return nil, errors.New("failed to hash password")
	}

	user := &domain.User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
		Role:     "user",
	}

	if err := s.userRepo.Create(user); err != nil {
		log.Error().Err(err).Str("email", email).Msg("Failed to save user in database")
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(email, password string) (string, *domain.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if user == nil {
		return "", nil, errors.New("invalid credentials")
	}

	if !hash.CheckPasswordHash(password, user.Password) {
		return "", nil, errors.New("invalid credentials")
	}

	token, err := jwt.GenerateToken(int64(user.ID), user.Role, s.jwtSecret)
	if err != nil {
		log.Error().Err(err).Int64("user_id", user.ID).Msg("Failed to generate JWT token")
		return "", nil, errors.New("failed to generate token")
	}

	return token, user, nil
}

func (s *authService) GetProfile(userID int64) (*domain.User, error) {
	user, err := s.userRepo.GetByID(int64(userID))
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (s *authService) ChangePassword(userID int64, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if !hash.CheckPasswordHash(oldPassword, user.Password) {
		return errors.New("old password is incorrect")
	}

	if oldPassword == newPassword {
		return errors.New("new password must be different from old password")
	}

	if len(newPassword) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	hashedPassword, err := hash.HashPassword(newPassword)
	if err != nil {
		log.Error().Err(err).Int64("user_id", userID).Msg("Failed to hash new password")
		return errors.New("failed to hash password")
	}

	user.Password = hashedPassword
	return s.userRepo.Update(user)
}

func (s *authService) ValidateToken(token string) (*jwt.Claims, error) {
	claims, err := jwt.ValidateToken(token, s.jwtSecret)
	if err != nil {
		return nil, errors.New("invalid or expired token")
	}

	if jwt.IsExpired(claims) {
		return nil, errors.New("token has expired")
	}

	return claims, nil
}
