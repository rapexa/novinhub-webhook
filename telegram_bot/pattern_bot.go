package main

import (
	"fmt"
	"log"
	"time"

	"novinhub-webhook/internal/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	BotToken = "8205935967:AAEI2jb_0y0-TlYZ_5gA2wF4cyIr7eaHYuU"
)

var appConfig *config.Config

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Panic("Failed to load configuration:", err)
	}
	appConfig = cfg

	// Initialize bot
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			handleMessage(bot, update.Message)
		} else if update.CallbackQuery != nil {
			handleCallbackQuery(bot, update.CallbackQuery)
		}
	}
}

func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	// Check if message is from admin (you can add admin user ID check here)

	switch message.Text {
	case "/start":
		sendMainMenu(bot, message.Chat.ID)
	case "📱 پترن امروز":
		showCurrentPattern(bot, message.Chat.ID)
	case "➡️ برو به پترن بعدی":
		nextPattern(bot, message.Chat.ID)
	case "📋 لیست پترن‌ها":
		showPatternsList(bot, message.Chat.ID)
	default:
		sendMainMenu(bot, message.Chat.ID)
	}
}

func sendMainMenu(bot *tgbotapi.BotAPI, chatID int64) {
	text := "🤖 ربات مدیریت پترن‌های SMS\n\n"
	text += "لطفاً یکی از گزینه‌های زیر را انتخاب کنید:"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📱 پترن امروز", "current_pattern"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➡️ برو به پترن بعدی", "next_pattern"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📋 لیست پترن‌ها", "list_patterns"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func showCurrentPattern(bot *tgbotapi.BotAPI, chatID int64) {
	pattern, index, groupName := appConfig.GetCurrentPatternInfo()

	text := fmt.Sprintf("📱 پترن فعلی:\n\n")
	text += fmt.Sprintf("🔹 گروه: %s\n", groupName)
	text += fmt.Sprintf("🔹 شماره: %d از 4\n", index)
	text += fmt.Sprintf("🔹 کد پترن: `%s`\n\n", pattern)
	text += fmt.Sprintf("⏰ آخرین تغییر: %s", time.Now().Format("2006-01-02 15:04:05"))

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func nextPattern(bot *tgbotapi.BotAPI, chatID int64) {
	pattern, index, groupName := appConfig.NextPattern()

	text := fmt.Sprintf("✅ پترن تغییر کرد!\n\n")
	text += fmt.Sprintf("🔹 گروه جدید: %s\n", groupName)
	text += fmt.Sprintf("🔹 شماره: %d از 4\n", index)
	text += fmt.Sprintf("🔹 کد پترن: `%s`\n\n", pattern)
	text += fmt.Sprintf("⏰ زمان تغییر: %s", time.Now().Format("2006-01-02 15:04:05"))

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func showPatternsList(bot *tgbotapi.BotAPI, chatID int64) {
	patterns := appConfig.GetPatternsList()

	text := "📋 لیست تمام پترن‌ها:\n\n"

	for _, p := range patterns {
		status := "❌"
		if p["is_current"].(bool) {
			status = "✅ فعلی"
		}

		text += fmt.Sprintf("%s %s (%d): `%s`\n",
			status,
			p["name"],
			p["index"],
			p["pattern"])
	}

	text += fmt.Sprintf("\n⏰ آخرین تغییر: %s", time.Now().Format("2006-01-02 15:04:05"))

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

// Handle callback queries for inline keyboards
func handleCallbackQuery(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	switch callbackQuery.Data {
	case "current_pattern":
		showCurrentPattern(bot, callbackQuery.Message.Chat.ID)
	case "next_pattern":
		nextPattern(bot, callbackQuery.Message.Chat.ID)
	case "list_patterns":
		showPatternsList(bot, callbackQuery.Message.Chat.ID)
	}

	// Answer callback query
	bot.Request(tgbotapi.NewCallback(callbackQuery.ID, ""))
}
