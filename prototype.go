// Copyright 2022 Adam Chalkley
//
// https://github.com/atc0005/go-teams-notify
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package goteamsnotify

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"
)

// msgPreparer is intended to cover MessageCard, AdaptiveCard, botapi.Message,
// etc.
type msgPreparer interface {
	Prepare(c teamsClient, webhookURL string) (*bytes.Buffer, error)
	// PrepareRequest(...) ?
	// ProcessResponse() ?
	// Validate(webhookURL string) error
	// String() ? - perhaps implement, but not add to this interface
}

// msgValidator is intended to cover MessageCard, AdaptiveCard,
// botapi.Message, etc.
type msgValidator interface {
	Validate() error
}

func (c teamsClient) sendWithContext(ctx context.Context, webhookURL string, message msgPreparer) error {
	// TODO: Do I need to implement String() method before this can be used?
	logger.Printf("sendWithContext: Webhook message received: %#v\n", message)

	messageBuffer, err := message.Prepare(c, webhookURL)
	if err != nil {
		return err
	}

	req, err := c.prepareRequest(ctx, webhookURL, messageBuffer)
	if err != nil {
		return err
	}

	// Submit message to endpoint.
	res, err := c.httpClient.Do(req)
	if err != nil {
		logger.Println(err)
		return err
	}

	// Make sure that we close the response body once we're done with it
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()

	responseText, err := processResponse(res)
	if err != nil {
		return err
	}

	logger.Printf("sendWithContext: Response string from Microsoft Teams API: %v\n", responseText)

	return nil
}

func (c teamsClient) sendWithRetry(ctx context.Context, webhookURL string, message msgPreparer, retries int, retriesDelay int) error {
	var result error

	// initial attempt + number of specified retries
	attemptsAllowed := 1 + retries

	// attempt to send message to Microsoft Teams, retry specified number of
	// times before giving up
	for attempt := 1; attempt <= attemptsAllowed; attempt++ {
		// the result from the last attempt is returned to the caller
		result = c.sendWithContext(ctx, webhookURL, message)

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
