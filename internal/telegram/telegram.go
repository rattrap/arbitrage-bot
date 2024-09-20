package telegram

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// TelegramService is a struct to hold the bot API and the chat ID
type TelegramService struct {
	isStarted bool
	bot       *tgbotapi.BotAPI
	chatID    int64
}

// NewTelegramService initializes a new TelegramService
func NewTelegramService(token string, chatID int64) *TelegramService {
	isStarted := true
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		isStarted = false
	}

	return &TelegramService{
		isStarted: isStarted,
		bot:       bot,
		chatID:    chatID,
	}
}

// SendMessage sends a message to the configured Telegram chat
func (ts *TelegramService) SendMessage(message string) error {
	if !ts.isStarted {
		return nil
	}
	msg := tgbotapi.NewMessage(ts.chatID, message)
	_, err := ts.bot.Send(msg)
	if err != nil {
		log.Printf("Failed to send message to Telegram: %v", err)
		return err
	}
	return nil
}

// FormatMessage formats a message with the given arguments
func FormatMessage(message string) string {
	const replacement = "\n"

	var replacer = strings.NewReplacer(
		", ", replacement,
	)
	return replacer.Replace(message)
}
