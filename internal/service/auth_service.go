package service

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/chillman2101/gits-catalogue/internal/model"
	"github.com/chillman2101/gits-catalogue/internal/repository"
	res_err "github.com/chillman2101/gits-catalogue/pkg/response"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type AuthService interface {
	Register(email, password string) error
	Login(email, password string) (*TokenPair, error)
	Refresh(refreshToken string) (*TokenPair, error)
	Logout(userID uint) error
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{repo}
}

func (s *authService) Register(email, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logrus.WithError(err).Error("failed to hash password")
		return errors.New("failed to hash password")
	}
	user := &model.User{Email: email, Password: string(hashed)}
	return s.repo.Create(user)
}

func (s *authService) Login(email, password string) (*TokenPair, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, res_err.ErrUserInvalidCredentials
	}
	return s.createUserToken(user)
}

func (s *authService) Refresh(refreshToken string) (*TokenPair, error) {
	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(os.Getenv("JWT_REFRESH_SECRET")), nil
	})
	if err != nil || !token.Valid {
		logrus.WithError(err).Warn("invalid refresh token")
		return nil, res_err.ErrUserInvalidToken
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))
	refreshTokenHash := claims["refresh_token_hash"].(string)

	user, err := s.repo.FindByIDAndRefreshTokenHash(userID, refreshTokenHash)
	if err != nil {
		logrus.WithError(err).WithField("user_id", userID).Error("user not found")
		return nil, err
	}

	if user.RefreshTokenHash != refreshTokenHash {
		logrus.WithField("user_id", userID).Warn("refresh token mismatch — possible token reuse")
		return nil, res_err.ErrUserInvalidRefreshToken
	}

	return s.createUserToken(user)
}

func (s *authService) Logout(userID uint) error {
	if err := s.repo.UpdateTokenHash(userID, ""); err != nil {
		logrus.WithError(err).WithField("user_id", userID).Error("failed to clear token hash")
		return err
	}
	if err := s.repo.UpdateRefreshTokenHash(userID, ""); err != nil {
		logrus.WithError(err).WithField("user_id", userID).Error("failed to clear refresh token hash")
		return err
	}
	return nil
}

func (s *authService) createUserToken(user *model.User) (*TokenPair, error) {
	tokenHash := fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%d%s%d", user.ID, user.Email, time.Now().Unix()))))
	refreshTokenHash := fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%d%s%d", user.ID, user.Email, time.Now().Unix()))))
	accessStr, err := signToken(jwt.MapClaims{
		"user_id":    user.ID,
		"email":      user.Email,
		"exp":        time.Now().Add(15 * time.Minute).Unix(),
		"token_hash": tokenHash,
	}, os.Getenv("JWT_SECRET"))
	if err != nil {
		logrus.WithError(err).Error("failed to sign access token")
		return nil, res_err.ErrUserSignToken
	}

	refreshStr, err := signToken(jwt.MapClaims{
		"user_id":            user.ID,
		"exp":                time.Now().Add(7 * 24 * time.Hour).Unix(),
		"refresh_token_hash": refreshTokenHash,
	}, os.Getenv("JWT_REFRESH_SECRET"))
	if err != nil {
		logrus.WithError(err).Error("failed to sign refresh token")
		return nil, res_err.ErrUserSignRefreshToken
	}

	if err := s.repo.UpdateTokenHash(user.ID, tokenHash); err != nil {
		logrus.WithError(err).WithField("user_id", user.ID).Error("failed to save access token hash")
		return nil, res_err.ErrUserSignToken
	}
	if err := s.repo.UpdateRefreshTokenHash(user.ID, refreshTokenHash); err != nil {
		logrus.WithError(err).WithField("user_id", user.ID).Error("failed to save refresh token hash")
		return nil, res_err.ErrUserSignRefreshToken
	}

	return &TokenPair{AccessToken: accessStr, RefreshToken: refreshStr}, nil
}

func signToken(claims jwt.MapClaims, secret string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}
