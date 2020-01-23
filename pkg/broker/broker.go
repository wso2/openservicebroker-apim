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

// Package broker holds the implementation of brokerapi.ServiceBroker interface.
package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/mitchellh/hashstructure"
	"github.com/pivotal-cf/brokerapi/domain"
	"github.com/pivotal-cf/brokerapi/domain/apiresponses"
	"github.com/pkg/errors"
	"github.com/wso2/service-broker-apim/pkg/apim"
	"github.com/wso2/service-broker-apim/pkg/client"
	"github.com/wso2/service-broker-apim/pkg/db"
	"github.com/wso2/service-broker-apim/pkg/log"
	"github.com/wso2/service-broker-apim/pkg/mapBrokerError"
	"github.com/wso2/service-broker-apim/pkg/model"
	"github.com/wso2/service-broker-apim/pkg/utils"
)

const (
	LogKeyAppID                      = "application-id"
	LogKeyServiceID                  = "service-id"
	LogKeyPlanID                     = "plan-id"
	LogKeyInstanceID                 = "instance-id"
	LogKeyBindID                     = "bind-id"
	LogKeyApplicationName            = "application-name"
	LogKeyPlatformApplicationName    = "platform-application-name"
	LogKeySpaceID                    = "cf-space-id"
	LogKeyOrgID                      = "cf-org-id"
	ApplicationDashboardURL          = "application-dashboard-url"
	ErrMsgUnableToStoreInstance      = "unable to store service instance in database"
	ErrActionStoreInstance           = "store service instance"
	ErrActionDelAPP                  = "delete Application"
	ErrActionDelInstance             = "delete service instance"
	ErrActionCreateAPIMResource      = "creating API-M resource"
	ErrActionUpdateAPIMResource      = "update API-M resource"
	ErrMsgUnableDelInstance          = "unable to delete service instance"
	ErrMsgUnableToGetBind            = "unable to retrieve Bind from the database"
	ErrMsgUnableToGenInputSchema     = "unable to generate %s plan input Schema"
	ErrMsgUnableToGenBindInputSchema = "unable to generate %s plan bind input Schema"
	ErrMsgInvalidPlanID              = "invalid plan id"
	ErrMsgUnableGenerateKeys         = "unable generate keys for application"
	ServiceID                        = "460F28F9-4D05-4889-970A-6BF5FB7D3CF8"
	ServiceName                      = "wso2apim-service"
	ServiceDescription               = "Manages WSO2 API Manager artifacts"
	ApplicationPlanID                = "00e851cd-ce8b-43eb-bc27-ac4d4fbb3204"
	ApplicationPlanName              = "app"
	ApplicationPlanDescription       = "Creates an Application with a set of subscription for a given set of APIs in WSO2 API Manager"
	DebugMsgDelInstance              = "delete instance"
	ApplicationPrefix                = "ServiceBroker_"
	StatusInstanceAlreadyExists      = "Instance already exists"
	StatusInstanceDoesNotExist       = "Instance does not exist"
)

var (
	ErrNotSupported                 = errors.New("not supported")
	ErrInvalidSVCPlan               = errors.New("invalid service or getServices")
	applicationPlanBindable         = true
	appPlanInputParameterSchema     map[string]interface{}
	appPlanBindInputParameterSchema map[string]interface{}
)

// APIM struct implements the interface brokerapi.ServiceBroker.
type APIM struct{}

// API struct represent an API.
type API struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ServiceParams represents the SVC create and update parameter.
type ServiceParams struct {
	APIs []API `json:"apis" hash:"set"`
}

// apimBrokerProvisionDetails represents the attributes retrieved from the provision request
type apimBrokerProvisionDetails struct {
	organizationalID  string
	spaceID           string
	serviceParameters ServiceParams
}

// Init method initialize the API-M broker. If there is an error it will cause a panic.
func (apimBroker *APIM) Init() {
	var err error
	appPlanInputParameterSchema, err = utils.GetJSONSchema(apim.AppPlanInputParameterSchemaRaw)
	if err != nil {
		log.HandleErrorAndExit(fmt.Sprintf(ErrMsgUnableToGenInputSchema, "app"), err)
	}
	appPlanBindInputParameterSchema, err = utils.GetJSONSchema(apim.AppPlanBindInputParameterSchemaRaw)
	if err != nil {
		log.HandleErrorAndExit(fmt.Sprintf(ErrMsgUnableToGenBindInputSchema, "app"), err)
	}
}

// Services returns the getServices offered(catalog) by this broker.
func (apimBroker *APIM) Services(ctx context.Context) ([]domain.Service, error) {
	return getServices()
}

