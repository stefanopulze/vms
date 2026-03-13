package telegram

type SendMsg struct {
	Text   string
	Markup any
}

type MessageResponse struct {
	OK     bool `json:"ok"`
	Result struct {
		MessageID int64 `json:"message_id"`
	} `json:"result"`
}

type UpdateResponse struct {
	OK     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	UpdateID      int            `json:"update_id"`
	Message       *Message       `json:"message"`
	CallbackQuery *CallbackQuery `json:"callback_query"`
}

func (u Update) ChatId() int64 {
	if u.Message != nil {
		return u.Message.Chat.ID
	} else if u.CallbackQuery != nil {
		return u.CallbackQuery.From.ID
	}

	return 0
}

func (u Update) MessageId() int64 {
	if u.Message != nil {
		return u.Message.ID
	} else if u.CallbackQuery != nil {
		return u.CallbackQuery.Message.ID
	}

	return 0
}

func (u Update) Username() string {
	if u.Message != nil {
		return u.Message.From.Username
	} else if u.CallbackQuery != nil {
		return u.CallbackQuery.From.Username
	}

	return ""
}

type Message struct {
	ID   int64  `json:"message_id"`
	Text string `json:"text"`
	From User   `json:"from"`
	Date int64  `json:"date"`
	Chat struct {
		ID int64 `json:"id"`
	} `json:"chat"`
}

type CallbackQuery struct {
	ID      string   `json:"id"`
	Data    string   `json:"data"`
	From    User     `json:"from"`
	Message *Message `json:"message"`
}

type User struct {
	ID        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
	Username  string `json:"username"`
}

type ReplyKeyboardMarkup struct {
	Keyboard              [][]KeyboardButton `json:"keyboard"`
	ResizeKeyboard        bool               `json:"resize_keyboard"`
	OneTimeKeyboard       bool               `json:"one_time_keyboard"`
	InputFieldPlaceholder string             `json:"input_field_placeholder"`
}

type KeyboardButton struct {
	Text  string `json:"text"`
	Style string `json:"style,omitempty"` // danger | success | primary
}

type InlineKeyboardMarkup struct {
	Keyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}
