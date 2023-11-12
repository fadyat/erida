package internal

import (
	"fmt"
	"github.com/emersion/go-smtp"
	"io"
	"log/slog"
	"time"
)

type backend struct {
	cfg *Config
}

func (b *backend) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	return &session{
		cfg: b.cfg,
	}, nil
}

type session struct {
	cfg  *Config
	from string
	body string
	to   []string
}

func (s *session) Reset() {
	clear(s.to)
}

func (s *session) Logout() error { return nil }

func (s *session) AuthPlain(_, _ string) error {
	// AuthPlain is ignored, because in our case, we are accepting
	// any username and password.
	//
	// Because it's launched in a local network, we don't care about
	// security.

	return nil
}

func (s *session) Mail(from string, _ *smtp.MailOptions) error {
	s.from = from
	return nil
}

func (s *session) Rcpt(to string, _ *smtp.RcptOptions) error {
	s.to = append(s.to, to)
	return nil
}

func (s *session) Data(body io.Reader) error {
	bodyRaw, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}
	s.body = string(bodyRaw)

	groups, err := groupByClientType(s.to)
	if err != nil {
		return fmt.Errorf("failed to group by client type: %w", err)
	}

	for clientType, recipients := range groups {
		c, e := selectClient(clientType, s.cfg)
		if e != nil {
			return fmt.Errorf("failed to select client: %w", e)
		}

		if err = c.proxify(s.from, recipients, s.body); err != nil {
			slog.Error(
				"failed to proxify",
				"error", err,
				"clientType", clientType,
			)
		}
	}

	return nil
}

type Server struct {
	*smtp.Server
}

func NewServer(
	cfg *Config,
) *Server {
	s := &Server{
		Server: smtp.NewServer(&backend{
			cfg: cfg,
		}),
	}

	s.ReadTimeout = 10 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.Addr = cfg.RelayAddr()
	s.Domain = cfg.ServerHost
	s.AllowInsecureAuth = cfg.Insecure()
	return s
}
