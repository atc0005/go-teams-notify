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

- webhook URL validation uses custom pattern
  - allows use of custom/private webhook URL endpoints
- other settings are the same as the basic example previously listed

*/

package main

import (
	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/atc0005/go-teams-notify/v2/messagecard"
)

func main() {
	_ = sendTheMessage()
}

func sendTheMessage() error {
	// init the client
	mstClient := goteamsnotify.NewTeamsClient()

	// setup webhook url
	webhookUrl := "https://my.domain.com/webhook/YOUR_WEBHOOK_URL_OF_TEAMS_CHANNEL"

	// Add a custom pattern for webhook URL validation
	mstClient.AddWebhookURLValidationPatterns(`^https://.*\.domain\.com/.*$`)
	// It's also possible to use multiple patterns with one call
	// mstClient.AddWebhookURLValidationPatterns(`^https://arbitrary\.example\.com/webhook/.*$`, `^https://.*\.domain\.com/.*$`)
	// To keep the default behavior and add a custom one, use something like the following:
	// mstClient.AddWebhookURLValidationPatterns(DefaultWebhookURLValidationPattern, `^https://.*\.domain\.com/.*$`)

	// setup message card
	msgCard := messagecard.NewMessageCard()
	msgCard.Title = "Hello world"
	msgCard.Text = "Here are some examples of formatted stuff like " +
		"<br> * this list itself  <br> * **bold** <br> * *italic* <br> * ***bolditalic***"
	msgCard.ThemeColor = "#DF813D"

	// send
	return mstClient.Send(webhookUrl, msgCard)
}
