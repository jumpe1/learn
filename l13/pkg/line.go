package pkg

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
)

type LineMessageAPI struct {
	channelToken string
	userID       string
}

func NewLineMessageAPI(channelToken, userID string) *LineMessageAPI {
	return &LineMessageAPI{
		channelToken: channelToken,
		userID:       userID,
	}
}

func (l *LineMessageAPI) Send(message string) error {
	const maxRetries = 3
	const backoffFactor = 2

	client := &http.Client{}
	bot, err := messaging_api.NewMessagingApiAPI(
		l.channelToken,
		messaging_api.WithHTTPClient(client),
	)
	if err != nil {
		return err
	}

	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		_, lastErr = bot.PushMessage(
			&messaging_api.PushMessageRequest{
				To: l.userID,
				Messages: []messaging_api.MessageInterface{
					messaging_api.TextMessage{
						Text: message,
					},
				},
			},
			"",
		)
		if lastErr == nil {
			return nil
		}

		log.Printf("Attempt %d/%d failed: %v", attempt, maxRetries, lastErr)
		time.Sleep(time.Duration(attempt*backoffFactor) * time.Second)
	}

	return errors.New("failed to send message after retries: " + lastErr.Error())
}
