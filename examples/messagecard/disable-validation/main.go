// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/go-teams-notify
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

/*

This example disables the validation webhook URLs, including the validation of
known prefixes so that custom/private webhook URL endpoints can be used (e.g.,
testing purposes).

Of note:

- webhook URL validation is **disabled**
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
	webhookUrl := "https://example.webhook.office.com/webhook/YOUR_WEBHOOK_URL_OF_TEAMS_CHANNEL"

	// Disable webhook URL validation
	mstClient.SkipWebhookURLValidationOnSend(true)

	// setup message card
	msgCard := messagecard.NewMessageCard()
	msgCard.Title = "Hello world"
	msgCard.Text = "Here are some examples of formatted stuff like " +
		"<br> * this list itself  <br> * **bold** <br> * *italic* <br> * ***bolditalic***"
	msgCard.ThemeColor = "#DF813D"

	// send
	return mstClient.Send(webhookUrl, msgCard)
}
