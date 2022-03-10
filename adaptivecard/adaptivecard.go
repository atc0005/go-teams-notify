// Copyright 2022 Adam Chalkley
//
// https://github.com/atc0005/go-teams-notify
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package adaptivecard

import (
	"bytes"
	"errors"
)

// TODO: Achieve feature parity with messagecard package (e.g., same
// equivalent functions, methods, constants).

// TODO: Spin off separate GH issues for known missing features, mention them
// here?

// TODO: Add one or more examples of using this package.

const (
	// MessageType is the type for an Adaptive Card Message.
	MessageType string = "message"

	// MentionType is the type for a user mention for a Adaptive Card Message.
	MentionType string = "mention"

	// MentionTextFormatTemplate is the expected format of the Mention.Text
	// field value.
	MentionTextFormatTemplate string = "<at>%s</at>"
)

// Valid types for an Adaptive Card element. Not all types are supported by
// Microsoft Teams.
//
// https://adaptivecards.io/explorer/AdaptiveCard.html
//
// TODO: Confirm whether all types are supported.
const (
	ElementTypeActionSet      string = "ActionSet"
	ElementTypeColumnSet      string = "ColumnSet"
	ElementTypeContainer      string = "Container"
	ElementTypeFactSet        string = "FactSet"
	ElementTypeImage          string = "Image"
	ElementTypeImageSet       string = "ImageSet"
	ElementTypeInputChoiceSet string = "Input.ChoiceSet"
	ElementTypeInputDate      string = "Input.Date"
	ElementTypeInputNumber    string = "Input.Number"
	ElementTypeInputText      string = "Input.Text"
	ElementTypeInputTime      string = "Input.Time"
	ElementTypeInputToggle    string = "Input.Toggle"
	ElementTypeMedia          string = "Media"
	ElementTypeRichTextBlock  string = "RichTextBlock"
	ElementTypeTable          string = "Table"
	ElementTypeTextBlock      string = "TextBlock"
)

var (
	// ErrInvalidType indicates that an invalid type was specified.
	ErrInvalidType = errors.New("invalid type value")

	// ErrInvalidFieldValue indicates that an invalid value was specified.
	ErrInvalidFieldValue = errors.New("invalid field value")

	// ErrMissingValue indicates that an expected value was missing.
	ErrMissingValue = errors.New("missing expected value")
)

// $ json2struct -f adaptive-card-with-mention.json
// https://docs.microsoft.com/en-us/microsoftteams/platform/webhooks-and-connectors/how-to/connectors-using?tabs=cURL#send-adaptive-cards-using-an-incoming-webhook
// https://docs.microsoft.com/en-us/microsoftteams/platform/task-modules-and-cards/cards/cards-format?tabs=adaptive-md%2Cconnector-html#mention-support-within-adaptive-cards
// https://stackoverflow.com/questions/50753072/microsoft-teams-webhook-generating-400-for-adaptive-card
type JSONToStruct struct {
	Attachments []struct {
		Content struct {
			Body []struct {
				Text string `json:"text"`
				Type string `json:"type"`
			} `json:"body"`
			Msteams struct {
				Entities []struct {
					Mentioned struct {
						Id   string `json:"id"`
						Name string `json:"name"`
					} `json:"mentioned"`
					Text string `json:"text"`
					Type string `json:"type"`
				} `json:"entities"`
			} `json:"msteams"`
			Schema  string `json:"schema"`
			Type    string `json:"type"`
			Version string `json:"version"`
		} `json:"content"`
		Contenttype string      `json:"contentType"`
		Contenturl  interface{} `json:"contentUrl"`
	} `json:"attachments"`
	Type string `json:"type"`
}

