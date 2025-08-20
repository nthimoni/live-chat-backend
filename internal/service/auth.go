package service

import (
	"errors"
	"fmt"
	"live-chat-backend/internal/models"
	"live-chat-backend/internal/repository"
	"log"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrPasswordPolicy = errors.New("invalid password policy")

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

func (s *AuthService) LoginUser(email, password string) (*models.User, error) {
	// check password policy to avoid useless call to the database
	err := s.ValidatePasswordPolicy(password)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	user, err := s.userRepo.FindByEmail(email)
	if err == gorm.ErrRecordNotFound {
		fmt.Println("user not found:", err)
		return nil, errors.New("invalid credentials")

	} else if err != nil {
		// Do not expose database errors to the user
		log.Println("database error", err)
		return nil, fmt.Errorf("unable to verify credentials, please try again later")
	}

	if !CheckPassword(user.Password, password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

func (s *AuthService) ValidatePasswordPolicy(password string) error {
	const minLength = 8
	const maxLength = 72 // maximum length handled by bcrypt
	var hasUpper, hasLower, hasDigit, hasSpecial bool

	if len(password) < minLength {
		return fmt.Errorf("%w: password must be at least %d characters", ErrPasswordPolicy, minLength)
	}
	if (len(password)) > maxLength {
		return fmt.Errorf("%w: password must be at most %d characters", ErrPasswordPolicy, maxLength)
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
		return fmt.Errorf("%w: password must contain at least one uppercase letter", ErrPasswordPolicy)
	}
	if !hasLower {
		return fmt.Errorf("%w: password must contain at least one lowercase letter", ErrPasswordPolicy)
	}
	if !hasDigit {
		return fmt.Errorf("%w: password must contain at least one digit", ErrPasswordPolicy)
	}
	if !hasSpecial {
		return fmt.Errorf("%w: password must contain at least one special character", ErrPasswordPolicy)
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

type JWTClaims struct {
	UserID uint `json:"sub"`
	jwt.RegisteredClaims
}

func (s *AuthService) GenerateUserJWT(userID uint) (string, error) {
	var now time.Time = time.Now()

	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)), // 24 hour validity
			Issuer:    "live-chat",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) ValidateTokenAndGetUser(tokenToValidate string) (*models.User, error) {
	token, err := jwt.ParseWithClaims(tokenToValidate, &JWTClaims{}, s.getTokenSigningKey)

	if err != nil || !token.Valid {
		log.Println("token error", err)
		return nil, errors.New("invalid or expired token")
	}

	claims := token.Claims.(*JWTClaims)

	user, err := s.userRepo.FindById(claims.UserID)
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("user not found")
	} else if err != nil {
		log.Println("database error", err)
		return nil, errors.New("unable to retrieve user")
	}

	return user, nil
}

func (s *AuthService) getTokenSigningKey(token *jwt.Token) (any, error) {
	// type assertion to verify the signing method
	_, isValidSigningMthod := token.Method.(*jwt.SigningMethodHMAC)
	if !isValidSigningMthod {
		return nil, errors.New("invalid signing method")
	}

	return s.jwtSecret, nil
}
