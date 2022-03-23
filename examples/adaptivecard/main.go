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

	// TODO: Remove these placeholders
	_ = mstClient
	_ = webhookUrl

	// setup message
	// msg := adaptivecard.NewSimpleMessage("")
	msg := adaptivecard.NewSimpleMessage("Hello there!")

	fmt.Printf("%+v\n", msg)

	err := msg.Prepare()
	if err != nil {
		fmt.Printf(
			"failed to prepare message: %v",
			err,
		)
		os.Exit(1)
	}

	fmt.Print(msg.PrettyPrint())

	// 	// add user mention
	// 	if err := msg.Mention("John Doe", "jdoe@example.com", true); err != nil {
	// 		fmt.Printf(
	// 			"failed to add user mention: %v",
	// 			err,
	// 		)
	// 	}
	//
	// 	// send message
	// 	if err := mstClient.Send(webhookUrl, msg); err != nil {
	// 		fmt.Printf(
	// 			"failed to send message: %v",
	// 			err,
	// 		)
	// 		os.Exit(1)
	// 	}
}
