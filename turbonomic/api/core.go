package api

import (
	"encoding/json"
	"math"

	log "github.com/foo/terraform-provider-utils/log"
)

// Constants used to denote the various class types
const (
	ClassNameContainer       = "Container"
	ClassNamePhysicalMachine = "PhysicalMachine"
	ClassNameStorage         = "Storage"
	ClassNameVirtualMachine  = "VirtualMachine"
)

type BaseApiDTO struct {
	ClassName   string `json:"className,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Links       []Link `json:"links,omitempty"`
	UUID        string `json:"uuid,omitempty"`
}

type Link struct {
	Href      string `json:"href,omitempty"`
	Rel       string `json:"rel,omitempty"`
	Templated bool   `json:"templated"`
}

type NameValueInputDTO struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

const (
	// Type constant for ResourceApiDTO corresponding to a storage resource
	ResourceTypeDisk = "disk"
)

type ResourceApiDTO struct {
	Provider BaseApiDTO   `json:"provider,omitempty"`
	Stats    []StatApiDTO `json:"stats,omitempty"`
	Template string       `json:"template,omitempty"`
	Type     string       `json:"type,omitempty"`
}

const (
	ValueInfinity = "Infinity"
)

type StatApiDTO struct {
	Capacity           StatValueApiDTO    `json:"capacity,omitempty"`
	ClassName          string             `json:"className,omitempty"`
	DisplayName        string             `json:"displayName,omitempty"`
	Filters            []StatFilterApiDTO `json:"filters,omitempty"`
	Links              []Link             `json:"links,omitempty"`
	Name               string             `json:"name,omitempty"`
	NumRelatedEntities int                `json:"numRelatedEntities,omitempty"`
	RelatedEntity      BaseApiDTO         `json:"relatedEntity,omitempty"`
	RelatedEntityType  string             `json:"relatedEntityType,omitempty"`
	Reserved           StatValueApiDTO    `json:"reserved,omitempty"`
	Units              string             `json:"units,omitempty"`
	UUID               string             `json:"uuid,omitempty"`
	Value              float64            `json:"value,omitempty"`
	Values             StatValueApiDTO    `json:"values,omitempty"`
}

func (obj *StatApiDTO) UnmarshalJSON(b []byte) error {
	log.Tracef("StatApiDTO#UnmarshalJSON")

	var jsonDecErr error

	// temporary struct to unmarshal embedded,complex types
	var tmpJSON struct {
		Capacity      StatValueApiDTO    `json:"capacity,omitempty"`
		Filters       []StatFilterApiDTO `json:"filters,omitempty"`
		Links         []Link             `json:"links,omitempty"`
		RelatedEntity BaseApiDTO         `json:"relatedEntity,omitempty"`
		Reserved      StatValueApiDTO    `json:"reserved,omitempty"`
		Values        StatValueApiDTO    `json:"values,omitempty"`
	}
	jsonDecErr = json.Unmarshal(b, &tmpJSON)
	if jsonDecErr != nil {
		return jsonDecErr
	}
	obj.Capacity = tmpJSON.Capacity
	obj.Filters = tmpJSON.Filters
	obj.Links = tmpJSON.Links
	obj.RelatedEntity = tmpJSON.RelatedEntity
	obj.Reserved = tmpJSON.Reserved
	obj.Values = tmpJSON.Values

	// Unmarshal the rest of the primitive attributes into a map
	var objMap map[string]interface{}
	jsonDecErr = json.Unmarshal(b, &objMap)
	if jsonDecErr != nil {
		return jsonDecErr
	}
	obj.ClassName, _ = objMap["className"].(string)
	obj.DisplayName, _ = objMap["displayName"].(string)
	obj.Name, _ = objMap["name"].(string)
	obj.NumRelatedEntities, _ = objMap["numRelatedEntities"].(int)
	obj.RelatedEntityType, _ = objMap["relatedEntityType"].(string)
	obj.Units, _ = objMap["units"].(string)
	obj.UUID, _ = objMap["uuid"].(string)

	// Unmarshal into temporary struct to handle API irregularities
	var objJSON struct {
		// Value is usually a float64, however in some instances it returns as
		// the value "Infinity" as a string. Capture either option, and convert
		// to float64
		Value interface{} `json:"value,omitempty"`
	}
	jsonDecErr = json.Unmarshal(b, &objJSON)
	if jsonDecErr != nil {
		return jsonDecErr
	}
	if value, ok := objJSON.Value.(float64); ok {
		obj.Value = value
	} else if value, ok := objJSON.Value.(string); ok {
		if value == ValueInfinity {
			obj.Value = math.Inf(1)
		}
	}

	return nil
}

type StatFilterApiDTO struct {
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}

type StatValueApiDTO struct {
	Avg   int `json:"avg,omitempty"`
	Max   int `json:"max,omitempty"`
	Min   int `json:"min,omitempty"`
	Total int `json:"total,omitempty"`
}
