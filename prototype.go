// Copyright 2022 Adam Chalkley
//
// https://github.com/atc0005/go-teams-notify
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package goteamsnotify

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// MessageSender describes the behavior of a baseline Microsoft Teams client.
type MessageSender interface {
	// validateInput(message MessageValidator, webhookURL string) error
	HTTPClient() *http.Client
	UserAgent() string
	ValidateWebhook(webhookURL string) error

	// TODO: Is this needed?
	//
	// A private method to prevent users implementing the interface so that
	// any future changes to it will not violate backwards compatibility.
	private()
}

// MessagePreparer is intended to cover MessageCard, AdaptiveCard,
// botapi.Message, etc.
type MessagePreparer interface {
	Prepare() (io.Reader, error)
}

// MessageValidator is intended to cover MessageCard, AdaptiveCard,
// botapi.Message, etc.
type MessageValidator interface {
	Validate() error
}

// Message is the interface shared by all supported message formats for
// submission to a Microsoft Teams channel.
type Message interface {
	MessagePreparer
	MessageValidator

	// TODO: Is this needed?
	//
	// A private method to prevent users implementing the interface so that
	// any future changes to it will not violate backwards compatibility.
	private()
}

func sendWithContext(ctx context.Context, client MessageSender, webhookURL string, message Message) error {
	logger.Printf("sendWithContext: Webhook message received: %#v\n", message)

	if err := client.ValidateWebhook(webhookURL); err != nil {
		return fmt.Errorf(
			"failed to validate webhook URL: %w",
			err,
		)
	}

	if err := message.Validate(); err != nil {
		return fmt.Errorf(
			"failed to validate message: %w",
			err,
		)
	}

	messageBuffer, err := message.Prepare()
	if err != nil {
		return fmt.Errorf(
			"failed to prepare message: %w",
			err,
		)
	}

	req, err := prepareRequest(ctx, client.UserAgent(), webhookURL, messageBuffer)
	if err != nil {
		return fmt.Errorf(
			"failed to prepare request: %w",
			err,
		)
	}

	// Submit message to endpoint.
	res, err := client.HTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf(
			"failed to submit message: %w",
			err,
		)
	}

	// Make sure that we close the response body once we're done with it
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()

	responseText, err := processResponse(res)
	if err != nil {
		return fmt.Errorf(
			"failed to process response: %w",
			err,
		)
	}

	logger.Printf("sendWithContext: Response string from Microsoft Teams API: %v\n", responseText)

	return nil
}

func sendWithRetry(ctx context.Context, client MessageSender, webhookURL string, message Message, retries int, retriesDelay int) error {
	var result error

	// initial attempt + number of specified retries
	attemptsAllowed := 1 + retries

	// attempt to send message to Microsoft Teams, retry specified number of
	// times before giving up
	for attempt := 1; attempt <= attemptsAllowed; attempt++ {
		// the result from the last attempt is returned to the caller
		result = sendWithContext(ctx, client, webhookURL, message)

		switch {
		case result != nil:

			logger.Printf(
				"sendWithRetry: Attempt %d of %d to send message failed: %v",
				attempt,
				attemptsAllowed,
				result,
			)

			if ctx.Err() != nil {
				errMsg := fmt.Errorf(
					"sendWithRetry: context cancelled or expired: %v; "+
						"aborting message submission after %d of %d attempts: %w",
					ctx.Err().Error(),
					attempt,
					attemptsAllowed,
					result,
				)

				logger.Println(errMsg)

				return errMsg
			}

			ourRetryDelay := time.Duration(retriesDelay) * time.Second

			logger.Printf(
				"sendWithRetry: Context not cancelled yet, applying retry delay of %v",
				ourRetryDelay,
			)
			time.Sleep(ourRetryDelay)

		default:
			logger.Printf(
				"sendWithRetry: successfully sent message after %d of %d attempts\n",
				attempt,
				attemptsAllowed,
			)

			// No further retries needed
			return nil
		}
	}

	return result
}
