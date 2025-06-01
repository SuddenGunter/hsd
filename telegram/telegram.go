package telegram

import (
	"fmt"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Notifier sends messages to a telegram chat.
type Notifier struct {
	chatID int64
	bot    *tgbotapi.BotAPI

	l *slog.Logger
}

// NewNotifier returns a new Notifier.
func NewNotifier(tgBotToken string, chatID int64, l *slog.Logger) (*Notifier, error) {
	bot, err := tgbotapi.NewBotAPI(tgBotToken)
	if err != nil {
		return nil, fmt.Errorf("new notifier: %w", err)
	}

	return &Notifier{chatID: chatID, bot: bot, l: l}, nil
}

// Notify sends a message to specific telegram chat about the alarm event.
func (n *Notifier) Notify(device, msg string) {
	payload := fmt.Sprintf("🚨 %s: %s", device, msg)
	tgMsg := tgbotapi.NewMessage(n.chatID, payload)

	_, err := n.bot.Send(tgMsg)
	if err != nil {
		n.l.Error("telegram message delivery failed", "err", err)
		return
	}
}
