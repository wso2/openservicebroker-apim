/*
 * Copyright (c) 2019 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http:www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

// Package client contains functions required to make HTTP calls.
package client

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/wso2/service-broker-apim/pkg/config"
	"github.com/wso2/service-broker-apim/pkg/log"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"time"
)

const (
	HeaderAuth                  = "Authorization"
	HeaderBear                  = "Bearer "
	HTTPContentType             = "Content-Type"
	ContentTypeApplicationJSON  = "application/json"
	ContentTypeURLEncoded       = "application/x-www-form-urlencoded; param=value"
	ErrMsgUnableToCreateReq     = "unable to create request"
	ErrMsgUnableToParseReqBody  = "unable to parse request body"
	ErrMsgUnableToParseRespBody = "unable to parse response body, context: %s "
	ErrMsgUnableInitiateReq     = "unable to initiate request: %s"
	ErrMsgUnsuccessfulAPICall   = "unsuccessful API call: %s response Code: %s URL: %s"
	ErrMsgUnableToCloseBody     = "unable to close the body"
)

var ErrInvalidParameters = errors.New("invalid parameters")

// RetryPolicy defines a function which validate the response and apply desired policy
// to determine whether to retry the particular request or not.
type RetryPolicy func(resp *http.Response) bool

// BackOffPolicy policy determines the duration between two retires
type BackOffPolicy func(min, max time.Duration, attempt int) time.Duration

// Client represent the state of the HTTP client.
type Client struct {
	httpClient    *http.Client
	checkForReTry RetryPolicy
	backOff       BackOffPolicy
	minBackOff    time.Duration
	maxBackOff    time.Duration
	maxRetry      int
}

// default client
var client = &Client{
	httpClient:    http.DefaultClient,
	checkForReTry: isErrorResponse,
	backOff:       calculateBackOff,
	minBackOff:    1 * time.Second,
	maxBackOff:    60 * time.Second,
	maxRetry:      3,
}

// HTTPRequest wraps the http.request and the Body.
// Body is wrapped with io.ReadSeeker which allows to reset the body buffer reader to initial state in retires.
type HTTPRequest struct {
	body    io.ReadSeeker
	httpReq *http.Request
}

// HTTPRequest returns the HTTP request.
func (r *HTTPRequest) HTTPRequest() *http.Request {
	return r.httpReq
}

// SetHeader method set the given header key and value to the HTTP request.
func (r *HTTPRequest) SetHeader(k, v string) {
	r.httpReq.Header.Set(k, v)
}

// Configure overrides the default client values. This function should be called before calling Invoke method.
func Configure(c *config.Client) {
	client = &Client{
		httpClient: &http.Client{
			Timeout: time.Duration(c.Timeout) * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: c.InsecureCon},
			},
		},
		minBackOff:    time.Duration(c.MinBackOff) * time.Second,
		maxBackOff:    time.Duration(c.MaxBackOff) * time.Second,
		maxRetry:      c.MaxRetries,
		backOff:       calculateBackOff,
		checkForReTry: isErrorResponse,
	}
}

// InvokeError wraps more information about the error.
type InvokeError struct {
	err        error
	StatusCode int
}

func (e *InvokeError) Error() string {
	return e.err.Error()
}

// B64BasicAuth returns a base64 encoded value of "u:p" string and any error encountered.
func B64BasicAuth(u, p string) (string, error) {
	if u == "" || p == "" {
		return "", ErrInvalidParameters
	}
	d := u + ":" + p
	return base64.StdEncoding.EncodeToString([]byte(d)), nil
}

// ParseBody parse response body into the given struct.
// Must send the pointer to the response body.
// Returns any error encountered.
func ParseBody(res *http.Response, v interface{}) error {
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(b, v); err != nil {
		return err
	}
	return nil
}

// Invoke the request and parse the response body to the given struct.
// context parameter is used to maintain the request context in the log.
// resCode parameter is used to determine the desired response code.
// Returns any error encountered.
func Invoke(context string, req *HTTPRequest, body interface{}, expectedRespCode int) error {
	resp, err := do(req)
	if err != nil {
		return errors.Wrapf(err, ErrMsgUnableInitiateReq, context)
	}
	if resp.StatusCode != expectedRespCode {
		return &InvokeError{
			err:        errors.Errorf(ErrMsgUnsuccessfulAPICall, context, resp.Status, req.httpReq.URL),
			StatusCode: resp.StatusCode,
		}
	}

	// If response has a body
	if body != nil {
		defer func() {
			ld := log.NewData().
				Add("context", context).
				Add("URL", req.httpReq.URL)
			if err := resp.Body.Close(); err != nil {
				log.Error(ErrMsgUnableToCloseBody, err, ld)
			}
		}()

		err = ParseBody(resp, body)
		if err != nil {
			return &InvokeError{
				err:        errors.Wrapf(err, ErrMsgUnableToParseRespBody, context),
				StatusCode: resp.StatusCode,
			}
		}
	}
	return nil
}

// CreateHTTPPOSTRequest returns a POST HTTP request with a Bearer token header with the content type to application/json
// and any error encountered.
func CreateHTTPPOSTRequest(token, url string, body io.ReadSeeker) (*HTTPRequest, error) {
	req, err := CreateHTTPRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, errors.Wrap(err, ErrMsgUnableToCreateReq)
	}
	req.SetHeader(HeaderAuth, HeaderBear+token)
	req.SetHeader(HTTPContentType, ContentTypeApplicationJSON)
	return req, nil
}

// CreateHTTPPUTRequest returns a PUT HTTP request with a Bearer token header with the content type to application/json
// and any error encountered.
func CreateHTTPPUTRequest(token, url string, body io.ReadSeeker) (*HTTPRequest, error) {
	req, err := CreateHTTPRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, errors.Wrap(err, ErrMsgUnableToCreateReq)
	}
	req.SetHeader(HeaderAuth, HeaderBear+token)
	req.SetHeader(HTTPContentType, ContentTypeApplicationJSON)
	return req, nil
}

// CreateHTTPGETRequest returns a GET HTTP request with a Bearer token header
// and any error encountered.
func CreateHTTPGETRequest(token, url string) (*HTTPRequest, error) {
	req, err := CreateHTTPRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, ErrMsgUnableToCreateReq)
	}
	req.SetHeader(HeaderAuth, HeaderBear+token)
	return req, nil
}

// CreateHTTPDELETERequest returns a DELETE HTTP request with a Bearer token header
// and any error encountered.
func CreateHTTPDELETERequest(token, url string) (*HTTPRequest, error) {
	req, err := CreateHTTPRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, ErrMsgUnableToCreateReq)
	}
	req.SetHeader(HeaderAuth, HeaderBear+token)
	return req, nil
}

// BodyReader returns the byte buffer representation of the provided struct and any error encountered.
func BodyReader(v interface{}) (io.ReadSeeker, error) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(v)
	if err != nil {
		return nil, errors.Wrap(err, ErrMsgUnableToParseReqBody)
	}
	return bytes.NewReader(buf.Bytes()), nil
}

// CreateHTTPRequest returns client.HTTPRequest struct which wraps the http.request, request Body, and any error encountered.
func CreateHTTPRequest(method, url string, body io.ReadSeeker) (*HTTPRequest, error) {
	var rcBody io.ReadCloser
	if body != nil {
		rcBody = ioutil.NopCloser(body)
	}
	req, err := http.NewRequest(method, url, rcBody)
	if err != nil {
		return nil, err
	}
	return &HTTPRequest{httpReq: req, body: body}, nil
}

// do invokes the request and returns the response and, an error if exists.
// If the request is failed it will retry according to the registered Retry policy and Back off policy.
func do(req *HTTPRequest) (resp *http.Response, err error) {
	i := 1
	for ok := true; ok; ok = i <= client.maxRetry {
		resp, err = client.httpClient.Do(req.httpReq)
		// This error occurs due to  network connectivity problem and not for non 2xx responses.
		if err != nil {
			return nil, err
		}
		if !client.checkForReTry(resp) {
			break
		}

		logData := log.NewData().
			Add("url", req.httpReq.URL).
			Add("response code", resp.StatusCode)
		if req.body != nil {
			// Reset the body reader
			log.Debug("resetting the request body", logData)
			if _, err := req.body.Seek(0, 0); err != nil {
				return nil, err
			}
		}
		bt := client.backOff(client.minBackOff, client.maxBackOff, i)
		logData.
			Add("back off time", bt.Seconds()).
			Add("attempt", i)
		log.Debug("retrying the request", logData)
		time.Sleep(bt)
		i++
	}
	return resp, nil
}

// isErrorResponse will retry the request if the response code is 4XX or 5XX.
func isErrorResponse(resp *http.Response) bool {
	return resp.StatusCode >= 400
}

// calculateBackOff waits until attempt^2 or (min,max).
func calculateBackOff(min, max time.Duration, attempt int) time.Duration {
	du := math.Pow(2, float64(attempt))
	sleep := time.Duration(du) * time.Second
	if sleep < min {
		return min
	}
	if sleep > max {
		return max
	}
	return sleep
}
