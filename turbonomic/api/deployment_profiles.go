package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/foo/terraform-provider-utils/log"
)

const (
	DeploymentProfilesPrefix = "deploymentprofiles"
)

// -----------------------------------------------------------------------------
// Struct Definition and Helpers
// -----------------------------------------------------------------------------

// -- Response Structs --

type DeploymentProfileApiDTO struct {
	Account          BaseApiDTO                      `json:"account,omitempty"`
	ClassName        string                          `json:"className,omitempty"`
	DeployParameters []DeploymentProfileTargetApiDTO `json:"deployParameters,omitempty"`
	DisplayName      string                          `json:"displayName,omitempty"`
	Links            []Link                          `json:"links,omitempty"`
	UUID             string                          `json:"uuid,omitempty"`
}

// DeploymentProfileApiInputDTO reads the attributes of the
// DeploymentProfileApiDTO and translates it to a DeploymentProfileApiInputDTO
// for use as inputs to create, update functions. Returns an error on failed
// conversion.
func (obj *DeploymentProfileApiDTO) DeploymentProfileApiInputDTO() (*DeploymentProfileApiInputDTO, error) {
	if obj == nil {
		return nil, nil
	}
	// quick conversion using JSON parsing as an intemediary since the
	// input DTO and the output DTO share an almost identical schema
	// definition
	//
	// Marshal to JSON byte string, then unmarshal back into the input DTO.
	// Handle any edge cases after.
	objJSONBytes, jsonEncErr := json.Marshal(obj)
	if jsonEncErr != nil {
		return nil, jsonEncErr
	}
	var inObj DeploymentProfileApiInputDTO
	jsonDecErr := json.Unmarshal(objJSONBytes, &inObj)
	if jsonDecErr != nil {
		return nil, jsonDecErr
	}
	// cascading conversions
	inObj.DeployParameters = make([]DeploymentProfileTargetApiInputDTO, len(obj.DeployParameters))
	for idx, val := range obj.DeployParameters {
		convObj, convErr := val.DeploymentProfileTargetApiInputDTO()
		if convErr != nil {
			return nil, convErr
		}
		inObj.DeployParameters[idx] = *convObj
	}
	// edge cases
	inObj.AccountId = obj.Account.UUID
	inObj.Name = obj.DisplayName
	return &inObj, nil
}

type DeploymentProfileParamApiDTO struct {
	ParameterType string              `json:"parameterType,omitempty"`
	Properties    []NameValueInputDTO `json:"properties,omitempty"`
}

type DeploymentProfileProviderApiDTO struct {
	Parameters []DeploymentProfileParamApiDTO `json:"parameters,omitempty"`
	Provider   BaseApiDTO                     `json:"provider,omitempty"`
}

// DeploymentProfileProviderApiInputDTO reads the attributes of the
// DeploymentProfileProviderApiDTO and translates it to a
// DeploymentProfileProviderApiInputDTO for use as inputs to create, update
// functions. Returns an error on failed conversion.
func (obj *DeploymentProfileProviderApiDTO) DeploymentProfileProviderApiInputDTO() (*DeploymentProfileProviderApiInputDTO, error) {
	if obj == nil {
		return nil, nil
	}
	// quick conversion using JSON parsing as an intemediary since the
	// input DTO and the output DTO share an almost identical schema
	// definition
	//
	// Marshal to JSON byte string, then unmarshal back into the input DTO.
	// Handle any edge cases after.
	objJSONBytes, jsonEncErr := json.Marshal(obj)
	if jsonEncErr != nil {
		return nil, jsonEncErr
	}
	var inObj DeploymentProfileProviderApiInputDTO
	jsonDecErr := json.Unmarshal(objJSONBytes, &inObj)
	if jsonDecErr != nil {
		return nil, jsonDecErr
	}
	// edge cases
	inObj.ProviderId = obj.Provider.UUID
	return &inObj, nil
}

type DeploymentProfileTargetApiDTO struct {
	Providers  []DeploymentProfileProviderApiDTO `json:"providers,omitempty"`
	TargetType string                            `json:"targetType,omitempty"`
}

// DeploymentProfileTargetApiInputDTO reads the attributes of the
// DeploymentProfileTargetApiDTO and translates it to a
// DeploymentProfileTargetApiInputDTO for use as inputs to create, update
// functions. Returns an error on failed conversion.
func (obj *DeploymentProfileTargetApiDTO) DeploymentProfileTargetApiInputDTO() (*DeploymentProfileTargetApiInputDTO, error) {
	if obj == nil {
		return nil, nil
	}
	// Copy same attributes
	inObj := DeploymentProfileTargetApiInputDTO{}
	inObj.TargetType = obj.TargetType
	inObj.Providers = make([]DeploymentProfileProviderApiInputDTO, len(obj.Providers))
	// cascading conversion of the deployment profile providers
	for idx, val := range obj.Providers {
		convObj, convErr := val.DeploymentProfileProviderApiInputDTO()
		if convErr != nil {
			return nil, convErr
		}
		inObj.Providers[idx] = *convObj
	}
	return &inObj, nil
}

// -- Input Structs --

// Deployment profile definition used as input to POST and PUT API endpoints
type DeploymentProfileApiInputDTO struct {
	// Business account related to the deployment profile
	AccountId string `json:"accountId,omitempty"`
	// Parameters
	DeployParameters []DeploymentProfileTargetApiInputDTO `json:"deployParameters,omitempty"`
	// Name of the deployment profile
	Name string `json:"name,omitempty"`
}

