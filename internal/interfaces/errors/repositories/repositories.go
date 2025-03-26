package repositories

import "errors"

var (
	ErrUserExists       = errors.New("user already exists")
	ErrUserNotFound     = errors.New("user not found")
	ErrUserIsNotUpdated = errors.New("user is not updated")
)
var (
	ErrActivationLinkIsNotFound   = errors.New("activation link is not found")
	ErrActivationLinkIsNotUpdated = errors.New("activation link is not updated")
	ErrGetActivationLink          = errors.New("error getting activation link")
	ErrSaveActivationLink         = errors.New("error saving activation link")
)
var (
	ErrTokenNotFound = errors.New("token not found")
)

var (
	ErrReferralNotFound = errors.New("referral not found")
)
