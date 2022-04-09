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

This example is a somewhat rambling exploration of available options for
generating user mentions. While functional, this example file does not
necessarily reflect optimal approaches for generating user mentions.

See the other Adaptive Card user mention examples for more targeted use cases
or the https://github.com/atc0005/send2teams project for a real-world CLI app
that uses this library to generate (optional) user mentions.

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
	"fmt"
	"os"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/atc0005/go-teams-notify/v2/adaptivecard"
)

func main() {

	// Initialize a new Microsoft Teams client.
	mstClient := goteamsnotify.NewTeamsClient()

	// Set webhook url.
	webhookUrl := "https://outlook.office.com/webhook/YOUR_WEBHOOK_URL_OF_TEAMS_CHANNEL"

	// Allow specifying webhook URL via environment variable, fall-back to
	// hard-coded value in this example file.
	expectedEnvVar := "WEBHOOK_URL"
	envWebhookURL := os.Getenv(expectedEnvVar)
	switch {
	case envWebhookURL != "":
		fmt.Printf(
			"Using webhook URL %q from environment variable %q\n\n",
			envWebhookURL,
			expectedEnvVar,
		)
		webhookUrl = envWebhookURL
	default:
		fmt.Println(expectedEnvVar, "environment variable not set.")
		fmt.Printf("Using hardcoded value %q as fallback\n\n", webhookUrl)
	}

	// Create, print & send simple message.
	simpleMsg, err := adaptivecard.NewSimpleMessage("Hello from NewSimpleMessage!", "", true)
	if err != nil {
		fmt.Printf(
			"failed to create message: %v",
			err,
		)
		os.Exit(1)
	}

	if err := simpleMsg.Prepare(); err != nil {
		fmt.Printf(
			"failed to prepare message: %v",
			err,
		)
		os.Exit(1)
	}

	fmt.Println(simpleMsg.PrettyPrint())

	if err := mstClient.Send(webhookUrl, simpleMsg); err != nil {
		fmt.Printf(
			"failed to send message: %v",
			err,
		)
		os.Exit(1)
	}

	// Create, print & send user mention message.
	mentionMsg, err := adaptivecard.NewMentionMessage(
		"John Doe",
		"jdoe@example.com",
		"New user mention message.",
	)
	if err != nil {
		fmt.Printf(
			"failed to create mention message: %v",
			err,
		)
		os.Exit(1)
	}

	if err := mentionMsg.Prepare(); err != nil {
		fmt.Printf(
			"failed to prepare message: %v",
			err,
		)
		os.Exit(1)
	}

	fmt.Println(mentionMsg.PrettyPrint())

	if err := mstClient.Send(webhookUrl, mentionMsg); err != nil {
		fmt.Printf(
			"failed to send message: %v",
			err,
		)
		os.Exit(1)
	}

	// Create simple message, then add a user mention to it.
	customMsg, err := adaptivecard.NewSimpleMessage("NewSimpleMessage.", "", true)
	if err != nil {
		fmt.Printf(
			"failed to create message: %v",
			err,
		)
		os.Exit(1)
	}

	if err := customMsg.Mention(
		false,
		"John Doe",
		"jdoe@example.com",
		"with a user mention added as a second step.",
	); err != nil {
		fmt.Printf(
			"failed to add user mention: %v",
			err,
		)
		os.Exit(1)
	}

	if err := customMsg.Prepare(); err != nil {
		fmt.Printf(
			"failed to prepare message: %v",
			err,
		)
		os.Exit(1)
	}

	fmt.Println(customMsg.PrettyPrint())

	if err := mstClient.Send(webhookUrl, customMsg); err != nil {
		fmt.Printf(
			"failed to send message: %v",
			err,
		)
		os.Exit(1)
	}

	// Create empty message, add a user mention to it.
	bareMsg := adaptivecard.NewMessage()
	err = bareMsg.Mention(
		false,
		"John Doe",
		"jdoe@example.com",
		"Testing Message.Mention() method on card with no prior Elements.",
	)
	if err != nil {
		fmt.Printf(
			"failed to add user mention: %v",
			err,
		)
		os.Exit(1)
	}

	if err := mstClient.Send(webhookUrl, bareMsg); err != nil {
		fmt.Printf(
			"failed to send message: %v",
			err,
		)
		os.Exit(1)
	}

}