func createCommonLogData(svcInstanceID, serviceID, planID string) *log.Data {
	return log.NewData().
		Add(LogKeyServiceID, serviceID).
		Add(LogKeyPlanID, planID).
		Add(LogKeyInstanceID, svcInstanceID)
}

func readProvisionDetails(serviceDetails *domain.ProvisionDetails, logData *log.Data) (*apimBrokerProvisionDetails, error) {
	apiParams, err := getServiceParamsIfExists(serviceDetails.RawParameters, logData)
	if err != nil {
		return nil, err
	}
	provDetails := &apimBrokerProvisionDetails{
		organizationalID:  serviceDetails.OrganizationGUID,
		spaceID:           serviceDetails.SpaceGUID,
		serviceParameters: apiParams,
	}
	return provDetails, nil
}

func getServiceParamsIfExists(rawParams json.RawMessage, logData *log.Data) (ServiceParams, error) { //check: name
	apiParams, err := unmarshalServiceParams(rawParams)
	if err != nil {
		return apiParams, err
	}

	if len(apiParams.APIs) == 0 {
		log.Error("no APIs defined", nil, logData)
		return apiParams, &mapBrokerError.ErrorEmptyAPIParameterSet{}
	}
	return apiParams, nil
}

func hasValidSpaceIDAndOrgID(spaceID string, orgID string) bool {
	return spaceID != "" && orgID != ""
}

func validateHashSpaceIDOrgID(svcInstance *model.ServiceInstance, paramHash string, spaceID string, orgID string) bool {
	return (svcInstance.ParameterHash == paramHash) && (svcInstance.SpaceID == spaceID) && (svcInstance.OrgID == orgID)
}

func generateHashForserviceParameters(appID string, svcParams ServiceParams, logData *log.Data) (string, error) {
	generatedHash, err := hashstructure.Hash(svcParams, nil)
	if err != nil {
		log.Error("unable to generate hash value for service parameters", err, logData)
		err := &mapBrokerError.ErrorUnableToGenerateHash{}
		return "", err
	}
	return strconv.FormatUint(generatedHash, 10), nil
}

func isSameInstanceWithDifferentAttrubutes(svcInstance *model.ServiceInstance, apimProvDetails *apimBrokerProvisionDetails, logData *log.Data) (bool, error) {
	parameterHash, err := generateHashForserviceParameters(svcInstance.ApplicationID, apimProvDetails.serviceParameters, logData)
	if err != nil {
		return false, err
	}
	if ok := validateHashSpaceIDOrgID(svcInstance, parameterHash, apimProvDetails.spaceID, apimProvDetails.organizationalID); ok {

		existingAPIs, err := getExistingAPIsForAppID(svcInstance.ApplicationID, logData)
		if err != nil {
			return false, err
		}
		if !isSameAPIs(existingAPIs, apimProvDetails.serviceParameters.APIs) {
			log.Debug("APIs does not match", logData)
			return true, nil
		}
		return false, nil
	}
	return true, nil
}

func deleteSubscriptions(removedSubsIds []string, svcInstanceID string) error {

	for _, sub := range removedSubsIds {
		err := apim.UnSubscribe(sub)
		if err != nil {
			return err
		}
		err = removeSubscription(sub, svcInstanceID)
		if err != nil {
			return err
		}
	}
	return nil
}

func removeSubscription(subId, svcInstId string) error {
	sub := &model.Subscription{
		ID:            subId,
		SVCInstanceID: svcInstId,
	}
	err := db.Delete(sub)
	if err != nil {
		return err
	}

	return nil
}

func getRemovedSubscriptionsIDs(applicationID string, existingAPIs, requestedAPIs []API, logData *log.Data) ([]string, error) {
	removedAPIs := getRemovedAPIs(existingAPIs, requestedAPIs)
	var removedSubsIDs []string
	for _, rAPI := range removedAPIs {
		rSub, err := getSubscriptionForAppAndAPI(applicationID, rAPI, logData)
		if err != nil {
			return nil, err
		}
		removedSubsIDs = append(removedSubsIDs, rSub.ID)
	}
	return removedSubsIDs, nil
}

func getRemovedAPIs(existingAPIs, requestedAPIs []API) []API {
	var removedAPIs []API
	for _, api := range existingAPIs {
		if !isArrayContainAPI(requestedAPIs, api) {
			removedAPIs = append(removedAPIs, api)
		}
	}
	return removedAPIs
}

func getAddedAPIs(existingAPIs, updatedAPIs []API, logData *log.Data) []API {
	var addedAPIs []API
	for _, api := range updatedAPIs {
		if !isArrayContainAPI(existingAPIs, api) {
			addedAPIs = append(addedAPIs, api)
		}
	}
	return addedAPIs
}

