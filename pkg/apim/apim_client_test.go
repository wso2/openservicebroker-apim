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

package apim

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/wso2/service-broker-apim/pkg/config"
)

const (
	publisherTestEndpoint         = "https://localhost:9443"
	StoreTestEndpoint             = "https://localhost:9443"
	StoreApplicationContext       = "/api/am/store/v0.14/applications"
	StoreSubscriptionContext      = "/api/am/store/v0.14/subscriptions"
	MultipleSubscriptionContext   = StoreSubscriptionContext + "/multiple"
	GenerateApplicationKeyContext = StoreApplicationContext + "/generate-keys"
	PublisherAPIContext           = "/api/am/publisher/v0.14/apis"
	successTestCase               = "success test case"
	failureTestCase               = "failure test case"
	ErrMsgTestIncorrectResult     = "expected value: %v but then returned value: %v"
)

type MockTokenManager struct {
}

func (m *MockTokenManager) Token() (string, error) {
	return "token", nil
}

func (m *MockTokenManager) Init(scopes []string) {

}

func init() {
	Init(&MockTokenManager{}, config.APIM{
		StoreEndpoint:                    StoreTestEndpoint,
		StoreApplicationContext:          StoreApplicationContext,
		StoreSubscriptionContext:         StoreSubscriptionContext,
		StoreMultipleSubscriptionContext: MultipleSubscriptionContext,
		PublisherAPIContext:              PublisherAPIContext,
		PublisherEndpoint:                publisherTestEndpoint,
		GenerateApplicationKeyContext:    GenerateApplicationKeyContext,
	})

}

func TestCreateApplication(t *testing.T) {
	t.Run(successTestCase, testCreateApplicationSuccessFunc())
	t.Run(failureTestCase, testCreateApplicationFailFunc())
}

func testCreateApplicationFailFunc() func(t *testing.T) {
	return func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		responder, err := httpmock.NewJsonResponder(http.StatusInternalServerError, nil)
		if err != nil {
			t.Error(err)
		}
		httpmock.RegisterResponder(http.MethodPost, StoreTestEndpoint+StoreApplicationContext, responder)

		_, err = CreateApplication(&ApplicationCreateReq{})
		if err == nil {
			t.Error("Expecting an error with code: " + strconv.Itoa(http.StatusInternalServerError))
		}
	}
}

func testCreateApplicationSuccessFunc() func(t *testing.T) {
	return func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		responder, err := httpmock.NewJsonResponder(http.StatusCreated, &AppCreateRes{ApplicationID: "1"})
		if err != nil {
			t.Error(err)
		}
		httpmock.RegisterResponder(http.MethodPost, StoreTestEndpoint+StoreApplicationContext, responder)

		id, err := CreateApplication(&ApplicationCreateReq{
			Name: "test",
		})
		if id != "1" {
			t.Errorf(ErrMsgTestIncorrectResult, "1", id)
		}
	}
}

func TestUpdateApplication(t *testing.T) {
	t.Run(successTestCase, testUpdateApplicationSuccessFunc())
	t.Run(failureTestCase, testUpdateApplicationFailFunc())
}

func testUpdateApplicationFailFunc() func(t *testing.T) {
	return func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		responder, err := httpmock.NewJsonResponder(http.StatusNotFound, nil)
		if err != nil {
			t.Error(err)
		}
		httpmock.RegisterResponder(http.MethodPut, StoreTestEndpoint+StoreApplicationContext+"/id", responder)

		err = UpdateApplication("id", &ApplicationCreateReq{
			Name: "test",
		})
		if err == nil {
			t.Error("Expecting an error with code: " + strconv.Itoa(http.StatusNotFound))
		}
	}
}

func testUpdateApplicationSuccessFunc() func(t *testing.T) {
	return func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		responder, err := httpmock.NewJsonResponder(http.StatusOK, nil)
		if err != nil {
			t.Error(err)
		}
		httpmock.RegisterResponder(http.MethodPut, StoreTestEndpoint+StoreApplicationContext+"/id", responder)

		err = UpdateApplication("id", &ApplicationCreateReq{
			Name: "test",
		})
		if err != nil {
			t.Error(err)
		}
	}
}

