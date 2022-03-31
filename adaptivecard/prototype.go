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
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
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
	if text == "" {
		return &Message{
			Type: TypeMessage,
		}
	}

	msg := Message{
		Type: TypeMessage,
	}

	textBlock := Element{
		Type: TypeElementTextBlock,
		Wrap: true,
		Text: text,
	}

	// TODO: Refactor further, make it easy to generate specific types of
	// simple cards.
	textCard := Card{
		Type:    TypeAdaptiveCard,
		Schema:  AdaptiveCardSchema,
		Version: fmt.Sprintf(AdaptiveCardVersionTmpl, AdaptiveCardMaxVersion),
		Body: []Element{
			textBlock,
		},
	}

	msg.Attach(&textCard)

	return &msg
}

// NewTextBlockCard creates and returns a new Card composed of a single
// TextBlock composed of the given text.
func NewTextBlockCard(text string) Card {
	textBlock := Element{
		Type: TypeElementTextBlock,
		Wrap: true,
		Text: text,
	}

	textCard := Card{
		Type:    TypeAdaptiveCard,
		Schema:  AdaptiveCardSchema,
		Version: fmt.Sprintf(AdaptiveCardVersionTmpl, AdaptiveCardMaxVersion),
		Body: []Element{
			textBlock,
		},
	}

	return textCard
}

// AddText appends given text to the message for delivery.
//
// TODO: What is needed to permit this to work?
// func (m *Message) AddText(text string) *Message {
// 	if text == "" {
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
//
// C# snippet:
//
// Display a carousel of all the rich card types.
// reply.AttachmentLayout = AttachmentLayoutTypes.Carousel;
// reply.Attachments.Add(Cards.CreateAdaptiveCardAttachment());
// reply.Attachments.Add(Cards.GetAnimationCard().ToAttachment());
// reply.Attachments.Add(Cards.GetAudioCard().ToAttachment());
// reply.Attachments.Add(Cards.GetHeroCard().ToAttachment());
// reply.Attachments.Add(Cards.GetOAuthCard().ToAttachment());
// reply.Attachments.Add(Cards.GetReceiptCard().ToAttachment());
// reply.Attachments.Add(Cards.GetSigninCard().ToAttachment());
// reply.Attachments.Add(Cards.GetThumbnailCard().ToAttachment());
// reply.Attachments.Add(Cards.GetVideoCard().ToAttachment());
func (a *Attachments) Add(attachment Attachment) *Attachments {
	*a = append(*a, attachment)

	return a
}

// Attach receives and adds one or more Card values to the Attachments
// collection for a Microsoft Teams message.
func (m *Message) Attach(cards ...*Card) {
	for _, card := range cards {
		attachment := Attachment{
			ContentType: AttachmentContentType,

			// Explicitly convert Card to TopLevelCard in order to assert that
			// TopLevelCard specific requirements are checked during
			// validation.
			Content: TopLevelCard{*card},
		}

		m.Attachments = append(m.Attachments, attachment)
	}
}

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

