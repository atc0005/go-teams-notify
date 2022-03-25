// Copyright 2022 Adam Chalkley
//
// https://github.com/atc0005/go-teams-notify
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package adaptivecard

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
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
// func (m *Message) AddText(text string) *Message {
// 	if strings.TrimSpace(text) == "" {
// 		return m
// 	}
//
// 	if len(m.Attachments) == 0 {
// 		// create new:
// 		// attachment
// 		// card
// 		// element
// 	}
//
// 	// PLACEHOLDER
// 	return m
//
// }

// Add appends an Attachment to the Attachments collection for a Microsoft
// Teams message.
//
// TODO: Is this useful for anything? We can just append directly to the
// Attachments field.
func (a *Attachments) Add(attachment Attachment) *Attachments {
	*a = append(*a, attachment)

	return a
}

// Attach receives and adds one or more Card values to the Attachments
// collection for a Microsoft Teams message.
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

	if m.Type != TypeMessage {
		return fmt.Errorf(
			"invalid message type %q; expected %q: %w",
			m.Type,
			TypeMessage,
			ErrInvalidType,
		)
	}

	for _, attachment := range m.Attachments {
		if err := attachment.Validate(); err != nil {
			return err
		}
	}

	for _, supportedValue := range supportedAttachmentLayoutValues() {
		if !strings.EqualFold(m.AttachmentLayout, supportedValue) {
			return fmt.Errorf(
				"invalid %s %q for Message; expected one of %v: %w",
				"AttachmentLayout",
				m.AttachmentLayout,
				supportedAttachmentLayoutValues(),
				ErrInvalidFieldValue,
			)
		}
	}

	return nil
}

//
// TODO: Create Validate() methods for all custom types that require specific
// type values.
//

// Validate asserts that required fields have valid values.
//
// TODO: Should we support user-specified ValidateFunc() here as well?
func (a Attachment) Validate() error {
	if a.ContentType != AttachmentContentType {
		return fmt.Errorf(
			"invalid attachment type %q; expected %q: %w",
			a.ContentType,
			AttachmentContentType,
			ErrInvalidType,
		)
	}

	return nil
}

