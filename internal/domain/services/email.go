package services

import (
	"context"
)

type EmailService interface {
	SendEmail(ctx context.Context, toAddress, subject, body string) error
	SendActivationLink(ctx context.Context, toAddress string, activationLink string) error
}
