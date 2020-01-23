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

// Package apim handles the interactions with APIM.
package apim

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"

	"github.com/pkg/errors"
	"github.com/wso2/service-broker-apim/pkg/client"
	"github.com/wso2/service-broker-apim/pkg/config"
	"github.com/wso2/service-broker-apim/pkg/log"
	"github.com/wso2/service-broker-apim/pkg/token"
	"github.com/wso2/service-broker-apim/pkg/utils"
)

const (
	CreateAPIContext                  = "create API"
	CreateApplicationContext          = "create application"
	CreateMultipleSubscriptionContext = "create multiple subscriptions"
	UpdateApplicationContext          = "update application"
	GenerateKeyContext                = "Generate application keys"
	UnSubscribeContext                = "unsubscribe api"
	ApplicationDeleteContext          = "delete application"
	APIDeleteContext                  = "delete API"
	APISearchContext                  = "search API"
	ApplicationSearchContext          = "search Application"
	ErrMsgAPPIDEmpty                  = "application id is empty"
)

var (
	publisherAPIEndpoint              string
	storeApplicationEndpoint          string
	storeSubscriptionEndpoint         string
	storeMultipleSubscriptionEndpoint string
	generateApplicationKeyEndpoint    string
	applicationDashBoardURLBase       string
	tokenManager                      token.Manager
	once                              sync.Once
)

// Init function initialize the API-M client. If there is an error, process exists with a panic.
func Init(manager token.Manager, conf config.APIM) {
	once.Do(func() {
		tokenManager = manager
		publisherAPIEndpoint = createEndpoint(conf.PublisherEndpoint, conf.PublisherAPIContext)
		storeApplicationEndpoint = createEndpoint(conf.StoreEndpoint, conf.StoreApplicationContext)
		storeSubscriptionEndpoint = createEndpoint(conf.StoreEndpoint, conf.StoreSubscriptionContext)
		storeMultipleSubscriptionEndpoint = createEndpoint(conf.StoreEndpoint, conf.StoreMultipleSubscriptionContext)
		generateApplicationKeyEndpoint = createEndpoint(conf.StoreEndpoint, conf.GenerateApplicationKeyContext)
		applicationDashBoardURLBase = createEndpoint(conf.StoreEndpoint, "/store/site/pages/application.jag")
	})
}

// createEndpoint returns a endpoint from the given paths.
func createEndpoint(paths ...string) string {
	endpoint, err := utils.ConstructURL(paths...)
	if err != nil {
		log.HandleErrorAndExit("cannot construct endpoint", err)
	}
	return endpoint
}

// CreateAPI function creates an API with the provided API spec.
// Returns the API ID and any error encountered.
func CreateAPI(reqBody *APIReqBody) (string, error) {
	req, err := creatHTTPPOSTAPIRequest(publisherAPIEndpoint, reqBody)
	var resBody APICreateResp
	err = send(CreateAPIContext, req, &resBody, http.StatusCreated)
	if err != nil {
		return "", err
	}
	return resBody.ID, nil
}

// GetAppDashboardURL returns DashBoard URL for the given Application.
func GetAppDashboardURL(appName string) string {
	q := url.Values{}
	q.Add("name", appName)
	return applicationDashBoardURLBase + "?" + q.Encode()
}

// CreateApplication creates an application with provided Application spec.
// Returns the Application ID and any error encountered.
func CreateApplication(reqBody *ApplicationCreateReq) (string, error) {
	req, err := creatHTTPPOSTAPIRequest(storeApplicationEndpoint, reqBody)
	if err != nil {
		return "", err
	}
	var resBody AppCreateRes
	err = send(CreateApplicationContext, req, &resBody, http.StatusCreated)
	if err != nil {
		return "", err
	}
	return resBody.ApplicationID, nil
}

// UpdateApplication updates an existing Application under the given ID with the provided Application spec.
// Returns any error encountered.
func UpdateApplication(id string, reqBody *ApplicationCreateReq) error {
	endpoint, err := utils.ConstructURL(storeApplicationEndpoint, id)
	if err != nil {
		return err
	}
	req, err := creatHTTPPUTAPIRequest(endpoint, reqBody)
	if err != nil {
		return err
	}
	err = send(UpdateApplicationContext, req, nil, http.StatusOK)
	if err != nil {
		return err
	}
	return nil
}