// Validate asserts that required fields have valid values.
//
// TODO: Should we support user-specified ValidateFunc() here as well?
func (c Card) Validate() error {
	if c.Type != TypeAdaptiveCard {
		return fmt.Errorf(
			"invalid card type %q; expected %q: %w",
			c.Type,
			TypeAdaptiveCard,
			ErrInvalidType,
		)
	}

	if strings.TrimSpace(c.Schema) != "" {
		if c.Schema != AdaptiveCardSchema {
			return fmt.Errorf(
				"invalid Schema value %q; expected %q: %w",
				c.Schema,
				AdaptiveCardSchema,
				ErrInvalidFieldValue,
			)
		}
	}

	// The Version field is required for top-level cards, optional for
	// Cards nested within an Action.ShowCard.
	//
	// TODO: Should we apply this check? Client code is highly unlikely to set
	// this value.
	//
	// TODO: Should we create a TopLevelCard type (embedding Card type) and
	// apply Version field validation to it instead?
	if !c.secondaryCard {
		if strings.TrimSpace(c.Version) == "" {
			return fmt.Errorf(
				"required field Version is empty for top-level Card: %w",
				ErrMissingValue,
			)
		}
	}

	for _, element := range c.Body {
		if err := element.Validate(); err != nil {
			return err
		}
	}

	for _, action := range c.Actions {
		if err := action.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// Validate asserts that required fields have valid values.
//
// TODO: Should we support user-specified ValidateFunc() here as well?
func (e Element) Validate() error {
	isValidValue := func(specifiedValue string, what string, supportedValues []string) error {
		for _, supportedValue := range supportedValues {
			if !strings.EqualFold(specifiedValue, supportedValue) {
				return fmt.Errorf(
					"invalid %s %q for element; expected one of %v: %w",
					what,
					specifiedValue,
					supportedValues,
					ErrInvalidFieldValue,
				)
			}
		}

		return nil
	}

	supportedElementTypes := supportedElementTypes()
	if err := isValidValue(e.Type, "type", supportedElementTypes); err != nil {
		return err
	}

	if strings.TrimSpace(e.Size) != "" {
		supportedSizeValues := supportedSizeValues()
		if err := isValidValue(e.Size, "size", supportedSizeValues); err != nil {
			return err
		}
	}

	if strings.TrimSpace(e.Weight) != "" {
		supportedWeightValues := supportedWeightValues()
		if err := isValidValue(e.Weight, "weight", supportedWeightValues); err != nil {
			return err
		}
	}

	if strings.TrimSpace(e.Color) != "" {
		supportedColorValues := supportedColorValues()
		if err := isValidValue(e.Color, "color", supportedColorValues); err != nil {
			return err
		}
	}

	if strings.TrimSpace(e.Spacing) != "" {
		supportedSpacingValues := supportedSpacingValues()
		if err := isValidValue(e.Spacing, "spacing", supportedSpacingValues); err != nil {
			return err
		}
	}

	for _, column := range e.Columns {
		if err := column.Validate(); err != nil {
			return err
		}
	}

	for _, action := range e.Actions {
		if err := action.Validate(); err != nil {
			return err
		}
	}

	for _, fact := range e.Facts {
		if err := fact.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// Validate asserts that required fields have valid values.
//
// TODO: Should we support user-specified ValidateFunc() here as well?
func (c Column) Validate() error {
	if c.Type != TypeColumn {
		return fmt.Errorf(
			"invalid column type %q; expected %q: %w",
			c.Type,
			TypeColumn,
			ErrInvalidType,
		)
	}

	if c.Width != nil {
		switch v := c.Width.(type) {
		// Assert fixed keyword values or valid pixel width.
		case string:
			v = strings.TrimSpace(v)

			switch v {
			case ColumnWidthAuto:
			case ColumnWidthStretch:
			default:
				matched, _ := regexp.MatchString(ColumnWidthPixelRegex, v)
				if !matched {
					return fmt.Errorf(
						"invalid pixel width %q; expected value in format %s: %w",
						v,
						ColumnWidthPixelWidthExample,
						ErrInvalidFieldValue,
					)
				}
			}

		// Number representing relative width of the column.
		case int:

		// Unsupported value.
		default:
			return fmt.Errorf(
				"invalid pixel width %q; "+
					"expected one of keywords %q, int value (e.g., %d) "+
					"or specific pixel width (e.g., %s): %w",
				v,
				strings.Join([]string{
					ColumnWidthAuto,
					ColumnWidthStretch,
				}, ","),
				1,
				ColumnWidthPixelWidthExample,
				ErrInvalidFieldValue,
			)
		}
	}

	for _, element := range c.Items {
		if err := element.Validate(); err != nil {
			return err
		}
	}

	if c.SelectAction != nil {
		return c.SelectAction.Validate()
	}

	return nil

	// return errors.New("error: Column.Validate() not implemented yet")
}
func (f Fact) Validate() error {
	return errors.New("error: Fact.Validate() not implemented yet")
}
func (i ISelectAction) Validate() error {
	return errors.New("error: ISelectAction.Validate() not implemented yet")
}
func (a Action) Validate() error {
	return errors.New("error: Action.Validate() not implemented yet")
}
func (m Mention) Validate() error {
	return errors.New("error: Mention.Validate() not implemented yet")
}
func (m Mentioned) Validate() error {
	return errors.New("error: Mentioned.Validate() not implemented yet")
}

/*

User mentions:

- need to add a mention entity for each mentioned person
- the text field of the mention entity has to be present elsewhere, presumably
  a TextBlock

Perhaps provide an AddMention() method to a TextBlock element type? Or, to a
Message type with an *Element method argument?

Provide a function and a method.

The method can call the function, passing in the pointer for the receiver it
was called against. Probably best to put the method on the Element type.

Not sure how to put it on the Message type, unless it tries to either create a
new TextBlock Element on the fly or finds the first one in the collection and
adds the mention there?

Perhaps create a standalone Mention() method that accepts sufficient arguments
to construct a Message with a TextBlock that generates a valid/minimal
user mention.

Because each attachment (Card) has its own msteams JSON object, we'll need a
pointer to the Card in addition to the Element, *unless* each Element knows
which Card it is attached to?

*/
