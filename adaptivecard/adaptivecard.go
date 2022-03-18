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
	// TypeAdaptiveCard is the supported type value for an Adaptive Card.
	TypeAdaptiveCard string = "AdaptiveCard"

	// TypeMessage is the type for an Adaptive Card Message.
	TypeMessage string = "message"

	// TypeMention is the type for a user mention for a Adaptive Card Message.
	TypeMention string = "mention"

	// MentionTextFormatTemplate is the expected format of the Mention.Text
	// field value.
	MentionTextFormatTemplate string = "<at>%s</at>"
)

const (
	// TypeColumn is the type for an Adaptive Card Column.
	TypeColumn string = "Column"

	// ColumnWidthAuto indicates that a column's width should be determined
	// automatically based on other columns in the column group.
	ColumnWidthAuto string = "auto"

	// ColumnWidthStretch indicates that a column's width should be stretched
	// to fill the enclosing column group.
	ColumnWidthStretch string = "stretch"

	// ColumnWidthPixelRegex is a regular expression pattern intended to match
	// specific pixel width values (e.g., 50px).
	ColumnWidthPixelRegex string = "^[0-9]+px"
)

// Text size for text within a TextBlock element.
const (
	SizeSmall      string = "small"
	SizeDefault    string = "default"
	SizeMedium     string = "medium"
	SizeLarge      string = "large"
	SizeExtraLarge string = "extraLarge"
)

// Text weight for TextBlock or TextRun elements.
const (
	WeightBolder  string = "Bolder"
	WeightLighter string = "Lighter"
	WeightDefault string = "Default"
)

// Supported colors for TextBlock elements.
const (
	ColorDefault   string = "Default"
	ColorDark      string = "Dark"
	ColorLight     string = "Light"
	ColorAccent    string = "Accent"
	ColorGood      string = "Good"
	ColorWarning   string = "Warning"
	ColorAttention string = "Attention"
)

// Supported spacing values for FactSet, Container and other container element
// types.
const (
	SpacingDefault    string = "default"
	SpacingNone       string = "none"
	SpacingSmall      string = "small"
	SpacingMedium     string = "medium"
	SpacingLarge      string = "large"
	SpacingExtraLarge string = "extraLarge"
	SpacingPadding    string = "padding"
)

// Supported width values for the msteams property used in in Adaptive Card
// messages sent via Microsoft Teams.
const (
	MSTeamsWidthFull string = "Full"
)

// Supported Actions
const (
	// TypeActionExecute is not supported in Microsoft Teams messages.
	// TypeActionExecute          string = "Action.Execute"
	TypeActionOpenURL          string = "Action.OpenUrl"
	TypeActionShowCard         string = "Action.ShowCard"
	TypeActionSubmit           string = "Action.Submit"
	TypeActionToggleVisibility string = "Action.ToggleVisibility"
)

// Valid types for an Adaptive Card element. Not all types are supported by
// Microsoft Teams.
//
// https://adaptivecards.io/explorer/AdaptiveCard.html
//
// TODO: Confirm whether all types are supported.
// NOTE: Based on current docs, version 1.3 is the latest supported at this
// time.
const (
	TypeElementActionSet      string = "ActionSet"
	TypeElementColumnSet      string = "ColumnSet"
	TypeElementContainer      string = "Container"
	TypeElementFactSet        string = "FactSet"
	TypeElementImage          string = "Image"
	TypeElementImageSet       string = "ImageSet"
	TypeElementInputChoiceSet string = "Input.ChoiceSet"
	TypeElementInputDate      string = "Input.Date"
	TypeElementInputNumber    string = "Input.Number"
	TypeElementInputText      string = "Input.Text"
	TypeElementInputTime      string = "Input.Time"
	TypeElementInputToggle    string = "Input.Toggle"
	TypeElementMedia          string = "Media"         // Introduced in version 1.1 (TODO: Is this supported in Teams message?)
	TypeElementRichTextBlock  string = "RichTextBlock" // Introduced in version 1.2
	TypeElementTextBlock      string = "TextBlock"
	TypeElementTextRun        string = "TextRun" // Introduced in version 1.2
)

var (
	// ErrInvalidType indicates that an invalid type was specified.
	ErrInvalidType = errors.New("invalid type value")

	// ErrInvalidFieldValue indicates that an invalid value was specified.
	ErrInvalidFieldValue = errors.New("invalid field value")

	// ErrMissingValue indicates that an expected value was missing.
	ErrMissingValue = errors.New("missing expected value")
)

// Message represents a Microsoft Teams message containing one or more
// Adaptive Cards.
type Message struct {
	// Type is required; must be set to "message".
	// TODO: Assert that this is present.
	Type string `json:"type"`

	// Attachments is a collection of one or more Adaptive Cards.
	Attachments []Attachment `json:"attachments"`

	// payload is a prepared Message in JSON format for submission or pretty
	// printing.
	payload *bytes.Buffer `json:"-"`
}

