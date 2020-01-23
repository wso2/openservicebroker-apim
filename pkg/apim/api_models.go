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

// APIMaxTps represents the max TPS(Transactions per second) for an API.
type APIMaxTps struct {
	Production int64 `json:"production,omitempty"`
	Sandbox    int64 `json:"sandbox,omitempty"`
}

// APIEndpointSecurity represents the endpoint security information.
type APIEndpointSecurity struct {
	Type     string `json:"type,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// Sequence represents a API sequence.
type Sequence struct {
	Name   string `json:"name"`
	Type   string `json:"type,omitempty"`
	ID     string `json:"id,omitempty"`
	Shared bool   `json:"shared,omitempty"`
}

// Label represents a API label.
type Label struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// APIBusinessInformation represents the  API business information.
type APIBusinessInformation struct {
	BusinessOwner       string `json:"businessOwner,omitempty"`
	BusinessOwnerEmail  string `json:"businessOwnerEmail,omitempty"`
	TechnicalOwner      string `json:"technicalOwner,omitempty"`
	TechnicalOwnerEmail string `json:"technicalOwnerEmail,omitempty"`
}

// APICorsConfiguration represents the CORS configuration for the API.
type APICorsConfiguration struct {
	CorsConfigurationEnabled      bool     `json:"corsConfigurationEnabled,omitempty"`
	AccessControlAllowOrigins     []string `json:"accessControlAllowOrigins,omitempty" hash:"set"`
	AccessControlAllowCredentials bool     `json:"accessControlAllowCredentials,omitempty"`
	AccessControlAllowHeaders     []string `json:"accessControlAllowHeaders,omitempty" hash:"set"`
	AccessControlAllowMethods     []string `json:"accessControlAllowMethods,omitempty" hash:"set"`
}

// APIReqBody represents the request of create "API" API call.
type APIReqBody struct {
	ID string `json:"id"`
	// Name of the API
	Name string `json:"name"`
	// A brief description about the API
	Description string `json:"description"`
	// A string that represents the context of the user's request
	Context string `json:"context"`
	// The version of the API
	Version string `json:"version"`
	// If the provider value is not given, the user invoking the API will be used as the provider.
	Provider string `json:"provider,omitempty"`
	// This describes in which status of the lifecycle the API is
	Status       string `json:"status"`
	ThumbnailURI string `json:"thumbnailUri,omitempty"`
	// Swagger definition of the API which contains details about URI templates and scopes
	APIDefinition string `json:"apiDefinition"`
	// WSDL URL if the API is based on a WSDL endpoint
	WsdlURI                 string `json:"wsdlUri,omitempty"`
	ResponseCaching         string `json:"responseCaching,omitempty"`
	CacheTimeout            int32  `json:"cacheTimeout,omitempty"`
	DestinationStatsEnabled bool   `json:"destinationStatsEnabled,omitempty"`
	IsDefaultVersion        bool   `json:"isDefaultVersion"`
	// The transport to be set. Accepted values are HTTP, WS
	Type string `json:"type"`
	// Supported transports for the API (http and/or https).
	Transport []string `json:"transport" hash:"set"`
	// Search keywords related to the API
	Tags []string `json:"tags,omitempty" hash:"set"`
	// The subscription tiers selected for the particular API
	Tiers []string `json:"tiers" hash:"set"`
	// The policy selected for the particular API
	APILevelPolicy string `json:"apiLevelPolicy,omitempty"`
	// Name of the Authorization header used for invoking the API. If it is not set, Authorization header name specified in tenant or system level will be used.
	AuthorizationHeader string     `json:"authorizationHeader,omitempty"`
	MaxTps              *APIMaxTps `json:"maxTps,omitempty"`
	// The visibility level of the API. Accepts one of the following. PUBLIC, PRIVATE, RESTRICTED OR CONTROLLED.
	Visibility string `json:"visibility"`
	// The user roles that are able to access the API
	VisibleRoles     []string             `json:"visibleRoles,omitempty" hash:"set"`
	VisibleTenants   []string             `json:"visibleTenants,omitempty" hash:"set"`
	EndpointConfig   string               `json:"endpointConfig"`
	EndpointSecurity *APIEndpointSecurity `json:"endpointSecurity,omitempty"`
	// Comma separated list of gateway environments.
	GatewayEnvironments string `json:"gatewayEnvironments,omitempty"`
	// Labels of micro-gateway environments attached to the API.
	Labels    []Label    `json:"labels,omitempty" hash:"set"`
	Sequences []Sequence `json:"sequences,omitempty" hash:"set"`
	// The subscription availability. Accepts one of the following. current_tenant, all_tenants or specific_tenants.
	SubscriptionAvailability     string   `json:"subscriptionAvailability,omitempty"`
	SubscriptionAvailableTenants []string `json:"subscriptionAvailableTenants,omitempty"`
	// Map of custom properties of API
	AdditionalProperties map[string]string `json:"additionalProperties,omitempty" hash:"set"`
	// Is the API is restricted to certain set of publishers or creators or is it visible to all the publishers and creators. If the accessControl restriction is none, this API can be modified by all the publishers and creators, if not it can only be viewable/modifiable by certain set of publishers and creators,  based on the restriction.
	AccessControl string `json:"accessControl,omitempty"`
	// The user roles that are able to view/modify as API publisher or creator.
	AccessControlRoles  []string                `json:"accessControlRoles,omitempty"`
	BusinessInformation *APIBusinessInformation `json:"businessInformation,omitempty"`
	CorsConfiguration   *APICorsConfiguration   `json:"corsConfiguration,omitempty"`
}

// APICreateResp represents the response of create "API" API call.
type APICreateResp struct {
	// UUID of the api registry artifact
	ID string `json:"id,omitempty"`
}

// ApplicationMetadata represents name, id and key of the generated application
type ApplicationMetadata struct {
	Name         string
	ID           string
	Keys         *ApplicationKeyResp
	DashboardURL string
}

// ApplicationCreateReq represents the response of create Application API call.
type ApplicationCreateReq struct {
	ThrottlingTier string `json:"throttlingTier"`
	Description    string `json:"description,omitempty"`
	Name           string `json:"name,omitempty"`
	CallbackURL    string `json:"callbackUrl,omitempty"`
}

// APIParam represents the structure for API plan parameters.
type APIParam struct {
	APISpec APIReqBody `json:"api"`
}

// ApplicationParam represents the structure for Application plan parameters.
type ApplicationParam struct {
	AppSpec ApplicationCreateReq `json:"app"`
}

// SubscriptionSpec represents the parameters for a Subscription.
type SubscriptionSpec struct {
	APIName          string `json:"apiName"`
	AppName          string `json:"appName"`
	SubscriptionTier string `json:"tier"`
}

// SubscriptionParam represents the structure for Subscription plan parameters.
type SubscriptionParam struct {
	SubsSpec SubscriptionSpec `json:"subs"`
}

// SubscriptionReq represents the APIM subscription create request.
type SubscriptionReq struct {
	Tier          string `json:"tier"`
	APIIdentifier string `json:"apiIdentifier"`
	ApplicationID string `json:"applicationId"`
}

// AppCreateReq represents the application creation request body.
type AppCreateReq struct {
	ThrottlingTier string `json:"throttlingTier"`
	Description    string `json:"description"`
	Name           string `json:"name"`
	CallbackURL    string `json:"callbackUrl"`
}

// AppCreateRes represents the application creation response body.
type AppCreateRes struct {
	ApplicationID string `json:"applicationId"`
}

// ApplicationKeyGenerateRequest represents the application key generation request.
type ApplicationKeyGenerateRequest struct {
	KeyType      string `json:"keyType"`
	ValidityTime string `json:"validityTime"`
	// The grant types that are supported by the application
	SupportedGrantTypes []string `json:"supportedGrantTypes,omitempty"`
	// Callback URL
	CallbackURL string `json:"callbackUrl,omitempty"`
	// Allowed domains for the access token
	AccessAllowDomains []string `json:"accessAllowDomains"`
	// Allowed scopes for the access token
	Scopes []string `json:"scopes,omitempty"`
	// Client ID for generating access token.
	ClientID string `json:"clientId,omitempty"`
	// Client secret for generating access token. This is given together with the client ID.
	ClientSecret string `json:"clientSecret,omitempty"`
}

// ApplicationKeyResp represents the Application generate keys API call response.
type ApplicationKeyResp struct {
	// The consumer key associated with the application and identifying the client
	ConsumerKey string `json:"consumerKey,omitempty"`
	// The client secret that is used to authenticate the client with the authentication server
	ConsumerSecret string `json:"consumerSecret,omitempty"`
	// The grant types that are supported by the application
	SupportedGrantTypes []string `json:"supportedGrantTypes,omitempty"`
	// Callback URL
	CallbackURL string `json:"callbackUrl,omitempty"`
	// Describes the state of the key generation.
	KeyState string `json:"keyState,omitempty"`
	// Describes to which endpoint the key belongs
	KeyType string `json:"keyType,omitempty"`
	// ApplicationConfig group id (if any).
	GroupID string `json:"groupId,omitempty"`
	Token   *Token `json:"token,omitempty"`
}

// Token represents an Application token.
type Token struct {
	// Access token
	AccessToken string `json:"accessToken,omitempty"`
	// Valid scopes for the access token
	TokenScopes []string `json:"tokenScopes,omitempty"`
	// Maximum validity time for the access token
	ValidityTime int64 `json:"validityTime,omitempty"`
}

// SubscriptionResp represents the response of create Subscription API call.
type SubscriptionResp struct {
	// The UUID of the subscription
	SubscriptionID string `json:"subscriptionId,omitempty"`
	// The UUID of the application
	ApplicationID string `json:"applicationId"`
	// The unique identifier of the API.
	APIIdentifier string `json:"apiIdentifier"`
	Tier          string `json:"tier"`
	Status        string `json:"status,omitempty"`
}

// APISearchInfo represents the API search information.
type APISearchInfo struct {
	Provider    string `json:"provider"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Name        string `json:"name"`
	Context     string `json:"context"`
	ID          string `json:"id"`
}

