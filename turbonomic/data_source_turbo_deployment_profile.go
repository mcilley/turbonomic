package turbonomic

import (
	"encoding/json"
	"fmt"

	autodoc "github.com/foo/terraform-provider-utils/autodoc"
	log "github.com/foo/terraform-provider-utils/log"
	"github.foo.com/shared/terraform-provider-turbonomic/turbonomic/api"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTurboDeploymentProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTurboDeploymentProfileRead,

		Schema: map[string]*schema.Schema{

			autodoc.MetaAttribute: &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
				Description: fmt.Sprintf(
					"%s Deployment profiles are used in conjunction with templates to "+
						"create reservations. Deployment profiles specify the physical "+
						"details about how to deploy VMs from a given template (ie: "+
						"the physical files that will be copied to the deployed workload "+
						"as well as optional placement limitations).",
					autodoc.MetaSummary,
				),
			},

			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: fmt.Sprintf(
					"The name of the deployment profile. "+
						"%s \"DEP-PCKR20181023141134_CENTOS751804_BO1_4.0\"",
					autodoc.MetaExample,
				),
			},

			"class_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Description: fmt.Sprintf(
					"The type/category of the deployment profile. "+
						"%s \"ServiceCatalogItem\"",
					autodoc.MetaExample,
				),
			},
		},
	}
}

// -----------------------------------------------------------------------------
// Conversion Helpers
// -----------------------------------------------------------------------------

// buildDeploymentProfile constructs a concrete DeploymentProfileApiDTO struct
// from the ResourceData reference
func buildDeploymentProfile(d *schema.ResourceData) *api.DeploymentProfileApiDTO {
	log.Tracef("buildDeploymentProfile")

	obj := api.DeploymentProfileApiDTO{}
	obj.UUID = d.Id()
	obj.DisplayName = d.Get("display_name").(string)

	var attr interface{}
	var ok bool

	if attr, ok = d.GetOk("class_name"); ok {
		obj.ClassName = attr.(string)
	}

	return &obj
}

// setResourceDataFromDeploymentProfile takes a DeploymentProfileApiDTO reference
// and writes back the state to the ResourceData reference
func setResourceDataFromDeploymentProfile(d *schema.ResourceData, obj *api.DeploymentProfileApiDTO) {
	log.Tracef("setResourceDataFromDeploymentProfile")
	d.SetId(obj.UUID)
	d.Set("display_name", obj.DisplayName)
	d.Set("class_name", obj.ClassName)
}

// -----------------------------------------------------------------------------
// CRUD Functions
// -----------------------------------------------------------------------------

func dataSourceTurboDeploymentProfileRead(d *schema.ResourceData, meta interface{}) error {
	log.Tracef("data_source_turbo_deployment_profile.go#Read")

	client := meta.(*api.Client)
	obj := buildTemplate(d)

	log.Debugf("DeploymentProfileApiDTO: [%+v]", obj)

	queryObjs, queryErr := client.DeploymentProfiles()
	if queryErr != nil {
		return queryErr
	}
	numQueryObjs := len(queryObjs)
	if numQueryObjs == 0 {
		return fmt.Errorf("Data source deployment profile returned 0 results")
	}

	queryMatches := make([]api.DeploymentProfileApiDTO, 0)
	for idx, queryObj := range queryObjs {
		queryObjJSON, _ := json.MarshalIndent(queryObj, "", "  ")
		log.Debugf("[%d] => [%s]", idx, queryObjJSON)
		if queryObj.DisplayName == obj.DisplayName {
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
	log.Debugf("Query DeploymentProfileApiDTO: [%s]", queryObjJSON)

	setResourceDataFromDeploymentProfile(d, queryObj)

	return nil
}
