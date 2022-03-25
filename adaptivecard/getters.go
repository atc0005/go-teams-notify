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

// supportedActionValues returns a list of valid Action types. This list is
// intended to be used for validation and display purposes.
//
// NOTE: See also the supportedISelectActionValues() function.
func supportedActionValues() []string {
	// https://adaptivecards.io/explorer/AdaptiveCard.html
	// https://docs.microsoft.com/en-us/adaptive-cards/authoring-cards/universal-action-model
	return []string{
		TypeActionExecute,
		TypeActionSubmit,
		TypeActionOpenURL,
		TypeActionShowCard,
		TypeActionToggleVisibility,
	}
}

// supportedISelectActionValues returns a list of valid ISelectAction types,
// which is a subset of the supported Action types. This list is intended to
// be used for validation and display purposes.
//
// NOTE: See also the supportedActionValues() function.
func supportedISelectActionValues() []string {
	return []string{
		TypeActionExecute,
		TypeActionSubmit,
		TypeActionOpenURL,
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
