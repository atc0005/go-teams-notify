// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/go-teams-notify
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

/*

This example demonstrates how to enable custom validation patterns for webhook
URLs.

Of note:

- message is in MessageCard format
- default timeout
- package-level logging is disabled by default
- webhook URL validation uses custom pattern
  - allows use of custom/private webhook URL endpoints
- simple message submitted to Microsoft Teams consisting of formatted body and
  title

*/

package main

import (
	"log"
	"os"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/atc0005/go-teams-notify/v2/messagecard"
)

func main() {

	// Initialize a new Microsoft Teams client.
	mstClient := goteamsnotify.NewTeamsClient()

	// Set webhook url.
	webhookUrl := "https://my.domain.com/webhook/YOUR_WEBHOOK_URL_OF_TEAMS_CHANNEL"

	// Add a custom pattern for webhook URL validation
	mstClient.AddWebhookURLValidationPatterns(`^https://.*\.domain\.com/.*$`)
	// It's also possible to use multiple patterns with one call
	// mstClient.AddWebhookURLValidationPatterns(`^https://arbitrary\.example\.com/webhook/.*$`, `^https://.*\.domain\.com/.*$`)
	// To keep the default behavior and add a custom one, use something like the following:
	// mstClient.AddWebhookURLValidationPatterns(DefaultWebhookURLValidationPattern, `^https://.*\.domain\.com/.*$`)

	// Setup message card.
	msgCard := messagecard.NewMessageCard()
	msgCard.Title = "Hello world"
	msgCard.Text = "Here are some examples of formatted stuff like " +
		"<br> * this list itself  <br> * **bold** <br> * *italic* <br> * ***bolditalic***"
	msgCard.ThemeColor = "#DF813D"

	// Send the message with default timeout/retry settings.
	if err := mstClient.Send(webhookUrl, msgCard); err != nil {
		log.Printf("failed to send message: %v", err)
		os.Exit(1)
	}
}
