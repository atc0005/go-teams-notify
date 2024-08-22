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
	"strings"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/atc0005/go-teams-notify/v2/adaptivecard"
)

func main() {

	// Initialize a new Microsoft Teams client.
	mstClient := goteamsnotify.NewTeamsClient()

	// Set webhook url.
	//
	// NOTE: This is for illustration purposes only. Best practice is to NOT
	// hardcode credentials of any kind.
	webhookUrl := "https://example.logic.azure.com:443/workflows/GUID_HERE/triggers/manual/paths/invoke?api-version=YYYY-MM-DD&sp=%2Ftriggers%2Fmanual%2Frun&sv=1.0&sig=SIGNATURE_HERE"

	// Allow specifying webhook URL via environment variable, fall-back to
	// hard-coded value in this example file.
	expectedEnvVar := "WEBHOOK_URL"
	envWebhookURL := os.Getenv(expectedEnvVar)
	switch {
	case envWebhookURL != "":
		log.Printf(
			"Using webhook URL %q from environment variable %q\n\n",
			envWebhookURL,
			expectedEnvVar,
		)
		webhookUrl = envWebhookURL
	default:
		log.Println(expectedEnvVar, "environment variable not set.")
		log.Printf("Using hardcoded value %q as fallback\n\n", webhookUrl)
	}

	// The title for message (first TextBlock element).
	msgTitle := "Hello world"

	// Formatted message body.
	msgText := "Here are some examples of formatted stuff like " +
		"\n * this list itself  \n * **bold** \n * *italic* \n * ***bolditalic***"

	// Create card using provided formatted title and text. We'll modify the
	// card and when finished use it to generate a message for delivery.
	card, err := adaptivecard.NewTextBlockCard(msgText, msgTitle, true)
	if err != nil {
		log.Printf(
			"failed to create card: %v",
			err,
		)
		os.Exit(1)
	}

	// See also https://yourbasic.org/golang/multiline-string/ for other
	// approaches to embedding formatted strings.
	codeSnippet := `
package main

import "log/slog"

func main() {
	slog.Info("hello, world")
}
`

	// If you want to strip leading/trailing whitespace.
	codeSnippet = strings.TrimSpace(codeSnippet)

	// Create CodeBlock using our snippet.
	codeBlock := adaptivecard.NewCodeBlock(codeSnippet, "Go", 1)

	// Add CodeBlock to our Card.
	if err := card.AddElement(false, codeBlock); err != nil {
		log.Printf(
			"failed to add codeblock to card: %v",
			err,
		)
		os.Exit(1)
	}

	// Create Message from Card
	msg, err := adaptivecard.NewMessageFromCard(card)
	if err != nil {
		log.Printf("failed to create message from card: %v", err)
		os.Exit(1)
	}

	// Send the message with default timeout/retry settings.
	if err := mstClient.Send(webhookUrl, msg); err != nil {
		log.Printf(
			"failed to send message: %v",
			err,
		)
		os.Exit(1)
	}
}
