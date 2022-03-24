// Copyright 2022 Adam Chalkley
//
// https://github.com/atc0005/go-teams-notify
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package adaptivecard

// supportedElementTypes returns a list of valid types for an Adaptive Card
// element used in Microsoft Teams messages. This list is intended to be used
// for validation and display purposes.
func supportedElementTypes() []string {
	// TODO: Confirm whether all types are supported.
	// NOTE: Based on current docs, version 1.4 is the latest supported at this
	// time.
	// https://docs.microsoft.com/en-us/microsoftteams/platform/task-modules-and-cards/cards/cards-reference#support-for-adaptive-cards
	// https://adaptivecards.io/explorer/AdaptiveCard.html
	return []string{
		TypeElementActionSet,
		TypeElementColumnSet,
		TypeElementContainer,
		TypeElementFactSet,
		TypeElementImage,
		TypeElementImageSet,
		TypeElementInputChoiceSet,
		TypeElementInputDate,
		TypeElementInputNumber,
		TypeElementInputText,
		TypeElementInputTime,
		TypeElementInputToggle,
		TypeElementMedia, // Introduced in version 1.1 (TODO: Is this supported in Teams message?)
		TypeElementRichTextBlock,
		TypeElementTextBlock,
		TypeElementTextRun,
	}
}
