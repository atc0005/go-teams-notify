// Copyright 2022 Adam Chalkley
//
// https://github.com/atc0005/go-teams-notify
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package adaptivecard

import (
	"fmt"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
)

// validater is the interface shared by all supported types which provide
// validation of their fields.
type validater interface {
	Validate() error
}

// validator is used to perform validation of given values. Each validation
// method for this type is designed to exit early in order to preserve any
// prior validation failure. If a previous validation check failure occurred,
// the most recent validation check result will
//
// After performing a validation check, the caller
// is responsible for checking the result to determine if further validation
// checks should be performed.
//
// Credit: https://stackoverflow.com/a/23960293/903870
type validator struct {
	err error
}

// MustBeGreaterThan asserts that the given value is greater than the
// specified highest value.
//
// A true value is returned if the validation step passed. A false value is
// returned false if this or a prior validation step failed.
// func (v *validator) MustBeGreaterThan(highest int, value int) bool {
// 	if v.err != nil {
// 		return false
// 	}
// 	if value <= highest {
// 		v.err = fmt.Errorf("Must be Greater than %d", highest)
// 		return false
// 	}
// 	return true
// }

// MustSelfValidate asserts that each given item can self-validate.
//
// A true value is returned if the validation step passed. A false value is
// returned false if this or a prior validation step failed.
func (v *validator) MustSelfValidate(items ...validater) bool {
	if v.err != nil {
		return false
	}
	for _, item := range items {
		if err := item.Validate(); err != nil {
			v.err = err
			return false
		}
	}
	return true
}

// MustBeNotEmptyValue asserts that fieldVal is not empty. fieldValDesc
// describes the field value being validated (e.g., "Type") and typeDesc
// describes the specific struct or value type whose field we are validating
// (e.g., "Element").
//
// A true value is returned if the validation step passed. A false value is
// returned if this or a prior validation step failed.
func (v *validator) MustBeNotEmptyValue(fieldVal string, fieldValDesc string, typeDesc string) bool {
	if v.err != nil {
		return false
	}
	if fieldVal == "" {
		v.err = fmt.Errorf(
			"required %s is empty for %s: %w",
			fieldValDesc,
			typeDesc,
			ErrMissingValue,
		)
		return false
	}
	return true
}

// func (v *validator) MustBeNotEmptyValYIfValXIs(x string,  y string) bool {
// 	if v.err != nil {
// 		return false
// 	}
// 	if x != "" {
// 		v.err = fmt.Errorf("value Must not be Empty")
// 		return false
// 	}
// 	return true
// }

// MustBeInListIfFieldValNotEmpty reports whether fieldVal is in validVals if
// fieldVal is not empty. fieldValDesc describes the field value being
// validated (e.g., "Type") and typeDesc describes the specific struct or
// value type whose field we are validating (e.g., "Element").
//
// A true value is returned if fieldVal is empty or is in validVals. A false
// value is returned if a prior validation step failed or if fieldVal is not
// empty and is not in validVals.
func (v *validator) MustBeInListIfFieldValNotEmpty(fieldVal string, fieldValDesc string, typeDesc string, validVals []string, baseErr error) bool {
	switch {
	case v.err != nil:
		return false

	case fieldVal != "" && !goteamsnotify.InList(fieldVal, validVals, false):
		v.err = fmt.Errorf(
			"invalid %s %q for %s; expected one of %v",
			fieldValDesc,
			fieldVal,
			typeDesc,
			validVals,
		)

		if baseErr != nil {
			v.err = fmt.Errorf(
				"invalid %s %q for %s; expected one of %v: %w",
				fieldValDesc,
				fieldVal,
				typeDesc,
				validVals,
				baseErr,
			)
		}

		return false

	// Validation is good.
	default:
		return true
	}
}

func (v *validator) MustBeNotEmptyCollection(fieldValueDesc string, typeDesc string, items ...interface{}) bool {
	if v.err != nil {
		return false
	}
	if len(items) == 0 {
		// v.err = fmt.Errorf("value Must not be Empty")
		v.err = fmt.Errorf(
			"required %s collection is empty for %s: %w",
			fieldValueDesc,
			typeDesc,
			ErrMissingValue,
		)
		return false
	}
	return true
}

// IsValid indicates whether validation checks performed thus far have all
// passed.
func (v *validator) IsValid() bool {
	return v.err != nil
}
