package turbonomic

import (
	"encoding/json"
	"fmt"
	"strings"

	autodoc "github.com/foo/terraform-provider-utils/autodoc"
	"github.com/foo/terraform-provider-utils/helper"
	log "github.com/foo/terraform-provider-utils/log"
	"github.foo.com/shared/terraform-provider-turbonomic/turbonomic/api"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTurboTemplate() *schema.Resource {
	// copy attributes from resource definition
	r := resourceTurboTemplate()
	ds := helper.DataSourceSchemaFromResourceSchema(r.Schema)

	// define searchable attributes for the data source
	ds["display_name"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		Description: fmt.Sprintf(
			"Name of the template. "+
				"%s \"PuREST VM\"",
			autodoc.MetaExample,
		),
	}
	ds["vcenter_server"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Description: fmt.Sprintf(
			"Hostname of a specific vcenter. "+
				"%s \"vcenter.host.foo.foo.com\"",
			autodoc.MetaExample,
		),
	}
	ds["has_deployment_profile"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
		Description: fmt.Sprintf(
			"If enabled, filter for only templates that have a deployment profile "+
				"associated with it. "+
				"%s ",
			autodoc.MetaExample,
		),
	}

	return &schema.Resource{
		Read: dataSourceTurboTemplateRead,
		// NOTE(ALL): See comments in the corresponding resource file
		Schema: ds,
	}
}

func dataSourceTurboTemplateRead(d *schema.ResourceData, meta interface{}) error {
	log.Tracef("data_source_turbo_template.go#Read")

	client := meta.(*api.Client)

	queryObjs, queryErr := client.Templates()
	if queryErr != nil {
		return queryErr
	}
	numQueryObjs := len(queryObjs)
	if numQueryObjs == 0 {
		return fmt.Errorf("Data source template returned 0 results")
	}

	displayName := d.Get("display_name").(string)
	hasDeploymentProfile := d.Get("has_deployment_profile").(bool)
	vcenter := ""

	var ok bool
	if vcenter, ok = d.Get("vcenter_server").(string); !ok {
		vcenter = ""
	}

	queryMatches := make([]api.TemplateApiDTO, 0)
	for idx, queryObj := range queryObjs {
		queryObjJSON, _ := json.MarshalIndent(queryObj, "", "  ")
		log.Debugf("[%d] => [%s]", idx, queryObjJSON)
		if queryObj.DisplayName == displayName &&
			((hasDeploymentProfile && queryObj.DeploymentProfile.UUID != "") ||
				(!hasDeploymentProfile && queryObj.DeploymentProfile.UUID == "")) &&
			(vcenter != "" && strings.Contains(queryObj.Model, vcenter)) {
			log.Debugf("  Matches!")
			queryMatches = append(queryMatches, queryObj)
		}
	}

	numQueryMatches := len(queryMatches)
	log.Debugf("numQueryMatches: [%d]", numQueryMatches)
	if numQueryMatches == 0 {
		return fmt.Errorf("Found [%d] templates matching the search criteria", numQueryMatches)
	} else if numQueryMatches > 1 {
		return fmt.Errorf("Found [%d] templates matching the search criteria", numQueryMatches)
	}

	queryObj := &queryMatches[0]
	queryObjJSON, _ := json.MarshalIndent(queryObj, "", "  ")
	log.Debugf("Query TemplateApiDTO: [%s]", queryObjJSON)

	setResourceDataFromTemplate(d, queryObj)

	return nil
}
