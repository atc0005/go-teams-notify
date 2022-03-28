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

	envWebhookURL := os.Getenv("WEBHOOK_URL")
	if envWebhookURL != "" {
		fmt.Println(envWebhookURL)
		webhookUrl = envWebhookURL
	}

	// setup message
	// msg := adaptivecard.NewSimpleMessage("")
	simpleMsg := adaptivecard.NewSimpleMessage("Hello there!")

	// fmt.Printf("%+v\n", simpleMsg)

	if err := simpleMsg.Prepare(); err != nil {
		fmt.Printf(
			"failed to prepare message: %v",
			err,
		)
		os.Exit(1)
	}

	fmt.Print(simpleMsg.PrettyPrint())

	if err := mstClient.Send(webhookUrl, simpleMsg); err != nil {
		fmt.Printf(
			"failed to send message: %v",
			err,
		)
		os.Exit(1)
	}

	mentionMsg, err := adaptivecard.NewMentionMsg(
		"Adam Chalkley",
		"atc0005@auburn.edu",
		"My spiffy new message!",
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

	fmt.Print(mentionMsg.PrettyPrint())

	if err := mstClient.Send(webhookUrl, mentionMsg); err != nil {
		fmt.Printf(
			"failed to send message: %v",
			err,
		)
		os.Exit(1)
	}

	// 	// add user mention
	// 	if err := msg.Mention("John Doe", "jdoe@example.com", true); err != nil {
	// 		fmt.Printf(
	// 			"failed to add user mention: %v",
	// 			err,
	// 		)
	// 	}
	//
	// 	// send message

}