func (apimBroker *APIM) GetBinding(ctx context.Context, svcInstanceID,
	bindingID string) (domain.GetBindingSpec, error) {
	return domain.GetBindingSpec{}, ErrNotSupported
}

func (apimBroker *APIM) GetInstance(ctx context.Context,
	svcInstanceID string) (domain.GetInstanceDetailsSpec, error) {
	return domain.GetInstanceDetailsSpec{}, ErrNotSupported
}

func (apimBroker *APIM) LastBindingOperation(ctx context.Context, svcInstanceID,
	bindingID string, details domain.PollDetails) (domain.LastOperation, error) {
	return domain.LastOperation{}, ErrNotSupported
}

// getServices returns an array of getServices offered by this service broker and any error encountered.
func getServices() ([]domain.Service, error) {
	return []domain.Service{
		{
			ID:                   ServiceID,
			Name:                 ServiceName,
			Description:          ServiceDescription,
			Bindable:             true,
			InstancesRetrievable: false,
			PlanUpdatable:        true,
			Plans: []domain.ServicePlan{
				{
					ID:          ApplicationPlanID,
					Name:        ApplicationPlanName,
					Description: ApplicationPlanDescription,
					Bindable:    &applicationPlanBindable,
					Schemas: &domain.ServiceSchemas{
						Instance: domain.ServiceInstanceSchema{
							Create: domain.Schema{
								Parameters: appPlanInputParameterSchema,
							},
							Update: domain.Schema{
								Parameters: appPlanInputParameterSchema,
							},
						},
						Binding: domain.ServiceBindingSchema{
							Create: domain.Schema{
								Parameters: appPlanBindInputParameterSchema,
							},
						},
					},
				},
			},
		},
	}, nil
}

// createApplication creates Subscription in API-M and returns App ID, App dashboard URL and an error if encountered.
func createApplication(appName string, logData *log.Data) (string, string, error) {
	req := &apim.ApplicationCreateReq{
		Name:           appName,
		ThrottlingTier: "Unlimited",
		Description:    "Application " + appName + " created by WSO2 APIM Service Broker",
	}
	appID, err := apim.CreateApplication(req)
	if err != nil {
		log.Error("unable to create application", err, logData)
		return "", "", handleAPIMResourceCreateError(err, appName, logData)
	}
	dashboardURL := apim.GetAppDashboardURL(appName)
	return appID, dashboardURL, nil
}

// handleAPIMResourceCreateError handles the API-M resource creation error. Returns an error mapped to apiresponses.FailureResponse.
func handleAPIMResourceCreateError(e error, resourceName string, logData *log.Data) error {
	invokeErr, ok := e.(*client.InvokeError)
	if ok && invokeErr.StatusCode == http.StatusConflict {
		log.Debug("API-M resource already exists", logData)
		return &mapBrokerError.ErrorAPIMResourceAlreadyExists{}
	}

	log.Debug("unable to create API-M resource", logData)
	return &mapBrokerError.ErrorUnableToCreateAPIMResource{}
}

// retriveServiceInstance function checks whether the given instance already exists.
// If the given instance exists then an initialized instance and
// If the given instance is unable to retrieve from database, an error is returned.
func retriveServiceInstance(svcInstanceID string, logData *log.Data) (*model.ServiceInstance, error) {
	instance := &model.ServiceInstance{
		ID: svcInstanceID,
	}
	exists, err := db.Retrieve(instance)
	if err != nil {
		log.Error("unable to retrieve the service instance from database", err, logData)
		return nil, &mapBrokerError.ErrorUnableToRetrieveServiceInstance{}
	}
	if !exists {
		log.Debug("instance doesn't exists", logData)
		return nil, nil
	}
	return instance, nil
}

// deleteInstance function deletes the given instance from database. An error type mapped to apiresponses.FailureResponse is returned.
func deleteInstance(i *model.ServiceInstance, logData *log.Data) error {
	err := db.Delete(i)
	if err != nil {
		log.Error("unable to delete the instance from the database", err, logData)
		err := &mapBrokerError.ErrorUnableToDeleteInstance{}
		return err
	}
	return nil
}

// handleAPIMResourceUpdateError handles the API-M resource update errors. Returns an error type mapped to apiresponses.FailureResponse.
func handleAPIMResourceUpdateError(err error, resourceName string) error {
	e, ok := err.(*client.InvokeError)
	if ok && e.StatusCode == http.StatusNotFound {
		return &mapBrokerError.ErrorAPIMResourceDoesNotExist{
			APIMResourceName: resourceName,
		}
	}
	return &mapBrokerError.ErrorUnableToUpdateAPIMResource{}
}

