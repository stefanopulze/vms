package notifier

type TelegramConfig struct {
	BotApiKey string `yaml:"botApiKey" env:"TELEGRAM_BOT_API_KEY"`
	ChatId    string `yaml:"chatId" env:"TELEGRAM_CHAT_ID"`
}
