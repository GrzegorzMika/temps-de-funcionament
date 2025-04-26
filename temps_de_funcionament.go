package main

import (
	"context"
	"log/slog"
	"net/http"
	"slices"
	"time"
)

const DEFAULT_TIMEOUT = 10 * time.Second

type TempsDeFuncionament struct {
	configuration       *Configuration
	notificationChannel chan Services
	httpClient          *http.Client
	notifiers           []Notifier
}

type Notifier interface {
	Notify(ctx context.Context, service Services)
}

func NewTempsDeFuncionament(c *Configuration, notifiers ...Notifier) *TempsDeFuncionament {
	httpClient := &http.Client{
		Timeout: time.Duration(DEFAULT_TIMEOUT),
	}
	return &TempsDeFuncionament{
		configuration:       c,
		notificationChannel: make(chan Services),
		httpClient:          httpClient,
		notifiers:           notifiers,
	}
}

func (t *TempsDeFuncionament) Notifier(ctx context.Context) {
	for {
		select {
		case service := <-t.notificationChannel:
			for _, notifier := range t.notifiers {
				go notifier.Notify(ctx, service)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (t *TempsDeFuncionament) Checker(ctx context.Context, service Services) {
	failures := 0
	ticker := time.NewTicker(time.Duration(service.PeriodSeconds) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			statusCode, err := t.request(ctx, service.URL)
			if err != nil {
				slog.With("error", err).Error("Error making request")
			}
			slog.With("statusCode", statusCode).Info("Service status code")
			if statusCode == 0 {
				slog.Error("Failed to get status code")
				continue
			}
			if !slices.Contains(service.AllowedCodes, statusCode) {
				failures++
			}
			if failures >= service.MaxFailures {
				t.notificationChannel <- service
				failures = 0
			}
		case <-ctx.Done():
			return
		}
	}

}

func (t *TempsDeFuncionament) request(ctx context.Context, url string) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	resp, err := t.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

func (t *TempsDeFuncionament) Start(ctx context.Context) {
	for _, service := range t.configuration.Services {
		go t.Checker(ctx, service)
	}
	go t.Notifier(ctx)
}