// Prepare handles tasks needed to construct a payload from a Message for
// delivery to an endpoint.
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

	// We need an attachment (containing one or more Adaptive Cards) in order
	// to generate a valid Message for Microsoft Teams delivery.
	if len(m.Attachments) == 0 {
		return fmt.Errorf(
			"required field Attachments is empty for Message: %w",
			ErrMissingValue,
		)
	}

	for _, attachment := range m.Attachments {
		if err := attachment.Validate(); err != nil {
			return err
		}
	}

	// Optional field, but only specific values permitted if set.
	if m.AttachmentLayout != "" {
		supportedValues := supportedAttachmentLayoutValues()
		if !goteamsnotify.InList(m.AttachmentLayout, supportedValues, false) {
			return fmt.Errorf(
				"invalid %s %q for Message; expected one of %v: %w",
				"AttachmentLayout",
				m.AttachmentLayout,
				supportedValues,
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
func (c Card) Validate() error {
	if c.Type != TypeAdaptiveCard {
		return fmt.Errorf(
			"invalid card type %q; expected %q: %w",
			c.Type,
			TypeAdaptiveCard,
			ErrInvalidType,
		)
	}

	if c.Schema != "" {
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

	// Both are optional fields, unless MinHeight is set in which case
	// VerticalContentAlignment is required.
	if c.MinHeight != "" && c.VerticalContentAlignment == "" {
		return fmt.Errorf(
			"field MinHeight is set, VerticalContentAlignment is not;"+
				" field VerticalContentAlignment is only optional when MinHeight"+
				" is not set: %w",
			ErrMissingValue,
		)
	}

	return nil
}

// Validate asserts that required fields have valid values.
func (tc TopLevelCard) Validate() error {
	// Validate embedded Card first as those validation requirements apply
	// here also.
	if err := tc.Card.Validate(); err != nil {
		return err
	}

	// The Version field is required for top-level cards (this one), optional
	// for Cards nested within an Action.ShowCard.
	switch {
	case strings.TrimSpace(tc.Version) == "":
		return fmt.Errorf(
			"required field Version is empty for top-level Card: %w",
			ErrMissingValue,
		)
	default:
		// Assert that Version value can be converted to the expected format.
		versionNum, err := strconv.ParseFloat(tc.Version, 64)
		if err != nil {
			return fmt.Errorf(
				"value %q incompatible with Version field: %w",
				tc.Version,
				ErrInvalidFieldValue,
			)
		}

		// This is a high confidence validation failure.
		if versionNum < AdaptiveCardMinVersion {
			return fmt.Errorf(
				"unsupported version %q;"+
					" expected minimum value of %0.1f: %w",
				tc.Version,
				AdaptiveCardMinVersion,
				ErrInvalidFieldValue,
			)
		}

		// This is *NOT* a high confidence validation failure; it is likely
		// that Microsoft Teams will gain support for future versions of the
		// Adaptive Card greater than the current recorded max configured
		// schema version. Because the max value constant is subject to fall
		// out of sync (at least briefly), this is a risky assertion to make.
		//
		// if versionNum < AdaptiveCardMinVersion || versionNum > AdaptiveCardMinVersion {
		// 	return fmt.Errorf(
		// 		"unsupported version %q;"+
		// 			" expected value between %0.1f and %0.1f: %w",
		// 		tc.Version,
		// 		AdaptiveCardMinVersion,
		// 		AdaptiveCardMaxVersion,
		// 		ErrInvalidFieldValue,
		// 	)
		// }
	}

	return nil
}

// WithSeparator indicates that a separating line should be drawn at the top
// of the element.
func (e *Element) WithSeparator() *Element {
	e.Separator = true
	return e
}

// Validate asserts that required fields have valid values.
func (e Element) Validate() error {
	supportedElementTypes := supportedElementTypes()
	if !goteamsnotify.InList(e.Type, supportedElementTypes, false) {
		return fmt.Errorf(
			"invalid %s %q for element; expected one of %v: %w",
			"Type",
			e.Type,
			supportedElementTypes,
			ErrInvalidType,
		)
	}

	if e.Size != "" {
		supportedSizeValues := supportedSizeValues()
		if !goteamsnotify.InList(e.Size, supportedSizeValues, false) {
			return fmt.Errorf(
				"invalid %s %q for element; expected one of %v: %w",
				"Size",
				e.Size,
				supportedSizeValues,
				ErrInvalidFieldValue,
			)
		}
	}

	if e.Weight != "" {
		supportedWeightValues := supportedWeightValues()
		if !goteamsnotify.InList(e.Weight, supportedWeightValues, false) {
			return fmt.Errorf(
				"invalid %s %q for element; expected one of %v: %w",
				"Weight",
				e.Weight,
				supportedWeightValues,
				ErrInvalidFieldValue,
			)
		}
	}

	if e.Color != "" {
		supportedColorValues := supportedColorValues()
		if !goteamsnotify.InList(e.Color, supportedColorValues, false) {
			return fmt.Errorf(
				"invalid %s %q for element; expected one of %v: %w",
				"Color",
				e.Color,
				supportedColorValues,
				ErrInvalidFieldValue,
			)
		}
	}

	if e.Spacing != "" {
		supportedSpacingValues := supportedSpacingValues()
		if !goteamsnotify.InList(e.Spacing, supportedSpacingValues, false) {
			return fmt.Errorf(
				"invalid %s %q for element; expected one of %v: %w",
				"Spacing",
				e.Spacing,
				supportedSpacingValues,
				ErrInvalidFieldValue,
			)
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
}

// Validate asserts that required fields have valid values.
func (f Fact) Validate() error {
	if f.Title == "" {
		return fmt.Errorf(
			"required field Title is empty for Fact: %w",
			ErrMissingValue,
		)
	}

	if f.Value == "" {
		return fmt.Errorf(
			"required field Value is empty for Fact: %w",
			ErrMissingValue,
		)
	}

	return nil
}

// Validate asserts that required fields have valid values.
func (m MSTeams) Validate() error {
	// If an optional width value is set, assert that it is a valid value.
	if m.Width != "" {
		supportedValues := supportedMSTeamsWidthValues()
		if !goteamsnotify.InList(m.Width, supportedValues, false) {
			return fmt.Errorf(
				"invalid %s %q for Action; expected one of %v: %w",
				"Width",
				m.Width,
				supportedValues,
				ErrInvalidFieldValue,
			)
		}
	}

	for _, mention := range m.Entities {
		if err := mention.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// Validate asserts that required fields have valid values.
func (i ISelectAction) Validate() error {
	supportedValues := supportedISelectActionValues()
	if !goteamsnotify.InList(i.Type, supportedValues, false) {
		return fmt.Errorf(
			"invalid %s %q for ISelectAction; expected one of %v: %w",
			"Type",
			i.Type,
			supportedValues,
			ErrInvalidType,
		)
	}

	if i.Type == TypeActionOpenURL {
		if i.URL == "" {
			return fmt.Errorf(
				"invalid URL for Action: %w",
				ErrMissingValue,
			)
		}
	}

	return nil
}

// Validate asserts that required fields have valid values.
func (a Action) Validate() error {
	supportedValues := supportedActionValues()
	if !goteamsnotify.InList(a.Type, supportedValues, false) {
		return fmt.Errorf(
			"invalid %s %q for Action; expected one of %v: %w",
			"Type",
			a.Type,
			supportedValues,
			ErrInvalidType,
		)
	}

	if a.Type == TypeActionOpenURL {
		if a.URL == "" {
			return fmt.Errorf(
				"invalid URL for Action: %w",
				ErrMissingValue,
			)
		}
	}

	// Optional, but only supported by the Action.ShowCard type.
	if a.Type != TypeActionShowCard && a.Card != nil {
		return fmt.Errorf(
			"error: specifying a Card is unsupported for Action type %q: %w",
			a.Type,
			ErrInvalidFieldValue,
		)
	}

	return nil
}

// Validate asserts that required fields have valid values.
//
// TODO: Need to assert that Text field of this Mention matches a portion of
func (m Mention) Validate() error {
	if m.Type != TypeMention {
		return fmt.Errorf(
			"invalid Mention type %q; expected %q: %w",
			m.Type,
			TypeMention,
			ErrInvalidType,
		)
	}

	if m.Text == "" {
		return fmt.Errorf(
			"required field Text is empty for Mention: %w",
			ErrMissingValue,
		)
	}

	// TODO: Need to assert that Text field of this Mention matches a portion
	// of a text field for a supported element in the same enclosing Card.
	//
	// This will require a "handle" to the enclosing Card in order to
	// loop over all elements in the body so that we can assert a text match.
	//
	// Expose the parent field as ParentCard or EnclosingCard and skip
	// recording the parent field as a MSTeams pointer. This will allow client
	// code to manage this directly if needed. For our purposes we can set the
	// EnclosingCard via Mention() methods:
	//
	// - method attached to a Card
	//
	// perhaps this method can look for the existing mention text and skip
	// adding it if found, otherwise add it.
	//
	// - method attached directly to an Element that requires a pointer to Card
	//
	// use pointer to a Card to apply the same logic?
	//
	// Perhaps require a pointer to the element as a mention method argument?
	// Is this feasible?
	//
	// We need to be able to validate a Mention type for a Message that has
	// been 100% generated by client code without assistance from this
	// package. This means that a Parent/Enclosing Card pointer is unlikely to
	// be set in those cases.
	//
	//
	// 	if m.parent != nil {
	// 		// If we have a pointer to the Card, we can evaluate supported
	// 		// elements of the Card body to assert that the required text string
	// 		// is present.
	// 		//
	// 		// https://docs.microsoft.com/en-us/microsoftteams/platform/task-modules-and-cards/cards/cards-format#mention-support-within-adaptive-cards
	// 		if m.parent.parent != nil {
	// 			var foundValidTextType bool
	// 			var foundExpectedTextString bool
	// 			for _, element := range m.parent.parent.Body {
	// 				// Look for valid text element types.
	// 				if element.Type == TypeElementTextBlock ||
	// 					element.Type == TypeElementFactSet {
	//
	// 					foundValidTextType = true
	//
	// 					// Look for the expected mention text.
	// 					if strings.Contains(element.Text, m.Text) {
	// 						foundExpectedTextString = true
	// 						break
	// 					}
	// 				}
	// 			}
	//
	// 			if !foundValidTextType {
	// 				// note that a supported text type wasn't found.
	// 			}
	//
	// 			if !foundExpectedTextString {
	// 				// note that the expected mention text was not found.
	// 			}
	// 		}
	// 	}

	return nil
}

// Validate asserts that required fields have valid values.
func (m Mentioned) Validate() error {
	if m.ID == "" {
		return fmt.Errorf(
			"required field ID is empty: %w",
			ErrMissingValue,
		)
	}

	if m.Name == "" {
		return fmt.Errorf(
			"required field Name is empty: %w",
			ErrMissingValue,
		)
	}

	return nil
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

// Mention uses the provided display name, ID and text values to add a new
// user Mention and TextBlock element to the first Card in the Message.
//
// If no Cards are yet attached to the Message, a new card is created using
// the Mention and TextBlock element. If specified, the new TextBlock element
// is added as the first element of the Card, otherwise it is added last. An
// error is returned if insufficient values are provided.
func (m *Message) Mention(displayName string, id string, msgText string, prependElement bool) error {
	switch {
	// If no existing cards, add a new one.
	case len(m.Attachments) == 0:
		mentionCard, err := NewMentionCard(displayName, id, msgText)
		if err != nil {
			return err
		}

		m.Attach(&mentionCard)

	// We have at least one Card already, use it.
	default:

		// Build mention.
		mention, err := NewMention(displayName, id)
		if err != nil {
			return fmt.Errorf(
				"add new Mention to Message: %w",
				err,
			)
		}

		textBlock := Element{
			Type: TypeElementTextBlock,
			Wrap: true,

			// The text block contains the mention text string (required) and
			// user-specified message text string. Use the mention text as a
			// "greeting" or lead-in for the user-specified message text.
			Text: mention.Text + " " + msgText,
		}

		switch {
		case prependElement:
			m.Attachments[0].Content.Body = append(
				[]Element{textBlock},
				m.Attachments[0].Content.Body...,
			)
		default:
			m.Attachments[0].Content.Body = append(
				m.Attachments[0].Content.Body,
				textBlock,
			)
		}

		m.Attachments[0].Content.MSTeams.Entities = append(
			m.Attachments[0].Content.MSTeams.Entities,
			mention,
		)
	}

	return nil
}

// Mention uses the given display name, ID and message text to add a new user
// Mention and TextBlock element to the Card. If specified, the new TextBlock
// element is added as the first element of the Card, otherwise it is added
// last. An error is returned if provided values are insufficient to create
// the user mention.
func (c *Card) Mention(displayName string, id string, msgText string, prependElement bool) error {
	if msgText == "" {
		return fmt.Errorf(
			"required msgText argument is empty: %w",
			ErrMissingValue,
		)
	}

	mention, err := NewMention(displayName, id)
	if err != nil {
		return err
	}

	textBlock := Element{
		Type: TypeElementTextBlock,
		Wrap: true,
		Text: mention.Text + " " + msgText,
	}

	switch {
	case prependElement:
		c.Body = append(c.Body, textBlock)
	default:
		c.Body = append([]Element{textBlock}, c.Body...)
	}

	return nil
}

// NewMention uses the given display name and ID to create a user Mention
// value for inclusion in a Card. An error is returned if provided values are
// insufficient to create the user mention.
func NewMention(displayName string, id string) (Mention, error) {
	switch {
	case displayName == "":
		return Mention{}, fmt.Errorf(
			"required name argument is empty: %w",
			ErrMissingValue,
		)

	case id == "":
		return Mention{}, fmt.Errorf(
			"required id argument is empty: %w",
			ErrMissingValue,
		)

	default:

		// Build mention.
		mention := Mention{
			Type: TypeMention,
			Text: fmt.Sprintf(MentionTextFormatTemplate, displayName),
			Mentioned: Mentioned{
				ID:   id,
				Name: displayName,
			},
		}

		return mention, nil
	}
}

// TODO: Accept values needed to create Mention, update Element receiver (if
// applicable type), then append to Card (using provided pointer).
//
// TODO: What happens if there is already a Mention in the entities
// collection? For example, one single text block with required mention.Text
// content in the textBlock.Text field, but two Mention values in the entities
// collection. Is this an error?
//
// A: No, I put 10 duplicate Mention values in the entities collection without
// issue. That said, if you have a unique Mention value in the entities
// collection and each Mention.Text value is not represented in the Text field
// of an applicable Element type then I suspect *that* would be an error
// scenario.

// AddMention adds one or more provided user mentions to the specified Card.
// The Text field for the specified TextBlock element is updated with the
// Mention Text. If specified, the Mention Text is prepended, otherwise
// appended. An error is returned if specified Mention values fail validation,
// or one of Card or Element pointers are null .
func AddMention(card *Card, textBlock *Element, prependText bool, mentions ...Mention) error {
	if card == nil {
		return fmt.Errorf(
			"specified pointer to Card is nil: %w",
			ErrMissingValue,
		)
	}

	if textBlock == nil {
		return fmt.Errorf(
			"specified pointer to TextBlock element is nil: %w",
			ErrMissingValue,
		)
	}

	if textBlock.Type != TypeElementTextBlock {
		return fmt.Errorf(
			"invalid element type %q; expected %q: %w",
			textBlock.Type,
			TypeElementTextBlock,
			ErrInvalidType,
		)
	}

	// Validate all user mentions before modifying Card or Element.
	for _, mention := range mentions {
		if err := mention.Validate(); err != nil {
			return err
		}
	}

	// Update TextBlock element text with required user mention text string.
	for _, mention := range mentions {
		switch prependText {
		case true:
			textBlock.Text = mention.Text + " " + textBlock.Text
		case false:
			textBlock.Text = textBlock.Text + " " + mention.Text
		}

		card.MSTeams.Entities = append(card.MSTeams.Entities, mention)
	}

	// The original text may have been sufficiently short to not be truncated,
	// but once we add the user mention text it likely would, so explicitly
	// indicate that we wish to disable wrapping.
	textBlock.Wrap = true

	return nil
}

// NewMentionMessage creates a new simple Message. Using the given message
// text, displayName and ID, a user Mention is also created and added to the
// new Message. An error is returned if provided values are insufficient to
// create the user mention.
func NewMentionMessage(displayName string, id string, msgText string) (*Message, error) {
	msg := Message{
		Type: TypeMessage,
	}

	mentionCard, err := NewMentionCard(displayName, id, msgText)
	if err != nil {
		return nil, err
	}

	msg.Attach(&mentionCard)

	return &msg, nil
}

// NewMentionCard creates a new Card with user Mention using the given message
// text, displayName and ID. An error is returned if provided values are
// insufficient to create the user mention.
func NewMentionCard(displayName string, id string, msgText string) (Card, error) {
	if msgText == "" {
		return Card{}, fmt.Errorf(
			"required msgText argument is empty: %w",
			ErrMissingValue,
		)
	}

	// Build mention.
	mention, err := NewMention(displayName, id)
	if err != nil {
		return Card{}, err
	}

	// Create basic card.
	textCard := NewTextBlockCard(msgText)

	// Update the text block so that it contains the mention text string
	// (required) and user-specified message text string. Use the mention
	// text as a "greeting" or lead-in for the user-specified message
	// text.
	textCard.Body[0].Text = mention.Text +
		" " + textCard.Body[0].Text

	textCard.MSTeams.Entities = append(
		textCard.MSTeams.Entities,
		mention,
	)

	return textCard, nil
}

// NewMessageFromCard is a helper function for creating a new Message based
// off of an existing Card value.
//
// TODO: Use pointer for Card?
func NewMessageFromCard(card Card) *Message {
	// TODO: Placeholder only
	return &Message{}
}
