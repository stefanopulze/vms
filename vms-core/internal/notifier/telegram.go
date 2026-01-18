package notifier

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

var _ Notifier = (*telegram)(nil)

type telegram struct {
	apiKey string
	chatId string
	botUrl string
}

func NewTelegram(cfg TelegramConfig) Notifier {
	return &telegram{
		botUrl: fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", cfg.BotApiKey),
		apiKey: cfg.BotApiKey,
		chatId: cfg.ChatId,
	}
}

func (t telegram) Name() string {
	return "telegram"
}

func (t telegram) Send(ctx context.Context, message string) error {
	data := url.Values{}
	data.Set("chat_id", t.chatId)
	data.Set("text", message)

	// Send POST request
	resp, err := http.PostForm(t.botUrl, data)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
