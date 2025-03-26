package payloads

type PayloadUserAccess struct {
	UserID       int64       `json:"user_id"`
	Name         string      `json:"name"`
	Email        string      `json:"email"`
	PhoneNumber  string      `json:"phone_number"`
	City         string      `json:"city"`
	Telegram     string      `json:"telegram"`
	ReferralCode string      `json:"referral_code"`
	Partners     []Partner   `json:"partners"`
	Statistics   []Statistic `json:"statistics"`
}
type PayloadUserRefresh struct {
	UserID int64 `json:"user_id"`
}

type Partner struct {
	ID          int
	City        string
	PhoneNumber string
	Status      bool
	Level       string
}

type Statistic struct {
	Name  string
	Value int
}
