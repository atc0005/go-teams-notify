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

	textCard := NewTextBlockCard(text, "")

	msg.Attach(textCard)

	return &msg
}

// NewTextBlockCard creates and returns a new Card composed of a single
// TextBlock composed of the given text.
// func NewTextBlockCard(text string) Card {
// 	textBlock := Element{
// 		Type: TypeElementTextBlock,
// 		Wrap: true,
// 		Text: text,
// 	}
//
// 	textCard := Card{
// 		Type:    TypeAdaptiveCard,
// 		Schema:  AdaptiveCardSchema,
// 		Version: fmt.Sprintf(AdaptiveCardVersionTmpl, AdaptiveCardMaxVersion),
// 		Body: []*Element{
// 			&textBlock,
// 		},
// 	}
//
// 	return textCard
// }

// NewTextBlockCard uses the specified text and optional title to create and
// return a new Card composed of a single TextBlock composed of the given
// text.
func NewTextBlockCard(text string, title string) Card {
	textBlock := Element{
		Type: TypeElementTextBlock,
		Wrap: true,
		Text: text,
	}

	card := Card{
		Type:    TypeAdaptiveCard,
		Schema:  AdaptiveCardSchema,
		Version: fmt.Sprintf(AdaptiveCardVersionTmpl, AdaptiveCardMaxVersion),
		Body: []Element{
			textBlock,
		},
	}

	if title != "" {
		titleTextBlock := Element{
			Type:  TypeElementTextBlock,
			Wrap:  true,
			Text:  title,
			Style: TextBlockStyleHeading,
		}

		card.Body = append([]Element{titleTextBlock}, card.Body...)
	}

	return card
}

