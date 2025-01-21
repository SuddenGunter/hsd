package telegram

import (
	"log"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Notifier sends messages to a telegram chat.
type Notifier struct {
	tgBotToken string
	l          *slog.Logger
}

// NewNotifier returns a new Notifier.
func NewNotifier(tgBotToken string, l *slog.Logger) *Notifier {
	return &Notifier{tgBotToken: tgBotToken, l: l}
}

// Notify sends a message to specific telegram chat about the alarm event.
func (n *Notifier) Notify(device, msg string) {
	n.l.Info("telegram message sent", "device", device, "msg", msg)
}

func (n *Notifier) Listen() {
	bot, err := tgbotapi.NewBotAPI(n.tgBotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	go func() {
		for update := range updates {
			if update.Message != nil { // If we got a message
				n.l.Debug("telegram message received", "user", update.Message.From.UserName, "msg", update.Message.Text)
			}
		}
	}()

	// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	// msg.ReplyToMessageID = update.Message.MessageID

	// bot.Send(msg)
}
