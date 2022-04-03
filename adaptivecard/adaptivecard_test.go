// Copyright 2022 Adam Chalkley
//
// https://github.com/atc0005/go-teams-notify
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package adaptivecard_test

import (
	"testing"

	"github.com/atc0005/go-teams-notify/v2/adaptivecard"
)

func TestNewTextBlockCardWithoutTitle(t *testing.T) {

	expectedText := "Tacos"

	// Create basic card.
	textCard := adaptivecard.NewTextBlockCard(expectedText, "")

	switch {
	// Assert that no more than one element is present.
	case len(textCard.Body) != 1:
		t.Fatalf(
			"want Card Body length of %d; got %d",
			1,
			len(textCard.Body),
		)

	// Assert specific type of element.
	case textCard.Body[0].Type != adaptivecard.TypeElementTextBlock:
		t.Fatalf(
			"want Card Body element type %q; got %q",
			adaptivecard.TypeElementTextBlock,
			textCard.Body[0].Type,
		)

	case textCard.Body[0].Text != expectedText:
		t.Fatalf(
			"want %q for Text field in element; got %q",
			expectedText,
			textCard.Body[0].Text,
		)
	}

}

func TestNewTextBlockCardWithTitle(t *testing.T) {

	expectedTitle := "Wonderful Food"
	expectedText := "Tacos"

	// Create basic card with title
	textCard := adaptivecard.NewTextBlockCard(expectedText, expectedTitle)

	// Assert that no more or less than two elements (one for title, one for
	// "body") are present.
	if len(textCard.Body) != 2 {
		t.Fatalf(
			"want Card Body length of %d; got %d",
			2,
			len(textCard.Body),
		)
	}

	// Assert specific type of element for title
	if textCard.Body[0].Type != adaptivecard.TypeElementTextBlock {
		t.Fatalf(
			"want Card Body element type %q; got %q",
			adaptivecard.TypeElementTextBlock,
			textCard.Body[0].Type,
		)
	}

	// Assert specific type of element for text "body"
	if textCard.Body[1].Type != adaptivecard.TypeElementTextBlock {
		t.Fatalf(
			"want Card Body element type %q; got %q",
			adaptivecard.TypeElementTextBlock,
			textCard.Body[1].Type,
		)
	}

	// Assert expected title "body" text
	if textCard.Body[0].Text != expectedTitle {
		t.Fatalf(
			"want %q for Text field in element %d; got %q",
			expectedTitle,
			0,
			textCard.Body[0].Text,
		)
	}

	// Assert expected "body" text
	if textCard.Body[1].Text != expectedText {
		t.Fatalf(
			"want %q for Text field in element %d; got %q",
			expectedText,
			1,
			textCard.Body[1].Text,
		)
	}

}
