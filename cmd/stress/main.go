package main

import (
	"context"
	"fmt"
	"github.com/emersion/go-smtp"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

var (
	eridaAddr       = os.Getenv("ERIDA_ADDR")
	from            = os.Getenv("STRESS_FROM")
	to              = parseRecipients(os.Getenv("STRESS_TO"))
	bodyPattern     = os.Getenv("STRESS_BODY_PATTERN")
	secondsInterval = os.Getenv("STRESS_SECONDS_INTERVAL")
)

func parseRecipients(recipients string) []string {
	var (
		parsed = strings.Split(recipients, ",")
		rcpt   = make([]string, 0, len(parsed))
	)
	for _, p := range parsed {
		if trim := strings.TrimSpace(p); trim != "" {
			rcpt = append(rcpt, trim)
		}
	}

	return parsed
}

// following code is a stress test for erida
// it will send 10 messages to erida server
//
// they will be launched in the same cluster
func main() {
	var (
		wg            sync.WaitGroup
		reqNumber     = 0
		interval, err = time.ParseDuration(secondsInterval)
	)
	if err != nil {
		interval = 5
		slog.Info("using default interval", "interval", interval)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case <-time.After(interval * time.Second):
				issueRequestInParallel(&wg, reqNumber)
			case <-ctx.Done():
				return
			}
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	<-sigCh

	cancel()
	wg.Wait()
}

func issueRequestInParallel(wg *sync.WaitGroup, reqNumber int) {
	var requests = 10

	for i := 0; i < requests; i++ {
		wg.Add(1)

		go func(rn int) {
			defer wg.Done()
			issueRequest(rn)
		}(reqNumber + i)
	}

	reqNumber += requests
}

func issueRequest(reqNumber int) {
	erida, err := smtp.Dial(eridaAddr)
	if err != nil {
		slog.Error("failed to dial erida server", "error", err)
		return
	}

	defer func() {
		if err = erida.Quit(); err != nil {
			slog.Error("failed to quit erida server", "error", err)
		}
	}()

	var body = strings.NewReader(fmt.Sprintf(bodyPattern, reqNumber))
	if err = erida.SendMail(from, to, body); err != nil {
		slog.Error("failed to send email", "error", err)
		return
	}

	slog.Info("email sent", "reqNumber", reqNumber)
}