func getSubscriptionsListForAppID(applicationID string, logData *log.Data) ([]model.Subscription, error) {
	subscription := &model.Subscription{
		ApplicationID: applicationID,
	}
	var subscriptionsList []model.Subscription
	hasSubscriptions, err := db.RetrieveList(subscription, &subscriptionsList)
	if err != nil {
		log.Error("unable to retrieve subscription", err, logData)
		return subscriptionsList, &mapBrokerError.ErrorUnableToRetrieveSubscriptionList{}
	}
	if !hasSubscriptions {
		return nil, nil
	}
	return subscriptionsList, nil
}

func getSubscriptionForAppAndAPI(applicationID string, api API, logData *log.Data) (*model.Subscription, error) {
	subscription := &model.Subscription{
		ApplicationID: applicationID,
		APIName:       api.Name,
		APIVersion:    api.Version,
	}
	hasSubscription, err := db.Retrieve(subscription)
	if err != nil {
		log.Error("unable to retrieve subscription", err, logData)
		return nil, &mapBrokerError.ErrorUnableToRetrieveSubscription{}
	}
	if !hasSubscription {
		log.Error("no subscription is available", err, logData)
		return nil, &mapBrokerError.ErrorNoSubscriptionAvailable{}
	}
	return subscription, nil
}

func getExistingAPIsForAppID(applicationID string, logData *log.Data) ([]API, error) {
	subscriptionsList, err := getSubscriptionsListForAppID(applicationID, logData)
	if err != nil {
		return nil, err
	}
	if subscriptionsList == nil {
		log.Error("no subscriptions are available", err, logData)
		return nil, &mapBrokerError.ErrorSubscriptionListUnavailable{}
	}
	var existingAPIs []API
	for _, sub := range subscriptionsList {
		api := API{
			Name:    sub.APIName,
			Version: sub.APIVersion,
		}
		existingAPIs = append(existingAPIs, api)
	}
	return existingAPIs, nil
}

func isSameAPIs(existingAPIs, requestedAPIs []API) bool {
	if len(existingAPIs) != len(requestedAPIs) {
		return false
	}
	for _, existingAPI := range existingAPIs { //TODO: try to optimize this (n*n too complex)
		if !isArrayContainAPI(requestedAPIs, existingAPI) {
			return false
		}
	}
	return true
}

func isArrayContainAPI(apis []API, api API) bool {
	for _, a := range apis {
		if a.Name == api.Name && a.Version == api.Version {
			return true
		}
	}
	return false
}

func unsubscribeMultipleAPIs(subs []model.Subscription, logData *log.Data) {
	for _, subscription := range subs {
		err := apim.UnSubscribe(subscription.ID)
		if err != nil {
			log.Error("Unable to unsubscribe APIs", err, logData)
		}
	}
}

func createAndStoreSubscriptions(instance *model.ServiceInstance, apis []API, logData *log.Data) error {
	subscriptions, err := createSubscriptions(instance, apis, logData)
	if err != nil {
		log.Error("unable to create subscriptions", err, logData)
		return err
	}
	err = storeSubscriptions(subscriptions)
	if err != nil {
		unsubscribeMultipleAPIs(subscriptions, logData)
		log.Error("unable to store subscriptions", err, logData)
	}

	return nil
}

func createServiceInstanceObject(ID, paramHash string, apimProvDetails *apimBrokerProvisionDetails, appData *apim.ApplicationMetadata) *model.ServiceInstance {
	svcInstance := &model.ServiceInstance{
		ID:              ID,
		ApplicationID:   appData.ID,
		ApplicationName: appData.Name,
		SpaceID:         apimProvDetails.spaceID,
		OrgID:           apimProvDetails.organizationalID,
		ConsumerKey:     appData.Keys.ConsumerKey,
		ConsumerSecret:  appData.Keys.ConsumerSecret,
		ParameterHash:   paramHash,
	}
	return svcInstance
}

func persistServiceInstance(svcInstance *model.ServiceInstance, logData *log.Data) error {
	err := storeServiceInstance(svcInstance, logData)
	if err != nil {
		return err
	}
	return nil
}

func createApplicationAndGenerateKeys(id string, logData *log.Data) (*apim.ApplicationMetadata, error) {
	appName := generateApplicationName(id)

	logData.Add(LogKeyApplicationName, appName)
	appID, appDashboardURL, err := createApplication(appName, logData)
	if err != nil {
		return nil, err
	}
	logData.Add(LogKeyAppID, appID).
		Add(ApplicationDashboardURL, appDashboardURL)

	keys, err := generateKeysForApplication(appID, logData)
	if err != nil {
		revertApplication(appID, logData)
		return nil, err
	}

	return &apim.ApplicationMetadata{
		Name:         appName,
		ID:           appID,
		Keys:         keys,
		DashboardURL: appDashboardURL,
	}, nil

}

