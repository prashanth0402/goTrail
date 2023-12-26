package Telegram

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// func TelegramApi() {
// 	log.Println("TelegramApi(+)")
// 	// Replace with your bot token
// 	bot, err := tgbotapi.NewBotAPI("Your:Token")
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	// Set up an update configuration
// 	u := tgbotapi.NewUpdate(0)
// 	u.Timeout = 60

// 	// Get updates from the Telegram API
// 	updates, err := bot.GetUpdatesChan(u)

// 	for update := range updates {
// 		if update.Message == nil {
// 			continue
// 		}
// 		// Extract the chat ID from the incoming message
// 		chatID := update.Message.Chat.ID

// 		// Extract the user's message
// 		userMessage := "Hai"

// 		// Process the user's message and respond accordingly
// 		responseText := processUserMessage(userMessage)

// 		// Create a new message to send as a response
// 		msg := tgbotapi.NewMessage(chatID, responseText)

// 		// Send the response message
// 		_, err := bot.Send(msg)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	}
// 	log.Println("Telegram (-)")
// }

// func processUserMessage(message string) string {
// 	// Process the user's message here
// 	// You can use the message content and chat ID as needed

// 	// In this example, we simply echo back the user's message
// 	return "You said: " + message
// }

func TelegramApi() {
	log.Println("TelegramApi (+)")
	// Replace with your Telegram Bot API token
	bot, err := tgbotapi.NewBotAPI("6178504153:AAGhq_LfIMd_YPQq2r3nBBVbKK5PmQ2ZNwI")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("bot", bot)
	bot.Debug = true // Enable debugging

	log.Printf("Authorized as %s", bot.Self.UserName)

	// Create a new message

	msg := tgbotapi.NewMessage(919025030926, "Hello, Telegram!")

	// Send the message
	_, err = bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}

	updates, err := bot.GetUpdates(tgbotapi.NewUpdate(0))
	if err != nil {
		log.Panic(err)
	}

	// Loop through the received updates
	for _, update := range updates {
		// The Chat ID of the sender
		chatID := update.Message.Chat.ID
		log.Printf("Received a message from chat ID: %d", chatID)
	}
	log.Println("TelegramApi (-)")

}
