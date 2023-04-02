package turbonomic

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	autodoc "github.com/foo/terraform-provider-utils/autodoc"
	log "github.com/foo/terraform-provider-utils/log"
	"github.foo.com/shared/terraform-provider-turbonomic/turbonomic/api"
)

func dataSourceTurboMarketPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTurboMarketPolicyRead,
		Schema: map[string]*schema.Schema{

			autodoc.MetaAttribute: &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
				Description: fmt.Sprintf(
					"%s Datasource for retrieving Turbonomic Market Policy information.\n",
					autodoc.MetaSummary,
				),
			},
			// -- Searchable Attributes --
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: fmt.Sprintf(
					"The name of Turbonomic Market Policy"+
						"%s \"BO1_vm_placement\"",
					autodoc.MetaExample,
				),
			},

			// -- Searchable Attributes --
			"market_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: fmt.Sprintf(
					"Market ID for searching Market Policy"+
						"%s \"BO1_vm_placement\"",
					autodoc.MetaExample,
				),
			},
		},
	}
}

func dataSourceTurboMarketPolicyRead(d *schema.ResourceData, meta interface{}) error {
	log.Tracef("data_source_turbo_market.go/#Read")

	client := meta.(*api.Client)

	readPolicy, readErr := client.ReadMarketPolicy(
		d.Get("display_name").(string),
		d.Get("market_id").(string),
	)
	if readErr != nil {
		return readErr
	}

	log.Debugf("Read Market Policy: [%s]", readPolicy)

	setResourceDataFromMarketPolicy(d, readPolicy)

	return nil
}

func setResourceDataFromMarketPolicy(d *schema.ResourceData, policy *api.TurboMarketPolicy) {
	log.Tracef("data_source_turbo_template.go/#setResourceDataFromMarket")

	d.SetId(policy.UUID)
	d.Set("display_name", policy.DisplayName)

}
