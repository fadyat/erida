package internal

import (
	"errors"
	"github.com/fadyat/erida/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSlackClientProxify(t *testing.T) {
	testcases := []struct {
		name    string
		from    string
		to      []string
		body    string
		expect  func(*mocks.SlackAPI)
		wantErr error
	}{
		{
			name: "success: personal message",
			from: "hellofrom",
			to:   []string{"personal.avfadeev@slack"},
			body: "hello",
			expect: func(api *mocks.SlackAPI) {
				api.On("PostMessage", "@avfadeev", mock.AnythingOfType("slack.MsgOption")).
					Return("", "", nil)
			},
		},
		{
			name: "success: channel message",
			from: "hellofrom",
			to:   []string{"channel.global@slack"},
			body: "hello",
			expect: func(api *mocks.SlackAPI) {
				api.On("PostMessage", "#global", mock.AnythingOfType("slack.MsgOption")).
					Return("", "", nil)
			},
		},
		{
			name: "success: multiple recipients",
			from: "hellofrom",
			to:   []string{"channel.global@slack", "personal.avfadeev@slack"},
			body: "hello",
			expect: func(api *mocks.SlackAPI) {
				api.On("PostMessage", "#global", mock.AnythingOfType("slack.MsgOption")).
					Return("", "", nil)
				api.On("PostMessage", "@avfadeev", mock.AnythingOfType("slack.MsgOption")).
					Return("", "", nil)
			},
		},
		{
			name: "failed: some requests failed",
			from: "hellofrom",
			to:   []string{"channel.global@slack", "personal.avfadeev@slack"},
			body: "hello",
			expect: func(api *mocks.SlackAPI) {
				api.On("PostMessage", "#global", mock.AnythingOfType("slack.MsgOption")).
					Return("", "", nil)
				api.On("PostMessage", "@avfadeev", mock.AnythingOfType("slack.MsgOption")).
					Return("", "", errors.New("something went wrong"))
			},
			wantErr: errors.New("failed to send messages: [something went wrong]"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			api := mocks.NewSlackAPI(t)
			c := newSlackClient(api)

			tc.expect(api)
			err := c.proxify(tc.from, tc.to, tc.body)
			if err != nil {
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			require.Nil(t, err)
			api.AssertExpectations(t)
		})
	}
}
