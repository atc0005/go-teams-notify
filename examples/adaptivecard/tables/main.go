// Copyright 2022 Adam Chalkley
//
// https://github.com/atc0005/go-teams-notify
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

/*
This is an example of a client application which uses this library to generate
a message with a table within a specific Microsoft Teams channel.

Of note:

- default timeout
- package-level logging is disabled by default
- validation of known webhook URL prefixes is *enabled*
- message is in Adaptive Card format
- text is unformatted
- a small table is added to the message

See https://docs.microsoft.com/en-us/adaptive-cards/authoring-cards/text-features
for the list of supported Adaptive Card text formatting options.
*/
package main

import (
	"fmt"
	"log"
	"os"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/atc0005/go-teams-notify/v2/adaptivecard"
)

func main() {

	// Initialize a new Microsoft Teams client.
	mstClient := goteamsnotify.NewTeamsClient()

	// Set webhook url.
	webhookUrl := "https://tigermailauburn.webhook.office.com/webhookb2/34ce6898-e165-4389-be14-e0cab541f0e1@ccb6deed-bd29-4b38-8979-d72780f62d3b/IncomingWebhook/a0d657ef349a4dfdabdb8bd219cd10f6/dfa896b8-6333-4ee8-a0a7-07962ebff8d7"

	// The title for message (first TextBlock element).
	msgTitle := "Hello world"

	// Formatted message body.
	msgText := "Here are some examples of formatted stuff like " +
		"\n * this list itself  \n * **bold** \n * *italic* \n * ***bolditalic***"

	_ = msgTitle
	_ = msgText

	card := adaptivecard.Card{
		Type:    adaptivecard.TypeAdaptiveCard,
		Schema:  adaptivecard.AdaptiveCardSchema,
		Version: fmt.Sprintf(adaptivecard.AdaptiveCardVersionTmpl, adaptivecard.AdaptiveCardMaxVersion),
		Body: []adaptivecard.Element{
			// {
			// 	Type: adaptivecard.TypeElementTextBlock,
			// 	Wrap: true,
			// 	Text: msgTitle,
			// },
			// {
			// 	Type: adaptivecard.TypeElementTextBlock,
			// 	Wrap: true,
			// 	Text: msgText,
			// },
			{
				Type:      adaptivecard.TypeElementTable,
				GridStyle: adaptivecard.ContainerStyleAccent,
				// ShowGridLines: func() *bool { show := true; return &show }(),
				FirstRowAsHeaders: func() *bool { show := true; return &show }(),
				Columns: []adaptivecard.Column{
					{
						Width: 1,
					},
					{
						Width: 1,
					},
					{
						Width: 1,
					},
				},
				Rows: []adaptivecard.TableRow{
					{
						Type: adaptivecard.TypeTableRow,
						Cells: []adaptivecard.TableCell{
							{
								Type: adaptivecard.TypeTableCell,
								Items: []*adaptivecard.Element{
									{
										Type: adaptivecard.TypeElementTextBlock,
										Wrap: true,
										Text: "Table cell test!",
									},
								},
							},
							{
								Type: adaptivecard.TypeTableCell,
								Items: []*adaptivecard.Element{
									{
										Type: adaptivecard.TypeElementTextBlock,
										Wrap: true,
										Text: "Table cell test!",
									},
								},
							},
							{
								Type: adaptivecard.TypeTableCell,
								Items: []*adaptivecard.Element{
									{
										Type: adaptivecard.TypeElementTextBlock,
										Wrap: true,
										Text: "Table cell test!",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	msg := &adaptivecard.Message{
		Type: adaptivecard.TypeMessage,
	}

	msg.Attach(card)

	fmt.Println("Testing")
	msg.Prepare()
	fmt.Println(msg.PrettyPrint())

	_ = mstClient
	_ = webhookUrl

	// Send the message with default timeout/retry settings.
	if err := mstClient.Send(webhookUrl, msg); err != nil {
		log.Printf(
			"failed to send message: %v",
			err,
		)
		os.Exit(1)
	}
}
