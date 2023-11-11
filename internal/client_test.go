package internal

import (
	"context"
	"fmt"
	smtpmock "github.com/mocktools/go-smtp-mock/v2"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTakeUsernames(t *testing.T) {
	testcases := []struct {
		name          string
		recipients    []string
		recipientType string
		expected      []string
	}{
		{
			name:          "success: slack",
			recipients:    []string{"personal.avfadeev@slack", "channel.global@slack"},
			recipientType: "slack",
			expected:      []string{"@avfadeev", "#global"},
		},
		{
			name:          "success: email",
			recipients:    []string{"avfadeev@gmail.com", "aboba@aboba.com"},
			recipientType: "gmail.com",
			expected:      []string{"avfadeev"},
		},
		{
			name: "success: skipping invalid recipients",
			recipients: []string{
				"personal.avfadeev@slack",
				"channel.global@slack",
				"avfadeev@gmail.com",
			},
			recipientType: "slack",
			expected:      []string{"@avfadeev", "#global"},
		},
		{
			name: "success: skipping because of invalid size",
			recipients: []string{
				"personal.avfadeev@slack",
				"channel.global@slack",
				"personal.avfadeev-slack",
			},
			recipientType: "slack",
			expected:      []string{"@avfadeev", "#global"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := takeUsernames(tc.recipients, tc.recipientType)
			require.Equal(t, tc.expected, got)
		})
	}
}

func TestConvertToRecipientWay(t *testing.T) {
	testcases := []struct {
		name          string
		username      string
		recipientType string
		expected      string
	}{
		{
			name:          "success: slack personal",
			username:      "personal.avfadeev",
			recipientType: "slack",
			expected:      "@avfadeev",
		},
		{
			name:          "success: slack channel",
			username:      "channel.global",
			recipientType: "slack",
			expected:      "#global",
		},
		{
			name:          "success: email",
			username:      "avfadeev@gmail.com",
			recipientType: "email",
			expected:      "avfadeev@gmail.com",
		},
		{
			name:          "success: unknown recipient type",
			username:      "something",
			recipientType: "unknown",
			expected:      "something",
		},
		{
			name:          "success: personal without prefix",
			username:      "personalavfadeev",
			recipientType: "slack",
			expected:      "personalavfadeev",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := convertToRecipientWay(tc.username, tc.recipientType)
			require.Equal(t, tc.expected, got)
		})
	}
}

func TestSelectClientType(t *testing.T) {
	testcases := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		{
			name:  "success: slack",
			input: "personal.avfadeev@slack",
			want:  "slack",
		},
		{
			name:  "success: email",
			input: "avfadeev@gmail.com",
			want:  "email",
		},
		{
			name:    "error: invalid recipient",
			input:   "avfadeev",
			wantErr: fmt.Errorf("invalid recipient: %s", "avfadeev"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := selectClientType(tc.input)
			if err != nil {
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			require.Equal(t, tc.want, got)
		})
	}
}

func TestSelectClient(t *testing.T) {
	testcases := []struct {
		name       string
		pre        func() context.Context
		clientType string
		wantErr    error
	}{
		{
			name:       "success: slack",
			clientType: "slack",
		},
		{
			name:       "success: email",
			clientType: "email",
			pre: func() context.Context {
				srv := smtpmock.New(smtpmock.ConfigurationAttr{
					HostAddress: "localhost",
					PortNumber:  2025,
				})

				go func() {
					require.NoError(t, srv.Start())
				}()

				ctx := context.Background()
				go func() {
					<-ctx.Done()
					require.NoError(t, srv.Stop())
				}()

				return ctx
			},
		},
		{
			name:       "error: invalid client type",
			clientType: "invalid",
			wantErr:    fmt.Errorf("unknown client type: %s", "invalid"),
		},
	}

	cfg := &Config{
		Host:       "localhost",
		Port:       "2025",
		SlackToken: "slack-token",
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.pre != nil {
				ctx := tc.pre()
				defer ctx.Done()
			}

			got, err := selectClient(tc.clientType, cfg)
			if err != nil {
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			require.NotNil(t, got)
		})
	}
}

func TestGroupByClientType(t *testing.T) {
	testcases := []struct {
		name     string
		input    []string
		expected map[string][]string
		wantErr  error
	}{
		{
			name:  "success: slack",
			input: []string{"personal.avfadeev@slack"},
			expected: map[string][]string{
				"slack": {"personal.avfadeev@slack"},
			},
		},
		{
			name:  "success: email",
			input: []string{"avfadeev@gmail.com"},
			expected: map[string][]string{
				"email": {"avfadeev@gmail.com"},
			},
		},
		{
			name:  "success: mixed",
			input: []string{"personal.avfadeev@slack", "channel.global@slack", "avfadeev@gmail.com"},
			expected: map[string][]string{
				"slack": {"personal.avfadeev@slack", "channel.global@slack"},
				"email": {"avfadeev@gmail.com"},
			},
		},
		{
			name:     "success: empty",
			input:    []string{},
			expected: map[string][]string{},
		},
		{
			name:    "error: invalid recipient",
			input:   []string{"avfadeev"},
			wantErr: fmt.Errorf("failed to select client type: invalid recipient: %s", "avfadeev"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := groupByClientType(tc.input)
			if err != nil {
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			require.Equal(t, tc.expected, got)
		})
	}
}
