// Copyright 2022 Adam Chalkley
//
// https://github.com/atc0005/go-teams-notify
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package botapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	// MessageType is the type for a BotAPI Message.
	MessageType string = "message"

	// MentionType is the type for a user mention for a BotAPI Message.
	MentionType string = "mention"

	// MentionTextFormatTemplate is the expected format of the Mention.Text
	// field value.
	MentionTextFormatTemplate string = "<at>%s</at>"
)

var (
	// ErrInvalidType indicates that an invalid type was specified.
	ErrInvalidType = errors.New("invalid type value")

	// ErrInvalidFieldValue indicates that an invalid value was specified.
	ErrInvalidFieldValue = errors.New("invalid field value")

	// ErrMissingValue indicates that an expected value was missing.
	ErrMissingValue = errors.New("missing expected value")
)

// Message is a minimal representation of the object used to mention one or
// more users in a Teams channel.
//
// https://docs.microsoft.com/en-us/microsoftteams/platform/bots/how-to/conversations/channel-and-group-conversations?tabs=json#add-mentions-to-your-messages
type Message struct {
	// Type is required; must be set to "message".
	Type string `json:"type"`

	// Text is required; mostly freeform content, but testing shows that the
	// "<at>Some User</at>" string (composed of Display Name value) is
	// required by Microsoft Teams for each Mention in the entities
	// collection.
	Text string `json:"text"`

	// Entities is required; a collection of Mention values, one per mentioned
	// individual.
	Entities []Mention `json:"entities"`

	// payload is a prepared Message in JSON format for submission or pretty
	// printing.
	payload *bytes.Buffer `json:"-"`
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

// NewMessage creates a new Message with required fields predefined.
func NewMessage() *Message {
	return &Message{
		Type: MessageType,
	}
}

// NewMessage creates a new Message using provided text with required fields
// predefined.
// func NewMessage(text string) *Message {
// 	return &Message{
// 		Type: MessageType,
// 		Text: text,
// 	}
// }

// AddText adds given text to the message for delivery. If specified, this
// method prepends given text instead of appending it.
//
// The caller may directly write to the exported Message Text field in order
// to overwrite existing Message text. The caller then takes responsibility
// for ensuring that any user mention placeholders are explicitly provided for
// the Message Text field in order to comply with API requirements.
// func (m *Message) AddText(text string, prepend bool) *Message {
// 	switch {
// 	case strings.TrimSpace(text) == "":
// 		// Passing an empty text string is effectively a NOOP.
// 	case prepend:
// 		m.Text = text + " " + m.Text
// 	default:
// 		m.Text += text
// 	}
//
// 	return m
// }

// AddText appends given text to the message for delivery.
//
// As an alternative to using this method, the caller may directly write to
// the exported Message Text field. If opting to use this approach, care
// should be taken by the caller to retain any previously added mention
// placeholders.
func (m *Message) AddText(text string) *Message {
	switch {
	case strings.TrimSpace(text) == "":
		// Passing an empty text string is effectively a NOOP.
	default:
		m.Text += text
	}

	return m
}

// PrettyPrint returns a formatted JSON payload of the Message if the
// Prepare() method has been called, or an empty string otherwise.
func (m *Message) PrettyPrint() string {
	if m.payload != nil {
		var prettyJSON bytes.Buffer

		// Validation is handled by the Message.Prepare() method.
		_ = json.Indent(&prettyJSON, m.payload.Bytes(), "", "\t")

		return prettyJSON.String()
	}

	return ""
}

// Validate performs basic validation of required field values.
func (m Message) Validate() error {
	if m.Text == "" {
		return fmt.Errorf(
			"required Text field is empty: %w",
			ErrInvalidFieldValue,
		)
	}

	if m.Type != MessageType {
		return fmt.Errorf(
			"got %s; wanted %s: %w",
			m.Type,
			MessageType,
			ErrInvalidType,
		)
	}

	// If we have any recorded user mentions, check each of them.
	if len(m.Entities) > 0 {
		for _, mention := range m.Entities {
			if err := mention.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

// Validate performs basic validation of required field values.
func (m Mention) Validate() error {
	if m.Type != MentionType {
		return fmt.Errorf(
			"got %s; wanted %s: %w",
			m.Type,
			MentionType,
			ErrInvalidType,
		)
	}

	if m.Text == "" {
		return fmt.Errorf(
			"required Text field is empty: %w",
			ErrInvalidFieldValue,
		)
	}

	if m.Mentioned.ID == "" {
		return fmt.Errorf(
			"required ID field is empty: %w",
			ErrInvalidFieldValue,
		)
	}

	if m.Mentioned.Name == "" {
		return fmt.Errorf(
			"required Name field is empty: %w",
			ErrInvalidFieldValue,
		)
	}

	return nil
}

// AddMention adds one or many Mention values to a Message.
//
// If specified, the Text field from each given Mention is prepended to the
// Text field of the Message in order to satisfy the API Message format
// requirements. If specified, the given separator is used, otherwise a space
// is assumed.
//
// If the caller opts to not update the Message Text field when adding a
// Mention, the caller is then responsible for ensuring that the Message Text
// field contains a valid match for each mentioned user.
//
// NOTE: Testing indicates that the expected format matches the DisplayName
// field for the user (e.g., "John Doe" instead of "John" or "Doe" or a custom
// format).
func (m *Message) AddMention(prependToText bool, separator string, mentions ...Mention) error {
	if len(mentions) == 0 {
		return fmt.Errorf(
			"func AddMention: missing value: %w",
			ErrMissingValue,
		)
	}

	for _, mention := range mentions {
		if err := mention.Validate(); err != nil {
			return fmt.Errorf(
				"func AddMention: validation failed: %w",
				err,
			)
		}

		m.Entities = append(m.Entities, mention)

		// Fallback to single space separator if user didn't specify one.
		if separator == "" {
			separator = " "
		}

		if prependToText {
			m.Text = mention.Text + separator + m.Text
		}
	}

	return nil
}

// Mention creates a new user Mention to be included in the Message entities
// collection.
//
// This method receives a user's DisplayName, ID and a boolean value used to
// indicate whether a leading text string of the format "<at>John Doe</at>"
// (i.e., a user "mention") should be prepended to the Message Text field.
//
// If the caller opts to not have this method update the Message Text field,
// then the caller will need to ensure that the Message Text field is updated
// to include a matching pattern for every Mention that is included in the
// entities collection for the Message.
//
// NOTE: Brief testing suggests that the user's display name (e.g., "John
// Doe") is required instead of a firstname (e.g., "John"), lastname ("Doe")
// or custom value (e.g., "JD") is required.
//
// The ID value can be an object ID (e.g.,
// 5e8b0f4d-2cd4-4e17-9467-b0f6a5c0c4d0) or a UserPrincipalName (e.g.,
// NewUser@contoso.onmicrosoft.com).
func (m *Message) Mention(displayName string, id string, prependToText bool) error {
	switch {
	case displayName == "":
		return fmt.Errorf(
			"func Mention: required name argument is empty: %w",
			ErrMissingValue,
		)

	case id == "":
		return fmt.Errorf(
			"func Mention: required id argument is empty: %w",
			ErrMissingValue,
		)

	default:
		mention := Mention{
			Type: MentionType,
			// Text: textVal,
			Text: fmt.Sprintf(MentionTextFormatTemplate, displayName),
			Mentioned: Mentioned{
				ID:   id,
				Name: displayName,
			},
		}

		m.Entities = append(m.Entities, mention)

		if prependToText {
			m.Text = mention.Text + " " + m.Text
		}
	}

	return nil
}

// Prepare handles tasks needed to prepare a given Message for delivery to an
// endpoint. If specified, tasks are repeated regardless of whether a previous
// Prepare call was made. Validation should be performed by the caller prior
// to calling this method.
func (m *Message) Prepare(recreate bool) error {
	if m.payload != nil && !recreate {
		return nil
	}

	jsonMessage, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf(
			"failed to prepare message: %w",
			err,
		)
	}

	m.payload = bytes.NewBuffer(jsonMessage)

	return nil
}

// Payload returns the prepared Message payload. The caller should call
// Prepare() prior to calling this method, results are undefined otherwise.
func (m *Message) Payload() io.Reader {
	return m.payload
}
