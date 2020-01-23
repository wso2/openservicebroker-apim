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

// Package token manages the Access token.
package token

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/wso2/service-broker-apim/pkg/client"
	"github.com/wso2/service-broker-apim/pkg/log"
	"github.com/wso2/service-broker-apim/pkg/utils"
)

const (
	SuffixSecond                    = "s"
	ErrMSGNotEnoughArgs             = "At least one scope should be present"
	ErrMSGUnableToGetClientCreds    = "Unable to get Client credentials"
	ErrMSGUnableToGetAccessToken    = "Unable to get access token for scopes: %v"
	ErrMSGUnableToParseExpireTime   = "Unable parse expiresIn time"
	ErrMsgUnableToParseRequestBody  = "unable to parse request body: %s"
	ErrMsgUnableToCreateRequestBody = "unable to create request body: %s"
	GenerateAccessToken             = "Generating access Token"
	DynamicClientRegMsg             = "Dynamic Client Reg"
	RefreshTokenContext             = "Refresh token"
	Context                         = "/token"
	UserName                        = "username"
	Password                        = "password"
	GrantPassword                   = "password"
	GrantRefreshToken               = "refresh_token"
	GrantType                       = "grant_type"
	Scope                           = "scope"
	RefreshToken                    = "refresh_token"
	ScopeSubscribe                  = "apim:subscribe"
	ScopeAPIView                    = "apim:api_view"
	LogKeyAT                        = "access-token"
	LogKeyRT                        = "refresh-token"
	LogKeyExpiresIn                 = "expires in"

	// CallBackURL is a dummy value
	CallBackURL = "www.dummy.com"

	// ClientName for dynamic client registration
	ClientName = "apim_service_broker"

	// DynamicClientRegGrantType for dynamic client registration
	DynamicClientRegGrantType = "password refresh_token"

	// Owner for dynamic client registration
	Owner = "admin"
)

// BasicCredentials represents the username and Password.
type BasicCredentials struct {
	Username string
	Password string
}

// DynamicClientRegReq represents the request for Dynamic client request body.
type DynamicClientRegReq struct {
	CallbackURL string `json:"callbackUrl"`
	ClientName  string `json:"clientName"`
	Owner       string `json:"owner"`
	GrantType   string `json:"grantType"`
	SaasApp     bool   `json:"saasApp"`
}

// DynamicClientRegResBody represents the message body for Dynamic client registration response body.
type DynamicClientRegResBody struct {
	CallbackURL       string `json:"callBackURL"`
	JSONString        string `json:"jsonString"`
	ClientName        string `json:"clientName"`
	ClientID          string `json:"clientId"`
	ClientSecret      string `json:"clientSecret"`
	IsSaasApplication bool   `json:"isSaasApplication"`
}