func removeServiceInstanceAndLogError(svcInstanceID string, logData *log.Data) {
	err := db.Delete(&model.ServiceInstance{
		ID: svcInstanceID,
	})
	if err != nil {
		log.Error(ErrMsgUnableDelInstance, err, logData)
	}
}

func storeSubscriptions(subscriptions []model.Subscription) error {
	var entities []model.Entity
	for _, val := range subscriptions {
		entities = append(entities, val)
	}
	err := db.BulkInsert(entities)
	if err != nil {
		log.Error("unable to store subscriptions", err, nil)
		return &mapBrokerError.ErrorUnableToStoreSubscriptions{}
	}
	return nil
}

func getSubscriptionList(svcInstanceID string, subsResponses []apim.SubscriptionResp) []model.Subscription {
	var subscriptions []model.Subscription
	for _, subsResponse := range subsResponses {
		apiIdentifier := strings.Split(subsResponse.APIIdentifier, "-")
		subs := model.Subscription{
			ID:            subsResponse.SubscriptionID,
			ApplicationID: subsResponse.ApplicationID,
			User:          apiIdentifier[0],
			APIName:       apiIdentifier[1],
			APIVersion:    apiIdentifier[2],
			SVCInstanceID: svcInstanceID,
		}
		subscriptions = append(subscriptions, subs)
	}
	return subscriptions
}

// storeServiceInstance stores service instance in the database returns an error type mapped to apiresponses.FailureResponse.
func storeServiceInstance(i *model.ServiceInstance, logData *log.Data) error {
	err := db.Store(i)
	if err != nil {
		log.Error(ErrMsgUnableToStoreInstance, err, logData)
		return &mapBrokerError.ErrorUnableToStoreServiceInstance{}
	}
	return nil
}

func createSubscriptions(svcInstance *model.ServiceInstance, apis []API, logData *log.Data) ([]model.Subscription, error) {

	var subscriptionRequests []apim.SubscriptionReq

	for _, api := range apis {
		apiID, err := apim.SearchAPIByNameVersion(api.Name, api.Version)
		if err != nil {
			return nil, &mapBrokerError.ErrorUnableToSearchAPIs{}
		}
		subReq := apim.SubscriptionReq{
			ApplicationID: svcInstance.ApplicationID,
			APIIdentifier: apiID,
			Tier:          "Unlimited",
		}
		subscriptionRequests = append(subscriptionRequests, subReq)
	}

	subscriptionCreateResp, err := apim.CreateMultipleSubscriptions(subscriptionRequests)
	if err != nil {
		log.Error("unable to create subscriptions", err, logData)
		return nil, &mapBrokerError.ErrorUnableToCreateSubscription{}
	}
	return getSubscriptionList(svcInstance.ID, subscriptionCreateResp), nil
}

func generateApplicationName(svcInstanceID string) string {
	return ApplicationPrefix + svcInstanceID
}

func (apimBroker *APIM) Provision(ctx context.Context, svcInstanceID string,
	provisionDetails domain.ProvisionDetails, asyncAllowed bool) (domain.ProvisionedServiceSpec, error) { //pdfProvisionDetails
	if !hasValidSpaceIDAndOrgID(provisionDetails.SpaceGUID, provisionDetails.OrganizationGUID) {
		return domain.ProvisionedServiceSpec{}, apiresponses.NewFailureResponse(errors.New("check space ID and org ID"), http.StatusBadRequest, "invalid parameters")
	}

	logData := createCommonLogData(svcInstanceID, provisionDetails.ServiceID, provisionDetails.PlanID)

	apimProvDetails, err := readProvisionDetails(&provisionDetails, logData)
	if err != nil {
		return domain.ProvisionedServiceSpec{}, mapBrokerError.MapBrokerErrors(err)
	}

	svcInstance, err := retriveServiceInstance(svcInstanceID, logData)
	if err != nil {
		return domain.ProvisionedServiceSpec{}, mapBrokerError.MapBrokerErrors(err)
	}

	if svcInstance != nil {
		confirm, err := isSameInstanceWithDifferentAttrubutes(svcInstance, apimProvDetails, logData)
		if err != nil {
			return domain.ProvisionedServiceSpec{}, mapBrokerError.MapBrokerErrors(err)
		}
		if confirm {
			return domain.ProvisionedServiceSpec{}, apiresponses.ErrInstanceAlreadyExists
		} else {
			return domain.ProvisionedServiceSpec{
				AlreadyExists: true,
			}, nil
		}

	}

	appMetadata, err := createApplicationAndGenerateKeys(svcInstanceID, logData)
	if err != nil {
		return domain.ProvisionedServiceSpec{}, mapBrokerError.MapBrokerErrors(err)
	}

	parameterHash, err := generateHashForserviceParameters(appMetadata.ID, apimProvDetails.serviceParameters, logData)
	if err != nil {
		return domain.ProvisionedServiceSpec{}, mapBrokerError.MapBrokerErrors(err)
	}

	svcInstance = createServiceInstanceObject(svcInstanceID, parameterHash, apimProvDetails, appMetadata)

	err = persistServiceInstance(svcInstance, logData)
	if err != nil {
		revertApplication(appMetadata.ID, logData)
		return domain.ProvisionedServiceSpec{}, mapBrokerError.MapBrokerErrors(err)
	}

	err = createAndStoreSubscriptions(svcInstance, apimProvDetails.serviceParameters.APIs, logData)
	if err != nil {
		revertApplication(appMetadata.ID, logData)
		removeServiceInstanceAndLogError(svcInstanceID, logData)
		return domain.ProvisionedServiceSpec{}, mapBrokerError.MapBrokerErrors(err)
	}

	return domain.ProvisionedServiceSpec{
		DashboardURL: appMetadata.DashboardURL,
	}, nil
}

