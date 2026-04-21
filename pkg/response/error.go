package response

import "github.com/pkg/errors"

var (
	ErrAuthorNotFound = errors.New("author not found")
	ErrAuthorConflict = errors.New("author already exists")
	ErrAuthorCache    = errors.New("author cache not found")

	ErrBookNotFound = errors.New("book not found")
	ErrBookConflict = errors.New("book already exists")
	ErrBookCache    = errors.New("book cache not found")

	ErrPublisherNotFound = errors.New("publisher not found")
	ErrPublisherConflict = errors.New("publisher already exists")
	ErrPublisherCache    = errors.New("publisher cache not found")

	ErrUserNotFound            = errors.New("user not found")
	ErrUserConflict            = errors.New("user already exists")
	ErrUserCache               = errors.New("user cache not found")
	ErrUserInvalidCredentials  = errors.New("invalid credentials")
	ErrUserInvalidToken        = errors.New("invalid or expired refresh token")
	ErrUserInvalidRefreshToken = errors.New("invalid or expired refresh token")
	ErrUserSignToken           = errors.New("failed to sign token")
	ErrUserSignRefreshToken    = errors.New("failed to sign refresh token")
)
