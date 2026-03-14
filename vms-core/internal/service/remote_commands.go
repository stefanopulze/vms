package service

import (
	"context"
	"fmt"
	"log/slog"
	"vms-core/internal/infrastructure/telegram"
	"vms-core/internal/service/commands"
)

func NewRemoteCommands(tc *telegram.Client, cmds ...commands.Command) *RemoteCommands {
	commandsMap := make(map[string]commands.Command)
	for _, cmd := range cmds {
		commandsMap[cmd.GetPattern()] = cmd
	}

	return &RemoteCommands{
		telegram:  tc,
		commands:  commandsMap,
		callbacks: make(map[int64]commands.Command),
	}
}

type RemoteCommands struct {
	telegram  *telegram.Client
	commands  map[string]commands.Command
	callbacks map[int64]commands.Command
}

func (rc RemoteCommands) HandleTelegramCommand(update telegram.Update) {
	chatId := update.ChatId()
	username := update.Username()

	slog.Debug(fmt.Sprintf("received telegram update on chatId: %d from: %s", chatId, username))

	if err := rc.telegram.ValidateMessage(chatId, username); err != nil {
		slog.Error(
			"someone tried to send a command to me without being in my chat",
			slog.Any("error", err),
			slog.Any("chat_id", chatId),
		)
		_ = rc.telegram.EditMessage(context.Background(), update.MessageId(), "You are not allowed to send commands to me")
		return
	}

	// User has sent a command
	if update.Message != nil {
		message := update.Message
		slog.Debug(fmt.Sprintf("received telegram command: %s", message.Text))
		cmdPattern, err := rc.telegram.ExtractCommand(update.Message.Text)
		if err != nil {
			slog.Error("cannot extract command", slog.Any("error", err))
			return
		}

		cmd, ok := rc.commands[cmdPattern]
		if !ok {
			slog.Debug(fmt.Sprintf("unknown command: %s", cmdPattern))
			return
		}

		msgId, msgErr := cmd.StartNewSession()
		if msgErr != nil {
			slog.Error("cannot start command", slog.Any("error", msgErr))
			return
		}

		if cmd.NeedCallback() {
			rc.callbacks[msgId] = cmd
		}

	} else if update.CallbackQuery != nil {
		slog.Info(fmt.Sprintf("received telegram callback query: %+v", update.CallbackQuery))
		msgId := update.CallbackQuery.Message.ID

		cmd, ok := rc.callbacks[msgId]
		if !ok {
			slog.Error(fmt.Sprintf("cannot find callback command for message: %d", msgId))
			_ = rc.telegram.EditMessage(context.Background(), msgId, "Please restart command")
			return
		}

		if err := cmd.HandleCallback(update.CallbackQuery); err != nil {
			slog.Error("cannot handle callback", slog.Any("error", err))
		}

		delete(rc.callbacks, msgId)
	}
}
