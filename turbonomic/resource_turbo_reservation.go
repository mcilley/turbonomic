package turbonomic

import (
	"fmt"
	"time"

	autodoc "github.com/foo/terraform-provider-utils/autodoc"

	log "github.com/foo/terraform-provider-utils/log"
	"github.foo.com/shared/terraform-provider-turbonomic/turbonomic/api"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

const (
	reservationOperationTimeout = 5 * time.Minute
)

func resourceTurboReservation() *schema.Resource {
	return &schema.Resource{

		Create: resourceTurboReservationCreate,
		Read:   resourceTurboReservationRead,
		Delete: resourceTurboReservationDelete,

		Schema: map[string]*schema.Schema{
			autodoc.MetaAttribute: &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
				Description: fmt.Sprintf(
					"%s Turbonomic Placement Recommendations.\n"+
						"NOTE: Once placement recommendation is retrieved, it only exists in state file. \n"+
						"\t Placement recommendatin is deleted by API itself when the call is completed "+
						"%s",
					autodoc.MetaSummary,
					autodoc.MetaImmutable,
				),
			},
			// REQUIRED
			"action": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"RESERVATION",
					"PLACEMENT",
				}, false),
				Description: fmt.Sprintf(
					"The intended action for the workload demand. "+
						"%s \"RESERVATION\"",
					autodoc.MetaExample,
				),
			},

			"entity_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: fmt.Sprintf(
					"Name of the instance to use for generating placement recommendation. "+
						"%s \"tftest.dev.foo.foo.com\"",
					autodoc.MetaExample,
				),
			},

			"template_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: fmt.Sprintf(
					"Template ID used for generating recommendation. "+
						"%s \"${data.turbonomic_template.template.id}\"",
					autodoc.MetaExample,
				),
			},

			// OPTIONAL
			"constraint_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Description: fmt.Sprintf(
					"List of constraint policies to use for reservation "+
						"%s [\"${data.turbonomic_market_policy.policy.id}\"]",
					autodoc.MetaExample,
				),
			},

			"deployment_profile_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: fmt.Sprintf(
					"ID of the deployment profile associated with the template. "+
						"%s \"${data.turbonomic_template.template.deployment_profile_id}\"",
					autodoc.MetaExample,
				),
			},

			"reservation_reserve_time": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: fmt.Sprintf(
					"timestamp of reservations for reservation/expiration times"+
						"%s \"1\"",
					autodoc.MetaExample,
				),
			},

			"reservation_expire_time": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: fmt.Sprintf(
					"timestamp of reservations for reservation/expiration times"+
						"%s \"1\"",
					autodoc.MetaExample,
				),
			},

			// bool value specifying whether we're creating a blocking/non-blocking request
			// false = nonblocking, true = blocking
			"reservation_blocking_req": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
				Description: fmt.Sprintf(
					"whether reservation requests will be blocking"+
						"%s",
					autodoc.MetaExample,
				),
			},

			//COMPUTED
			"compute_provider": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Generated recommendation for compute placememt",
			},

			"storage_provider": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Generated recommendation for storage placememt",
			},

			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the placement recommendation",
			},
		},
	}
}

func resourceTurboReservationCreate(d *schema.ResourceData, meta interface{}) error {
	log.Tracef("resource_turbo_reservation.go#Create")

	client := meta.(*api.Client)

	resCreate := api.ReservationCreate{}
	resDepParam := api.ReservationDeploymentParameter{}
	resPlcParam := api.ReservationPlacementParameter{}

	var err error
	var attr interface{}
	var ok bool

	if attr, ok = d.GetOk("action"); ok {
		resCreate.Action = attr.(string)
	}

	if attr, ok = d.GetOk("deployment_profile_id"); ok {
		resDepParam.DeploymentProfileID = attr.(string)
	}

	if attr, ok = d.GetOk("entity_name"); ok {
		entityNames := make([]string, 1)
		entityNames[0] = attr.(string)
		resPlcParam.EntityNames = entityNames
		resPlcParam.Count = 1
		resCreate.DemandName = attr.(string)
	}

	if attr, ok = d.GetOk("template_id"); ok {
		resPlcParam.TemplateID = attr.(string)
	}

	if attr, ok = d.GetOk("reservation_reserve_time"); ok {
		resCreate.ReserveDateTime = attr.(string)
	}

	if attr, ok = d.GetOk("reservation_expire_time"); ok {
		resCreate.ExpireDateTime = attr.(string)
	}

	if attr, ok = d.GetOk("constraint_ids"); ok {
		resPlcParam.ConstraintIDs = convertStringSet(attr.(*schema.Set))
	}

	resParams := make([]api.ReservationParameter, 1)
	resParams[0].PlacementParameters = resPlcParam
	resParams[0].DeploymentParameters = resDepParam

	resCreate.Parameters = resParams

	res, err := client.CreateReservation(&resCreate, d.Get("reservation_blocking_req").(bool))
	if err != nil {
		return err
	}

	d.SetId(res.UUID)

	//Wait for the Reservation resource to be ready
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"IN_PROGRESS", "LOADING", "RETRYING", "FUTURE", "UNFULFILLED"},
		Target:     []string{"PLACEMENT_SUCCEEDED", "RESERVED"},
		Refresh:    refreshReservation(d, meta),
		Timeout:    2 * time.Minute,
		MinTimeout: 3 * time.Second,
		Delay:      5 * time.Second,
	}

	_, err = stateConf.WaitForState()

	if err != nil {
		_ = client.DeleteReservation(res.UUID)
		d.SetId("")
		return fmt.Errorf("error waiting for turbonomic reservation id: %s Error: %s", res.UUID, err)
	}

	return nil
}

