// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/go-teams-notify
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

/*

This example illustrates adding an OpenUri Action to a message card. When
used, this action triggers opening a URI in a separate browser or application.


Of note:

- default timeout
- package-level logging is disabled by default
- validation of known webhook URL prefixes is *enabled*
- message submitted to Microsoft Teams consisting of formatted body, title and
  one OpenUri Action

See also:

- https://docs.microsoft.com/en-us/outlook/actionable-messages/message-card-reference#actions

*/

package main

import (
	"log"

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
	webhookUrl := "https://outlook.office.com/webhook/YOUR_WEBHOOK_URL_OF_TEAMS_CHANNEL"

	// destination for OpenUri action
	targetURL := "https://github.com/atc0005/go-teams-notify"
	targetURLDesc := "Project Homepage"

	// setup message card
	msgCard := messagecard.NewMessageCard()
	msgCard.Title = "Hello world"
	msgCard.Text = "Here are some examples of formatted stuff like " +
		"<br> * this list itself  <br> * **bold** <br> * *italic* <br> * ***bolditalic***"
	msgCard.ThemeColor = "#DF813D"

	// setup Action for message card
	pa, err := messagecard.NewPotentialAction(
		messagecard.PotentialActionOpenURIType,
		targetURLDesc,
	)

	if err != nil {
		log.Fatal("error encountered when creating new action:", err)
	}

	pa.PotentialActionOpenURI.Targets =
		[]messagecard.PotentialActionOpenURITarget{
			{
				OS:  "default",
				URI: targetURL,
			},
		}

	// add the Action to the message card
	if err := msgCard.AddPotentialAction(pa); err != nil {
		log.Fatal("error encountered when adding action to message card:", err)
	}

	// send
	return mstClient.Send(webhookUrl, msgCard)
}
