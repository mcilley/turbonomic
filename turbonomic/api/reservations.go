package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	_ "time"

	log "github.com/foo/terraform-provider-utils/log"
)

const (
	//ReservationsPrefix - API endpoint for managing reservations
	ReservationsPrefix = "reservations"
)

// -----------------------------------------------------------------------------
// Struct Definitions for Creating Reservations
// -----------------------------------------------------------------------------
//
// Below structs are intended to represent this json model:
// {
// 	"action": "PLACEMENT",
// 	"demandName": "string",
// 	"deployDateTime": "string",
// 	"expireDateTime": "string",
// 	"parameters": [
// 	  {
// 		"deploymentParameters": {
// 		  "deploymentProfileID": "string",
// 		  "highAvailability": true,
// 		  "priority": "string"
// 		},
// 		"placementParameters": {
// 		  "constraintIDs": [
// 			"string"
// 		  ],
// 		  "count": 0,
// 		  "entityNames": [
// 			"string"
// 		  ],
// 		  "geographicRedundancy": true,
// 		  "templateID": "string"
// 		}
// 	  }
// 	],
// 	"reserveDateTime": "string"
//   }
type ReservationCreate struct {
	Action          string                 `json:"action,omitempty"`
	DemandName      string                 `json:"demandName,omitempty"`
	DeployDateTime  string                 `json:"deployDateTime,omitempty"`
	ExpireDateTime  string                 `json:"expireDateTime,omitempty"`
	Parameters      []ReservationParameter `json:"parameters,omitempty"`
	ReserveDateTime string                 `json:"reserveDateTime,omitempty"`
}

type ReservationParameter struct {
	DeploymentParameters ReservationDeploymentParameter `json:"deploymentParameters,omitempty"`
	PlacementParameters  ReservationPlacementParameter  `json:"placementParameters,omitempty"`
}

type ReservationDeploymentParameter struct {
	DeploymentProfileID string `json:"deploymentProfileID,omitempty"`
	HighAvailability    bool   `json:"highAvailability,omitempty"`
	Priority            bool   `json:"priority,omitempty"`
}

type ReservationPlacementParameter struct {
	ConstraintIDs        []string `json:"constraintIDs,omitempty"`
	Count                int      `json:"count,omitempty"`
	EntityNames          []string `json:"entityNames,omitempty"`
	GeographicRedundancy bool     `json:"geographicRedundancy,omitempty"`
	TemplateID           string   `json:"templateID,omitempty"`
}

// -----------------------------------------------------------------------------
// Struct Definitions for Create Reservation Response
// -----------------------------------------------------------------------------
type ReservationResponse struct {
	UUID            string         `json:"uuid,omitempty"`
	DisplayName     string         `json:"displayName,omitempty"`
	Count           int            `json:"count,omitempty"`
	Status          string         `json:"status,omitempty"`
	ReserveDateTime string         `json:"reserveDateTime,omitempty"`
	ExpireDateTime  string         `json:"expireDateTime,omitempty"`
	DeployDateTIme  string         `json:"deployDateTime,omitempty"`
	ReserveCount    int            `json:"reserveCount,omitempty"`
	DeployCount     int            `json:"deployCount,omitempty"`
	DemandEntities  []DemandEntity `json:"demandEntities,omitempty"`
}

type DemandEntity struct {
	UUID              string     `json:"uuid,omitempty"`
	DisplayName       string     `json:"displayName,omitempty"`
	Template          Identifier `json:"template,omitempty"`
	DeploymentProfile Identifier `json:"deploymentProfile,omitempty"`
	Placements        Placement  `json:"placements,omitempty"`
}

type Placement struct {
	ComputeResources []ComputeResource `json:"computeResources,omitempty"`
	StorageResources []StorageResource `json:"storageResources,omitempty"`
}

type ComputeResource struct {
	Provider Identifier `json:"provider,omitempty"`
	Stats    []stat     `json:"stats,omitempty"`
}

type StorageResource struct {
	Provider Identifier `json:"provider,omitempty"`
	Stats    []stat     `json:"stats,omitempty"`
	Type     string     `json:"type,omitempty"`
}

