package tgbot

import (
	"context"
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

type Client struct {
	bot *gotgbot.Bot
}

func NewClient(token string) (*Client, error) {
	bot, err := gotgbot.NewBot(token, nil)
	if err != nil {
		return nil, fmt.Errorf("new tg bot: %w", err)
	}

	return &Client{
		bot: bot,
	}, nil
}

func (c *Client) SetWebhook(webhookURL string) error {
	_, err := c.bot.SetWebhook(webhookURL, nil)
	if err != nil {
		return fmt.Errorf("set tg bot webhook: %w", err)
	}

	return nil
}

func (c *Client) SendTextMessage(_ context.Context, chatID int64, text string) error {
	if _, err := c.bot.SendMessage(chatID, text, nil); err != nil {
		return fmt.Errorf("send text message: %w", err)
	}

	return nil
}
