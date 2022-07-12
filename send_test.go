// Copyright 2020 Enrico Hoffmann
// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/go-teams-notify
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package goteamsnotify

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	assert.IsType(t, &teamsClient{}, client)
}

func TestTeamsClientSend(t *testing.T) {
	// THX@Hassansin ... http://hassansin.github.io/Unit-Testing-http-client-in-Go
	simpleMsgCard := NewMessageCard()
	simpleMsgCard.Text = "Hello World"
	var tests = []struct {
		reqURL                string
		reqMsg                MessageCard
		resBody               string // httpClient response body text
		resError              error  // httpClient error
		error                 error  // method error
		validationURLPatterns []string
		skipURLVal            bool // whether webhook URL validation is applied (e.g., GH-68)
		resStatus             int  // httpClient response status
	}{
		// invalid webhookURL - url.Parse error
		{
			reqURL:     "ht\ttp://",
			reqMsg:     simpleMsgCard,
			resStatus:  0,
			resBody:    "invalid",
			resError:   nil,
			error:      &url.Error{},
			skipURLVal: false,
		},
		// invalid webhookURL - missing prefix in webhook URL
		{
			reqURL:     "",
			reqMsg:     simpleMsgCard,
			resStatus:  0,
			resBody:    "invalid",
			resError:   nil,
			error:      ErrWebhookURLUnexpected,
			skipURLVal: false,
		},
		// invalid httpClient.Do call
		{
			reqURL:     "https://outlook.office.com/webhook/xxx",
			reqMsg:     simpleMsgCard,
			resStatus:  200,
			resBody:    http.StatusText(http.StatusOK),
			resError:   errors.New("pling"),
			error:      &url.Error{},
			skipURLVal: false,
		},
		// invalid httpClient.Do call
		{
			reqURL:     "https://outlook.office365.com/webhook/xxx",
			reqMsg:     simpleMsgCard,
			resStatus:  200,
			resBody:    http.StatusText(http.StatusOK),
			resError:   errors.New("pling"),
			error:      &url.Error{},
			skipURLVal: false,
		},
		// invalid response status code
		{
			reqURL:     "https://outlook.office.com/webhook/xxx",
			reqMsg:     simpleMsgCard,
			resStatus:  400,
			resBody:    http.StatusText(http.StatusBadRequest),
			resError:   nil,
			error:      errors.New(""),
			skipURLVal: false,
		},
		// invalid response status code
		{
			reqURL:     "https://outlook.office365.com/webhook/xxx",
			reqMsg:     simpleMsgCard,
			resStatus:  400,
			resBody:    http.StatusText(http.StatusBadRequest),
			resError:   nil,
			error:      errors.New(""),
			skipURLVal: false,
		},
		// valid
		{
			reqURL:     "https://outlook.office.com/webhook/xxx",
			reqMsg:     simpleMsgCard,
			resStatus:  200,
			resBody:    ExpectedWebhookURLResponseText,
			resError:   nil,
			error:      nil,
			skipURLVal: false,
		},
		// valid
		{
			reqURL:     "https://outlook.office365.com/webhook/xxx",
			reqMsg:     simpleMsgCard,
			resStatus:  200,
			resBody:    ExpectedWebhookURLResponseText,
			resError:   nil,
			error:      nil,
			skipURLVal: false,
		},
		// custom webhook domain (e.g., GH-68) without disabling validation
		//
		// This test case should not result in an actual client request going
		// out as validation failure should occur.
		{
			reqURL:     "https://example.webhook.office.com/webhook/xxx",
			reqMsg:     simpleMsgCard,
			resStatus:  0,
			resBody:    "",
			resError:   nil,
			error:      ErrWebhookURLUnexpected,
			skipURLVal: false,
		},
		// custom webhook domain (e.g., GH-68) with validation disabled
		//
		// This is expected to succeed, provided that the actual webhook URL
		// is valid. GH-68 indicates that private webhook endpoints exist, but
		// without knowing the names or valid patterns, this is about all we
		// can do for now?
		{
			reqURL:     "https://example.webhook.office.com/webhook/xxx",
			reqMsg:     simpleMsgCard,
			resStatus:  200,
			resBody:    ExpectedWebhookURLResponseText,
			resError:   nil,
			error:      nil,
			skipURLVal: true,
		},
		// custom webhook domain with custom validation patterns
		{
			reqURL:                "https://arbitrary.domain.com/webhook/xxx",
			reqMsg:                simpleMsgCard,
			resStatus:             200,
			resBody:               ExpectedWebhookURLResponseText,
			resError:              nil,
			error:                 nil,
			skipURLVal:            false,
			validationURLPatterns: []string{DefaultWebhookURLValidationPattern, "arbitrary.domain.com"},
		},
		// custom webhook domain with custom validation patterns not matching requirements
		{
			reqURL:                "https://arbitrary.test.com/webhook/xxx",
			reqMsg:                simpleMsgCard,
			resStatus:             200,
			resBody:               ExpectedWebhookURLResponseText,
			resError:              nil,
			error:                 ErrWebhookURLUnexpected,
			skipURLVal:            false,
			validationURLPatterns: []string{DefaultWebhookURLValidationPattern, "arbitrary.domain.com"},
		},
		// custom webhook domain with complex custom validation pattern matching requirements
		{
			reqURL:                "https://foo.domain.com/webhook/xxx",
			reqMsg:                simpleMsgCard,
			resStatus:             200,
			resBody:               ExpectedWebhookURLResponseText,
			resError:              nil,
			error:                 nil,
			skipURLVal:            false,
			validationURLPatterns: []string{`^https://.*\.domain\.com/.*$`},
		},
	}
	for idx, test := range tests {
		// Create range scoped var for use within closure
		test := test

		client := NewTestClient(func(req *http.Request) (*http.Response, error) {
			// Test request parameters
			assert.Equal(t, req.URL.String(), test.reqURL)

			// GH-46; fix contributed by @davecheney (thank you!)
			//
			// The RoundTripper documentation notes that nil must be returned
			// as the error value if a response is received. A non-nil error
			// should be returned for failure to obtain a response. Failure to
			// obtain a response is indicated by the test table response
			// error, so we represent that failure to obtain a response by
			// returning nil and the test table response error explaining why
			// a response could not be retrieved.
			if test.resError != nil {
				return nil, test.resError
			}

			// GH-46 (cont) If no table test response errors are provided,
			// then the response was retrieved (provided below), so we are
			// required to return nil as the error value along with the
			// response.
			return &http.Response{
				StatusCode: test.resStatus,

				// Send response to be tested
				Body: ioutil.NopCloser(bytes.NewBufferString(test.resBody)),

				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}, nil
		})
		c := &teamsClient{httpClient: client}
		c.AddWebhookURLValidationPatterns(test.validationURLPatterns...)

		// Disable webhook URL prefix validation if specified by table test
		// entry. See GH-68 for additional details.
		if test.skipURLVal {
			t.Log("Calling SkipWebhookURLValidationOnSend")
			c.SkipWebhookURLValidationOnSend(true)
		}

		err := c.Send(test.reqURL, test.reqMsg)
		if err != nil {

			// Assert that returned error chain contains specified table test
			// error.
			if !errors.Is(err, test.error) {

				// Debugging; trying to figure out how wrapping is applied.

				t.Log("calling initial errors.Unwrap")
				inner := errors.Unwrap(err)
				for {
					t.Logf("inner error: %T", inner)
					if inner == nil {
						break
					}

					t.Log("calling errors.Unwrap")
					inner = errors.Unwrap(inner)
				}

				t.Logf("test %d; got %s", idx, err)

				t.Fatalf(
					"FAIL: test %d; got %T, want %T",
					idx,
					err,
					test.error,
				)
			} else {
				t.Logf(
					"OK: test %d; test.error is of type %T, err is of type %T",
					idx,
					test.error,
					err,
				)
				t.Logf(
					"OK: test %d; test.error has value '%s'",
					idx,
					test.error.Error(),
				)
				t.Logf(
					"OK: test %d; error response has value '%s'",
					idx,
					err.Error(),
				)
			}
		} else {
			t.Logf("OK: test %d; no error; response body: '%s'", idx, test.resBody)
		}

		switch {

		// An error occurred, but table test entry indicates no error expected.
		case err != nil && test.error == nil:
			t.Logf("FAIL: test %d; error occurred, but none expected!", idx)
			t.Fatalf(
				"FAIL: test %d; got '%v', want '%v'",
				idx,
				err,
				test.error,
			)

		// No error occurred, but table test entry indicates one expected.
		case err == nil && test.error != nil:
			t.Logf("FAIL: test %d; no error occurred, but one was expected!", idx)
			t.Fatalf(
				"FAIL: test %d; got '%v', want '%v'",
				idx,
				err,
				test.error,
			)

		// No error occurred and table test entry indicates no error expected.
		case err == nil && test.error == nil:
			t.Logf("OK: test %d; no error occurred and table test entry indicates no error expected.", idx)

		// Error occurred and table test entry indicates one expected.
		case err != nil && test.error != nil:
			t.Logf("OK: test %d; error occurred and table test entry indicates one expected.", idx)

		}

	}
}

// helper for testing --------------------------------------------------------------------------------------------------

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) (*http.Response, error)

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// NewTestClient returns *http.API with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}
