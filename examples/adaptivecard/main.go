// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/go-teams-notify
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

/*

TODO: Fix this example

This is an example of a simple client application which uses this library to
generate a user mention within a specific Microsoft Teams channel.

Of note:

- default timeout
- package-level logging is disabled by default
- validation of known webhook URL prefixes is *enabled*
- simple message submitted to Microsoft Teams consisting of plain text message
  (formatting is allowed, just not shown here) with a specific user mention

*/

package main

import (
	"fmt"
	"os"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/atc0005/go-teams-notify/v2/adaptivecard"
)

func main() {

	// init the client
	mstClient := goteamsnotify.NewTeamsClient()

	webhookUrl := "https://outlook.office.com/webhook/YOUR_WEBHOOK_URL_OF_TEAMS_CHANNEL"

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

	// Test handling of incomplete message
	stubMsg := adaptivecard.NewSimpleMessage("")
	if err := stubMsg.Validate(); err != nil {
		fmt.Printf("test message fails validation: %v\n", err)
		// os.Exit(1)
	} else {
		if err := mstClient.Send(webhookUrl, stubMsg); err != nil {
			fmt.Printf(
				"failed to send message: %v",
				err,
			)
			os.Exit(1)
		}
	}

	// Create, print & send simple message.
	simpleMsg := adaptivecard.NewSimpleMessage("Hello from NewSimpleMessage!")
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
		"Adam Chalkley",
		"atc0005@auburn.edu",
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
	customMsg := adaptivecard.NewSimpleMessage("NewSimpleMessage.")
	err = customMsg.Mention(
		"Adam Chalkley",
		"atc0005@auburn.edu",
		"with a user mention added as a second step.",
		false,
	)
	if err != nil {
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
		"Adam Chalkley",
		"atc0005@auburn.edu",
		"Testing Message.Mention() method on card with no prior Elements.",
		false,
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