// Resp represents the message body of the token api response.
type Resp struct {
	Scope        string `json:"scope"`
	TokenTypes   string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

// token represent the Access token & Refresh token for a particular scope.
type token struct {
	lock         sync.RWMutex // ensures atomic writes to the following fields.
	accessToken  string
	refreshToken string
	expiresIn    time.Time
}

// PasswordRefreshTokenGrantManager is used to manage Access token using password and refresh_token grant type.
// Holds access token for the given scopes and regenerate token using refresh token if the access token is expired.
type PasswordRefreshTokenGrantManager struct {
	once                             sync.Once
	token                            *token
	clientID                         string
	clientSec                        string
	TokenEndpoint                    string
	DynamicClientEndpoint            string
	DynamicClientRegistrationContext string
	UserName                         string
	Password                         string
}

// Manager interface manages the token for a set of given scopes.
type Manager interface {
	// Init initialize the Token Manager. Generate token for the given scopes.
	// Must run before using the Token Manager.
	Init(scopes []string)

	// Token method returns an access token and any error occurred.
	Token() (string, error)
}

// Init initialize the Token Manager. Generate token for the given scopes.
// Must run before using the Token Manager.
func (m *PasswordRefreshTokenGrantManager) Init(scopes []string) {
	m.once.Do(func() {
		if len(scopes) == 0 {
			log.HandleErrorAndExit(ErrMSGNotEnoughArgs, nil)
		}
		err := m.registerDynamicClient(defaultClientRegBody())
		if err != nil {
			log.HandleErrorAndExit(ErrMSGUnableToGetClientCreds, err)
		}

		data := m.createAccessTokenReq(scopes)
		aT, rT, validPeriod, err := m.generateToken(data, GenerateAccessToken)
		if err != nil {
			log.HandleErrorAndExit(fmt.Sprintf(ErrMSGUnableToGetAccessToken, scopes), err)
		}
		// Handling the expire time of the access token
		duration, err := time.ParseDuration(strconv.Itoa(validPeriod) + "s")
		if err != nil {
			log.HandleErrorAndExit(ErrMSGUnableToParseExpireTime, err)
		}

		expiresIn := time.Now().Add(duration)
		m.token = &token{
			accessToken:  aT,
			refreshToken: rT,
			expiresIn:    expiresIn,
		}
		ld := log.NewData().
			Add(LogKeyAT, aT).
			Add(LogKeyExpiresIn, expiresIn).
			Add(LogKeyRT, rT)
		log.Debug(fmt.Sprintf("generated a token for scopes %v", scopes), ld)
	})
}

// createAccessTokenReq method returns access token request body for the given scopes.
// This functions assumes scopes array contains at least one element.
func (m *PasswordRefreshTokenGrantManager) createAccessTokenReq(scopes []string) url.Values {
	data := url.Values{}
	data.Set(UserName, m.UserName)
	data.Add(Password, m.Password)
	data.Add(GrantType, GrantPassword)
	var scopeVal = ""
	if scopes != nil && len(scopes) != 0 {
		scopeVal = scopes[0]
		for i := 1; i < len(scopes); i++ {
			scopeVal = scopeVal + " " + scopes[i]
		}
	}
	data.Add(Scope, scopeVal)
	return data
}

// createRefreshTokenReq method returns refresh token request body.
func createRefreshTokenReq(rT string) url.Values {
	data := url.Values{}
	data.Add(RefreshToken, rT)
	data.Add(GrantType, GrantRefreshToken)
	return data
}

// isExpired method returns true if the difference between the current time and given time is 10s.
func isExpired(expiresIn time.Time) bool {
	if time.Now().Sub(expiresIn) > (10 * time.Second) {
		return true
	}
	return false
}

// Token method returns an access token and any error occurred.
func (m *PasswordRefreshTokenGrantManager) Token() (string, error) {
	m.token.lock.RLock()
	ld := log.NewData().
		Add(LogKeyAT, m.token.accessToken).
		Add(LogKeyExpiresIn, m.token.expiresIn.String()).
		Add(LogKeyRT, m.token.refreshToken)
	if !isExpired(m.token.expiresIn) {
		log.Debug("access token is not expired", ld)
		aT := m.token.accessToken
		m.token.lock.RUnlock()
		return aT, nil
	}
	m.token.lock.RUnlock()
	m.token.lock.Lock()
	if !isExpired(m.token.expiresIn) {
		m.token.lock.Unlock()
		return m.token.accessToken, nil
	}
	ld = log.NewData().
		Add(LogKeyAT, m.token.accessToken).
		Add(LogKeyExpiresIn, m.token.expiresIn.String()).
		Add(LogKeyRT, m.token.refreshToken)
	log.Debug("access token is expired, re-generating", ld)
	aT, rT, expiresIn, err := m.generateRefreshToken(m.token.refreshToken)
	if err != nil {
		m.token.lock.Unlock()
		return "", err
	}
	//Parse time to type time.Duration
	duration, err := time.ParseDuration(strconv.Itoa(expiresIn) + SuffixSecond)
	ld = log.NewData().
		Add(LogKeyAT, aT).
		Add(LogKeyRT, rT).
		Add(LogKeyExpiresIn, expiresIn)
	log.Debug("new access token is generated", ld)
	m.token.refreshToken = rT
	m.token.accessToken = aT
	m.token.expiresIn = time.Now().Add(duration)

	m.token.lock.Unlock()
	return aT, nil
}

// generateRefreshToken method generates a new Access token and a Refresh token.
func (m *PasswordRefreshTokenGrantManager) generateRefreshToken(rTNow string) (aT, newRT string, expiresIn int, err error) {
	data := createRefreshTokenReq(rTNow)
	aT, rT, expiresIn, err := m.generateToken(data, RefreshTokenContext)
	if err != nil {
		return "", "", 0, err
	}
	return aT, rT, expiresIn, nil
}

// generateToken method returns an Access token and a Refresh token from given params.
func (m *PasswordRefreshTokenGrantManager) generateToken(reqBody url.Values, context string) (aT, rT string, expiresIn int, err error) {
	u, err := utils.ConstructURL(m.TokenEndpoint, Context)
	if err != nil {
		return "", "", 0, errors.Wrap(err, "cannot construct, token endpoint")
	}
	req, err := client.CreateHTTPRequest(http.MethodPost, u, bytes.NewReader([]byte(reqBody.Encode())))
	if err != nil {
		return "", "", 0, errors.Wrapf(err, ErrMsgUnableToCreateRequestBody,
			context)
	}
	req.HTTPRequest().SetBasicAuth(m.clientID, m.clientSec)
	req.SetHeader(client.HTTPContentType, client.ContentTypeURLEncoded)
	var resBody Resp
	if err := client.Invoke(context, req, &resBody, http.StatusOK); err != nil {
		return "", "", 0, err
	}
	return resBody.AccessToken, resBody.RefreshToken, resBody.ExpiresIn, nil
}

// registerDynamicClient method gets the Client ID and Client Secret using the given Dynamic client registration request.
func (m *PasswordRefreshTokenGrantManager) registerDynamicClient(reqBody *DynamicClientRegReq) error {
	bodyReader, err := client.BodyReader(reqBody)
	if err != nil {
		return errors.Wrapf(err, ErrMsgUnableToParseRequestBody, DynamicClientRegMsg)
	}
	dynamicClientRegistrationEndpoint, err := utils.ConstructURL(m.DynamicClientEndpoint, m.DynamicClientRegistrationContext)
	if err != nil {
		return errors.Wrap(err, "cannot construct, token endpoint")
	}
	req, err := client.CreateHTTPRequest(http.MethodPost, dynamicClientRegistrationEndpoint, bodyReader)
	if err != nil {
		return errors.Wrapf(err, ErrMsgUnableToCreateRequestBody, DynamicClientRegMsg)
	}
	req.HTTPRequest().SetBasicAuth(m.UserName, m.Password)
	req.SetHeader(client.HTTPContentType, client.ContentTypeApplicationJSON)

	var resBody DynamicClientRegResBody
	if err := client.Invoke(DynamicClientRegMsg, req, &resBody, http.StatusOK); err != nil {
		return err
	}
	m.clientID = resBody.ClientID
	m.clientSec = resBody.ClientSecret
	return nil
}

// defaultClientRegBody function returns an initialized dynamic client registration request body.
func defaultClientRegBody() *DynamicClientRegReq {
	return &DynamicClientRegReq{
		CallbackURL: CallBackURL,
		ClientName:  ClientName,
		GrantType:   DynamicClientRegGrantType,
		Owner:       Owner,
		SaasApp:     true,
	}
}
