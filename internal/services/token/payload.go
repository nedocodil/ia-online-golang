package token

import "ia-online-golang/internal/dto"

type PayloadUserAccess struct {
	UserID       int64             `json:"user_id"`
	Roles        []string          `json:"roles"`
	Name         string            `json:"name"`
	Email        string            `json:"email"`
	PhoneNumber  string            `json:"phone_number"`
	City         string            `json:"city"`
	Telegram     string            `json:"telegram"`
	ReferralCode string            `json:"referral_code"`
	Referrals    []dto.ReferralDTO `json:"referrals"`
	Statistic    dto.UserStatistic `json:"statistic"`
}
type PayloadUserRefresh struct {
	UserID int64 `json:"user_id"`
}
