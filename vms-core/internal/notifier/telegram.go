package notifier

import (
	"context"
	"vms-core/internal/infrastructure/telegram"
)

var _ Notifier = (*telegramNotifier)(nil)

type telegramNotifier struct {
	client *telegram.Client
}

func NewTelegram(client *telegram.Client) Notifier {
	return &telegramNotifier{
		client: client,
	}
}

func (t telegramNotifier) Send(ctx context.Context, message string) error {
	_, err := t.client.Send(ctx, message)
	return err
}
