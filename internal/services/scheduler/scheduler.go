package scheduler

import (
	"context"
	"ia-online-golang/internal/services/referral"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type SchedulerService struct {
	log             *logrus.Logger
	ReferralService referral.ReferralServiceI
	cron            *cron.Cron
}

type SchedulerServiceI interface {
	Run()
	Stop()
}

func New(log *logrus.Logger, referralService referral.ReferralServiceI) *SchedulerService {
	return &SchedulerService{
		log:             log,
		ReferralService: referralService,
		cron:            cron.New(),
	}
}

func (s *SchedulerService) Run() {
	op := "SchedulerService.Run"

	// Пример: запуск каждый день в 3:00 ночи
	_, err := s.cron.AddFunc("*/10 * * * * *", func() {
		ctx := context.Background()
		err := s.ReferralService.UpdateActiveReferrals(ctx)
		if err != nil {
			s.log.Errorf("%s:%v", op, err)
		} else {
			s.log.Infof("%s: успешно обновлены активные рефералы", op)
		}
	})

	if err != nil {
		s.log.Fatalf("%s:%v", op, err)
	}

	s.cron.Start()
	s.log.Info("⏱️ Планировщик запущен")
}

func (s *SchedulerService) Stop() {
	s.cron.Stop()
	s.log.Info("🛑 Планировщик остановлен")
}
