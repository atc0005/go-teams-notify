// Copyright 2022 Adam Chalkley
//
// https://github.com/atc0005/go-teams-notify
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package botapi

import (
	"errors"
	"fmt"
)

// import (
// 	"github.com/atc0005/go-teams-notify/v2"
// )

const (
	// MentionType is the type for a user mention for a BotAPI Message.
	MentionType string = "mention"

	// MessageType is the type for a BotAPI Message.
	MessageType string = "message"
)

var (
	// ErrInvalidType indicates that an invalid type was specified.
	ErrInvalidType = errors.New("invalid type value")

	// ErrInvalidValue indicates that an invalid value was specified.
	ErrInvalidValue = errors.New("invalid field value")
)

/*

curl -X POST -H "Content-type: application/json" -d '{
    "type": "message",
    "text": "Hey <at>Some User</at> check out this message",
    "entities": [
        {
            "type":"mention",
            "mentioned":{
                "id":"some.user@company.com",
                "name":"Some User"
            },
            "text": "<at>Some User</at>"
        }
    ]
}' <webhook_url>

*/

// Message is a minimal representation of the object used to mention one or
// more users in a Teams channel.
//
// https://docs.microsoft.com/en-us/microsoftteams/platform/bots/how-to/conversations/channel-and-group-conversations?tabs=json#add-mentions-to-your-messages
type Message struct {
	// Type is required; must be set to "message".
	Type string `json:"type"`

	// Text is required; mostly freeform content, but testing shows that the
	// "<at>Some User</at>" string is required by Microsoft Teams.
	//
	// TODO: A unique "<at>Some User</at>" string is believed to be required
	// for each Entity value in the Entities collection.
	Text string `json:"text"`

	// Entities is required; a collection of Mention values, one per mentioned
	// individual.
	Entities []Mention `json:"entities"`
}

// Mention represents a mention in the message for a specific user.
type Mention struct {
	// Type is required; must be set to "mention".
	Type string `json:"type"`

	// Text must match a portion of the message text field. If it does not,
	// the mention is ignored.
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

	// Name is the name of the user mentioned.
	Name string `json:"name"`
}

// NewMessage creates a new Message with required fields predefined.
func NewMessage() Message {
	return Message{
		Type: MessageType,
		// Text?
		// Entities?
	}
}

// Validate implements goteamsnotify.msgValidator by performing basic
// validation of required field values.
func (m Message) Validate() error {
	if m.Text == "" {
		return fmt.Errorf(
			"required Text field is empty: %w",
			ErrInvalidValue,
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
	//
	// TODO: Are we working with a valid BotAPI Message if we do not have any
	// mentions? Presumably yes?
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
			ErrInvalidValue,
		)
	}

	if m.Mentioned.ID == "" {
		return fmt.Errorf(
			"required ID field is empty: %w",
			ErrInvalidValue,
		)
	}

	if m.Mentioned.Name == "" {
		return fmt.Errorf(
			"required Name field is empty: %w",
			ErrInvalidValue,
		)
	}

	return nil
}

func (m Message) AddMention(mention ...Mention) error {
	return errors.New("TODO: Implement this")
}

// func (m Message) Prepare(c teamsClient, )
