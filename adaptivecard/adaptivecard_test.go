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

func TestNewTextBlockCard(t *testing.T) {

	expectedText := "Tacos"

	// Create basic card.
	textCard := adaptivecard.NewTextBlockCard(expectedText)

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
