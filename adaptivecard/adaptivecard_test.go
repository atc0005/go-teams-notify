// Copyright 2022 Adam Chalkley
//
// https://github.com/atc0005/go-teams-notify
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package adaptivecard

import (
	"errors"
	"testing"
)

/*

	TODO:

	Test Column.Validate()

		- nil pointer SelectAction
			- c.SelectAction.Validate()

	Test ColumnItems.Validate()

		- slice of pointers, some potentially nil
			- include valid pointers
			- include some nil pointers

			What happens for the value receiver?
			Does the Validate() method receive a zero value Element type?
			Perhaps modify Element.Validate() to emit the address of Element
				values it has been asked to validate?


*/

func TestColumnItems_ValidatePanicExpected(t *testing.T) {
	/*
		A ColumnSet contains Column values.
		A Column contains an Items value which is a []*Element.

		The []*Element could *potentially* contain a nil pointer. Validation
		should catch this without initiating a panic?

		The []*Element could *potentially* contain a zero value Element.
		Validation should catch this without initiating a panic?

		This specific test ignores the parent Column and its parent ColumnSet
		and instead focuses just on the ColumnItems validation behavior.
	*/

	// []*Element
	columnItems := make([]*Element, 0, 10)

	// A zero value Element.
	element1 := Element{}

	// A properly filled out Element. We opt to use a TextBlock Element.
	element2 := Element{
		Type: TypeElementTextBlock,

		// Not required, but we go ahead and fill it in.
		Text: "placeholder",
	}

	element3 := &Element{}

	columnItems = append(columnItems, &element1)
	columnItems = append(columnItems, &element2)
	columnItems = append(columnItems, nil) // Problem entry.
	columnItems = append(columnItems, element3)

	// Run validation using item "copy"
	for i, item := range columnItems {
		// --- FAIL: TestColumnItems_Validate (0.00s)
		// panic: runtime error: invalid memory address or nil pointer dereference [recovered]
		//         panic: runtime error: invalid memory address or nil pointer dereference
		// [signal 0xc0000005 code=0x0 addr=0x0 pc=0x100621e]
		//
		// github.com/atc0005/go-teams-notify/v2/adaptivecard.TestColumnItems_Validate(0xc00014c1a0)
		//         T:/github/go-teams-notify/adaptivecard/adaptivecard_test.go:77 +0x1be
		//
		//
		// NOTE: Because Validate() has a value receiver, it cannot be
		// responsible for ensuring that the pointer is not nil (this isn't
		// visible) to the method. Instead, the caller of that method has to
		// ensure that the attempt to dereference the pointer (which happens
		// when you call the method) can succeed.
		if err := item.Validate(); err != nil {
			t.Errorf("failed to validate item %d: %v", i, item)
		} else {
			t.Logf("successfully validated item %d: %v", i, item)
		}
	}

	// Run validation using original slice member. This has the same problem
	// as noted above; you deference the pointer when you call the Validate()
	// method (which uses a value receiver).
	for i := range columnItems {
		if err := columnItems[i].Validate(); err != nil {
			t.Error("failed to validate item:", columnItems[i])
		}
	}

}

func TestColumnItems_Validate(t *testing.T) {
	/*
		A ColumnSet contains Column values.
		A Column contains an Items value which is a []*Element.

		The []*Element could *potentially* contain a nil pointer. Validation
		should catch this without initiating a panic?

		The []*Element could *potentially* contain a zero value Element.
		Validation should catch this without initiating a panic?

		This specific test ignores the parent Column and its parent ColumnSet
		and instead focuses just on the ColumnItems validation behavior.
	*/

	// []*Element
	// columnItems := make(ColumnItems, 0, 10)
	columnItems := make([]*Element, 0, 10)

	// A zero value Element.
	element1 := Element{}

	// A properly filled out Element. We opt to use a TextBlock Element.
	element2 := Element{
		Type: TypeElementTextBlock,

		// Not required, but we go ahead and fill it in.
		Text: "placeholder",
	}

	element3 := &Element{}

	columnItems = append(columnItems, &element1)
	columnItems = append(columnItems, &element2)
	columnItems = append(columnItems, nil) // Problem entry.
	columnItems = append(columnItems, element3)

	// Explicitly convert []*Element so that we can ensure that the values are
	// handled properly.
	got := ColumnItems(columnItems).Validate()

	// We expect an error due to the nil value we've added to the collection.
	want := ErrMissingValue

	switch {
	case errors.Is(got, want):
		t.Log("validation failed for Column Items")
		t.Logf("got %q; want %q", got, want)
	case got == nil:
		t.Log("validation successful for Column Items")
		t.Errorf("got %q; want %q", got, want)
	default:
		t.Log("validation failed for Column Items")
		t.Errorf("got %q; want %q", got, want)
	}
}
