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
		name                  string
		reqURL                string
		reqMsg                MessageCard
		resBody               string // httpClient response body text
		resError              error  // httpClient error
		error                 error  // method error
		validationURLPatterns []string
		skipURLVal            bool // whether webhook URL validation is applied (e.g., GH-68)
		resStatus             int  // httpClient response status
	}{
		{
			name:       "invalid webhookURL - url.Parse error",
			reqURL:     "ht\ttp://",
			reqMsg:     simpleMsgCard,
			resStatus:  0,
			resBody:    "invalid",
			resError:   nil,
			error:      &url.Error{},
			skipURLVal: false,
		},
		{
			name:       "invalid webhookURL - missing prefix in webhook URL",
			reqURL:     "",
			reqMsg:     simpleMsgCard,
			resStatus:  0,
			resBody:    "invalid",
			resError:   nil,
			error:      ErrWebhookURLUnexpected,
			skipURLVal: false,
		},
		{
			name:       "invalid httpClient.Do call using outlook.office.com URL",
			reqURL:     "https://outlook.office.com/webhook/xxx",
			reqMsg:     simpleMsgCard,
			resStatus:  200,
			resBody:    http.StatusText(http.StatusOK),
			resError:   errors.New("pling"),
			error:      &url.Error{},
			skipURLVal: false,
		},
		{
			name:       "invalid httpClient.Do call using outlook.office365.com URL",
			reqURL:     "https://outlook.office365.com/webhook/xxx",
			reqMsg:     simpleMsgCard,
			resStatus:  200,
			resBody:    http.StatusText(http.StatusOK),
			resError:   errors.New("pling"),
			error:      &url.Error{},
			skipURLVal: false,
		},
		{
			name:       "invalid response status code using outlook.office.com URL",
			reqURL:     "https://outlook.office.com/webhook/xxx",
			reqMsg:     simpleMsgCard,
			resStatus:  400,
			resBody:    http.StatusText(http.StatusBadRequest),
			resError:   nil,
			error:      errors.New(""),
			skipURLVal: false,
		},
		{
			name:       "invalid response status code using outlook.office365.com URL",
			reqURL:     "https://outlook.office365.com/webhook/xxx",
			reqMsg:     simpleMsgCard,
			resStatus:  400,
			resBody:    http.StatusText(http.StatusBadRequest),
			resError:   nil,
			error:      errors.New(""),
			skipURLVal: false,
		},
		{
			name:       "valid values using outlook.office.com URL",
			reqURL:     "https://outlook.office.com/webhook/xxx",
			reqMsg:     simpleMsgCard,
			resStatus:  200,
			resBody:    ExpectedWebhookURLResponseText,
			resError:   nil,
			error:      nil,
			skipURLVal: false,
		},
		{
			name:       "valid values using outlook.office365.com URL",
			reqURL:     "https://outlook.office365.com/webhook/xxx",
			reqMsg:     simpleMsgCard,
			resStatus:  200,
			resBody:    ExpectedWebhookURLResponseText,
			resError:   nil,
			error:      nil,
			skipURLVal: false,
		},
		{
			// This test case should not result in an actual client request
			// going out as validation failure should occur.
			name:       "custom webhook domain without disabling validation",
			reqURL:     "https://example.webhook.office.com/webhook/xxx",
			reqMsg:     simpleMsgCard,
			resStatus:  0,
			resBody:    "",
			resError:   nil,
			error:      ErrWebhookURLUnexpected,
			skipURLVal: false,
		},
		{
			// This is expected to succeed, provided that the actual webhook
			// URL is valid. GH-68 indicates that private webhook endpoints
			// exist, but without knowing the names or valid patterns, this is
			// about all we can do for now?
			name:       "custom webhook domain with validation disabled",
			reqURL:     "https://example.webhook.office.com/webhook/xxx",
			reqMsg:     simpleMsgCard,
			resStatus:  200,
			resBody:    ExpectedWebhookURLResponseText,
			resError:   nil,
			error:      nil,
			skipURLVal: true,
		},
		{
			name:                  "custom webhook domain with custom validation patterns matching requirements",
			reqURL:                "https://arbitrary.domain.com/webhook/xxx",
			reqMsg:                simpleMsgCard,
			resStatus:             200,
			resBody:               ExpectedWebhookURLResponseText,
			resError:              nil,
			error:                 nil,
			skipURLVal:            false,
			validationURLPatterns: []string{DefaultWebhookURLValidationPattern, "arbitrary.domain.com"},
		},
		{
			name:                  "custom webhook domain with custom validation patterns not matching requirements",
			reqURL:                "https://arbitrary.test.com/webhook/xxx",
			reqMsg:                simpleMsgCard,
			resStatus:             200,
			resBody:               ExpectedWebhookURLResponseText,
			resError:              nil,
			error:                 ErrWebhookURLUnexpected,
			skipURLVal:            false,
			validationURLPatterns: []string{DefaultWebhookURLValidationPattern, "arbitrary.domain.com"},
		},
		{
			name:                  "custom webhook domain with complex custom validation pattern matching requirements",
			reqURL:                "https://foo.domain.com/webhook/xxx",
			reqMsg:                simpleMsgCard,
			resStatus:             200,
			resBody:               ExpectedWebhookURLResponseText,
			resError:              nil,
			error:                 nil,
			skipURLVal:            false,
			validationURLPatterns: []string{`^https://.*\.domain\.com/.*$`},
		},
		{
			name:                  "ms teams workflow webhook request with default patterns",
			reqURL:                "https://default123xxxx.de.environment.api.powerplatform.com:443/powerautomate/automations/direct/workflows/1x2x3x4x5/triggers/manual/paths/invoke?api-version=1&sp=%2Ftriggers%2Fmanual%2Frun&sv=1.0&sig=xxxx",
			reqMsg:                simpleMsgCard,
			resStatus:             200,
			resBody:               ExpectedWebhookURLResponseText,
			resError:              nil,
			error:                 nil,
			skipURLVal:            false,
			validationURLPatterns: []string{DefaultWebhookURLValidationPattern, WorkflowURLBaseDomain},
		},
	}
	for idx, test := range tests {
		// Create range scoped var for use within closure
		test := test

		t.Run(test.name, func(t *testing.T) {

			client := NewTestClient(func(req *http.Request) (*http.Response, error) {
				// Test request parameters
				assert.Equal(t, req.URL.String(), test.reqURL)

				// GH-46; fix contributed by @davecheney (thank you!)
				//
				// The RoundTripper documentation notes that nil must be
				// returned as the error value if a response is received. A
				// non-nil error should be returned for failure to obtain a
				// response. Failure to obtain a response is indicated by the
				// test table response error, so we represent that failure to
				// obtain a response by returning nil and the test table
				// response error explaining why a response could not be
				// retrieved.
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

			// Disable webhook URL prefix validation if specified by table
			// test entry. See GH-68 for additional details.
			if test.skipURLVal {
				t.Log("Calling SkipWebhookURLValidationOnSend")
				c.SkipWebhookURLValidationOnSend(true)
			}

			err := c.Send(test.reqURL, test.reqMsg)
			switch {

			// An error occurred, but table test entry indicates no error expected.
			case err != nil && test.error == nil:
				t.Logf("FAIL: test %d; error occurred, but none expected!", idx)
				t.Fatalf(
					"\nFAIL: test %d\ngot %v\nwant %v",
					idx,
					err,
					test.error,
				)

			// No error occurred, but table test entry indicates one expected.
			case err == nil && test.error != nil:
				t.Logf("FAIL: test %d; no error occurred, but one was expected!", idx)
				t.Fatalf(
					"\nFAIL: test %d\ngot %v\nwant %v",
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

		})

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
