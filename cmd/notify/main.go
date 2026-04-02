package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Sender delivers a notification message.
type Sender interface {
	Send(chatID, message string) error
}

func main() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr, os.Getenv, nil))
}

// Run is the testable CLI entry point.
func Run(args []string, stdout, stderr io.Writer, getenv func(string) string, sender Sender) int {
	fs := flag.NewFlagSet("notify", flag.ContinueOnError)
	fs.SetOutput(stderr)

	token := fs.String("telegram-bot-token", "", "Telegram bot token (or TELEGRAM_BOT_TOKEN env)")
	chatID := fs.String("telegram-chat-id", "", "Telegram chat ID (or TELEGRAM_CHAT_ID env)")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if *token == "" {
		*token = getenv("TELEGRAM_BOT_TOKEN")
	}
	if *chatID == "" {
		*chatID = getenv("TELEGRAM_CHAT_ID")
	}

	if *token == "" {
		fmt.Fprintln(stderr, "error: Telegram bot token is required (--telegram-bot-token or TELEGRAM_BOT_TOKEN)")
		return 1
	}
	if *chatID == "" {
		fmt.Fprintln(stderr, "error: Telegram chat ID is required (--telegram-chat-id or TELEGRAM_CHAT_ID)")
		return 1
	}

	message := strings.TrimSpace(strings.Join(fs.Args(), " "))
	if message == "" {
		fmt.Fprintln(stderr, "error: notification message is required")
		return 1
	}

	if sender == nil {
		sender = newTelegramSender(http.DefaultClient, *token)
	}

	if err := sender.Send(*chatID, message); err != nil {
		fmt.Fprintf(stderr, "error: %s\n", err)
		return 2
	}

	fmt.Fprintln(stdout, "notification sent")
	return 0
}

type telegramSender struct {
	client  *http.Client
	baseURL string
	token   string
}

func newTelegramSender(client *http.Client, token string) *telegramSender {
	if client == nil {
		client = http.DefaultClient
	}
	return &telegramSender{
		client:  client,
		baseURL: "https://api.telegram.org",
		token:   token,
	}
}

func (s *telegramSender) Send(chatID, message string) error {
	form := url.Values{}
	form.Set("chat_id", chatID)
	form.Set("text", message)

	resp, err := s.client.PostForm(fmt.Sprintf("%s/bot%s/sendMessage", strings.TrimRight(s.baseURL, "/"), s.token), form)
	if err != nil {
		return fmt.Errorf("sending Telegram message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("sending Telegram message: unexpected status %s", resp.Status)
	}

	return nil
}
