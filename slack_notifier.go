package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const SlackPostMessageUrl = "https://slack.com/api/chat.postMessage"

const (
	HTTPTimeoutSeconds    = 10
	RequestTimeoutSeconds = 60
)

type SlackNotifier struct {
	client *http.Client
	token  string
}

func NewSlackNotifier(token string) *SlackNotifier {
	slog.Info("Creating Slack notifier")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}
	return &SlackNotifier{
		client: &http.Client{
			Timeout:   time.Duration(HTTPTimeoutSeconds) * time.Second,
			Transport: tr,
		},
		token: token,
	}
}

func (d *SlackNotifier) Notify(ctx context.Context, service Services) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(RequestTimeoutSeconds)*time.Second)
	defer cancel()

	slackMessage := d.messageBody(service.Name, service.SlackChannel)
	slog.With("slackMessage", slackMessage).Info("Sending slack message")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, SlackPostMessageUrl, strings.NewReader(slackMessage))
	if err != nil {
		slog.With("error", err).Error("failed to create slack request")
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", d.token))
	response, err := d.client.Do(req)
	if err != nil {
		slog.With("error", err).Error("failed to send slack request")
		return
	}
	if response.StatusCode > http.StatusOK {
		slog.With("response", response).Error("error sending slack message")
	}
	err = response.Body.Close()
	if err != nil {
		slog.With("error", err).Error("failed to close slack response body")
		return
	}
}

func (d *SlackNotifier) messageBody(serviceName string, channel string) string {
	data := url.Values{}
	data.Set("channel", channel)
	data.Set("text", fmt.Sprintf("Service %s is not functioning as expected. Please investigate!", serviceName))
	return data.Encode()
}
