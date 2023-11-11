package internal

import (
	"crypto/tls"
	"fmt"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"log/slog"
	"strings"
)

type mailClient struct {
	*smtp.Client
	cfg *Config
}

func newMailClient(cfg *Config) (*mailClient, error) {
	c, err := smtp.Dial(cfg.SMTPAddr())
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	return &mailClient{
		Client: c,
		cfg:    cfg,
	}, nil
}

func (c *mailClient) proxify(from string, to []string, body string) error {
	defer func() {
		if err := c.Quit(); err != nil {
			slog.Error("failed to quit", "error", err)
		}
	}()

	if err := c.handshake(); err != nil {
		return fmt.Errorf("failed to handshake: %w", err)
	}

	slog.Info("sending message to", "recipients", to)
	if err := c.SendMail(from, to, strings.NewReader(body)); err != nil {
		return fmt.Errorf("failed to send mail: %w", err)
	}

	return nil
}

func (c *mailClient) handshake() error {
	if err := c.startTLS(); err != nil {
		return fmt.Errorf("failed to start tls: %w", err)
	}

	return c.auth()
}

func (c *mailClient) startTLS() error {
	if !c.cfg.StartTLS {
		return nil
	}

	if err := c.StartTLS(&tls.Config{
		InsecureSkipVerify: c.cfg.Insecure(),
		ServerName:         c.cfg.Host,
	}); err != nil {
		return fmt.Errorf("failed to start tls: %w", err)
	}

	return nil
}

func (c *mailClient) auth() error {
	if !c.cfg.Auth {
		return nil
	}

	if err := c.Auth(
		sasl.NewLoginClient(c.cfg.User, c.cfg.Pass),
	); err != nil {
		return fmt.Errorf("failed to auth: %w", err)
	}

	return nil
}