func TestGenerateKeys(t *testing.T) {
	t.Run(successTestCase, testGenerateKeysSuccessFunc())
	t.Run(failureTestCase, testGenerateKeysFailFunc())
}

func testGenerateKeysFailFunc() func(t *testing.T) {
	return func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		responder, err := httpmock.NewJsonResponder(http.StatusInternalServerError, nil)
		if err != nil {
			t.Error(err)
		}
		httpmock.RegisterResponder(http.MethodPost, StoreTestEndpoint+GenerateApplicationKeyContext, responder)

		_, err = GenerateKeys("")
		if err.Error() != ErrMsgAPPIDEmpty {
			t.Error("Expecting an error : " + ErrMsgAPPIDEmpty + " got: " + err.Error())
		}
		_, err = GenerateKeys("")
		if err == nil {
			t.Error("Expecting an error with code: " + strconv.Itoa(http.StatusInternalServerError))
		}
	}
}

func testGenerateKeysSuccessFunc() func(t *testing.T) {
	return func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		responder, err := httpmock.NewJsonResponder(http.StatusOK, &ApplicationKeyResp{
			Token: &Token{AccessToken: "abc"},
		})
		if err != nil {
			t.Error(err)
		}
		httpmock.RegisterResponder(http.MethodPost, StoreTestEndpoint+GenerateApplicationKeyContext, responder)

		got, err := GenerateKeys("123")
		if err != nil {
			t.Error(err)
		}
		if got.Token.AccessToken != "abc" {
			t.Errorf(ErrMsgTestIncorrectResult, "abc", got.Token.AccessToken)
		}
	}
}

func TestCreateMultipleSubscription(t *testing.T) {
	t.Run(successTestCase, testCreateMultipleSubscriptionSuccessFunc())
	t.Run(failureTestCase, testCreateMultipleSubscriptionFailFunc())
}

func testCreateMultipleSubscriptionFailFunc() func(t *testing.T) {
	return func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		responder, err := httpmock.NewJsonResponder(http.StatusInternalServerError, nil)
		if err != nil {
			t.Error(err)
		}
		httpmock.RegisterResponder(http.MethodPost, StoreTestEndpoint+MultipleSubscriptionContext, responder)

		_, err = CreateMultipleSubscriptions([]SubscriptionReq{
			{
				APIIdentifier: "a",
				ApplicationID: "b",
				Tier:          "c",
			},
		})
		if err == nil {
			t.Error("Expecting an error with code: " + strconv.Itoa(http.StatusInternalServerError))
		}
	}
}

func testCreateMultipleSubscriptionSuccessFunc() func(t *testing.T) {
	return func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		responder, err := httpmock.NewJsonResponder(http.StatusOK, []SubscriptionResp{
			{
				SubscriptionID: "abc",
			},
			{
				SubscriptionID: "abc1",
			},
		})
		if err != nil {
			t.Error(err)
		}
		httpmock.RegisterResponder(http.MethodPost, StoreTestEndpoint+MultipleSubscriptionContext, responder)
		got, err := CreateMultipleSubscriptions([]SubscriptionReq{
			{
				APIIdentifier: "a",
				ApplicationID: "b",
				Tier:          "c",
			},
			{
				APIIdentifier: "a1",
				ApplicationID: "b1",
				Tier:          "c1",
			},
		})
		if err != nil {
			t.Error(err)
		}
		if got[0].SubscriptionID != "abc" {
			t.Errorf(ErrMsgTestIncorrectResult, "abc", got)
		}
		if got[1].SubscriptionID != "abc1" {
			t.Errorf(ErrMsgTestIncorrectResult, "abc", got)
		}
	}
}

func TestUnSubscribe(t *testing.T) {
	t.Run(successTestCase, testUnSubscribeSuccessFunc())
	t.Run(failureTestCase, testUnSubscribeFailFunc())
}

