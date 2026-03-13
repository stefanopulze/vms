package commands

import (
	"context"
	"errors"
	"time"
	"vms-core/internal/infrastructure/telegram"
	"vms-core/internal/voltronic"
)

var _ Command = (*UpdateSourcePriority)(nil)

func NewUpdateSourcePriority(tc *telegram.Client, inverter *voltronic.Client) *UpdateSourcePriority {
	return &UpdateSourcePriority{
		telegram: tc,
		inverter: inverter,
	}

}

type UpdateSourcePriority struct {
	telegram *telegram.Client
	inverter *voltronic.Client
}

func (c *UpdateSourcePriority) StartNewSession() (int64, error) {
	keyboard := telegram.InlineKeyboardMarkup{
		Keyboard: [][]telegram.InlineKeyboardButton{
			{
				{Text: "SBU", CallbackData: "sbu"},
				{Text: "SUB", CallbackData: "sub"},
				{Text: "USB", CallbackData: "usb"},
			},
		},
	}

	return c.telegram.SendWithMarkup(context.Background(), "Choose output source", keyboard)
}

func (c *UpdateSourcePriority) GetPattern() string {
	return "/update_source_priority"
}

func (c *UpdateSourcePriority) HandleCallback(callback *telegram.CallbackQuery) error {
	// TODO move into remote_commands.go? Is useful for all callbacks?
	if time.Unix(callback.Message.Date, 0).Add(5 * time.Second).Before(time.Now()) {
		c.telegram.AnswerCallback(context.Background(), callback.ID, "Timeout")
		_ = c.telegram.EditMessage(context.Background(), callback.Message.ID, "Message is too old. Please try again.")
		return errors.New("message is too old")
	}

	if !isValidSourcePriority(callback.Data) {
		return errors.New("invalid source priority")
	}

	if err := c.inverter.UpdateSourcePriority(callback.Data); err != nil {
		c.telegram.AnswerCallback(context.Background(), callback.ID, "Something went wrong")
		_ = c.telegram.EditMessage(context.Background(), callback.Message.ID, "Something went wrong. Please try again.")
		return err
	}

	c.telegram.AnswerCallback(context.Background(), callback.ID, "Source priority updated")
	return c.telegram.EditMessage(context.Background(), callback.Message.ID, "Source priority updated")
}

func isValidSourcePriority(v string) bool {
	return v == "sbu" || v == "sub" || v == "usb"
}

func (c *UpdateSourcePriority) NeedCallback() bool {
	return true
}