// Attachment represents an attached Adaptive Card for a Microsoft Teams
// message.
type Attachment struct {

	// ContentType is required; must be set to
	// "application/vnd.microsoft.card.adaptive".
	// TODO: Assert that this is present.
	ContentType string `json:"contentType"`

	// ContentURL appears to be related to support for tabs. Most examples
	// have this value set to null.
	//
	// TODO: Update this description with confirmed details.
	ContentURL NullString `json:"contentUrl"`

	// Content represents the content of an Adaptive Card.
	Content Card `json:"content"`
}

// Card represents the content of an Adaptive Card.
// https://adaptivecards.io/explorer/
type Card struct {

	// Type is required; must be set to "AdaptiveCard"
	//
	// TODO: Assert that this is present.
	Type string `json:"type"`

	// Schema represents the URI of the Adaptive Card schema.
	//
	// TODO: Assert "http://adaptivecards.io/schemas/adaptive-card.json".
	Schema string `json:"$schema"`

	// Version is required for top-level cards; the schema version that the
	// content for an Adaptive Card requires.
	//
	// Version 1.3 is the highest supported for user-generated cards.
	// https://docs.microsoft.com/en-us/microsoftteams/platform/task-modules-and-cards/cards/cards-reference#support-for-adaptive-cards
	// https://adaptivecards.io/designer
	//
	// Version 1.4 is when Action.Execute was introduced.
	//
	// Per this doc:
	// https://docs.microsoft.com/en-us/microsoftteams/platform/task-modules-and-cards/cards/cards-reference
	//
	// the "Action.Execute" action is supported:
	//
	// "For Adaptive Cards in Incoming Webhooks, all native Adaptive Card
	// schema elements, except Action.Submit, are fully supported. The
	// supported actions are Action.OpenURL, Action.ShowCard,
	// Action.ToggleVisibility, and Action.Execute."
	//
	// TODO: Assert that this is present for top-level cards only?
	Version string `json:"version"`

	// FallbackText is the text shown when the client doesn't support the
	// version specified (may contain markdown).
	FallbackText string `json:"fallbackText,omitempty"`

	// Body represents the body of an Adaptive Card. The body is made up of
	// building-blocks known as elements. Elements can be composed to create
	// many types of cards. These elements are shown in the primary card
	// region.
	//
	// NOTE: If we make this an interface type then the fields of the Element
	// won't be exposed to client code. Perhaps it's better to create
	// constructors for each supported Element type so that required fields
	// are populated and unneeded fields are skipped.
	Body []Element `json:"body"`

	// Actions is a collection of actions to show in the card's action bar.
	// TODO: Should this be a pointer?
	Actions []Action `json:"actions,omitempty"`

	// MSTeams is a container for properties specific to Microsoft Teams
	// messages, including formatting properties and user mentions.
	MSTeams MSTeams `json:"msteams"`

	// VerticalContentAlignment defines how the content should be aligned
	// vertically within the container. Only relevant for fixed-height cards,
	// or cards with a minHeight specified.
	VerticalContentAlignment string `json:"verticalContentAlignment"`
}

