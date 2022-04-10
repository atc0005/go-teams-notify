// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/go-teams-notify
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

/*

This is an example of a client application which uses this library to process
(pretend) user display names and IDs and generate multiple user mentions
within a specific Microsoft Teams channel.

Of note:

- default timeout
- package-level logging is disabled by default
- validation of known webhook URL prefixes is *enabled*
- message is in Adaptive Card format
- text is formatted
- multiple user mentions are added to the message

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

type userMention struct {
	DisplayName string
	ID          string
}

func main() {

	// Initialize a new Microsoft Teams client.
	mstClient := goteamsnotify.NewTeamsClient()

	// Set webhook url.
	webhookUrl := "https://outlook.office.com/webhook/YOUR_WEBHOOK_URL_OF_TEAMS_CHANNEL"

	// The title for message (first TextBlock element).
	msgTitle := "Hello world"

	// Formatted message body.
	msgText := "Here are some examples of formatted stuff like " +
		"\n * this list itself  \n * **bold** \n * *italic* \n * ***bolditalic***"

	// Create Card using provided formatted title and text. We'll attach user
	// mentions to this Card and then later generate a valid Message for
	// delivery using the Card as input.
	card, err := adaptivecard.NewTextBlockCard(msgText, msgTitle, true)
	if err != nil {
		log.Printf("failed to create card: %v", err)
		os.Exit(1)
	}

	// We pretend that the user has submitted these values via command line
	// flags or some other input source and we have stored them in a struct
	// with two fields for conversion to user mentions in our Microsoft Teams
	// message.
	usersToMention := []userMention{
		{
			DisplayName: "John Doe",
			ID:          "jdoe@example.com",
		},
		{
			DisplayName: "Harry Dresden",
			ID:          "hdresden@example.com",
		},
	}

	if len(usersToMention) > 0 {
		// Process user mention details specified by user, create user mention
		// values that we can attach to the card.
		userMentions := make([]adaptivecard.Mention, 0, len(usersToMention))

		for _, user := range usersToMention {
			userMention, err := adaptivecard.NewMention(user.DisplayName, user.ID)
			if err != nil {
				log.Printf("failed to process user mention: %v\n", err)
				os.Exit(1)
			}
			userMentions = append(userMentions, userMention)
		}

		// Add user mention collection to card.
		if err := card.AddMention(true, userMentions...); err != nil {
			log.Printf("failed to add user mentions to message: %v\n", err)
		}
	}

	// Create new Message using Card as input.
	msg, err := adaptivecard.NewMessageFromCard(card)
	if err != nil {
		log.Printf("failed to create message from card: %v", err)
		os.Exit(1)
	}

	// Send the message with default timeout/retry settings.
	if err := mstClient.Send(webhookUrl, msg); err != nil {
		log.Printf("failed to send message: %v", err)
		os.Exit(1)
	}
}
