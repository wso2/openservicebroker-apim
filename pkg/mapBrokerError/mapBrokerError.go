package mapBrokerError

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/pivotal-cf/brokerapi/domain/apiresponses"
	"github.com/wso2/service-broker-apim/pkg/apim"
	"github.com/wso2/service-broker-apim/pkg/log"
)

const (
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
)

type BindingNotExistsError struct{}
type InstanceDoesNotExistError struct{}
type InstanceConflictError struct{}
type ErrorUnableToRetrieveServiceInstance struct{}
type ErrorUnableToDeleteInstance struct{}
type ErrorUnableToRetrieveSubscriptionList struct{}
type ErrorSubscriptionListUnavailable struct{}
type ErrorUnableToRetrieveSubscription struct{}
type ErrorNoSubscriptionAvailable struct{}
type ErrorUnableToGenerateHash struct{}
type ErrorUnableToStoreSubscriptions struct{}
type ErrorUnableToStoreServiceInstance struct{}
type ErrorUnableToSearchAPIs struct{}
type ErrorUnableToCreateSubscription struct{}
type ErrorUnableToRetrieveBind struct{}
type ErrorUnableToStoreBind struct{}
type ErrorUnableToDeleteBind struct{}
type ErrorUnableToGenerateKeys struct{}
type ErrorBindDoesNotExist struct{}
type ErrorUnableToCreateAPIMResource struct{}
type ErrorAPIMResourceAlreadyExists struct{}
type ErrorEmptyAPIParameterSet struct{}
type ErrorUnableToUpdateAPIMResource struct{}
type ErrorAPIMResourceDoesNotExist struct {
	APIMResourceName string
}

func (e *ErrorUnableToRetrieveServiceInstance) Error() string {
	return "unable to get Service instance from database"
}

func (e *ErrorUnableToDeleteInstance) Error() string {
	return ErrMsgUnableDelInstance
}

func (e *ErrorUnableToRetrieveSubscriptionList) Error() string {
	return "unable to get subscriptions list"
}

func (e *ErrorUnableToRetrieveSubscription) Error() string {
	return "unable to retrieve subscriptions"
}

func (e *ErrorNoSubscriptionAvailable) Error() string {
	return "no subscription is available"
}

func (e *ErrorSubscriptionListUnavailable) Error() string {
	return "unable to get subscriptions list"
}

func (e *ErrorUnableToGenerateHash) Error() string {
	return "unable to generate hash value for service parameters"
}

func (e *ErrorUnableToStoreSubscriptions) Error() string {
	return "unable to store subscriptions in the database"
}

func (e *ErrorUnableToStoreServiceInstance) Error() string {
	return "unable to store subscriptions in the database"
}

func (e *ErrorUnableToSearchAPIs) Error() string {
	return "unable to search API's"
}

func (e *ErrorUnableToCreateSubscription) Error() string {
	return "unable to create subscriptions"
}

func (e *ErrorUnableToRetrieveBind) Error() string {
	return ErrMsgUnableToGetBind
}

func (e *ErrorUnableToStoreBind) Error() string {
	return "unable to store Bind"
}

func (e *ErrorUnableToDeleteBind) Error() string {
	return "unable to delete Bind"
}

func (e *ErrorBindDoesNotExist) Error() string {
	return "bind does not exist"
}

func (e *ErrorUnableToGenerateKeys) Error() string {
	return ErrMsgUnableGenerateKeys
}
func (e *ErrorUnableToCreateAPIMResource) Error() string {
	return "unable to create the API-M resource"
}

func (e *ErrorAPIMResourceAlreadyExists) Error() string {
	return "API-M resource already exists"
}

func (e *ErrorAPIMResourceDoesNotExist) Error() string {
	return fmt.Sprintf("API-M resource %s not found !", e.APIMResourceName)
}

func (e *ErrorUnableToUpdateAPIMResource) Error() string {
	return "unable to update the API-M resource"
}

func (e *ErrorEmptyAPIParameterSet) Error() string {
	return "No APIs Defined"
}

func returnInternalServerResponse(errMsg, loggerAction string) error {
	return apiresponses.NewFailureResponse(errors.New(errMsg), http.StatusInternalServerError, loggerAction)
}

func returnBadRequestResponsee(errMsg, loggerAction string) error {
	return apiresponses.NewFailureResponse(errors.New(errMsg), http.StatusBadRequest, loggerAction)
}

func revertApplication(appID string, logData *log.Data) {
	err := apim.DeleteApplication(appID)
	if err != nil {
		log.Error("unable to delete application", err, logData)
	}
	log.Debug("Delete Application", logData)

}

func MapBrokerErrors(err error) error { //TODO: mapBrokerError

	switch err.(type) {
	case *ErrorUnableToRetrieveServiceInstance:
		return returnInternalServerResponse("unable to get Service instance from database", "get instance from the database")
	case *ErrorUnableToDeleteInstance:
		return returnInternalServerResponse(ErrMsgUnableDelInstance, ErrActionDelInstance)
	case *ErrorUnableToRetrieveSubscriptionList:
		return returnInternalServerResponse("unable to retrieve subscriptions", "retrive subscription")
	case *ErrorSubscriptionListUnavailable:
		return returnInternalServerResponse("no subscriptions are available", "retrive subscription")
	case *ErrorUnableToRetrieveSubscription:
		return returnInternalServerResponse("unable to query database", "retrieve subscription")
	case *ErrorNoSubscriptionAvailable:
		return returnInternalServerResponse("no subscription is available", "retrieve subscription")
	case *ErrorUnableToGenerateHash:
		return returnInternalServerResponse("unable to generate hash value for service parameters", "generate hash for service parameter")
	case *ErrorUnableToStoreSubscriptions:
		return returnInternalServerResponse("unable to store subscriptions in the database", "store subscriptions in the database")
	case *ErrorUnableToStoreServiceInstance:
		return returnInternalServerResponse(ErrMsgUnableToStoreInstance, ErrActionStoreInstance)
	case *ErrorUnableToSearchAPIs:
		return returnInternalServerResponse("unable to search API's", "get API ID's")
	case *ErrorUnableToCreateSubscription:
		return returnInternalServerResponse("unable to create subscriptions", "get API ID's")
	case *ErrorUnableToRetrieveBind:
		return returnInternalServerResponse(ErrMsgUnableToGetBind, "retrieve bind")
	case *ErrorUnableToStoreBind:
		return returnInternalServerResponse("unable to store Bind", "store Bind")
	case *ErrorUnableToDeleteBind:
		return returnInternalServerResponse("unable to unbind", "delete bind")
	case *ErrorUnableToGenerateKeys:
		return returnInternalServerResponse(ErrMsgUnableGenerateKeys, "generate keys for application")
	case *ErrorUnableToCreateAPIMResource:
		return returnInternalServerResponse("unable to create the API-M resource", ErrActionCreateAPIMResource)
	case *ErrorAPIMResourceAlreadyExists:
		return returnInternalServerResponse("API-M resource already exists", ErrActionCreateAPIMResource)
	case *ErrorAPIMResourceDoesNotExist:
		return returnInternalServerResponse("API-M resource does not exist", ErrActionUpdateAPIMResource)
	case *ErrorUnableToUpdateAPIMResource:
		return returnInternalServerResponse("unable to update the API-M resource", ErrActionUpdateAPIMResource)
	case *ErrorEmptyAPIParameterSet:
		return returnBadRequestResponsee("No APIs Defined", "get service parameters")
	default:
		return err
	}

}
