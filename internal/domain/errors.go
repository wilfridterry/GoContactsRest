package domain

import "errors"

var (
	ErrContactNotFound = errors.New("contact not found")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
)
