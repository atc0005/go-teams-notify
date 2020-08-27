// Copyright 2020 Enrico Hoffmann
// Copyright 2020 Adam Chalkley
//
// https:#github.com/atc0005/go-teams-notify
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
		reqURL    string
		reqMsg    MessageCard
		resStatus int    // httpClient response status
		resBody   string // httpClient response body text
		resError  error  // httpClient error
		error     error  // method error
	}{
		// invalid webhookURL - url.Parse error
		{
			reqURL:    "ht\ttp://",
			reqMsg:    simpleMsgCard,
			resStatus: 0,
			resBody:   "invalid",
			resError:  nil,
			error:     &url.Error{},
		},
		// invalid webhookURL - missing prefix in webhook URL
		{
			reqURL:    "",
			reqMsg:    simpleMsgCard,
			resStatus: 0,
			resBody:   "invalid",
			resError:  nil,
			error:     errors.New(""),
		},
		// invalid httpClient.Do call
		{
			reqURL:    "https://outlook.office.com/webhook/xxx",
			reqMsg:    simpleMsgCard,
			resStatus: 200,
			resBody:   http.StatusText(http.StatusOK),
			resError:  errors.New("pling"),
			error:     &url.Error{},
		},
		// invalid httpClient.Do call
		{
			reqURL:    "https://outlook.office365.com/webhook/xxx",
			reqMsg:    simpleMsgCard,
			resStatus: 200,
			resBody:   http.StatusText(http.StatusOK),
			resError:  errors.New("pling"),
			error:     &url.Error{},
		},
		// invalid response status code
		{
			reqURL:    "https://outlook.office.com/webhook/xxx",
			reqMsg:    simpleMsgCard,
			resStatus: 400,
			resBody:   http.StatusText(http.StatusBadRequest),
			resError:  nil,
			error:     errors.New(""),
		},
		// invalid response status code
		{
			reqURL:    "https://outlook.office365.com/webhook/xxx",
			reqMsg:    simpleMsgCard,
			resStatus: 400,
			resBody:   http.StatusText(http.StatusBadRequest),
			resError:  nil,
			error:     errors.New(""),
		},
		// valid
		{
			reqURL:    "https://outlook.office.com/webhook/xxx",
			reqMsg:    simpleMsgCard,
			resStatus: 200,
			resBody:   http.StatusText(http.StatusOK),
			resError:  nil,
			error:     nil,
		},
		// valid
		{
			reqURL:    "https://outlook.office365.com/webhook/xxx",
			reqMsg:    simpleMsgCard,
			resStatus: 200,
			resBody:   http.StatusText(http.StatusOK),
			resError:  nil,
			error:     nil,
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

		err := c.Send(test.reqURL, test.reqMsg)

		// BUG: This does not handle comparing wrapped errors (GH-23).
		// assert.IsType(t, test.error, err)

		// FIX: Current master branch of stretchr/testing provides
		// assert.ErrorAs(), but this isn't officially available until v1.7.x.
		//
		// Note: This provides a wrapper around Go 1.13+ stdlib error wrapping
		// functionality.
		//
		// if err != nil {
		// 	assert.ErrorAs(t, err, &test.error)
		// }

		// FIX: Use Go 1.13 stdlib errors.As() as a replacement for
		// assert.IsType() in order to provide a type assertion of wrapped
		// errors to table test errors.
		if err != nil {
			// FIXME: This won't work if the test.reqURL is well-formed, but
			// does not contain one of the two known valid prefixes.
			if !errors.As(err, &test.error) {
				t.Fatalf(
					"FAIL: test %d; got %T, want %T",
					idx,
					errors.Unwrap(err),
					test.error,
				)
			} else {
				t.Logf(
					"OK: test %d; test.error is of type %T, err is of type %T",
					idx,
					test.error,
					err,
				)
			}
		} else {
			t.Logf("OK: test %d; no error", idx)
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