// NewCard creates and returns an empty Card.
func NewCard() Card {
	return Card{
		Type:    TypeAdaptiveCard,
		Schema:  AdaptiveCardSchema,
		Version: fmt.Sprintf(AdaptiveCardVersionTmpl, AdaptiveCardMaxVersion),
	}
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
//
// NOTE: Including multiple cards in the attachments collection *without*
// attachmentLayout set to "carousel" hides cards after the first. Not sure if
// this is a bug, or if it's intentional.
func (m *Message) Attach(cards ...Card) {
	for _, card := range cards {
		attachment := Attachment{
			ContentType: AttachmentContentType,

			// Explicitly convert Card to TopLevelCard in order to assert that
			// TopLevelCard specific requirements are checked during
			// validation.
			Content: TopLevelCard{card},
		}

		m.Attachments = append(m.Attachments, attachment)
	}
}

// Carousel sets the Message Attachment layout to Carousel display mode.
func (m *Message) Carousel() *Message {
	m.AttachmentLayout = AttachmentLayoutCarousel
	return m
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

// Validate asserts that fields have valid values.
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

// Validate asserts that fields have valid values.
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

	// If there are recorded user mentions, we need to assert that
	// Mention.Text is contained (substring match) within an applicable
	// field of a supported Element of the Card Body.
	//
	// At present, this includes the Text field of a TextBlock Element or
	// the Title or Value fields of a Fact from a FactSet.
	//
	// https://docs.microsoft.com/en-us/microsoftteams/platform/task-modules-and-cards/cards/cards-format#mention-support-within-adaptive-cards
	if len(c.MSTeams.Entities) > 0 {
		hasMentionText := func(elements []Element, m Mention) bool {
			for _, element := range elements {
				if element.HasMentionText(m) {
					return true
				}
			}
			return false
		}

		// User mentions recorded, but no elements in Card Body to potentially
		// contain required text string.
		if len(c.Body) == 0 {
			return fmt.Errorf(
				"user mention text not found in empty Card Body: %w",
				ErrMissingValue,
			)
		}

		// For every user mention, we require at least one match in an
		// applicable Element in the Card Body.
		for _, mention := range c.MSTeams.Entities {
			if !hasMentionText(c.Body, mention) {
				// Card Body contains no applicable elements with required
				// Mention text string.
				return fmt.Errorf(
					"user mention text not found in elements of Card Body: %w",
					ErrMissingValue,
				)
			}
		}
	}

	return nil
}

// Validate asserts that fields have valid values.
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
		// if versionNum < AdaptiveCardMinVersion || versionNum > AdaptiveCardMaxVersion {
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
//
// TODO: Are there any element types which do not support this?
func (e *Element) WithSeparator() *Element {
	e.Separator = true
	return e
}

// Validate asserts that fields have valid values.
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

	// The Text field is required by TextBlock and TextRun elements, but an
	// empty string appears to be permitted. Because of this, we do not have
	// to assert that a value is present for the field.

	if e.Type == TypeElementImage {
		// URL is required for Image element type.
		// https://adaptivecards.io/explorer/Image.html
		if e.URL == "" {
			return fmt.Errorf(
				"required URL is empty for %s: %w",
				e.Type,
				ErrMissingValue,
			)
		}
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

	if e.Style != "" {
		// Valid Style field values differ based on type. For example, a
		// Container element supports Container styles whereas a TextBlock
		// supports a different and more limited set of style values. We use a
		// helper function to retrieve valid style values for evaluation.
		supportedStyleValues := supportedStyleValues(e.Type)

		switch {
		case len(supportedStyleValues) == 0:
			return fmt.Errorf(
				"invalid %s %q for element; %s values not supported for element: %w",
				"Style",
				e.Style,
				"Style",
				ErrInvalidFieldValue,
			)

		case !goteamsnotify.InList(e.Style, supportedStyleValues, false):
			return fmt.Errorf(
				"invalid %s %q for element; expected one of %v: %w",
				"Style",
				e.Style,
				supportedStyleValues,
				ErrInvalidFieldValue,
			)
		}
	}

	if e.Type == TypeElementContainer {
		// Items collection is required for Container element type.
		// https://adaptivecards.io/explorer/Container.html
		if len(e.Items) == 0 {
			return fmt.Errorf(
				"required Items collection is empty for %s: %w",
				e.Type,
				ErrMissingValue,
			)
		}

		for _, item := range e.Items {
			if err := item.Validate(); err != nil {
				return err
			}
		}
	}

	// Used by ColumnSet type, but not required.
	for _, column := range e.Columns {
		if err := column.Validate(); err != nil {
			return err
		}
	}

	if e.Type == TypeElementActionSet {
		// Actions collection is required for ActionSet element type.
		// https://adaptivecards.io/explorer/ActionSet.html
		if len(e.Actions) == 0 {
			return fmt.Errorf(
				"required Actions collection is empty for %s: %w",
				e.Type,
				ErrMissingValue,
			)
		}

		for _, action := range e.Actions {
			if err := action.Validate(); err != nil {
				return err
			}
		}
	}

	if e.Type == TypeElementFactSet {
		// Facts collection is required for FactSet element type.
		// https://adaptivecards.io/explorer/FactSet.html
		if len(e.Facts) == 0 {
			return fmt.Errorf(
				"required Facts collection is empty for %s: %w",
				e.Type,
				ErrMissingValue,
			)
		}

		for _, fact := range e.Facts {
			if err := fact.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

// Validate asserts that fields have valid values.
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

// Validate asserts that fields have valid values.
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

// Validate asserts that fields have valid values.
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

// Validate asserts that fields have valid values.
func (i ISelectAction) Validate() error {
	// Some supportedISelectActionValues are restricted to later Adaptive Card
	// schema versions.
	supportedValues := supportedISelectActionValues(AdaptiveCardMaxVersion)
	if !goteamsnotify.InList(i.Type, supportedValues, false) {
		return fmt.Errorf(
			"invalid %s %q for ISelectAction; expected one of %v: %w",
			"Type",
			i.Type,
			supportedValues,
			ErrInvalidType,
		)
	}

	if i.Fallback != "" {
		supportedValues := supportedISelectActionFallbackValues(AdaptiveCardMaxVersion)
		if !goteamsnotify.InList(i.Fallback, supportedValues, false) {
			return fmt.Errorf(
				"invalid %s %q for ISelectAction; expected one of %v: %w",
				"Fallback",
				i.Fallback,
				supportedValues,
				ErrInvalidFieldValue,
			)
		}
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

// Validate asserts that fields have valid values.
func (a Action) Validate() error {

	// Some Actions are restricted to later Adaptive Card schema versions.
	supportedValues := supportedActionValues(AdaptiveCardMaxVersion)
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

	if a.Fallback != "" {
		supportedValues := supportedActionFallbackValues(AdaptiveCardMaxVersion)
		if !goteamsnotify.InList(a.Fallback, supportedValues, false) {
			return fmt.Errorf(
				"invalid %s %q for Action; expected one of %v: %w",
				"Fallback",
				a.Fallback,
				supportedValues,
				ErrInvalidFieldValue,
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

// Validate asserts that fields have valid values.
//
// Element.Validate() asserts that required Mention.Text content is found for
// each recorded user mention the Card..
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

	return nil
}

// Validate asserts that fields have valid values.
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

		m.Attach(mentionCard)

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

// AddMention adds one or more provided user mentions to the associated Card
// along with a new TextBlock element as the first element in the Card body.
// The Text field for the new TextBlock element is updated with the Mention
// Text. This effectively creates a dedicated TextBlock that acts as a
// "lead-in" or "announcement block" for other elements in the Card.
//
// An error is returned if specified Mention values fail validation.
func (c *Card) AddMention(mentions ...Mention) error {
	textBlock := Element{
		Type: TypeElementTextBlock,
		Wrap: true,
	}

	// Whether the mention text is prepended or appended doesn't matter since
	// the TextBlock element we are adding is empty.
	if err := AddMention(c, &textBlock, true, mentions...); err != nil {
		return err
	}

	// Insert new TextBlock with modified Mention.Text included as the first
	// element.
	c.Body = append([]Element{textBlock}, c.Body...)

	return nil
}

// AddElement adds one or more provided Elements to the Body of the associated
// Card. If specified, the Element values are prepended to the Card Body (as a
// contiguous set retaining current order), otherwise appended to the Card
// Body.
//
// An error is returned if specified Element values fail validation.
func (c *Card) AddElement(prepend bool, elements ...Element) error {
	// Validate first before adding to Card Body.
	for _, element := range elements {
		if err := element.Validate(); err != nil {
			return err
		}
	}

	switch prepend {
	case true:
		c.Body = append(elements, c.Body...)
	case false:
		c.Body = append(c.Body, elements...)
	}

	return nil
}

// AddFactSet adds one or more provided FactSet elements to the Body of the
// associated Card. If specified, the FactSet values are prepended to the Card
// Body (as a contiguous set retaining current order), otherwise appended to
// the Card Body.
//
// An error is returned if specified FactSet values fail validation.
//
// TODO: Is this needed? Should we even have a separate FactSet type that is
// so difficult to work with?
func (c *Card) AddFactSet(prepend bool, factsets ...FactSet) error {
	// Convert to base Element type
	factsetElements := make([]Element, 0, len(factsets))
	for _, factset := range factsets {
		element := Element(factset)
		factsetElements = append(factsetElements, element)
	}

	// Validate first before adding to Card Body.
	for _, element := range factsetElements {
		if err := element.Validate(); err != nil {
			return err
		}
	}

	switch prepend {
	case true:
		c.Body = append(factsetElements, c.Body...)
	case false:
		c.Body = append(c.Body, factsetElements...)
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

// AddMention adds one or more provided user mentions to the specified Card.
// The Text field for the specified TextBlock element is updated with the
// Mention Text. If specified, the Mention Text is prepended, otherwise
// appended. An error is returned if specified Mention values fail validation,
// or one of Card or Element pointers are null.
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

	msg.Attach(mentionCard)

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
	textCard := NewTextBlockCard(msgText, "")

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
// TODO: Require Card pointer?
func NewMessageFromCard(card Card) *Message {
	msg := Message{
		Type: TypeMessage,
	}

	msg.Attach(card)

	return &msg
}

// NewContainer creates an empty Container.
func NewContainer() Container {
	container := Container{
		Type: TypeElementContainer,
	}

	return container
}

// NewTextBlock creates a new TextBlock element using the optional user
// specified Text.
func NewTextBlock(text string) Element {
	textBlock := Element{
		Type: TypeElementTextBlock,
		Wrap: true,
		Text: text,
	}

	return textBlock
}

// NewFactSet creates an empty FactSet.
func NewFactSet() FactSet {
	factSet := FactSet{
		Type: TypeElementFactSet,
	}

	return factSet
}

// AddFact adds one or many Fact values to a FactSet. An error is returned if
// the Fact fails validation or if AddFact is called on an unsupported Element
// type.
func (fs *FactSet) AddFact(facts ...Fact) error {
	// Fail early if called on the wrong Element type.
	if fs.Type != TypeElementFactSet {
		return fmt.Errorf(
			"unsupported element type %s; expected %s: %w",
			fs.Type,
			TypeElementFactSet,
			ErrInvalidType,
		)
	}

	// Validate all Fact values before adding them to the collection.
	for _, fact := range facts {
		if err := fact.Validate(); err != nil {
			return err
		}
	}

	fs.Facts = append(fs.Facts, facts...)

	return nil
}

// HasMentionText asserts that a supported Element type contains the required
// Mention text string necessary to link a user mention to a specific Element.
func (e Element) HasMentionText(m Mention) bool {
	switch {
	case e.Type == TypeElementTextBlock:
		if strings.Contains(e.Text, m.Text) {
			return true
		}
		return false

	case e.Type == TypeElementFactSet:
		for _, fact := range e.Facts {
			if strings.Contains(fact.Title, m.Text) ||
				strings.Contains(fact.Value, m.Text) {

				return true
			}
		}
		return false

	default:
		return false
	}
}

// NewActionOpenURL creates a new Action.OpenURL value using the provided URL
// and title. An error is returned if invalid values are supplied.
func NewActionOpenURL(url string, title string) (Action, error) {
	// Accept the user-specified values as-is, use Validate() method to do the
	// heavy lifting.
	action := Action{
		Type:  TypeActionOpenURL,
		Title: title,
		URL:   url,
	}

	err := action.Validate()
	if err != nil {
		return Action{}, err
	}

	return action, nil
}

// AddElement adds the given Element to the collection of Element values in
// the container. If specified, the Element is inserted at the beginning of
// the collection, otherwise appended to the end.
func (c *Container) AddElement(prepend bool, element Element) error {
	if err := element.Validate(); err != nil {
		return err
	}

	switch prepend {
	case true:
		c.Items = append([]Element{element}, c.Items...)
	case false:
		c.Items = append(c.Items, element)
	}

	return nil
}

// AddContainer adds the given Container Element to the collection of Element
// values for the Card. If specified, the Container Element is inserted at the
// beginning of the collection, otherwise appended to the end.
func (c *Card) AddContainer(prepend bool, container Container) error {
	element := Element(container)

	if err := element.Validate(); err != nil {
		return err
	}

	switch prepend {
	case true:
		c.Body = append([]Element{element}, c.Body...)
	case false:
		c.Body = append(c.Body, element)
	}

	return nil
}
