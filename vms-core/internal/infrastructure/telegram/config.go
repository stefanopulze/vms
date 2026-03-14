package telegram

type Config struct {
	BotName        string   `yaml:"botName" env:"TELEGRAM_BOT_NAME"`
	BotApiKey      string   `yaml:"botApiKey" env:"TELEGRAM_BOT_API_KEY"`
	ChatId         string   `yaml:"chatId" env:"TELEGRAM_CHAT_ID"`
	EnableCommands bool     `yaml:"enabledCommands" env:"TELEGRAM_ENABLE_COMMANDS"`
	EnabledUsers   []string `yaml:"enabledUsers"`
}