// resourceTurboReservationRead - placeholder read function(required by resource schema)
// Reservation doesn't exist in Turbonomic.
// It's only stored in state file. State file is populated at create time
func resourceTurboReservationRead(d *schema.ResourceData, meta interface{}) error {
	log.Tracef("resource_turbo_reservation.go#Read")

	return nil
}

// resourceTurboReservationDelete - reservation delete function(required by resource schema)
// Turbonomic reservations/placements are ephemeral and aren't likely in Turbonomic.
// Though we want to handle the scenario where reservations have not expired, and
// we're looking to destroy this resource.
func resourceTurboReservationDelete(d *schema.ResourceData, meta interface{}) error {
	log.Tracef("resource_turbo_reservation.go#Delete")

	client := meta.(*api.Client)
	_, readErr := client.ReadReservation(d.Id())

	if readErr == nil {
		return client.DeleteReservation(d.Id())
	}

	return nil
}

func setResourceDataFromReservation(d *schema.ResourceData, r *api.ReservationResponse) error {
	log.Tracef("resource_turbo_reservation.go#setResourceDataFromReservation")

	//check that demand ents is = count
	if len(r.DemandEntities) < 1 {
		return fmt.Errorf(
			"Reservation response is not in the expected format. \n"+
				"DemandEnitites is not equal to 1. \n"+
				"Response: [%+v]",
			r)
	}

	d.Set("compute_provider", r.DemandEntities[0].Placements.ComputeResources[0].Provider.DisplayName)
	d.Set("storage_provider", r.DemandEntities[0].Placements.StorageResources[0].Provider.DisplayName)
	d.Set("status", r.Status)
	return nil
}

func convertStringSet(set *schema.Set) []string {
	s := make([]string, 0, set.Len())
	for _, v := range set.List() {
		s = append(s, v.(string))
	}
	return s
}

// Poll status of reservation to verify if placement has succeeded for reservation
// A Successful reservation is when we recieve "PLACEMENT_SUCCEEDED" or "RESERVED" status
func refreshReservation(d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Debugf("Refreshing reservation state")

		client := meta.(*api.Client)

		resDetail, err := client.ReadReservation(d.Id())

		if err != nil {
			return nil, "Failed", err
		}

		switch resDetail.Status {
		case "IN_PROGRESS":
			return resDetail, resDetail.Status, nil
		case "RETRYING":
			return resDetail, resDetail.Status, nil
		case "FUTURE":
			return resDetail, resDetail.Status, nil
		case "LOADING":
			return resDetail, resDetail.Status, nil
		case "UNFULFILLED":
			return resDetail, resDetail.Status, nil
		case "PLACEMENT_SUCCEEDED":
			setResourceDataFromReservation(d, resDetail)
			return resDetail, resDetail.Status, nil
		case "RESERVED":
			setResourceDataFromReservation(d, resDetail)
			return resDetail, resDetail.Status, nil
		case "PLACEMENT_FAILED":
			return resDetail, resDetail.Status, fmt.Errorf("reservation unfulfilled, environment does not have resources to place the workload")
		default:
			return resDetail, resDetail.Status, fmt.Errorf("%s, unrecognized reservation status", resDetail.Status)
		}
	}
}
