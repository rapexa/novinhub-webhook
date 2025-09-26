package main

import (
	"log"
	"strconv"
	"time"

	"novinhub-webhook/internal/config"
	"novinhub-webhook/internal/server"
	"novinhub-webhook/pkg/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	BotToken = "8205935967:AAEI2jb_0y0-TlYZ_5gA2wF4cyIr7eaHYuU"
)

// Admin IDs - لیست ادمین‌های مجاز
var AdminIDs = map[int64]string{
	76599340:  "Admin Original",      // ادمین اصلی
	110435852: "MahYaR (@Saeidpour)", // ادمین جدید
}

// isAdmin بررسی می‌کند که آیا کاربر ادمین است یا نه
func isAdmin(userID int64) (bool, string) {
	if name, exists := AdminIDs[userID]; exists {
		return true, name
	}
	return false, ""
}

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize logger
	logger := logger.New()

	// Create server
	srv := server.New(cfg, logger)

	// Start webhook server in a goroutine
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatal("Server failed to start:", err)
		}
	}()

	// Start Telegram bot
	startTelegramBot(cfg, logger)

	// Wait forever
	select {}
}

func startTelegramBot(cfg *config.Config, logger *logger.Logger) {
	// Initialize bot
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		logger.Error("Failed to initialize Telegram bot", "error", err)
		return
	}

	bot.Debug = false
	logger.Info("Telegram bot authorized", "username", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Handle updates in a goroutine
	go func() {
		for update := range updates {
			if update.Message != nil {
				handleMessage(bot, update.Message, cfg, logger)
			} else if update.CallbackQuery != nil {
				handleCallbackQuery(bot, update.CallbackQuery, cfg, logger)
			}
		}
	}()
}

func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, cfg *config.Config, logger *logger.Logger) {
	// Check if message is from admin
	isAdminUser, adminName := isAdmin(message.From.ID)
	if !isAdminUser {
		logger.Info("🚫 Unauthorized access attempt",
			"user_id", message.From.ID,
			"username", message.From.UserName,
			"first_name", message.From.FirstName,
			"available_admins", len(AdminIDs))
		return
	}

	// Log admin access
	logger.Info("✅ Admin access granted",
		"user_id", message.From.ID,
		"username", message.From.UserName,
		"first_name", message.From.FirstName,
		"admin_name", adminName,
		"message", message.Text)

	switch message.Text {
	case "/start":
		sendMainMenu(bot, message.Chat.ID)
	case "📱 پترن امروز":
		showCurrentPattern(bot, message.Chat.ID, cfg)
	case "➡️ برو به پترن بعدی":
		nextPattern(bot, message.Chat.ID, cfg)
	case "📋 لیست پترن‌ها":
		showPatternsList(bot, message.Chat.ID, cfg)
	case "👥 لیست ادمین‌ها":
		showAdminsList(bot, message.Chat.ID)
	default:
		sendMainMenu(bot, message.Chat.ID)
	}
}

func handleCallbackQuery(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery, cfg *config.Config, logger *logger.Logger) {
	// Check if callback is from admin
	isAdminUser, adminName := isAdmin(callbackQuery.From.ID)
	if !isAdminUser {
		logger.Info("🚫 Unauthorized callback attempt",
			"user_id", callbackQuery.From.ID,
			"username", callbackQuery.From.UserName,
			"first_name", callbackQuery.From.FirstName,
			"available_admins", len(AdminIDs))
		// Answer callback query with error
		bot.Request(tgbotapi.NewCallback(callbackQuery.ID, "❌ دسترسی غیرمجاز"))
		return
	}

	// Log admin callback access
	logger.Info("✅ Admin callback access granted",
		"user_id", callbackQuery.From.ID,
		"username", callbackQuery.From.UserName,
		"first_name", callbackQuery.From.FirstName,
		"admin_name", adminName,
		"callback_data", callbackQuery.Data)

	switch callbackQuery.Data {
	case "current_pattern":
		showCurrentPattern(bot, callbackQuery.Message.Chat.ID, cfg)
	case "next_pattern":
		nextPattern(bot, callbackQuery.Message.Chat.ID, cfg)
	case "list_patterns":
		showPatternsList(bot, callbackQuery.Message.Chat.ID, cfg)
	case "list_admins":
		showAdminsList(bot, callbackQuery.Message.Chat.ID)
	}

	// Answer callback query
	bot.Request(tgbotapi.NewCallback(callbackQuery.ID, ""))
}

func sendMainMenu(bot *tgbotapi.BotAPI, chatID int64) {
	text := "🤖 ربات مدیریت پترن‌های SMS\n"
	text += "👋 سلام ! خوش آمدید\n"
	text += "🔒 دسترسی امنیتی فعال\n\n"
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
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👥 لیست ادمین‌ها", "list_admins"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func showCurrentPattern(bot *tgbotapi.BotAPI, chatID int64, cfg *config.Config) {
	pattern, index, groupName := cfg.GetCurrentPatternInfo()

	text := "📱 پترن فعلی:\n\n"
	text += "🔹 گروه: " + groupName + "\n"
	text += "🔹 شماره: " + strconv.Itoa(index) + " از 4\n"
	text += "🔹 کد پترن: `" + pattern + "`\n\n"
	text += "⏰ آخرین تغییر: " + time.Now().Format("2006-01-02 15:04:05")

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func nextPattern(bot *tgbotapi.BotAPI, chatID int64, cfg *config.Config) {
	pattern, index, groupName := cfg.NextPattern()

	text := "✅ پترن تغییر کرد!\n\n"
	text += "🔹 گروه جدید: " + groupName + "\n"
	text += "🔹 شماره: " + strconv.Itoa(index) + " از 4\n"
	text += "🔹 کد پترن: `" + pattern + "`\n\n"
	text += "⏰ زمان تغییر: " + time.Now().Format("2006-01-02 15:04:05")

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func showPatternsList(bot *tgbotapi.BotAPI, chatID int64, cfg *config.Config) {
	patterns := cfg.GetPatternsList()

	text := "📋 لیست تمام پترن‌ها:\n\n"

	for _, p := range patterns {
		status := "❌"
		if p["is_current"].(bool) {
			status = "✅ فعلی"
		}

		text += status + " " + p["name"].(string) + " (" + strconv.Itoa(p["index"].(int)) + "): `" + p["pattern"].(string) + "`\n"
	}

	text += "\n⏰ آخرین تغییر: " + time.Now().Format("2006-01-02 15:04:05")

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func showAdminsList(bot *tgbotapi.BotAPI, chatID int64) {
	text := "👥 لیست ادمین‌های سیستم:\n\n"

	for userID, name := range AdminIDs {
		text += "🔹 " + name + "\n"
		text += "   ID: `" + strconv.FormatInt(userID, 10) + "`\n\n"
	}

	text += "📊 تعداد کل ادمین‌ها: " + strconv.Itoa(len(AdminIDs))

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}
