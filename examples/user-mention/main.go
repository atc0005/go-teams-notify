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
- simple message submitted to Microsoft Teams consisting of plain text message
  (formatting is allowed, just not shown here) with a specific user mention

*/

package main

import (
	"fmt"
	"os"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/atc0005/go-teams-notify/v2/botapi"
)

func main() {

	// init the client
	mstClient := goteamsnotify.NewTeamsClient()

	webhookUrl := "https://outlook.office.com/webhook/YOUR_WEBHOOK_URL_OF_TEAMS_CHANNEL"

	// setup message
	msg := botapi.NewMessage().AddText("Hello there!")

	// add user mention
	if err := msg.Mention("John Doe", "jdoe@example.com", true); err != nil {
		fmt.Printf(
			"failed to add user mention: %v",
			err,
		)
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
