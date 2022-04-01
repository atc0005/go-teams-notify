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

// supportedSizeValues returns a list of valid Size values for applicable
// Element types. This list is intended to be used for validation and display
// purposes.
func supportedSizeValues() []string {
	// https://adaptivecards.io/explorer/TextBlock.html
	return []string{
		SizeSmall,
		SizeDefault,
		SizeMedium,
		SizeLarge,
		SizeExtraLarge,
	}
}

// supportedWeightValues returns a list of valid Weight values for text in
// applicable Element types. This list is intended to be used for validation
// and display purposes.
func supportedWeightValues() []string {
	// https://adaptivecards.io/explorer/TextBlock.html
	return []string{
		WeightBolder,
		WeightLighter,
		WeightDefault,
	}
}

// supportedColorValues returns a list of valid Color values for text in
// applicable Element types. This list is intended to be used for validation
// and display purposes.
func supportedColorValues() []string {
	// https://adaptivecards.io/explorer/TextBlock.html
	return []string{
		ColorDefault,
		ColorDark,
		ColorLight,
		ColorAccent,
		ColorGood,
		ColorWarning,
		ColorAttention,
	}
}

// supportedSpacingValues returns a list of valid Spacing values for Element
// types. This list is intended to be used for validation and display
// purposes.
func supportedSpacingValues() []string {
	// https://adaptivecards.io/explorer/TextBlock.html
	return []string{
		SpacingDefault,
		SpacingNone,
		SpacingSmall,
		SpacingMedium,
		SpacingLarge,
		SpacingExtraLarge,
		SpacingPadding,
	}
}

// supportedActionValues accepts a value indicating the maximum Adaptive Card
// schema version supported and returns a list of valid Action types. This
// list is intended to be used for validation and display purposes.
//
// NOTE: See also the supportedISelectActionValues() function. See ref links
// for unsupported Action types.
func supportedActionValues(version float64) []string {
	// https://adaptivecards.io/explorer/AdaptiveCard.html
	// https://docs.microsoft.com/en-us/adaptive-cards/authoring-cards/universal-action-model
	// https://docs.microsoft.com/en-us/microsoftteams/platform/task-modules-and-cards/cards/cards-reference
	supportedValues := []string{
		// TypeActionSubmit,
		TypeActionOpenURL,
		TypeActionShowCard,
		TypeActionToggleVisibility,
	}

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
	if version >= 1.4 {
		supportedValues = append(supportedValues, TypeActionExecute)
	}

	return supportedValues
}

// supportedISelectActionValues returns a list of valid ISelectAction types,
// which is a subset of the supported Action types. This list is intended to
// be used for validation and display purposes.
//
// NOTE: See also the supportedActionValues() function. See ref links for
// unsupported Action types.
func supportedISelectActionValues() []string {
	// https://adaptivecards.io/explorer/Column.html
	// https://adaptivecards.io/explorer/TableCell.html
	// https://docs.microsoft.com/en-us/microsoftteams/platform/task-modules-and-cards/cards/cards-reference
	return []string{
		TypeActionExecute,
		// TypeActionSubmit,
		TypeActionOpenURL,
		// TypeActionShowCard,
		TypeActionToggleVisibility,
	}
}

// supportedAttachmentLayoutValues returns a list of valid AttachmentLayout
// values for Message type. This list is intended to be used for validation
// and display purposes.
//
// NOTE: See also the supportedActionValues() function.
func supportedAttachmentLayoutValues() []string {
	return []string{
		AttachmentLayoutList,
		AttachmentLayoutCarousel,
	}
}

// supportedMSTeamsWidthValues returns a list of valid Width field values for
// MSTeams type. This list is intended to be used for validation and display
// purposes.
func supportedMSTeamsWidthValues() []string {
	// https://docs.microsoft.com/en-us/microsoftteams/platform/task-modules-and-cards/cards/cards-format#full-width-adaptive-card
	return []string{
		MSTeamsWidthFull,
	}
}