// https://docs.microsoft.com/en-us/microsoftteams/platform/webhooks-and-connectors/how-to/connectors-using?tabs=cURL#send-adaptive-cards-using-an-incoming-webhook
// https://docs.microsoft.com/en-us/microsoftteams/platform/task-modules-and-cards/cards/cards-format?tabs=adaptive-md%2Cconnector-html#mention-support-within-adaptive-cards
// https://stackoverflow.com/questions/50753072/microsoft-teams-webhook-generating-400-for-adaptive-card
// https://docs.microsoft.com/en-us/adaptive-cards/authoring-cards/getting-started
type Message struct {
	// Type is required; must be set to "message".
	Type string `json:"type"`

	// Attachments is a collection of card objects.
	Attachments []Card `json:"attachments"`

	Content Content `json:"content"`

	// payload is a prepared Message in JSON format for submission or pretty
	// printing.
	payload *bytes.Buffer `json:"-"`
}

// Card represents an Adaptive Card.
type Card struct {

	// ContentType is required; must be set to
	// "application/vnd.microsoft.card.adaptive".
	ContentType string `json:"contentType"`

	// ContentURL appears to be related to support for tabs. Most examples
	// have this value set to null.
	//
	// TODO: Update this description with confirmed details.
	ContentURL NullString `json:"contentUrl"`

	// Content represents the content of an Adaptive Card.
	Content Content `json:"content"`
}

// Content represents the content of an Adaptive Card.
// https://adaptivecards.io/explorer/
type Content struct {

	// Body represents the body of an Adaptive Card. The body is made up of
	// building-blocks known as elements. Elements can be composed to create
	// many types of cards. These elements are shown in the primary card
	// region.
	Body []Element `json:"body"`

	// MSTeams is a container for user mentions.
	MSTeams MSTeams `json:"msteams"`

	// Schema is required; schema represents the URI of the Adaptive Card
	// schema.
	//
	// TODO: Assert "http://adaptivecards.io/schemas/adaptive-card.json".
	Schema string `json:"schema"`

	// Type is required; must be set to "AdaptiveCard"
	//
	// TODO: Assert that this is present.
	Type string `json:"type"`

	// Version is the schema version that the content for an Adaptive Card
	// requires.
	//
	// TODO: Test & determine the minimum supported version that we can use.
	//
	// This doc:
	// https://docs.microsoft.com/en-us/microsoftteams/platform/webhooks-and-connectors/how-to/connectors-using?tabs=cURL#send-adaptive-cards-using-an-incoming-webhook
	// uses "1.2" as the version string.
	Version string `json:"version"`
}

// Element is a "building block" for the body of an Adaptive Card and is shown
// in the primary card region.
type Element struct {

	// Type is required; the type of the element used in the body of an
	// Adaptive Card.
	//
	// TODO: Assert that this is present.
	Type string `json:"type"`

	// Text is used by supported element types to display text. A subset of
	// markdown is supported.
	//
	// https://docs.microsoft.com/en-us/adaptive-cards/authoring-cards/text-features
	Text string `json:"text,omitempty"`

	// Type is required and indicates the type of element used.
	// https://adaptivecards.io/explorer/AdaptiveCard.html
	//

	Size    string   `json:"size,omitempty"`
	Weight  string   `json:"weight,omitempty"`
	Color   string   `json:"color,omitempty"`
	Wrap    bool     `json:"wrap,omitempty"`
	Columns []Column `json:"columns,omitempty"`
}

// MSTeams represents a container for a collection of user mentions.
type MSTeams struct {

	// Entities is a collection of user mentions.
	Entities []Mention `json:"entities"`
}

// Mention represents a mention in the message for a specific user.
type Mention struct {
	// Type is required; must be set to "mention".
	Type string `json:"type"`

	// Text must match a portion of the message text field. If it does not,
	// the mention is ignored.
	//
	// Brief testing indicates that this needs to wrap a name/value in <at>NAME
	// HERE</at> tags.
	Text string `json:"text"`

	// Mentioned represents a user that is mentioned.
	Mentioned Mentioned `json:"mentioned"`
}

// Mentioned represents the user id and name of a user that is mentioned.
type Mentioned struct {
	// ID is the unique identifier for a user that is mentioned. This value
	// can be an object ID (e.g., 5e8b0f4d-2cd4-4e17-9467-b0f6a5c0c4d0) or a
	// UserPrincipalName (e.g., NewUser@contoso.onmicrosoft.com).
	ID string `json:"id"`

	// Name is the DisplayName of the user mentioned.
	Name string `json:"name"`
}
