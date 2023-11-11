//go:build integration

package main

import (
	"bytes"
	"context"
	"github.com/emersion/go-smtp"
	"github.com/fadyat/erida/internal"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/stretchr/testify/require"
	"log/slog"
	"testing"
	"time"
)

const (
	testConfigPath = "../test.env"
)

func readTestConfig() (*internal.Config, error) {
	var cfg internal.Config
	if err := cleanenv.ReadConfig(testConfigPath, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// TestRelayFlow is an integration test, that will start a server and send a
// message to it.
//
// Need to run it with `go test -tags=integration ./...`
//
// Real Google GMail SMTP server and Slack API are used.
// Current test doesn't check the correctness of the message, that was sent.
// It only checks, that the server is able to send a message to the real services.
func TestRelayFlow(t *testing.T) {
	cfg, err := readTestConfig()
	require.NoError(t, err)

	srv := internal.NewServer(cfg)
	go func() {
		if err = srv.ListenAndServe(); err != nil {
			slog.Error("failed to start server: ", err)
		}
	}()

	c, err := smtp.Dial(srv.Addr)
	require.NoError(t, err)

	var (
		from = "erida@erida.com"
		to   = []string{
			"avfadeev@gmail.com",
			"personal.fadyat@slack",
			"channel.empty@slack",
		}
		body = []byte(`Subject: Integration test
Content-Type: text/plain; charset=UTF-8

Hello, this Erida!
We are testing the process of sending messages to Slack and GMail,
using this SMTP server as a relay.

Best regards, Erida
`)
	)

	require.NoError(t, c.SendMail(from, to, bytes.NewReader(body)))

	// sleeping for 3 seconds to let the server process the message
	time.Sleep(3 * time.Second)

	require.NoError(t, c.Quit())
	require.NoError(t, srv.Shutdown(context.Background()))
}
