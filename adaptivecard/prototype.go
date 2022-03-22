package adaptivecard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

/*

Mocking package API


msg := botapi.NewMessage().AddText("Hello there!")

This would require (likely in reverse order):

- creating the Message
- setting the Type field
- creating an Attachment
- setting ContentType
- creating a Card
- setting Type
- setting Schema
- setting Version
- creating an Element (of a specific type?)
- appending Element to the Body slice of Card
- attaching Card to the Attachment slice
- appending Attachment to the Attachments slice of the Message

AddText() could operate on the *Message, appending to the Text field of the
first Element identified.

*/

/*

msg := adaptivecard.NewMessage().AddText("Hello there!")


*/

// NewMessage creates a new Message with required fields predefined.
func NewMessage() *Message {
	return &Message{
		Type: TypeMessage,
	}
}

// NewSimpleMessage creates a new simple Message using given text. If given an
// empty string a minimal Message is returned.
func NewSimpleMessage(text string) *Message {
	if strings.TrimSpace(text) == "" {
		return &Message{
			Type: TypeMessage,
		}
	}

	msg := Message{
		Type: TypeMessage,

		// TODO: Add Attachments type with an Add method that accepts an
		// attachment?
		Attachments: []Attachment{
			{
				ContentType: AttachmentContentType,
				Content: Card{
					Type:    TypeAdaptiveCard,
					Schema:  SchemaAdaptiveCard,
					Version: VersionAdaptiveCardMax,
					Body: []Element{
						{
							Type: TypeElementTextBlock,
							Text: text,
						},
					},
				},
			},
		},
	}

	return &msg
}

// AddText appends given text to the message for delivery.
//
// TODO: What is needed to permit this to work?
func (m *Message) AddText(text string) *Message {
	if strings.TrimSpace(text) == "" {
		return m
	}

	if len(m.Attachments) == 0 {
		// create new:
		// attachment
		// card
		// element
	}

	// PLACEHOLDER
	return m

}

/*

C# snippet:

// Display a carousel of all the rich card types.
reply.AttachmentLayout = AttachmentLayoutTypes.Carousel;
reply.Attachments.Add(Cards.CreateAdaptiveCardAttachment());
reply.Attachments.Add(Cards.GetAnimationCard().ToAttachment());
reply.Attachments.Add(Cards.GetAudioCard().ToAttachment());
reply.Attachments.Add(Cards.GetHeroCard().ToAttachment());
reply.Attachments.Add(Cards.GetOAuthCard().ToAttachment());
reply.Attachments.Add(Cards.GetReceiptCard().ToAttachment());
reply.Attachments.Add(Cards.GetSigninCard().ToAttachment());
reply.Attachments.Add(Cards.GetThumbnailCard().ToAttachment());
reply.Attachments.Add(Cards.GetVideoCard().ToAttachment());

Using that API makes sense: msg.Attachments.Add(...)

*/

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
