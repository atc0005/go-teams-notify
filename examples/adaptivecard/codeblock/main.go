// Copyright 2022 Adam Chalkley
//
// https://github.com/atc0005/go-teams-notify
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

/*

This is an example of a client application which uses this library to generate
a Microsoft Teams message containing a codeblock in Adaptive Card format.

Of note:

- default timeout
- package-level logging is disabled by default
- validation of known webhook URL formats is *enabled*
- message submitted to Microsoft Teams consisting of title, formatted message
  body and embedded codeblock

See these links for Adaptive Card text formatting options:

- https://docs.microsoft.com/en-us/adaptive-cards/authoring-cards/text-features
- https://learn.microsoft.com/en-us/microsoftteams/platform/task-modules-and-cards/cards/cards-format?tabs=adaptive-md%2Cdesktop%2Cconnector-html#codeblock-in-adaptive-cards


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
	webhookUrl := "https://example.logic.azure.com:443/workflows/GUID_HERE/triggers/manual/paths/invoke?api-version=YYYY-MM-DD&sp=%2Ftriggers%2Fmanual%2Frun&sv=1.0&sig=SIGNATURE_HERE"

	// The title for message (first TextBlock element).
	msgTitle := "Hello world"

	// Formatted message body.
	msgText := "Here are some examples of formatted stuff like " +
		"\n * this list itself  \n * **bold** \n * *italic* \n * ***bolditalic***"

	// Create message using provided formatted title and text.
	msg, err := adaptivecard.NewSimpleMessage(msgText, msgTitle, true)
	if err != nil {
		log.Printf(
			"failed to create message: %v",
			err,
		)
		os.Exit(1)
	}

	codeSnippet := `
	package main

	import "log/slog"

	func main() {
		slog.Info("hello, world")
	}
	`

	// Create codeblock using our snippet.
	codeBlock := adaptivecard.NewCodeBlock(codeSnippet, "Go", 1)

	// TODO: Attach codeblock
	//
	// Q: How? I'd like to append it directly after the message I've prepared.
	// How do I do that in a way that makes sense?
	//
	// We need to limit the cards to just one so that we don't require
	// Carousel mode to be enabled (with presumably limited support).
	//
	// The API surface suggests that we use Message.Attach() to attach a new
	// card, but instead we need to attach a new element to the existing Card
	// which is already "attached" to the Message.
	//
	// Card.AddElement() could be useful here?
	// What provides access to the Card from the Message?
	//
	// The Message.Mention() method uses this logic:
	//
	// switch {
	// case prependElement:
	// 	m.Attachments[0].Content.Body = append(
	// 		[]Element{textBlock},
	// 		m.Attachments[0].Content.Body...,
	// 	)
	// default:
	// 	m.Attachments[0].Content.Body = append(
	// 		m.Attachments[0].Content.Body,
	// 		textBlock,
	// 	)
	// }
	//
	// This field is of type `[]adaptivecard.Element`:
	// msg.Attachments[0].Content.Body
	//
	// We need to handle that by retrieving what is currently there and
	// appending a new Element value.
	//
	// Maybe `Message.AddElementToFirstCard()` or similar?
	//
	// NOTE: We know that there is an entry in the msg.Attachments collection
	// due to the use of the `adaptivecard.NewSimpleMessage()` factory
	// function taking care of that for us.
	//
	// msg.Attachments[0].Content.Body = append(
	// 	msg.Attachments[0].Content.Body,
	// 	codeBlock,
	// )
	//
	// This is a little more ergonomic:
	msg.Attachments[0].Content.AddElement(false, codeBlock)
	//
	// but it would probably be even more ergonomic to provide:
	// msg.AddElement(false, codeBlock)
	//
	// with the method adding the given Element to the first card in the
	// collection or an error if a Card is not already present.

	// Send the message with default timeout/retry settings.
	if err := mstClient.Send(webhookUrl, msg); err != nil {
		log.Printf(
			"failed to send message: %v",
			err,
		)
		os.Exit(1)
	}
}
