package api

import (
	"fmt"
	"net/http"

	log "github.com/foo/terraform-provider-utils/log"
)

const (
	//MarketsPrefixPrefix - API endpoint for managing devices
	MarketsPrefix = "markets"
)

// -----------------------------------------------------------------------------
// Struct Definition and Helpers
// -----------------------------------------------------------------------------
// TurboMarket - struct definning a Turbo marker json response
type TurboMarket struct {
	Links            []Link `json:"links,omitempty"`
	UUID             string `json:"uuid,omitempty"`
	DisplayName      string `json:"displayName,omitempty"`
	ClassName        string `json:"className,omitempty"`
	State            string `json:"state,omitempty"`
	UnplacedEntities bool   `json:"unplacedEntities,omitempty"`
	EnvironmentType  string `json:"environmentType,omitempty"`
}

type CriteriaList struct {
	ExpVal        string `json:"expVal,omitempty"`
	ExpType       string `json:"expType,omitempty"`
	FilterType    string `json:"filterType,omitempty"`
	CaseSensitive bool   `json:"caseSensitive,omitempty"`
}

type ProviderConsumerGroup struct {
	Links           []Link         `json:"links,omitempty"`
	UUID            string         `json:"uuid,omitempty"`
	DisplayName     string         `json:"displayName,omitempty"`
	ClassName       string         `json:"className,omitempty"`
	EntitiesCount   int            `json:"entitiesCount,omitempty"`
	MembersCount    int            `json:"membersCount,omitempty"`
	GroupType       string         `json:"groupType,omitempty"`
	Severity        string         `json:"severity,omitempty"`
	IsStatic        bool           `json:"isStatic,omitempty"`
	LogicalOperator string         `json:"logicalOperator,omitempty"`
	CriteriaList    []CriteriaList `json:"criteriaList,omitempty"`
	EnvironmentType string         `json:"environmentType,omitempty"`
}

type TurboMarketPolicy struct {
	DisplayName   string                `json:"displayName,omitempty"`
	UUID          string                `json:"uuid,omitempty"`
	Links         []Link                `json:"links,omitempty"`
	Type          string                `json:"type,omitempty"`
	Name          string                `json:"name,omitempty"`
	Enabled       bool                  `json:"enabled,omitempty"`
	Capacity      float64               `json:"capacity,omitempty"`
	CommodityType string                `json:"commodityType,omitempty"`
	ConsumerGroup ProviderConsumerGroup `json:"consumerGroup,omitempty"`
	ProviderGroup ProviderConsumerGroup `json:"providerGroup,omitempty"`
}

// -----------------------------------------------------------------------------
// CRUD Implementation
// -----------------------------------------------------------------------------

// ReadMarket - reads the attributes of a TURBO Market
// identified by the supplied name
func (c *Client) ReadMarket(marketName string) (*TurboMarket, error) {
	log.Tracef("turbonomic/api/markets.go#ReadMarket")

	reqEndpoint := fmt.Sprintf("/%s", MarketsPrefix)

	req, reqErr := c.NewRequest(
		http.MethodGet,
		reqEndpoint,
		nil,
	)
	if reqErr != nil {
		return nil, reqErr
	}

	var markets []TurboMarket
	sendErr := c.SendAndParse(req, &markets)
	if sendErr != nil {
		return nil, sendErr
	}

	log.Debugf("markets: [%+v]", markets)

	for _, e := range markets {
		if e.DisplayName == marketName {
			return &e, nil
		}
	}

	return nil, fmt.Errorf(
		"Market [%s] is not found",
		marketName,
	)
}

// ReadMarketPolicy- reads the attributes of a TURBO Market
// identified by the supplied name
func (c *Client) ReadMarketPolicy(policyName string, marketUUID string) (*TurboMarketPolicy, error) {
	log.Tracef("turbonomic/api/markets.go#ReadMarketPolicy")

	reqEndpoint := fmt.Sprintf("/%s/%s/policies", MarketsPrefix, marketUUID)

	req, reqErr := c.NewRequest(
		http.MethodGet,
		reqEndpoint,
		nil,
	)
	if reqErr != nil {
		return nil, reqErr
	}

	var policies []TurboMarketPolicy
	sendErr := c.SendAndParse(req, &policies)
	if sendErr != nil {
		return nil, sendErr
	}

	log.Debugf("policies: [%+v]", policies)

	for _, e := range policies {
		if e.DisplayName == policyName {
			return &e, nil
		}
	}
	return nil, fmt.Errorf(
		"Market policy [%s] is not found",
		policyName,
	)
}
