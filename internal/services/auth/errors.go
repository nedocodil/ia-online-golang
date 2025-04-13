package auth

import "errors"

var (
	ErrActiveLinkAlreadyExists = errors.New("active link already exists")
	ErrActiveLinkNotExists     = errors.New("active link not exists")
	ErrInvalidActiveLink       = errors.New("invalid active link")
	ErrActiveLinkExpired       = errors.New("activation link already expired")
)

var (
	ErrReferralIdNotFound = errors.New("referral code not found")
)
