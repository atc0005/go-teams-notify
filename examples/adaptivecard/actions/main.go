// Copyright 2022 Adam Chalkley
//
// https://github.com/atc0005/go-teams-notify
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

/*

This is an example of a client application which uses this library to generate
a basic Microsoft Teams message in Adaptive Card format.

Of note:

- default timeout
- package-level logging is disabled by default
- validation of known webhook URL prefixes is *enabled*
- simple message submitted to Microsoft Teams consisting of title and
  formatted message body

See https://docs.microsoft.com/en-us/adaptive-cards/authoring-cards/text-features
for the list of supported Adaptive Card text formatting options.

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

	// Create card using provided formatted title and text. We'll modify the
	// card and when finished use it to generate a message for delivery.
	card, err := adaptivecard.NewTextBlockCard("Simple message with OpenURL action", "Hello World", true)
	if err != nil {
		log.Printf(
			"failed to create message: %v",
			err,
		)
		os.Exit(1)
	}

	// Destination for OpenURL action.
	targetURL := "https://github.com/atc0005/go-teams-notify"
	targetURLDesc := "Project Homepage"

	urlAction, err := adaptivecard.NewActionOpenURL(targetURL, targetURLDesc)
	if err != nil {
		log.Printf(
			"failed to create action for card: %v",
			err,
		)
		os.Exit(1)
	}

	if err := card.AddAction(true, urlAction); err != nil {
		log.Printf(
			"failed to add action to card: %v",
			err,
		)
		os.Exit(1)
	}

	// Create Message from Card
	msg, err := adaptivecard.NewMessageFromCard(card)
	if err != nil {
		log.Printf("failed to create message from card: %v", err)
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