// MSTeams represents a container for properties specific to Microsoft Teams
// messages, including formatting properties and user mentions.
type MSTeams struct {

	// Width controls the width of Adaptive Cards within a Microsoft Teams
	// messages.
	// https://docs.microsoft.com/en-us/microsoftteams/platform/task-modules-and-cards/cards/cards-format#full-width-adaptive-card
	//
	// TODO: assert specific values
	// TODO: Research supported values, add as MSTeamsWidthXYZ constants.
	Width string `json:"width,omitempty"`

	// Wrap indicates whether text is ...
	//
	// TODO: Research specific purpose of this field and how interacts with a
	// value set on a specific element of an Adaptive Card.
	//
	// TODO: Confirm that this is a value field.
	// https://github.com/MicrosoftDocs/msteams-docs/issues/5003
	Wrap bool `json:"wrap,omitempty"`

	// AllowExpand controls whether images can be displayed in stage view
	// selectively.
	//
	// https://docs.microsoft.com/en-us/microsoftteams/platform/task-modules-and-cards/cards/cards-format#stage-view-for-images-in-adaptive-cards
	AllowExpand bool `json:"allowExpand,omitempty"`

	// Entities is a collection of user mentions.
	// TODO: Should this be a pointer?
	Entities []Mention `json:"entities,omitempty"`
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

// Element is a "building block" for an Adaptive Card. Elements are shown
// within the primary card region (aka, "body"), columns and other container
// types. Not all fields of this Go struct type are supported by all Adaptive
// Card element types.
type Element struct {

	// Type is required and indicates the type of the element used in the body
	// of an Adaptive Card.
	// https://adaptivecards.io/explorer/AdaptiveCard.html
	//
	// TODO: Assert that this is present.
	Type string `json:"type"`

	// Text is used by supported element types to display text. A subset of
	// markdown is supported for text used in TextBlock elements, but no
	// formatting is permitted in text used in TextRun elements.
	//
	// https://docs.microsoft.com/en-us/adaptive-cards/authoring-cards/text-features
	// https://adaptivecards.io/explorer/TextBlock.html
	// https://adaptivecards.io/explorer/TextRun.html
	Text string `json:"text,omitempty"`

	// Size controls the size of text within a TextBlock element.
	//
	// TODO: Assert specific values
	Size string `json:"size,omitempty"`

	// Weight controls the weight of text in TextBlock or TextRun elements.
	//
	// TODO: Assert specific values
	Weight string `json:"weight,omitempty"`

	// Color controls the color of TextBlock elements or text used in TextRun
	// elements.
	//
	// TODO: Assert specific values
	Color string `json:"color,omitempty"`

	// Spacing controls the amount of spacing between this element and the
	// preceding element.
	// TODO: Assert specific values
	Spacing string `json:"spacing,omitempty"`

	// Columns is a container used by a ColumnSet element type which contains
	// one or more elements.
	Columns []Column `json:"columns,omitempty"`

	// Actions is a collection of actions to show.
	// TODO: Should this be a pointer?
	Actions []Action `json:"actions,omitempty"`

	// Facts is a collection of Fact values that are part of a FactSet element
	// type. Each Fact value is a key/value pair displayed in tabular form.
	// TODO: Should this be a pointer?
	Facts []Fact `json:"facts,omitempty"`

	// Wrap controls whether text is allowed to wrap or is clipped for
	// TextBlock elements.
	Wrap bool `json:"wrap,omitempty"`

	// Separator, when true, indicates that a separating line shown should
	// drawn at the top of the element.
	Separator bool `json:"separator,omitempty"`
}

// Column is a container used by a ColumnSet element type. Each container
// may contain one or more elements.
//
// https://adaptivecards.io/explorer/Column.html
type Column struct {

	// Type is required; must be set to "Column".
	Type string `json:"type"`

	// Width represents the width of a column in the column group. Valid
	// values consist of fixed strings OR a number representing the relative
	// width.
	//
	// "auto", "stretch", a number representing relative width of the column
	// in the column group, or in version 1.1 and higher, a specific pixel
	// width, like "50px".
	//
	// TODO: Assert fixed string constants, integer type OR pixel regex (use
	// ColumnWidthPixelRegex)
	Width interface{} `json:"width"`

	// Items are the card elements that should be rendered inside of the
	// column.
	// TODO: Should this be a pointer?
	Items []Element `json:"items"`
}

// Fact represents a Fact in a FactSet as a key/value pair.
type Fact struct {
	// Title is required; the title of the fact.
	Title string `json:"title"`

	// Value is required; the value of the fact.
	Value string `json:"value"`
}

/*



TODO:

Look at creating separate Action variants similar to what was done for the
MessageCard format by Nicolas Maupu (potential actions). That approach has
separate types for each action variant and has each field set as `omitempty`
which effectively excludes the action from rendered JSON unless ...




*/

// Action represents an action that a user may take on a card. Actions
// typically get rendered in an "action bar" at the bottom of a card.
//
// https://adaptivecards.io/explorer/ActionSet.html
// https://adaptivecards.io/explorer/AdaptiveCard.html
// https://docs.microsoft.com/en-us/microsoftteams/platform/task-modules-and-cards/cards/cards-reference
type Action struct {

	// Type is required; specific values are supported.
	// TODO: Assert that this is present for each action.
	//
	// For Adaptive Cards in Incoming Webhooks, all native Adaptive Card
	// schema elements, except Action.Submit, are fully supported.
	//
	// The supported actions are Action.OpenURL, Action.ShowCard,
	// Action.ToggleVisibility, and Action.Execute.
	//
	// See also https://docs.microsoft.com/en-us/microsoftteams/platform/task-modules-and-cards/cards/cards-reference
	Type string `json:"type"`

	// Title is a label for the button or link that represents this action.
	Title string `json:"title"`

	// URL is required for the Action.OpenUrl type, optional for other action
	// types.
	// TODO: Assert that this is present for Action.OpenUrl type.
	URL string `json:"url,omitempty"`

	// Card property is used by Action.ShowCard type.
	//
	// NOTE: Based on a review of JSON content, it looks like `ActionCard` is
	// really just a `Card` type.
	//
	// refs https://github.com/matthidinger/ContosoScubaBot/blob/master/Cards/SubscriberNotification.JSON
	Card ActionCard `json:"card,omitempty"`
}

type ActionCard struct {
	Type string `json:"type"`
	// Body []ActionCardBody `json:"body"`
	Body []Card `json:"body"`
}

// ActionCardBody appears to be a slice of Card.
// TODO: Duplicate?
// type ActionCardBody struct {
// 	Type string `json:"type"`
// 	Text string `json:"text"`
// 	Wrap bool   `json:"wrap"`
// }
