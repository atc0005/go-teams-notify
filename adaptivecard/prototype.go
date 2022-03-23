package adaptivecard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

	// 	msg := Message{
	// 		Type: TypeMessage,
	//
	// 		// TODO: Add Attachments type with an Add method that accepts an
	// 		// attachment?
	// 		Attachments: []Attachment{
	// 			{
	// 				ContentType: AttachmentContentType,
	// 				Content: Card{
	// 					Type:    TypeAdaptiveCard,
	// 					Schema:  AdaptiveCardSchema,
	// 					Version: AdaptiveCardMaxVersion,
	// 					Body: []Element{
	// 						{
	// 							Type: TypeElementTextBlock,
	// 							Text: text,
	// 						},
	// 					},
	// 				},
	// 			},
	// 		},
	// 	}
	msg := Message{
		Type: TypeMessage,
	}

	// TODO: Refactor further, make it easy to generate specific types of
	// simple cards.
	textCard := Card{
		Type:    TypeAdaptiveCard,
		Schema:  AdaptiveCardSchema,
		Version: AdaptiveCardMaxVersion,
		Body: []Element{
			{
				Type: TypeElementTextBlock,
				Text: text,
			},
		},
	}

	msg.Attach(&textCard)

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

// TODO: Is this useful for anything? We can just append directly to the
// Attachments field?
func (a *Attachments) Add(attachment Attachment) *Attachments {
	*a = append(*a, attachment)

	return a
}

func (m *Message) Attach(cards ...*Card) *Message {

	if len(cards) == 0 {
		return m
	}

	for _, card := range cards {
		attachment := Attachment{
			ContentType: AttachmentContentType,
			Content:     *card,
		}

		m.Attachments = append(m.Attachments, attachment)
	}

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
		_ = json.Indent(&prettyJSON, m.payload.Bytes(), "", "\t")

		return prettyJSON.String()
	}

	return ""
}

// Prepare handles tasks needed to prepare a given Message for delivery to an
// endpoint. Validation should be performed by the caller prior to calling
// this method.
func (m *Message) Prepare() error {
	jsonMessage, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf(
			"error marshalling Message to JSON: %w",
			err,
		)
	}

	switch {
	case m.payload == nil:
		m.payload = &bytes.Buffer{}
	default:
		m.payload.Reset()
	}

	_, err = m.payload.Write(jsonMessage)
	if err != nil {
		return fmt.Errorf(
			"error updating JSON payload for Message: %w",
			err,
		)
	}

	return nil
}

// Payload returns the prepared Message payload. The caller should call
// Prepare() prior to calling this method, results are undefined otherwise.
func (m *Message) Payload() io.Reader {
	return m.payload
}

// Validate performs validation for Message using ValidateFunc if defined,
// otherwise applying default validation.
func (m Message) Validate() error {
	if m.ValidateFunc != nil {
		return m.ValidateFunc()
	}

	// if m.Text == "" {
	// 	return fmt.Errorf(
	// 		"required Text field is empty: %w",
	// 		ErrInvalidFieldValue,
	// 	)
	// }

	fmt.Printf("\n\nFIXME: Message.Validate() is INCOMPLETE\n\n")

	if m.Type != TypeMessage {
		return fmt.Errorf(
			"invalid message type %q; expected %q: %w",
			m.Type,
			TypeMessage,
			ErrInvalidType,
		)
	}

	// // If we have any recorded user mentions, check each of them.
	// if len(m.Entities) > 0 {
	// 	for _, mention := range m.Entities {
	// 		if err := mention.Validate(); err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	return nil
}