func (apimBroker *APIM) Deprovision(ctx context.Context, svcInstanceID string,
	serviceDetails domain.DeprovisionDetails, asyncAllowed bool) (domain.DeprovisionServiceSpec, error) {
	logData := createCommonLogData(svcInstanceID, serviceDetails.ServiceID, serviceDetails.PlanID)

	svcInstance, err := retriveServiceInstance(svcInstanceID, logData)
	if err != nil {
		return domain.DeprovisionServiceSpec{}, mapBrokerError.MapBrokerErrors(err)
	}
	if svcInstance == nil {
		log.Debug("instance doesn't exists", logData)
		return domain.DeprovisionServiceSpec{}, apiresponses.ErrInstanceDoesNotExist
	}

	logData.
		Add(LogKeyAppID, svcInstance.ApplicationID).
		Add(LogKeyApplicationName, svcInstance.ApplicationName)

	log.Debug("delete the application", logData)
	err = apim.DeleteApplication(svcInstance.ApplicationID)
	if err != nil {
		log.Error("unable to delete the Application", err, logData)
		return domain.DeprovisionServiceSpec{}, apiresponses.NewFailureResponse(errors.New(ErrMsgUnableDelInstance), http.StatusInternalServerError, ErrActionDelAPP) //TODO: consistancy ?
	}

	log.Debug(DebugMsgDelInstance, logData)

	err = deleteInstance(&model.ServiceInstance{ID: svcInstanceID}, logData)
	if err != nil {
		return domain.DeprovisionServiceSpec{}, mapBrokerError.MapBrokerErrors(err)
	}

	return domain.DeprovisionServiceSpec{}, nil
}

func unmarshalServiceParams(rawParam json.RawMessage) (ServiceParams, error) {
	var serviceParams ServiceParams
	err := json.Unmarshal(rawParam, &serviceParams)
	if err != nil {
		return serviceParams, err
	}
	return serviceParams, nil
}

// retrieveServiceBind returns the initialized Bind struct if successfull and
// returns a nil object and an error for any error encountered.
func retrieveServiceBind(bindingID string, logData *log.Data) (*model.Bind, error) {
	bind := &model.Bind{
		ID: bindingID,
	}
	exists, err := db.Retrieve(bind)
	if err != nil {
		log.Error(ErrMsgUnableToGetBind, err, logData)
		return nil, &mapBrokerError.ErrorUnableToRetrieveBind{}
	}
	if !exists {
		return nil, nil
	}
	return bind, nil
}

// isBindWithSameAttributes returns true of the Bind is already exists and attached with the given instance ID,attributes.
func isBindWithSameAttributes(bind *model.Bind, svcInstanceID string, bindResource *domain.BindResource) bool {
	var isSameAttributes = svcInstanceID == bind.SVCInstanceID
	if !isOriginatedFromCreateServiceKey(bindResource) {
		isSameAttributes = isSameAttributes && (bindResource.AppGuid == bind.PlatformAppID)
	}
	return isSameAttributes
}

