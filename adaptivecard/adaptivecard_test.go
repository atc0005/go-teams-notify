// Copyright 2022 Adam Chalkley
//
// https://github.com/atc0005/go-teams-notify
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package adaptivecard

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
