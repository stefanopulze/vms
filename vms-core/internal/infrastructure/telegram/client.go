package telegram

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"
)

func NewClient(cfg Config) *Client {
	chatId, _ := strconv.ParseInt(cfg.ChatId, 10, 64)

	return &Client{
		botName:         cfg.BotName,
		apiKey:          cfg.BotApiKey,
		chatId:          chatId,
		enabledCommands: cfg.EnableCommands,
		enabledUsers:    cfg.Users,
		baseBotUrl:      fmt.Sprintf("https://api.telegram.org/bot%s", cfg.BotApiKey),
		lastUpdateId:    -1,
	}
}

type Client struct {
	botName         string
	apiKey          string
	chatId          int64
	baseBotUrl      string
	lastUpdateId    int
	enabledCommands bool
	enabledUsers    []string
}

func (tc Client) send(_ context.Context, msg string, replyMarkup any) (int64, error) {
	payload := url.Values{}
	payload.Set("chat_id", fmt.Sprintf("%d", tc.chatId))
	payload.Set("text", msg)

	if replyMarkup != nil {
		rmj, err := json.Marshal(replyMarkup)
		if err != nil {
			return 0, errors.Join(errors.New("failed to marshal reply markup"), err)
		}
		payload.Set("reply_markup", string(rmj))
	}

	// Send POST request
	sendUrl := fmt.Sprintf("%s/sendMessage", tc.baseBotUrl)
	resp, err := http.PostForm(sendUrl, payload)
	if err != nil {
		return 0, err
	}
	defer func() { _ = resp.Body.Close() }()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected network error: %d", resp.StatusCode)
	}

	telegramResponse := new(MessageResponse)
	if err = json.NewDecoder(resp.Body).Decode(telegramResponse); err != nil {
		return 0, err
	}
	if !telegramResponse.OK {
		return 0, fmt.Errorf("unable to send message: %s", telegramResponse)
	}

	return telegramResponse.Result.MessageID, nil
}

func (tc Client) SendMsg(ctx context.Context, msg SendMsg) (int64, error) {
	return tc.send(ctx, msg.Text, msg.Markup)
}

func (tc Client) Send(ctx context.Context, message string) (int64, error) {
	return tc.send(ctx, message, nil)
}

func (tc Client) SendWithMarkup(ctx context.Context, msg string, markup any) (int64, error) {
	return tc.send(ctx, msg, markup)
}

func (tc Client) EditMessage(_ context.Context, messageID int64, msg string) error {
	payload := url.Values{}
	payload.Set("chat_id", fmt.Sprintf("%d", tc.chatId))
	payload.Set("message_id", fmt.Sprintf("%d", messageID))
	payload.Set("text", msg)

	apiUrl := fmt.Sprintf("%s/editMessageText", tc.baseBotUrl)
	resp, err := http.PostForm(apiUrl, payload)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	return nil
}

func (tc Client) GetUpdates(ctx context.Context, handler func(Update)) {
	client := &http.Client{
		Timeout: 40 * time.Second,
	}

	for {
		select {
		case <-ctx.Done():
			slog.Info("telegram: context done")
			return

		default:
			updates, err := tc.fetchUpdates(client, tc.lastUpdateId+1)
			if err != nil {
				slog.Error("telegram: failed to fetch updates", "error", err)
				time.Sleep(5 * time.Second)
				continue
			}

			for i := 0; i < len(updates); i++ {
				tc.lastUpdateId = updates[i].UpdateID
				go handler(updates[i])
			}
		}
	}
}

func (tc Client) fetchUpdates(client *http.Client, offset int) ([]Update, error) {
	path := fmt.Sprintf("%s/getUpdates?offset=%d&timeout=30", tc.baseBotUrl, offset)

	resp, err := client.Get(path)
	if err != nil {
		return nil, err
	}
	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(resp.Body)

	var response UpdateResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Result, nil
}

func (tc Client) ValidateMessage(chatId int64, username string) error {
	if tc.chatId != chatId {
		return fmt.Errorf("chat id %d is not valid", chatId)
	}

	if len(tc.enabledUsers) == 0 {
		return nil
	}

	if !slices.Contains(tc.enabledUsers, username) {
		return fmt.Errorf("user %s is not authorized", username)
	}

	return nil
}

func (tc Client) ExtractCommand(cmd string) (string, error) {
	if len(cmd) == 0 {
		return "", fmt.Errorf("empty command")
	}

	if !strings.Contains(cmd, tc.botName) {
		return "", fmt.Errorf("command is not for this bot")
	}

	return strings.TrimSuffix(cmd, tc.botName), nil
}

func (tc Client) AnswerCallback(_ context.Context, callbackQueryId string, toastMsg string) {
	apiUrl := fmt.Sprintf("%s/answerCallbackQuery", tc.baseBotUrl)
	payload := url.Values{}
	payload.Set("callback_query_id", callbackQueryId)
	payload.Set("text", toastMsg)
	payload.Set("show_alert", "false")

	response, err := http.PostForm(apiUrl, payload)
	if err != nil {
		slog.Error("telegram: failed to answer callback query", "error", err)
		return
	}

	defer func() { _ = response.Body.Close() }()
}
