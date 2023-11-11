package internal

import (
	"fmt"
	"github.com/slack-go/slack"
	"log/slog"
	"sync"
)

//go:generate mockery --name=slackAPI --output=../mocks --exported
type slackAPI interface {
	PostMessage(channelID string, options ...slack.MsgOption) (string, string, error)
}

type slackClient struct {
	slackAPI
}

func newSlackClient(sc slackAPI) *slackClient {
	return &slackClient{
		slackAPI: sc,
	}
}

func (c *slackClient) proxify(from string, to []string, body string) error {
	return c.sendMessageConcurrent(
		takeUsernames(to, slackRecipient),
		slack.MsgOptionBlocks(
			slack.NewHeaderBlock(
				slack.NewTextBlockObject(
					slack.PlainTextType,
					fmt.Sprintf("New message from %s", from),
					false,
					false,
				),
			),
			slack.NewContextBlock(
				"",
				slack.NewTextBlockObject(
					slack.PlainTextType,
					body,
					false,
					false,
				),
			),
		),
	)
}

func (c *slackClient) sendMessageConcurrent(
	recipients []string,
	msgOpts ...slack.MsgOption,
) error {
	var (
		wg         sync.WaitGroup
		errorsChan = make(chan error)
		msg        = slack.MsgOptionCompose(msgOpts...)
	)

	slog.Info("sending message to", "recipients", recipients)
	for _, user := range recipients {
		wg.Add(1)

		go func(u string) {
			defer wg.Done()

			_, _, err := c.PostMessage(u, msg)
			if err != nil {
				errorsChan <- err
			}
		}(user)
	}

	go func() {
		defer close(errorsChan)
		wg.Wait()
	}()

	var errors []error
	for err := range errorsChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to send messages: %s", errors)
	}

	return nil
}