// Deployment profile provider definition used as input to POST and PUT API
// endpoints
type DeploymentProfileProviderApiInputDTO struct {
	// Provider parameters
	Parameters []DeploymentProfileParamApiDTO `json:"parameters,omitempty"`
	// Provider entity ID
	ProviderId string `json:"providerId,omitempty"`
}

// Deployment profile target definition used as input to POST and PUT API
// endpoints
type DeploymentProfileTargetApiInputDTO struct {
	// Provider's entity parameters
	Providers []DeploymentProfileProviderApiInputDTO `json:"providers,omitempty"`
	// Type of deployment target, Ex: "vCenter", "AWS", "Softlayer"
	TargetType string `json:"targetType,omitempty"`
}

// -----------------------------------------------------------------------------
// CRUD Implementation
// -----------------------------------------------------------------------------

// CreateTemplate creates a new Turbonomic deployment profile with the provided
// DeploymentProfileApiInputDTO. It returns a reference to the
// DeploymentProfileApiInputDTO that was returned, representing the created
// deployment profile or an error if encountered.
func (c *Client) CreateDeploymentProfile(obj *DeploymentProfileApiInputDTO) (*DeploymentProfileApiDTO, error) {
	log.Tracef("turbonomic/api/deployment_profiles.go#CreateDeploymentProfile")

	reqEndpoint := fmt.Sprintf("/%s", DeploymentProfilesPrefix)

	objJSONBytes, jsonEncErr := json.Marshal(obj)
	if jsonEncErr != nil {
		return nil, jsonEncErr
	}

	req, reqErr := c.NewRequest(
		http.MethodPost,
		reqEndpoint,
		bytes.NewBuffer(objJSONBytes),
	)
	if reqErr != nil {
		return nil, reqErr
	}

	var createObj DeploymentProfileApiDTO
	sendErr := c.SendAndParse(req, &createObj)
	if sendErr != nil {
		return nil, sendErr
	}

	return &createObj, nil
}

// ReadTemplate returns a DeploymentProfileApiDTO representing the deployment
// profile identified by the supplied UUID or an error if encountered.
func (c *Client) ReadDeploymentProfile(uuid string) (*DeploymentProfileApiDTO, error) {
	log.Tracef("turbonomic/api/deployment_profiles.go#ReadDeploymentProfile")

	reqEndpoint := fmt.Sprintf("/%s/%s", DeploymentProfilesPrefix, uuid)

	req, reqErr := c.NewRequest(
		http.MethodGet,
		reqEndpoint,
		nil,
	)
	if reqErr != nil {
		return nil, reqErr
	}

	var readObj DeploymentProfileApiDTO
	sendErr := c.SendAndParse(req, &readObj)
	if sendErr != nil {
		return nil, sendErr
	}

	return &readObj, nil
}

// UpdateTemplate updates a Turbonomic deployment profile identified by the
// given UUID.  The properties of the deployment profile will be set the values
// of the supplied DeploymentProfileApiInputDTO. This function returns a
// reference to the deployment profile with the updated properties or an error
// if encountered.
func (c *Client) UpdateDeploymentProfile(uuid string, obj *DeploymentProfileApiInputDTO) (*DeploymentProfileApiDTO, error) {
	log.Tracef("turbonomic/api/deployment_profiles.go#UpdateDeploymentProfile")

	reqEndpoint := fmt.Sprintf("/%s/%s", DeploymentProfilesPrefix, uuid)

	objJSONBytes, jsonEncErr := json.Marshal(obj)
	if jsonEncErr != nil {
		return nil, jsonEncErr
	}

	req, reqErr := c.NewRequest(
		http.MethodPut,
		reqEndpoint,
		bytes.NewBuffer(objJSONBytes),
	)
	if reqErr != nil {
		return nil, reqErr
	}

	var updateObj DeploymentProfileApiDTO
	sendErr := c.SendAndParse(req, &updateObj)
	if sendErr != nil {
		return nil, sendErr
	}

	return &updateObj, nil
}

// DeleteDeploymentProfile deletes the Deployment Profile identitied by the
// given UUID. It will return an error if encountered.
func (c *Client) DeleteDeploymentProfile(uuid string) error {
	log.Tracef("turbonomic/api/deployment_profiles.go#DeleteDeploymentProfile")

	reqEndpoint := fmt.Sprintf("/%s/%s", DeploymentProfilesPrefix, uuid)

	req, reqErr := c.NewRequest(
		http.MethodDelete,
		reqEndpoint,
		nil,
	)
	if reqErr != nil {
		return reqErr
	}

	return c.SendAndParse(req, nil)
}

// DeploymentProfiles returns all Deployment Profile objects in Turbonomic or
// an error if one encountered.
func (c *Client) DeploymentProfiles() ([]DeploymentProfileApiDTO, error) {
	log.Tracef("turbonomic/api/deployment_profiles.go#DeploymentProfiles")

	reqEndpoint := fmt.Sprintf("/%s", DeploymentProfilesPrefix)

	req, reqErr := c.NewRequest(
		http.MethodGet,
		reqEndpoint,
		nil,
	)
	if reqErr != nil {
		return nil, reqErr
	}

	var profiles []DeploymentProfileApiDTO
	sendErr := c.SendAndParse(req, &profiles)
	if sendErr != nil {
		return nil, sendErr
	}
	return profiles, nil
}
