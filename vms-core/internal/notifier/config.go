package notifier

type TelegramConfig struct {
	BotApiKey string `env:"TELEGRAM_BOT_API_KEY"`
	ChatId    string `env:"TELEGRAM_CHAT_ID"`
}
