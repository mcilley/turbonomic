package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/foo/terraform-provider-utils/log"
)

const (
	//TemplatesPrefix - API endpoint for managing devices
	TemplatesPrefix = "templates"
)

// -----------------------------------------------------------------------------
// Struct Definition and Helpers
// -----------------------------------------------------------------------------

// Template struct definition used for input to the POST and PUT API endpoints.
type TemplateApiInputDTO struct {
	// Type of template, this should be one of the ClassNameXxx constants
	ClassName string `json:"className,omitempty"`
	// Command with arguments, used for ClassNameContainer templates
	CmdWithArgs string `json:"cmdWithArgs,omitempty"`
	// Compute resources: Number of CPU, CPU speed, memory size, etc
	ComputeResources []ResourceApiDTO `json:"computeResources,omitempty"`
	// Ids of the Deployment Profiles associated with this template. In order to
	// set up VM workloads and deploy to a reservation, the VM template must have
	// a deployment profile mapped to it. This is the UUID of the
	// deployment profile.
	DeploymentProfileId string `json:"deploymentProfileId,omitempty"`
	// Description of the template
	Description string `json:"description,omitempty"`
	// Name of the template
	DisplayName string `json:"displayName,omitempty"`
	// Profile image, used for ClassNameContainer templates
	Image string `json:"image,omitempty"`
	// Profile image tag, used for ClassNameContainer templates
	ImageTag string `json:"imageTag,omitempty"`
	// Infrastructure resources: Power, size, cooling, etc
	InfrastructureResources []ResourceApiDTO `json:"infrastructureResources,omitempty"`
	// Network resources
	NetworkResources []ResourceApiDTO `json:"networkResources,omitempty"`
	// Cost price associated with this template when performing market analysis
	Price float64 `json:"price,omitempty"`
	// Storage resources: Disk I/O, disk size, percentage of disk consumed, etc
	StorageResources []ResourceApiDTO `json:"storageResources,omitempty"`
	// Hardware, software vendor
	Vendor string `json:"vendor,omitempty"`
}

// Template struct definition used as output from the API endpoints
type TemplateApiDTO struct {
	// Type of template, this should be one of the ClassNameXxx constants
	ClassName string `json:"className,omitempty"`
	// Command with arguments, used for ClassNameContainer templates
	CmdWithArgs string `json:"cmdWithArgs,omitempty"`
	// Compute resources: Number of CPU, CPU speed, memory size, etc
	ComputeResources []ResourceApiDTO `json:"computeResources,omitempty"`
	// Database edition, used for Database templates
	DbEdition string `json:"dbEdition,omitempty"`
	// Database engine, used for Database templates
	DbEngine string `json:"dbEngine,omitempty"`
	// Deployment profile associated with this template
	DeploymentProfile DeploymentProfileApiDTO `json:"deploymentProfile,omitempty"`
	// Description of the template
	Description string `json:"description,omitempty"`
	// Whether or not the template is discovered or manually created
	Discovered bool `json:"discovered"`
	// Name of the template
	DisplayName string `json:"displayName,omitempty"`
	// Profile image, used for ClassNameContainer templates
	Image string `json:"image,omitempty"`
	// Profile image tag, used for ClassNameContainer templates
	ImageTag string `json:"imageTag,omitempty"`
	// Infrastructure resources: Power, size, cooling, etc
	InfrastructureResources []ResourceApiDTO `json:"infrastructureResources,omitempty"`
	// Links associated with the template
	Links []Link `json:"links,omitempty"`
	// API model URI
	Model string `json:"model,omitempty"`
	// Network resources
	NetworkResources []ResourceApiDTO `json:"networkResources,omitempty"`
	// Cost price associated with this template when performing market analysis
	Price float64 `json:"price,omitempty"`
	// Storage resources: Disk I/O, disk size, percentage of disk consumed, etc
	StorageResources []ResourceApiDTO `json:"storageResources,omitempty"`
	// Unique identifier for this template
	UUID string `json:"uuid,omitempty"`
	// Hardware, software vendor
	Vendor string `json:"vendor,omitempty"`
}