func testUnSubscribeFailFunc() func(t *testing.T) {
	return func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		responder, err := httpmock.NewJsonResponder(http.StatusInternalServerError, nil)
		if err != nil {
			t.Error(err)
		}
		httpmock.RegisterResponder(http.MethodDelete, StoreTestEndpoint+StoreSubscriptionContext+"/abc", responder)

		err = UnSubscribe("abc")
		if err == nil {
			t.Error("Expecting an error with code: " + strconv.Itoa(http.StatusInternalServerError))
		}
	}
}

func testUnSubscribeSuccessFunc() func(t *testing.T) {
	return func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		responder, err := httpmock.NewJsonResponder(http.StatusOK, nil)
		if err != nil {
			t.Error(err)
		}
		httpmock.RegisterResponder(http.MethodDelete, StoreTestEndpoint+StoreSubscriptionContext+"/abc", responder)

		err = UnSubscribe("abc")
		if err != nil {
			t.Error(err)
		}
	}
}

func TestDeleteApplication(t *testing.T) {
	t.Run(successTestCase, testDeleteApplicationSuccessFunc())
	t.Run(failureTestCase, testDeleteApplicationFailFunc())
}

func testDeleteApplicationFailFunc() func(t *testing.T) {
	return func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		responder, err := httpmock.NewJsonResponder(http.StatusInternalServerError, nil)
		if err != nil {
			t.Error(err)
		}
		httpmock.RegisterResponder(http.MethodDelete, StoreTestEndpoint+StoreApplicationContext+"/abc", responder)

		err = DeleteApplication("abc")
		if err == nil {
			t.Error("Expecting an error with code: " + strconv.Itoa(http.StatusInternalServerError))
		}
	}
}

func testDeleteApplicationSuccessFunc() func(t *testing.T) {
	return func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		responder, err := httpmock.NewJsonResponder(http.StatusOK, nil)
		if err != nil {
			t.Error(err)
		}
		httpmock.RegisterResponder(http.MethodDelete, StoreTestEndpoint+StoreApplicationContext+"/abc", responder)

		err = DeleteApplication("abc")
		if err != nil {
			t.Error(err)
		}
	}
}

func TestDeleteAPI(t *testing.T) {
	t.Run(successTestCase, testDeleteAPISuccessFunc())
	t.Run(failureTestCase, testDeleteAPIFailFunc())
}

func testDeleteAPIFailFunc() func(t *testing.T) {
	return func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		responder, err := httpmock.NewJsonResponder(http.StatusInternalServerError, nil)
		if err != nil {
			t.Error(err)
		}
		httpmock.RegisterResponder(http.MethodDelete, publisherTestEndpoint+PublisherAPIContext+"/abc", responder)

		err = DeleteAPI("abc")
		if err == nil {
			t.Error("Expecting an error with code: " + strconv.Itoa(http.StatusInternalServerError))
		}
	}
}

func testDeleteAPISuccessFunc() func(t *testing.T) {
	return func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		responder, err := httpmock.NewJsonResponder(http.StatusOK, nil)
		if err != nil {
			t.Error(err)
		}
		httpmock.RegisterResponder(http.MethodDelete, publisherTestEndpoint+PublisherAPIContext+"/abc", responder)

		err = DeleteAPI("abc")
		if err != nil {
			t.Error(err)
		}
	}
}

func TestSearchAPIByNameVersion(t *testing.T) {
	t.Run(successTestCase, testSearchAPIByNameVersionSuccessFunc())
	t.Run("failure test case 1", testSearchAPIByNameVersionFail1Func())
	t.Run("failure test case 2", testSearchAPIByNameVersionFail2Func())
}

func testSearchAPIByNameVersionSuccessFunc() func(t *testing.T) {
	return func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		responder, err := httpmock.NewJsonResponder(http.StatusOK, &APISearchResp{
			Count: 1,
			List: []APISearchInfo{{
				ID: "111-111",
			}},
		})
		if err != nil {
			t.Error(err)
		}
		httpmock.RegisterResponder(http.MethodGet, publisherTestEndpoint+PublisherAPIContext+"?query=name%3ATest+version%3Av1", responder)
		apiID, err := SearchAPIByNameVersion("Test", "v1")
		if err != nil {
			t.Error(err)
		}
		if apiID != "111-111" {
			t.Errorf(ErrMsgTestIncorrectResult, "111-111", apiID)
		}
	}
}

