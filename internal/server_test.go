package internal

import (
	"bytes"
	"fmt"
	"github.com/emersion/go-smtp"
	smtpmock "github.com/mocktools/go-smtp-mock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strconv"
	"strings"
	"testing"
)

const (
	relayHost = "localhost"
	relayPort = 1025

	host = "localhost"
	port = 1026
)

func TestServer(t *testing.T) {
	mockedServer := smtpmock.New(smtpmock.ConfigurationAttr{
		HostAddress: host,
		PortNumber:  port,
		LogToStdout: true,
	})

	require.NoError(t, mockedServer.Start())
	server := NewServer(&Config{
		ServerHost: relayHost,
		ServerPort: strconv.Itoa(relayPort),
		Host:       host,
		Port:       strconv.Itoa(port),

		// Mocked server doesn't support TLS and AUTH, disable them.
		// https://github.com/mocktools/go-smtp-mock/issues/76
		// https://github.com/mocktools/go-smtp-mock/issues/84
		StartTLS: false,
		Auth:     false,
	})
	go func() { require.NoError(t, server.ListenAndServe()) }()

	c, err := smtp.Dial(server.Addr)
	require.NoError(t, err)

	var (
		from = "avfadeev@gmail.com"
		to   = []string{"to@icloud.com"}
		body = []byte(`Subject: Hello
Content-Type: text/plain; charset=UTF-8

World!
`)
	)

	require.NoError(t, c.SendMail(from, to, bytes.NewReader(body)))
	msgs := mockedServer.Messages()
	require.Equal(t, 1, len(msgs))

	msg := msgs[0]
	assert.Equal(t, fmt.Sprintf("MAIL FROM:<%s>", from), msg.MailfromRequest())
	assert.Equal(t, "250 Received", msg.MailfromResponse())

	rcpt := msg.RcpttoRequestResponse()[0]
	assert.Equal(t, fmt.Sprintf("RCPT TO:<%s>", to[0]), rcpt[0])
	assert.Equal(t, "250 Received", rcpt[1])

	assert.Equal(t, "DATA", msg.DataRequest())
	assert.Equal(
		t,
		"354 Ready for receive message. End data with <CR><LF>.<CR><LF>",
		msg.DataResponse(),
	)

	assert.Equal(
		t,
		strings.ReplaceAll(string(body), "\n", "\r\n"),
		msg.MsgRequest(),
	)
	assert.Equal(t, "250 Received", msg.MsgResponse())

	require.NoError(t, mockedServer.Stop())
	require.NoError(t, server.Close())
}