// Bind method creates a Bind between given Service instance and the App.
func (apimBroker *APIM) Bind(ctx context.Context, svcInstanceID, bindingID string,
	bindDetails domain.BindDetails, asyncAllowed bool) (domain.Binding, error) {

	logData := createCommonLogData(svcInstanceID, bindDetails.ServiceID, bindDetails.PlanID)
	logData.Add(LogKeyBindID, bindingID)

	bind, err := retrieveServiceBind(bindingID, logData)
	if err != nil {
		return domain.Binding{}, mapBrokerError.MapBrokerErrors(err)
	}

	log.Debug("retrieve instance", logData)
	svcInstance, err := retriveServiceInstance(svcInstanceID, logData)
	if err != nil {
		return domain.Binding{}, mapBrokerError.MapBrokerErrors(err)
	}
	if svcInstance == nil {
		log.Debug("instance doesn't exists", logData)
		return domain.Binding{}, apiresponses.ErrInstanceDoesNotExist
	}

	credentialsMap := credentialsMap(svcInstance.ApplicationName, svcInstance.ConsumerKey, svcInstance.ConsumerSecret)

	var isWithSameAttr = false
	if bind != nil {
		isWithSameAttr = isBindWithSameAttributes(bind, svcInstanceID, bindDetails.BindResource)
		if !isWithSameAttr {
			return domain.Binding{}, apiresponses.ErrBindingAlreadyExists
		}
		return domain.Binding{
			Credentials:   credentialsMap,
			AlreadyExists: true,
		}, nil
	}

	platformAppID := getPlatformAppID(bindDetails.BindResource)
	logData.Add(LogKeyPlatformApplicationName, platformAppID)

	bind = &model.Bind{
		ID:            bindingID,
		PlatformAppID: platformAppID,
		SVCInstanceID: svcInstanceID,
	}
	err = storeBind(bind, logData)
	if err != nil {
		return domain.Binding{}, mapBrokerError.MapBrokerErrors(err)
	}
	log.Debug("successfully stored the Bind", logData)
	return domain.Binding{
		Credentials: credentialsMap,
	}, nil
}

func isApplicationPlan(planID string) bool {
	return planID == ApplicationPlanID
}

func getPlatformAppID(b *domain.BindResource) string {
	var cfAppID string
	if isOriginatedFromCreateServiceKey(b) {
		cfAppID = ""
	} else {
		cfAppID = b.AppGuid
	}
	return cfAppID
}

func revertApplication(appID string, logData *log.Data) {
	err := apim.DeleteApplication(appID)
	if err != nil {
		log.Error("unable to delete application", err, logData)
	}
	log.Debug("Delete Application", logData)
}

// generateKeysForApplication function generates keys for the given Subscription.
// Returns generated keys and an error type mapped to apiresponses.FailureResponse if encountered.
func generateKeysForApplication(appID string, logData *log.Data) (*apim.ApplicationKeyResp, error) {
	appKeys, err := apim.GenerateKeys(appID)
	if err != nil {
		log.Error(ErrMsgUnableGenerateKeys, err, logData)
		return appKeys, &mapBrokerError.ErrorUnableToGenerateKeys{}
	}
	return appKeys, nil
}

// storeBind function stores the given Bind in the database.
// Return an error type mapped to apiresponses.FailureResponse error if encountered.
func storeBind(b *model.Bind, logData *log.Data) error {
	err := db.Store(b)
	if err != nil {
		log.Error("unable to store bind", err, logData)
		return &mapBrokerError.ErrorUnableToStoreBind{}
	}
	return nil
}

// createServiceKey check whether the command is a "create-service-key".
// BindResources or BindResource.AppGuid is nil only if the it is a "create-service-key" command.
func isOriginatedFromCreateServiceKey(b *domain.BindResource) bool {
	return b == nil || b.AppGuid == ""
}

func credentialsMap(appName, consumerKey, consumerSecret string) map[string]interface{} {
	return map[string]interface{}{
		"ApplicationName": appName,
		"ConsumerKey":     consumerKey,
		"ConsumerSecret":  consumerSecret,
	}
}

// Unbind deletes the Bind from database and returns domain.UnbindSpec struct and any error encountered.
func (apimBroker *APIM) Unbind(ctx context.Context, svcInstanceID, bindingID string,
	unbindDetails domain.UnbindDetails, asyncAllowed bool) (domain.UnbindSpec, error) {

	logData := createCommonLogData(svcInstanceID, unbindDetails.ServiceID, unbindDetails.PlanID)

	if !isApplicationPlan(unbindDetails.PlanID) {
		log.Error(ErrMsgInvalidPlanID, ErrInvalidSVCPlan, logData)
		return domain.UnbindSpec{}, apiresponses.NewFailureResponse(errors.New("unbinding"), http.StatusBadRequest, "invalid planID") //TODO: consistant?
	}

	bind, err := retrieveServiceBind(bindingID, logData)
	if err != nil {
		return domain.UnbindSpec{}, mapBrokerError.MapBrokerErrors(err)
	}
	if bind == nil {
		return domain.UnbindSpec{}, apiresponses.ErrBindingDoesNotExist // TODO: consist?
	}

	logData.Add("cf-app-id", bind.PlatformAppID)

	err = deleteBind(bind, logData)
	if err != nil {
		return domain.UnbindSpec{}, mapBrokerError.MapBrokerErrors(err)
	}

	return domain.UnbindSpec{}, nil
}

