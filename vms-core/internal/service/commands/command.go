package commands

import "vms-core/internal/infrastructure/telegram"

type Command interface {
	GetPattern() string
	StartNewSession() (int64, error)
	HandleCallback(*telegram.CallbackQuery) error
	NeedCallback() bool
}
