package main

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/cucumber/godog"
)

type fakeSender struct {
	called  bool
	chatID  string
	message string
	err     error
}

func (f *fakeSender) Send(chatID, message string) error {
	f.called = true
	f.chatID = chatID
	f.message = message
	return f.err
}

type notifyCtx struct {
	envVars  map[string]string
	message  string
	sender   *fakeSender
	exitCode int
	stdout   string
	stderr   string
}

func (nc *notifyCtx) reset() {
	nc.envVars = make(map[string]string)
	nc.message = ""
	nc.sender = &fakeSender{}
	nc.exitCode = -1
	nc.stdout = ""
	nc.stderr = ""
}

func (nc *notifyCtx) telegramCredentialsViaEnvironmentVariables() error {
	nc.envVars["TELEGRAM_BOT_TOKEN"] = "test-token"
	nc.envVars["TELEGRAM_CHAT_ID"] = "123456"
	return nil
}

func (nc *notifyCtx) envVarIsSetTo(name, value string) error {
	nc.envVars[name] = value
	return nil
}

func (nc *notifyCtx) aNotificationMessage(message string) error {
	nc.message = message
	return nil
}

func (nc *notifyCtx) theTelegramSenderWillFailWith(errMsg string) error {
	nc.sender.err = fmt.Errorf("%s", errMsg)
	return nil
}

func (nc *notifyCtx) theNotificationCLIRuns() error {
	args := []string{}
	if nc.message != "" {
		args = append(args, nc.message)
	}

	getenv := func(name string) string {
		return nc.envVars[name]
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	nc.exitCode = Run(args, &stdoutBuf, &stderrBuf, getenv, nc.sender)
	nc.stdout = stdoutBuf.String()
	nc.stderr = stderrBuf.String()
	return nil
}

func (nc *notifyCtx) theTelegramSenderIsCalledWithChatID(chatID string) error {
	if !nc.sender.called {
		return fmt.Errorf("expected Telegram sender to be called")
	}
	if nc.sender.chatID != chatID {
		return fmt.Errorf("expected chat ID %q, got %q", chatID, nc.sender.chatID)
	}
	return nil
}

func (nc *notifyCtx) theTelegramSenderIsCalledWithMessage(message string) error {
	if !nc.sender.called {
		return fmt.Errorf("expected Telegram sender to be called")
	}
	if nc.sender.message != message {
		return fmt.Errorf("expected message %q, got %q", message, nc.sender.message)
	}
	return nil
}

func (nc *notifyCtx) theTelegramSenderIsNotCalled() error {
	if nc.sender.called {
		return fmt.Errorf("expected Telegram sender not to be called")
	}
	return nil
}

func (nc *notifyCtx) theExitCodeIs(code int) error {
	if nc.exitCode != code {
		return fmt.Errorf("expected exit code %d, got %d\nstdout: %s\nstderr: %s", code, nc.exitCode, nc.stdout, nc.stderr)
	}
	return nil
}

func (nc *notifyCtx) stderrContains(expected string) error {
	if !strings.Contains(nc.stderr, expected) {
		return fmt.Errorf("expected stderr to contain %q, got %q", expected, nc.stderr)
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	nc := &notifyCtx{}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		nc.reset()
		return ctx, nil
	})

	ctx.Step(`^Telegram credentials via environment variables$`, nc.telegramCredentialsViaEnvironmentVariables)
	ctx.Step(`^([A-Z_]+) is set to "([^"]*)"$`, nc.envVarIsSetTo)
	ctx.Step(`^a notification message "([^"]*)"$`, nc.aNotificationMessage)
	ctx.Step(`^the Telegram sender will fail with "([^"]*)"$`, nc.theTelegramSenderWillFailWith)

	ctx.Step(`^the notification CLI runs$`, nc.theNotificationCLIRuns)

	ctx.Step(`^the Telegram sender is called with chat ID "([^"]*)"$`, nc.theTelegramSenderIsCalledWithChatID)
	ctx.Step(`^the Telegram sender is called with message "([^"]*)"$`, nc.theTelegramSenderIsCalledWithMessage)
	ctx.Step(`^the Telegram sender is not called$`, nc.theTelegramSenderIsNotCalled)
	ctx.Step(`^stderr contains "([^"]*)"$`, nc.stderrContains)
	ctx.Step(`^the exit code is (\d+)$`, nc.theExitCodeIs)
}

// BDD: specs/notify.feature :: Feature: Notification CLI
func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../../specs/notify.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