// TemplateApiInputDTO reads the attributes of the TemplateApiDTO and
// translates it to a TemplateApiInputDTO for use as inputs to create, update
// functions. Returns an error on failed conversion.
func (obj *TemplateApiDTO) TemplateApiInputDTO() (*TemplateApiInputDTO, error) {
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
	var inObj TemplateApiInputDTO
	jsonDecErr := json.Unmarshal(objJSONBytes, &inObj)
	if jsonDecErr != nil {
		return nil, jsonDecErr
	}
	// edge cases
	inObj.DeploymentProfileId = obj.DeploymentProfile.UUID
	return &inObj, nil
}

// -----------------------------------------------------------------------------
// CRUD Implementation
// -----------------------------------------------------------------------------

// CreateTemplate creates a new Turbonomic template with the
// provided TemplateApiInputDTO. It returns a reference to the TemplateApiDTO
// that was returned, representing the created template or an error if
// encountered.
func (c *Client) CreateTemplate(obj *TemplateApiInputDTO) (*TemplateApiDTO, error) {
	log.Tracef("turbonomic/api/templates.go#CreateTemplate")

	reqEndpoint := fmt.Sprintf("/%s", TemplatesPrefix)

	objJSONBytes, jsonEncErr := json.Marshal(obj)

	log.Debugf("templateCreate: [%s]", objJSONBytes)

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

	var createObj TemplateApiDTO
	sendErr := c.SendAndParse(req, &createObj)
	if sendErr != nil {
		return nil, sendErr
	}

	return &createObj, nil
}

// ReadTemplate returns a TemplateApiDTO representing the template identified
// by the supplied UUID or an error if encountered.
func (c *Client) ReadTemplate(uuid string) (*TemplateApiDTO, error) {
	log.Tracef("turbonomic/api/templates.go#ReadTemplate")

	reqEndpoint := fmt.Sprintf("/%s/%s", TemplatesPrefix, uuid)

	req, reqErr := c.NewRequest(
		http.MethodGet,
		reqEndpoint,
		nil,
	)
	if reqErr != nil {
		return nil, reqErr
	}

	var readObjs []TemplateApiDTO
	sendErr := c.SendAndParse(req, &readObjs)
	if sendErr != nil {
		return nil, sendErr
	}

	lenReadObjs := len(readObjs)
	if lenReadObjs == 0 {
		return nil, fmt.Errorf("Could not find template by UUID [%s]", uuid)
	} else if lenReadObjs > 1 {
		return nil, fmt.Errorf("Matched multiple templates by UUID [%s]", uuid)
	}

	return &readObjs[0], nil
}

// UpdateTemplate updates a Turbonomic template identified by the given UUID.
// The properties of the template will be set the values of the supplied
// TemplateApiInputDTO. This function returns a reference to the template with
// the updated properties or an error if encountered.
func (c *Client) UpdateTemplate(uuid string, obj *TemplateApiInputDTO) (*TemplateApiDTO, error) {
	log.Tracef("turbonomic/api/templates.go#UpdateTemplate")

	reqEndpoint := fmt.Sprintf("/%s/%s", TemplatesPrefix, uuid)

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

	var updateObj TemplateApiDTO
	sendErr := c.SendAndParse(req, &updateObj)
	if sendErr != nil {
		return nil, sendErr
	}

	return &updateObj, nil
}

// DeleteTemplate deletes the Turbonomic template identitied by the given
// UUID. It will return an error if encountered.
func (c *Client) DeleteTemplate(uuid string) error {
	log.Tracef("turbonomic/api/templates.go#DeleteTemplate")

	reqEndpoint := fmt.Sprintf("/%s/%s", TemplatesPrefix, uuid)

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

// Templates returns all Template objects in Turbonomic or an error
// if one encountered.
func (c *Client) Templates() ([]TemplateApiDTO, error) {
	log.Tracef("turbonomic/api/templates.go#Templates")

	reqEndpoint := fmt.Sprintf("/%s", TemplatesPrefix)

	req, reqErr := c.NewRequest(
		http.MethodGet,
		reqEndpoint,
		nil,
	)
	if reqErr != nil {
		return nil, reqErr
	}

	var templates []TemplateApiDTO
	sendErr := c.SendAndParse(req, &templates)
	if sendErr != nil {
		return nil, sendErr
	}

	return templates, nil
}
