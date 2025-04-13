package referral

import (
	"context"
	"errors"
	"fmt"
	"ia-online-golang/internal/dto"
	"ia-online-golang/internal/storage"

	"github.com/sirupsen/logrus"
)

type ReferralService struct {
	log                *logrus.Logger
	ReferralRepository storage.ReferralRepositoryI
}

type ReferralServiceI interface {
	ReferralsUser(ctx context.Context, referral_code string) ([]dto.ReferralDTO, error)
	UpdateActiveReferrals(ctx context.Context) error
}

func New(log *logrus.Logger, referralRepository storage.ReferralRepositoryI) *ReferralService {
	return &ReferralService{
		log:                log,
		ReferralRepository: referralRepository,
	}
}

func (r *ReferralService) ReferralsUser(ctx context.Context, referral_code string) ([]dto.ReferralDTO, error) {
	op := "ReferralService.Referrals"

	referrals, err := r.ReferralRepository.ReferralsUser(ctx, referral_code)
	if err != nil {
		if errors.Is(err, storage.ErrReferralsNotFound) {
			return []dto.ReferralDTO{}, nil
		}
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	referralsDTO := []dto.ReferralDTO{}

	for _, ref := range referrals {
		referralDTO := dto.ReferralDTO{
			ID:          ref.Referral.ID,
			Name:        ref.User.Name,
			PhoneNumber: ref.User.PhoneNumber,
			City:        ref.User.City,
			Active:      ref.Referral.Active,
		}
		referralsDTO = append(referralsDTO, referralDTO)
	}

	return referralsDTO, nil
}

func (r *ReferralService) UpdateActiveReferrals(ctx context.Context) error {
	op := "ReferralService.UpdateActive"

	referrals, err := r.ReferralRepository.GetInactiveReferralsWithReadyLeads(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, referral := range referrals {
		err = r.ReferralRepository.UpdateActive(ctx, referral.ID, true)
		if err != nil {
			r.log.Error(err)
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}
