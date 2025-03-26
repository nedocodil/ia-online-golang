package services

import "errors"

var (
	ErrSaveRefreshToken          = errors.New("error saving refresh token")
	ErrRefreshTokenAlreadyExists = errors.New("refreshing token already exists")
	ErrRefreshTokenNotExists     = errors.New("refreshing token not exists")
	ErrInvalidToken              = errors.New("token is invalid")
	ErrInvalidRefreshToken       = errors.New("refresh token is invalid")
	ErrInvalidAccessToken        = errors.New("access token is invalid")
	ErrExpiredToken              = errors.New("expired token")
	ErrExpiredRefreshToken       = errors.New("expired access token")
	ErrExpiredAccessToken        = errors.New("expired refresh token")
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotActivated  = errors.New("user not activated")
	ErrUserNotFound      = errors.New("user not found")
	ErrIncorrectPassword = errors.New("incorrect password")
)

var (
	ErrActiveLinkAlreadyExists = errors.New("active link already exists")
	ErrActiveLinkNotExists     = errors.New("active link not exists")
	ErrInvalidActiveLink       = errors.New("invalid active link")
	ErrActiveLinkExpired       = errors.New("activation link already expired")
)

var (
	ErrReferralAlreadyUsed = errors.New("referral not used")
)