// Identifier - generic object idenitfier block
type Identifier struct {
	UUID        string `json:"uuid,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	ClassName   string `json:"className,omitempty"`
}

type stat struct {
	Name  string  `json:"name,omitempty"`
	Value float32 `json:"value,omitempty"`
}

// {
// 	"uuid": "_58zFA7GKEeiLOK9ARx98hA",
// 	"displayName": "ishashchuk-test",
// 	"count": 1,
// 	"status": "PLACEMENT_SUCCEEDED",
// 	"placementExpirationDateTime": "2018-09-06T00:16:19Z",
// 	"demandEntities": [
// 	  {
// 		"uuid": "_581hQbGKEeiLOK9ARx98hA",
// 		"displayName": "test.foo.com",
// 		"className": "VirtualMachine",
// 		"template": {
// 		  "uuid": "T564dbefa-8a99-1aad-542d-a1e751c5beba",
// 		  "displayName": "vcsa.host.foo.foo.com::TMP-PCKR20180605140114_CENTOS751804_SE1_2.0",
// 		  "className": "VirtualMachineProfile"
// 		},
// 		"deploymentProfile": {
// 		  "uuid": "_-qz6MKZDEeiLOK9ARx98hA",
// 		  "displayName": "DEP-PCKR20180605140114_CENTOS751804_SE1_2.0",
// 		  "className": "ServiceCatalogItem"
// 		},
// 		"placements": {
// 		  "computeResources": [
// 			{
// 			  "stats": [
// 				{
// 				  "name": "numOfCpu",
// 				  "value": 1
// 				},
// 				{
// 				  "name": "cpuSpeed",
// 				  "value": 2594
// 				},
// 				{
// 				  "name": "cpuConsumedFactor",
// 				  "value": 0.5
// 				},
// 				{
// 				  "name": "memorySize",
// 				  "value": 1048576
// 				},
// 				{
// 				  "name": "memoryConsumedFactor",
// 				  "value": 0.75
// 				},
// 				{
// 				  "name": "ioThroughput",
// 				  "value": 0
// 				},
// 				{
// 				  "name": "networkThroughput",
// 				  "value": 0
// 				}
// 			  ],
// 			  "provider": {
// 				"uuid": "48deb3f3-cff0-e711-0001-00000000003e",
// 				"displayName": "psc01n06.esx.foo.foo.com",
// 				"className": "PhysicalMachine"
// 			  }
// 			}
// 		  ],
// 		  "storageResources": [
// 			{
// 			  "stats": [
// 				{
// 				  "name": "diskSize",
// 				  "value": 27522.389
// 				},
// 				{
// 				  "name": "diskIops",
// 				  "value": 0
// 				}
// 			  ],
// 			  "provider": {
// 				"uuid": "5a62034b-b247a364-f836-0025b511a1df",
// 				"displayName": "lun01_psc01_vmax_foo",
// 				"className": "Storage"
// 			  }
// 			}
// 		  ]
// 		}
// 	  }
// 	]
//   }

// -----------------------------------------------------------------------------
// CRUD Implementation
// -----------------------------------------------------------------------------

// DeleteReservation deletes the Turbonomic reservation identitied by the given
// UUID. It will return an error if encountered.
func (c *Client) DeleteReservation(uuid string) error {
	log.Tracef("turbonomic/api/reservations.go#DeleteReservation")

	reqEndpoint := fmt.Sprintf("/%s/%s", ReservationsPrefix, uuid)

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

// ReadReservation reads the Turbonomic reservation identitied by the given
// UUID. It will return an error if encountered.
func (c *Client) ReadReservation(uuid string) (*ReservationResponse, error) {
	log.Tracef("turbonomic/api/reservations.go#ReadReservation")

	reqEndPoint := fmt.Sprintf("/%s/%s", ReservationsPrefix, uuid)

	req, reqErr := c.NewRequest(
		http.MethodGet,
		reqEndPoint,
		nil,
	)
	if reqErr != nil {
		return nil, reqErr
	}

	reqQuery := req.URL.Query()
	req.URL.RawQuery = reqQuery.Encode()

	var response ReservationResponse
	sendErr := c.SendAndParse(req, &response)
	if sendErr != nil {
		return nil, sendErr
	}

	log.Debugf("response: [%+v]", response)

	return &response, nil
}

// CreateReservation creates a new Turbonomic reservation with the
// provided ReservationCreate object. It returns a reference to the ReservationResponse object
// that was returned, representing the created template or an error if
// encountered.
func (c *Client) CreateReservation(rCreate *ReservationCreate, blocking bool) (*ReservationResponse, error) {
	log.Tracef("turbonomic/api/reservations.go#ReservationCreate")

	reqEndPoint := fmt.Sprintf("/%s", ReservationsPrefix)

	resJSON, jsonEncErr := json.Marshal(rCreate)
	if jsonEncErr != nil {
		return nil, jsonEncErr
	}

	log.Debugf("rCreate: [%s]", resJSON)

	req, reqErr := c.NewRequest(
		http.MethodPost,
		reqEndPoint,
		bytes.NewBuffer(resJSON),
	)
	if reqErr != nil {
		return nil, reqErr
	}

	reqQuery := req.URL.Query()

	if blocking {
		reqQuery.Set("apiCallBlock", "true")

	} else {
		reqQuery.Set("apiCallBlock", "false")
	}

	//reqQuery.Set("apiCallBlock", "false")
	req.URL.RawQuery = reqQuery.Encode()

	var response ReservationResponse
	sendErr := c.SendAndParse(req, &response)
	if sendErr != nil {
		return nil, sendErr
	}

	log.Debugf("response: [%+v]", response)

	return &response, nil
}
