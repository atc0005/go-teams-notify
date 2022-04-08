// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/go-teams-notify
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

/*

This is an example of a simple client application which uses this library to
generate a user mention within a specific Microsoft Teams channel.

Of note:

- default timeout
- package-level logging is disabled by default
- validation of known webhook URL prefixes is *enabled*
- message is in Adaptive Card format
- text is unformatted (formatting is allowed, just not shown in this example)
- a specific user mention is added to the message

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

	// setup empty message
	msg := adaptivecard.NewMessage()

	// add user mention
	if err := msg.Mention(true, "John Doe", "jdoe@example.com", "Hello there!"); err != nil {
		fmt.Printf(
			"failed to add user mention: %v",
			err,
		)
		os.Exit(1)
	}

	// send message
	if err := mstClient.Send(webhookUrl, msg); err != nil {
		fmt.Printf(
			"failed to send message: %v",
			err,
		)
		os.Exit(1)
	}
}
