// Copyright 2022 Adam Chalkley
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

- message is in Adaptive Card format
- default timeout
- package-level logging is disabled by default
- webhook URL validation is **disabled**
  - allows use of custom/private webhook URL endpoints
- simple message submitted to Microsoft Teams consisting of formatted body and
  title

*/

package main

import (
	"log"
	"os"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/atc0005/go-teams-notify/v2/adaptivecard"
)

func main() {

	// Initialize a new Microsoft Teams client.
	mstClient := goteamsnotify.NewTeamsClient()

	// Set webhook url.
	webhookUrl := "https://outlook.office.com/webhook/YOUR_WEBHOOK_URL_OF_TEAMS_CHANNEL"

	// Disable webhook URL validation
	mstClient.SkipWebhookURLValidationOnSend(true)

	// The title for message (first TextBlock element).
	msgTitle := "Hello world"

	// Formatted message body.
	msgText := "Here are some examples of formatted stuff like " +
		"\n * this list itself  \n * **bold** \n * *italic* \n * ***bolditalic***"

	// Create message using provided formatted title and text.
	msg, err := adaptivecard.NewSimpleMessage(msgText, msgTitle, true)
	if err != nil {
		log.Printf(
			"failed to create message: %v",
			err,
		)
		os.Exit(1)
	}

	// Send the message with default timeout/retry settings.
	if err := mstClient.Send(webhookUrl, msg); err != nil {
		log.Printf(
			"failed to send message: %v",
			err,
		)
		os.Exit(1)
	}
}