func testSearchAPIByNameVersionFail1Func() func(t *testing.T) {
	return func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		responder, err := httpmock.NewJsonResponder(http.StatusOK, &APISearchResp{
			Count: 0,
		})
		if err != nil {
			t.Error(err)
		}
		httpmock.RegisterResponder(http.MethodGet, publisherTestEndpoint+PublisherAPIContext+"?query=name%3ATest+version%3Av1", responder)
		_, err = SearchAPIByNameVersion("Test", "v1")
		if err == nil {
			t.Error("Expecting an error")
		}
		if err.Error() != "couldn't find the API Test" {
			t.Error("Expecting the error 'couldn't find the API Test' but got " + err.Error())
		}
	}
}
func testSearchAPIByNameVersionFail2Func() func(t *testing.T) {
	return func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		responder, err := httpmock.NewJsonResponder(http.StatusOK, &APISearchResp{
			Count: 2,
		})
		if err != nil {
			t.Error(err)
		}
		httpmock.RegisterResponder(http.MethodGet, publisherTestEndpoint+PublisherAPIContext+"?query=name%3ATest+version%3Av1", responder)
		_, err = SearchAPIByNameVersion("Test", "v1")
		if err == nil {
			t.Error("Expecting an error")
		}
		if err.Error() != "returned more than one API for API Test" {
			t.Error("Expecting the error 'returned more than one API for API Test' but got " + err.Error())
		}
	}
}

func TestSearchApplication(t *testing.T) {
	t.Run(successTestCase, testSearchApplicationSuccessFunc())
	t.Run("failure test case 1", testSearchApplicationFail1Func())
	t.Run("failure test case 2", testSearchApplicationFail2Func())
}

func testSearchApplicationSuccessFunc() func(t *testing.T) {
	return func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		responder, err := httpmock.NewJsonResponder(http.StatusOK, &ApplicationSearchResp{
			Count: 1,
			List: []ApplicationSearchInfo{{
				ApplicationID: "111-111",
			}},
		})
		if err != nil {
			t.Error(err)
		}
		httpmock.RegisterResponder(http.MethodGet, StoreTestEndpoint+StoreApplicationContext+"?query=Test", responder)
		apiID, err := SearchApplication("Test")
		if err != nil {
			t.Error(err)
		}
		if apiID != "111-111" {
			t.Errorf(ErrMsgTestIncorrectResult, "Test", apiID)
		}
	}
}

func testSearchApplicationFail1Func() func(t *testing.T) {
	return func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		responder, err := httpmock.NewJsonResponder(http.StatusOK, &APISearchResp{
			Count: 0,
		})
		if err != nil {
			t.Error(err)
		}
		httpmock.RegisterResponder(http.MethodGet, StoreTestEndpoint+StoreApplicationContext+"?query=Test", responder)
		_, err = SearchApplication("Test")
		if err == nil {
			t.Error("Expecting an error")
		}
		if err.Error() != "couldn't find the Application Test" {
			t.Error("Expecting the error 'couldn't find the Application Test' but got " + err.Error())
		}
	}
}

func testSearchApplicationFail2Func() func(t *testing.T) {
	return func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		responder, err := httpmock.NewJsonResponder(http.StatusOK, &ApplicationSearchResp{
			Count: 2,
		})
		if err != nil {
			t.Error(err)
		}
		httpmock.RegisterResponder(http.MethodGet, StoreTestEndpoint+StoreApplicationContext+"?query=Test", responder)
		_, err = SearchApplication("Test")
		if err == nil {
			t.Error("Expecting an error")
		}
		if err.Error() != "returned more than one Application for Test" {
			t.Error("Expecting the error 'returned more than one Application for Test' but got " + err.Error())
		}
	}
}