// APISearchResp represents the response of search "API" by name API call.
type APISearchResp struct {
	Previous string          `json:"previous"`
	List     []APISearchInfo `json:"list"`
	Count    int             `json:"count"`
	Next     string          `json:"next"`
}

// ApplicationSearchInfo represents the Application search information.
type ApplicationSearchInfo struct {
	GroupID        string `json:"groupId"`
	Subscriber     string `json:"subscriber"`
	ThrottlingTier string `json:"throttlingTier"`
	ApplicationID  string `json:"applicationId"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Status         string `json:"status"`
}

// ApplicationSearchResp represents the response of search Application by name API call.
type ApplicationSearchResp struct {
	Previous string                  `json:"previous"`
	List     []ApplicationSearchInfo `json:"list"`
	Count    int                     `json:"count"`
	Next     string                  `json:"next"`
}

var AppPlanBindInputParameterSchemaRaw = `{
  "$schema": "http://json-schema.org/draft-04/schema#"
}`

var AppPlanInputParameterSchemaRaw = `{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "type": "object",
  "properties": {
    "apis": {
      "type": "array",
      "items": [
        {
          "type": "object",
          "properties": {
            "name": {
              "type": "string"
            },
            "version": {
              "type": "string"
            }
          },
          "required": [
            "name",
            "version"
          ]
        }
      ]
    }
  },
  "required": [
    "apis"
  ]
}`
