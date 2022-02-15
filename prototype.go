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

// validateMessage acts as a validation "router" all message types, calling
// out to appropriate helper functions as required.
func (c teamsClient) validateMessage(message interface{}, webhookURL string) error {

	// PLACEHOLDER; FIXME
	return fmt.Errorf(
		"PLACEHOLDER",
	)
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
