package internal

import (
	"fmt"
	"github.com/slack-go/slack"
	"strings"
)

const (
	emailRecipient = "email"
	slackRecipient = "slack"

	personalMessage = "personal."
	channelMessage  = "channel."
)

var (
	operators = map[string]map[string]string{
		slackRecipient: {
			personalMessage: "@",
			channelMessage:  "#",
		},
	}
)

type client interface {

	// proxify is a function, that will be used as a proxy for the data.
	proxify(from string, to []string, body string) error
}

func takeUsernames(recipients []string, recipientType string) []string {
	var usernames = make([]string, 0, len(recipients))

	for _, r := range recipients {
		at := strings.Split(r, "@")
		if len(at) != 2 {
			continue
		}

		if at[1] != recipientType {
			continue
		}

		usernames = append(usernames, convertToRecipientWay(at[0], recipientType))
	}

	return usernames
}

func convertToRecipientWay(username, recipientType string) string {
	if _, ok := operators[recipientType]; !ok {
		return username
	}

	var msgType string
	switch {
	case strings.HasPrefix(username, personalMessage):
		msgType = personalMessage
	case strings.HasPrefix(username, channelMessage):
		msgType = channelMessage
	}

	operator := operators[recipientType][msgType]
	return operator + strings.TrimPrefix(username, msgType)
}

func selectClientType(recipient string) (string, error) {
	at := strings.Split(recipient, "@")
	if len(at) != 2 {
		return "", fmt.Errorf("invalid recipient: %s", recipient)
	}

	domain := at[1]
	if domain == slackRecipient {
		return slackRecipient, nil
	}

	return emailRecipient, nil
}

func selectClient(clientType string, cfg *Config) (client, error) {
	switch clientType {
	case emailRecipient:
		return newMailClient(cfg)
	case slackRecipient:
		return newSlackClient(slack.New(cfg.SlackToken)), nil
	default:
		return nil, fmt.Errorf("unknown client type: %s", clientType)
	}
}

func groupByClientType(recipients []string) (map[string][]string, error) {
	groups := make(map[string][]string)

	for _, recipient := range recipients {
		clientType, err := selectClientType(recipient)
		if err != nil {
			return nil, fmt.Errorf("failed to select client type: %w", err)
		}

		groups[clientType] = append(groups[clientType], recipient)
	}

	return groups, nil
}
