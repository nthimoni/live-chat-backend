package service

import (
	"fmt"
	"live-chat-backend/internal/models"
	"live-chat-backend/internal/repository"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	jwtSecret []byte
	userRepo  repository.UserRepository
}

func NewAuthService(jwtSecret []byte, userRepo repository.UserRepository) *AuthService {
	return &AuthService{jwtSecret: jwtSecret, userRepo: userRepo}
}

func (s *AuthService) RegisterUser(username, email, password string) (*models.User, error) {
	user := &models.User{
		Username: username,
		Email:    email,
		Password: password,
	}

	err := s.ValidatePasswordPolicy(password)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

type JWTClaims struct {
	UserID uint   `json:"sub"`
	Email  string `json:"email"`
	Role   string `json:"role,omitempty"`
	jwt.RegisteredClaims
}

func (s *AuthService) GenerateUserJWT(userID uint, email string) (string, error) {
	var now time.Time = time.Now()

	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)), // 24 hour validity
			Issuer:    "live-chat",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) ValidatePasswordPolicy(password string) error {
	const minLength = 8
	var hasUpper, hasLower, hasDigit, hasSpecial bool

	if len(password) < minLength {
		return fmt.Errorf("password must be at least %d characters", minLength)
	}

	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsDigit(c):
			hasDigit = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !hasDigit {
		return fmt.Errorf("password must contain at least one digit")
	}
	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
