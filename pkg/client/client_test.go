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
package client

import (
	"encoding/json"
	"github.com/jarcoal/httpmock"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"testing"
)

type testVal struct {
	ID   int
	Name string
}

const (
	Host                      = "localhost"
	HTTPMockEndpoint          = "https://" + Host + "/api"
	Context                   = "testing"
	Token                     = "Token"
	PayloadID                 = 1
	PayloadName               = "test"
	PayloadString             = `{"id": 1, "name": "test"}`
	ErrMsgTestIncorrectResult = "expected value: %v but then returned value: %v"
)

var payload = testVal{
	ID:   PayloadID,
	Name: PayloadName,
}

func TestB64BasicAuth(t *testing.T) {
	_, err := B64BasicAuth("", "")
	if err == nil {
		t.Errorf("Expected error didn't occur")
	}
	re1, err := B64BasicAuth("admin", "admin")
	if err != nil {
		t.Error(err)
	}
	exp := "YWRtaW46YWRtaW4="
	if re1 != exp {
		t.Errorf(ErrMsgTestIncorrectResult, exp, re1)
	}
}

func TestPostHTTPRequestWrapper(t *testing.T) {
	b, err := BodyReader(payload)
	if err != nil {
		t.Error(err)
	}
	req, err := CreateHTTPPOSTRequest(Token, HTTPMockEndpoint, b)
	if err != nil {
		t.Error(err)
	}
	if req.httpReq.Method != http.MethodPost {
		t.Errorf(ErrMsgTestIncorrectResult, http.MethodPost, req.httpReq.Method)
	}
	if req.httpReq.Host != Host {
		t.Errorf(ErrMsgTestIncorrectResult, Host, req.httpReq.Host)
	}
	if req.httpReq.Header.Get(HeaderAuth) != (HeaderBear + Token) {
		t.Errorf(ErrMsgTestIncorrectResult, HeaderBear+Token, req.httpReq.Header.Get(HeaderAuth))
	}
	var val testVal
	err = json.NewDecoder(req.httpReq.Body).Decode(&val)
	if err != nil {
		t.Error(err)
	}
	if val.ID != PayloadID {
		t.Errorf(ErrMsgTestIncorrectResult, PayloadID, val.ID)
	}
}

func TestDeleteHTTPRequestWrapper(t *testing.T) {
	req, err := CreateHTTPDELETERequest(Token, "http://"+Host)
	if err != nil {
		t.Error(err)
	}
	if req.httpReq.Method != http.MethodDelete {
		t.Errorf(ErrMsgTestIncorrectResult, http.MethodDelete, req.httpReq.Method)
	}
	if req.httpReq.Host != Host {
		t.Errorf(ErrMsgTestIncorrectResult, Host, req.httpReq.Host)
	}
	if req.httpReq.Header.Get(HeaderAuth) != (HeaderBear + Token) {
		t.Errorf(ErrMsgTestIncorrectResult, HeaderBear+Token, req.httpReq.Header.Get(HeaderAuth))
	}
}

func TestInvoke(t *testing.T) {
	t.Run("success test case", testInvokeSuccessFunc())
	t.Run("failure test case", testInvokeFailFunc())
}

func testInvokeSuccessFunc() func(t *testing.T) {
	return func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("POST", HTTPMockEndpoint,
			httpmock.NewStringResponder(200, PayloadString))

		buf, err := BodyReader(PayloadString)
		if err != nil {
			t.Error(err)
		}

		req, err := CreateHTTPPOSTRequest(Token, HTTPMockEndpoint, buf)
		if err != nil {
			t.Error(err)
		}
		var body testVal
		err = Invoke(Context, req, &body, http.StatusOK)
		if err != nil {
			t.Error(err)
		}
		if body.ID != PayloadID {
			t.Errorf(ErrMsgTestIncorrectResult, PayloadID, body.ID)
		}
	}
}

func testInvokeFailFunc() func(t *testing.T) {
	return func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("POST", HTTPMockEndpoint,
			httpmock.NewStringResponder(http.StatusNotFound, ""))
		req, err := CreateHTTPPOSTRequest(Token, HTTPMockEndpoint, nil)
		if err != nil {
			t.Error(err)
		}
		err = Invoke(Context, req, nil, http.StatusOK)
		if err == nil {
			t.Errorf(ErrMsgTestIncorrectResult, "error response with code: "+
				strconv.Itoa(http.StatusNotFound), "reponse code: "+strconv.Itoa(http.StatusOK))
		} else {
			if e, ok := err.(*InvokeError); ok {
				if e.StatusCode != http.StatusNotFound {
					t.Errorf(ErrMsgTestIncorrectResult, "response code = "+
						strconv.Itoa(http.StatusNotFound), "reponse code = "+strconv.Itoa(http.StatusOK))
				}
			} else {
				t.Errorf(ErrMsgTestIncorrectResult, "type = "+reflect.TypeOf(InvokeError{}).Name(),
					"type = "+reflect.TypeOf(err).Name())
			}

		}

	}
}

func TestBodyReader(t *testing.T) {
	r, err := BodyReader(payload)
	if err != nil {
		t.Error(err)
	}
	var val testVal
	err = json.NewDecoder(r).Decode(&val)
	if err != nil {
		t.Error(err)
	}
	if val.ID != PayloadID {
		t.Errorf(ErrMsgTestIncorrectResult, PayloadID, val.ID)
	}
}

func TestParseBody(t *testing.T) {
	buf, err := BodyReader(payload)
	if err != nil {
		t.Error(err)
	}
	resp := &http.Response{
		Body: ioutil.NopCloser(buf),
	}
	var body testVal
	err = ParseBody(resp, &body)
	if err != nil {
		t.Error(err)
	}
	if body.ID != PayloadID {
		t.Errorf(ErrMsgTestIncorrectResult, PayloadID, body.ID)
	}
}
