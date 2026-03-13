package notifier

import (
	"context"
	"fmt"
	"log/slog"
)

type Notifier interface {
	Send(ctx context.Context, message string) error
}

func NewNotify(notifiers ...Notifier) *Notify {
	return &Notify{
		notifiers: notifiers,
	}
}

type Notify struct {
	notifiers []Notifier
}

func (n Notify) Send(ctx context.Context, message string) error {
	for _, notifier := range n.notifiers {
		go func() {
			if err := notifier.Send(ctx, message); err != nil {
				slog.Error(fmt.Sprintf("notifier: %s", err))
			}
		}()
	}

	return nil
}
