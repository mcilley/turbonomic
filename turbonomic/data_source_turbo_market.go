package turbonomic

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	autodoc "github.com/foo/terraform-provider-utils/autodoc"
	log "github.com/foo/terraform-provider-utils/log"
	"github.foo.com/shared/terraform-provider-turbonomic/turbonomic/api"
)

func dataSourceTurboMarket() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTurboMarketRead,
		Schema: map[string]*schema.Schema{

			autodoc.MetaAttribute: &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
				Description: fmt.Sprintf(
					"%s Datasource for retrieving Turbonomic Market information.\n",
					autodoc.MetaSummary,
				),
			},
			// -- Searchable Attributes --
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Market",
					"Market_projection",
					"Market_Capacity",
				}, false),
				Default:     "Market",
				Description: "The name of Turbonomic Market. DEFAULT: `Market`",
			},
		},
	}
}

func dataSourceTurboMarketRead(d *schema.ResourceData, meta interface{}) error {
	log.Tracef("data_source_turbo_market.go/#Read")

	client := meta.(*api.Client)

	readMarket, readErr := client.ReadMarket(d.Get("display_name").(string))
	if readErr != nil {
		return readErr
	}

	mktObjJSON, _ := json.MarshalIndent(readMarket, "", "  ")
	log.Debugf("Read Market: [%s]", mktObjJSON)

	setResourceDataFromMarket(d, readMarket)

	return nil
}

func setResourceDataFromMarket(d *schema.ResourceData, market *api.TurboMarket) {
	log.Tracef("data_source_turbo_template.go/#setResourceDataFromMarket")

	d.SetId(market.UUID)
	d.Set("display_name", market.DisplayName)

}