func deleteBind(bind *model.Bind, logData *log.Data) error {
	err := db.Delete(bind)
	if err != nil {
		log.Error("unable to delete the bind from the database", err, logData)
		return &mapBrokerError.ErrorUnableToDeleteBind{}
	}
	return nil
}

// LastOperation ...
// If the broker provisions asynchronously, the Cloud Controller will poll this endpoint
// for the status of the provisioning operation.
func (apimBroker *APIM) LastOperation(ctx context.Context, svcInstanceID string,
	details domain.PollDetails) (domain.LastOperation, error) {
	return domain.LastOperation{}, errors.New("not supported")
}

func updateServiceForAddedAPIs(existingAPIs, updatedAPIs []API, svcInstance *model.ServiceInstance, logData *log.Data) ([]API, error) {

	addedAPIs := getAddedAPIs(existingAPIs, updatedAPIs, logData)
	if len(addedAPIs) == 0 {
		log.Debug("No new APIs found", logData)
		return addedAPIs, nil
	}

	err := createAndStoreSubscriptions(svcInstance, addedAPIs, logData)
	if err != nil {
		return nil, err
	}

	return addedAPIs, nil
}

func updateServiceForRemovedAPIs(existingAPIs []API, paramAPIs []API, svcInstance *model.ServiceInstance, logData *log.Data) error {

	removeSubscriptionIDs, err := getRemovedSubscriptionsIDs(svcInstance.ApplicationID, existingAPIs, paramAPIs, logData) //updatesvcforremovedapis
	if err != nil {
		return err
	}

	err = deleteSubscriptions(removeSubscriptionIDs, svcInstance.ID)
	if err != nil {
		return err
	}

	return nil
}

func revertAddedAPIs(appID, instanceID string, apis []API, logData *log.Data) {
	log.Debug("remove previously added APIs", logData)
	var removedSubsIDs []string
	for _, rAPI := range apis {
		rSub, err := getSubscriptionForAppAndAPI(appID, rAPI, logData)
		if err != nil {
			log.Error("unable to get subscriptions", err, logData)
		}
		removedSubsIDs = append(removedSubsIDs, rSub.ID)
	}
	err := deleteSubscriptions(removedSubsIDs, instanceID)
	if err != nil {
		log.Error("unable to delete subscriptions", err, logData)
	}
}

func (apimBroker *APIM) Update(cxt context.Context, svcInstanceID string,
	updateDetails domain.UpdateDetails, asyncAllowed bool) (domain.UpdateServiceSpec, error) {

	logData := createCommonLogData(svcInstanceID, updateDetails.ServiceID, updateDetails.PlanID)
	log.Debug("update service instance", logData)

	svcParams, err := getServiceParamsIfExists(updateDetails.RawParameters, logData)
	if err != nil {
		return domain.UpdateServiceSpec{}, mapBrokerError.MapBrokerErrors(err)
	}

	svcInstance, err := retriveServiceInstance(svcInstanceID, logData)
	if err != nil {
		return domain.UpdateServiceSpec{}, mapBrokerError.MapBrokerErrors(err)
	}
	if svcInstance == nil {
		log.Debug("instance doesn't exists", logData)
		return domain.UpdateServiceSpec{}, apiresponses.ErrInstanceDoesNotExist
	}

	existingAPIs, err := getExistingAPIsForAppID(svcInstance.ApplicationID, logData)
	if err != nil {
		return domain.UpdateServiceSpec{}, mapBrokerError.MapBrokerErrors(err)
	}

	addedAPIs, err := updateServiceForAddedAPIs(existingAPIs, svcParams.APIs, svcInstance, logData)
	if err != nil {
		return domain.UpdateServiceSpec{}, mapBrokerError.MapBrokerErrors(err)
	}

	err = updateServiceForRemovedAPIs(existingAPIs, svcParams.APIs, svcInstance, logData)
	if err != nil {
		revertAddedAPIs(svcInstance.ApplicationID, svcInstanceID, addedAPIs, logData)
		return domain.UpdateServiceSpec{}, mapBrokerError.MapBrokerErrors(err)
	}

	log.Debug("Instace successfully updated", logData)

	return domain.UpdateServiceSpec{}, nil
}
