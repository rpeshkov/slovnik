package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/rpeshkov/slovnik"

	"github.com/pkg/errors"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type slovnikClient struct {
	client  *http.Client
	baseURL *url.URL
}

// NewSlovnikClient creates new client for accessing slovnik web server
func NewSlovnikClient(baseURL string, httpClient *http.Client) (*slovnikClient, error) {
	var c *http.Client

	if httpClient == nil {
		c = &http.Client{Timeout: 10 * time.Second}
	} else {
		c = http.DefaultClient
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	return &slovnikClient{c, u}, nil
}

// Translate word
func (c *slovnikClient) Translate(word string, language slovnik.Language) ([]*slovnik.Word, error) {
	const methodURL = "/api/translate/"
	u := *c.baseURL
	u.Path = path.Join(u.Path, methodURL, word)
	r, err := c.client.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Got bad status (%d) from server", r.StatusCode)
	}

	w := []*slovnik.Word{}
	json.NewDecoder(r.Body).Decode(&w)
	return w, nil
}

// Bot type aggregate all bot logic
type Bot struct {
	api *tgbotapi.BotAPI

	// updates is a channel that's used to communicate with telegram server
	updates tgbotapi.UpdatesChannel

	templates *Template

	translator slovnik.Translator
}

// NewBot creates and initializes new bot
func NewBot(config *Config, templates *Template) (*Bot, error) {
	botAPI, err := tgbotapi.NewBotAPI(config.BotID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init bot")
	}

	slovnikClient, err := NewSlovnikClient(config.SlovnikURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init slovnikClient")
	}

	var updates tgbotapi.UpdatesChannel

	if config.IsWebhook() {
		log.Printf("WebhookHost is set. URL: %s. Using webhooks\n", config.WebhookURL)

		webHook := tgbotapi.NewWebhook(config.WebhookURL)
		if _, err = botAPI.SetWebhook(webHook); err != nil {
			return nil, errors.Wrap(err, "webhook set failed")
		}

		updates = botAPI.ListenForWebhook("/bot" + botAPI.Token)

		go http.ListenAndServe("0.0.0.0:8080", nil)
	} else {
		log.Println("WebhookURL environment variable is not set. Using polling.")

		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		updates, err = botAPI.GetUpdatesChan(u)
		if err != nil {
			return nil, errors.Wrap(err, "unable to get updates channel")
		}
	}

	return &Bot{botAPI, updates, templates, slovnikClient}, nil
}

// Listen start listening on message updates and calling provided handler for processing incoming messages
func (bot *Bot) Listen() {
	for update := range bot.updates {
		if update.Message != nil {
			bot.handleMessage(&update)
		} else if update.CallbackQuery != nil {
			bot.handleCallbackQuery(&update)
		}
	}
}

func (bot *Bot) handleMessage(update *tgbotapi.Update) {
	words, err := bot.translator.Translate(update.Message.Text, slovnik.Cz)
	if err != nil {
		bot.respondError(update.Message.Chat.ID, "Something bad happened :(")
		log.Println(err)
		return
	}

	messageText := bot.templates.Translation(words)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
	msg.ParseMode = tgbotapi.ModeMarkdown

	hasPhrases := len(words) == 1 && len(words[0].Samples) > 0
	if hasPhrases {
		keyboard := bot.addMessageKeyboard(words)
		if keyboard != nil {
			msg.ReplyMarkup = keyboard
		}
	}

	_, err = bot.api.Send(msg)

	if err != nil {
		log.Println(err)
	}
}

// respondError writes an error to the chat
func (bot *Bot) respondError(updateID int64, text string) {
	msg := tgbotapi.NewMessage(updateID, text)
	_, err := bot.api.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func (bot *Bot) handleCallbackQuery(update *tgbotapi.Update) {
	callbackData := update.CallbackQuery.Data
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID

	if strings.HasPrefix(callbackData, "phrases:") {
		w := strings.TrimPrefix(callbackData, "phrases:")

		words, err := bot.translator.Translate(w, slovnik.Cz)
		if err != nil {
			bot.respondError(chatID, "Error occured when I tried to get phrases :(")
			log.Println(err)
			return
		}

		messageText := bot.templates.Phrases(words)

		msg := tgbotapi.NewMessage(chatID, messageText)
		msg.ParseMode = tgbotapi.ModeMarkdown

		_, err = bot.api.Send(msg)
		if err != nil {
			log.Println(err)
			return
		}

		editMsg := tgbotapi.NewEditMessageText(chatID, messageID, bot.templates.Translation(words))
		editMsg.ReplyMarkup = nil
		editMsg.ParseMode = tgbotapi.ModeMarkdown
		_, err = bot.api.Send(editMsg)

		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (bot *Bot) addMessageKeyboard(words []*slovnik.Word) *tgbotapi.InlineKeyboardMarkup {
	if words == nil || len(words) > 1 || len(words[0].Samples) <= 0 {
		return nil
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Фразы", "phrases:"+words[0].Word),
		),
	)

	return &keyboard
}