// GenerateKeys generates keys for the given application.
// Returns generated keys and any error encountered.
func GenerateKeys(appID string) (*ApplicationKeyResp, error) {
	if appID == "" {
		return nil, errors.New(ErrMsgAPPIDEmpty)
	}
	reqBody := defaultApplicationKeyGenerateReq()
	req, err := creatHTTPPOSTAPIRequest(generateApplicationKeyEndpoint, reqBody)
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Add("applicationId", appID)
	req.HTTPRequest().URL.RawQuery = q.Encode()

	var resBody ApplicationKeyResp
	err = send(GenerateKeyContext, req, &resBody, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return &resBody, nil
}

// CreateMultipleSubscriptions creates the given subscriptions.
// Returns list of SubscriptionResp and any error encountered.
func CreateMultipleSubscriptions(subs []SubscriptionReq) ([]SubscriptionResp, error) {
	req, err := creatHTTPPOSTAPIRequest(storeMultipleSubscriptionEndpoint, subs)
	if err != nil {
		return nil, err
	}
	resBody := make([]SubscriptionResp, 0)
	err = send(CreateMultipleSubscriptionContext, req, &resBody, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return resBody, nil
}

// UnSubscribe method removes the given subscription.
// Returns any error encountered.
func UnSubscribe(subscriptionID string) error {
	endpoint, err := utils.ConstructURL(storeSubscriptionEndpoint, subscriptionID)
	if err != nil {
		return err
	}
	req, err := creatHTTPDELETEAPIRequest(endpoint)
	if err != nil {
		return err
	}
	err = send(UnSubscribeContext, req, nil, http.StatusOK)
	if err != nil {
		return err
	}
	return nil
}

// DeleteApplication method deletes the given application.
// Returns any error encountered.
func DeleteApplication(applicationID string) error {
	endpoint, err := utils.ConstructURL(storeApplicationEndpoint, applicationID)
	if err != nil {
		return err
	}
	req, err := creatHTTPDELETEAPIRequest(endpoint)
	if err != nil {
		return err
	}
	err = send(ApplicationDeleteContext, req, nil, http.StatusOK)
	if err != nil {
		return err
	}
	return nil
}

// DeleteAPI method deletes the given API.
// Returns any error encountered.
func DeleteAPI(apiID string) error {
	endpoint, err := utils.ConstructURL(publisherAPIEndpoint, apiID)
	if err != nil {
		return err
	}
	req, err := creatHTTPDELETEAPIRequest(endpoint)
	if err != nil {
		return err
	}
	err = client.Invoke(APIDeleteContext, req, nil, http.StatusOK)
	if err != nil {
		return err
	}
	return nil
}

// send sends the given HTTP request, initialize the given response body if it is expected response code.
// Returns any error encountered.
func send(context string, req *client.HTTPRequest, resBody interface{}, expectedRespCode int) error {
	err := client.Invoke(context, req, resBody, expectedRespCode)
	if err != nil {
		return err
	}
	return nil
}

// getBodyReaderAndToken returns a token, a Reader for the given HTTP request body and any error encountered.
func getBodyReaderAndToken(reqBody interface{}) (string, io.ReadSeeker, error) {
	aT, err := tokenManager.Token()
	if err != nil {
		return "", nil, err
	}
	var bodyReader io.ReadSeeker
	if reqBody != nil {
		bodyReader, err = client.BodyReader(reqBody)
		if err != nil {
			return "", nil, err
		}
	}
	return aT, bodyReader, nil
}

func creatHTTPPOSTAPIRequest(endpoint string, reqBody interface{}) (*client.HTTPRequest, error) {
	aT, bodyReader, err := getBodyReaderAndToken(reqBody)
	if err != nil {
		return nil, err
	}
	req, err := client.CreateHTTPPOSTRequest(aT, endpoint, bodyReader)
	if err != nil {
		return nil, err
	}
	return req, err
}

func creatHTTPDELETEAPIRequest(endpoint string) (*client.HTTPRequest, error) {
	aT, err := tokenManager.Token()
	if err != nil {
		return nil, err
	}
	req, err := client.CreateHTTPDELETERequest(aT, endpoint)
	if err != nil {
		return nil, err
	}
	return req, err
}

// creatAPIMSearchHTTPRequest returns a API-M resource search request and any error encountered.
func creatAPIMSearchHTTPRequest(endpoint, query string) (*client.HTTPRequest, error) {
	aT, err := tokenManager.Token()
	if err != nil {
		return nil, err
	}
	req, err := client.CreateHTTPGETRequest(aT, endpoint)
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Add("query", query)
	req.HTTPRequest().URL.RawQuery = q.Encode()
	return req, err
}

func creatHTTPPUTAPIRequest(endpoint string, reqBody interface{}) (*client.HTTPRequest, error) {
	aT, bodyReader, err := getBodyReaderAndToken(reqBody)
	if err != nil {
		return nil, err
	}
	req, err := client.CreateHTTPPUTRequest(aT, endpoint, bodyReader)
	if err != nil {
		return nil, err
	}
	return req, err
}

// SearchAPIByNameVersion method returns API ID of the Given API.
// An error is returned if the number of result for the search is not equal to 1.
// Returns API ID and any error encountered.
func SearchAPIByNameVersion(apiName, version string) (string, error) {
	query := "name:" + apiName + " version:" + version
	req, err := creatAPIMSearchHTTPRequest(publisherAPIEndpoint, query)
	if err != nil {
		return "", err
	}
	var resp APISearchResp
	err = send(APISearchContext, req, &resp, http.StatusOK)
	if err != nil {
		return "", err
	}
	if resp.Count == 0 {
		return "", errors.New(fmt.Sprintf("couldn't find the API %s", apiName))
	}
	if resp.Count > 1 {
		return "", errors.New(fmt.Sprintf("returned more than one API for API %s", apiName))
	}
	return resp.List[0].ID, nil
}

// SearchApplication method returns Application ID of the Given Application.
// An error is returned if the number of result for the search is not equal to 1.
// Returns Application ID and any error encountered.
func SearchApplication(appName string) (string, error) {
	req, err := creatAPIMSearchHTTPRequest(storeApplicationEndpoint, appName)
	if err != nil {
		return "", err
	}
	var resp ApplicationSearchResp
	err = send(ApplicationSearchContext, req, &resp, http.StatusOK)
	if err != nil {
		return "", err
	}
	if resp.Count == 0 {
		return "", errors.New(fmt.Sprintf("couldn't find the Application %s", appName))
	}
	if resp.Count > 1 {
		return "", errors.New(fmt.Sprintf("returned more than one Application for %s", appName))
	}
	return resp.List[0].ApplicationID, nil
}

func defaultApplicationKeyGenerateReq() *ApplicationKeyGenerateRequest {
	return &ApplicationKeyGenerateRequest{
		ValidityTime:       "3600",
		KeyType:            "PRODUCTION",
		AccessAllowDomains: []string{"ALL"},
		Scopes:             []string{"am_application_scope", "default"},
		SupportedGrantTypes: []string{"urn:ietf:params:oauth:grant-type:saml2-bearer", "iwa:ntlm", "refresh_token",
			"client_credentials", "password"},
	}
}
