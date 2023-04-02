package turbonomic

import (
	"fmt"
	"strings"

	autodoc "github.com/foo/terraform-provider-utils/autodoc"
	log "github.com/foo/terraform-provider-utils/log"
	"github.foo.com/shared/terraform-provider-turbonomic/turbonomic/api"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceTurboTemplate() *schema.Resource {
	return &schema.Resource{

		Create: resourceTurboTemplateCreate,
		Read:   resourceTurboTemplateRead,
		Update: resourceTurboTemplateUpdate,
		Delete: resourceTurboTemplateDelete,

		Schema: map[string]*schema.Schema{
			autodoc.MetaAttribute: &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
				Description: fmt.Sprintf(
					"%s Tubonomic uses templates to reserve resource and deploy workloads "+
						"into the environment, calculate supply/demand changes in a plan, "+
						"and to calculate workloads for cloud environments.",
					autodoc.MetaSummary,
				),
			},

			// -- Required Arguments --

			"class_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: fmt.Sprintf(
					"The type of template. Values include: `\"Container\"`, `\"PhysicalMachine\"`, "+
						"`\"Storage\"`, `\"VirtualMachine\"`. "+
						"%s `\"VirtualMachine\"`",
					autodoc.MetaExample,
				),
				ValidateFunc: validation.StringInSlice([]string{
					api.ClassNameContainer,
					api.ClassNamePhysicalMachine,
					api.ClassNameStorage,
					api.ClassNameVirtualMachine,
					// NOTE(ALL): do not ignore case when comparing values
				}, false),
				DiffSuppressFunc: diffSuppressTemplateClassName,
			},

			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: fmt.Sprintf(
					"Name of the template. "+
						"%s \"PuREST VM\"",
					autodoc.MetaExample,
				),
			},

			// -- Optional Arguments --
			// Template resources

			"compute_resource": &schema.Schema{
				Type:     schema.TypeSet,
				Set:      schema.HashResource(resourceTurboResource()),
				Elem:     resourceTurboResource(),
				Optional: true,
				Description: fmt.Sprintf(
					"Set of compute resource statistics such as number of "+
						"CPU, CPU speed, memory size, etc. "+
						"%s {\n"+
						"    name  = \"cpuSpeed\"\n"+
						"    units = \"MHz\"\n"+
						"    value = 500\n"+
						"  }",
					autodoc.MetaExample,
				),
			},

			"infrastructure_resource": &schema.Schema{
				Type:     schema.TypeSet,
				Set:      schema.HashResource(resourceTurboResource()),
				Elem:     resourceTurboResource(),
				Optional: true,
				Description: fmt.Sprintf(
					"Set of infrastructure resource statistics such as "+
						"power size, space size, cooling, etc. "+
						"%s {\n"+
						"    name  = \"powerSize\"\n"+
						"    value = 1\n"+
						"  }",
					autodoc.MetaExample,
				),
			},

			"network_resource": &schema.Schema{
				Type:     schema.TypeSet,
				Set:      schema.HashResource(resourceTurboResource()),
				Elem:     resourceTurboResource(),
				Optional: true,
				Description: "Set of network resource statistics such as " +
					"network throughput, etc.",
			},

			"storage_resource": &schema.Schema{
				Type:     schema.TypeSet,
				Set:      schema.HashResource(resourceTurboResource()),
				Elem:     resourceTurboResource(),
				Optional: true,
				Description: fmt.Sprintf(
					"Set of storage resource statistics such as "+
						"disk I/O, disk size, percentage of disk consumed, etc. "+
						"%s {\n"+
						"    name  = \"diskSize\"\n"+
						"    units = \"GB\"\n"+
						"    value = 20\n"+
						"  }",
					autodoc.MetaExample,
				),
			},

			// -- Optional Arguments --

			"deployment_profile_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "ID of the deployment profile associated with this " +
					"template. In order to set up VM workloads and deploy to a " +
					"reservation, the VM template must have a deployment profile mapped " +
					"to it. This is the UUID of the deployment profile.",
			},

			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the template",
			},

			"price": &schema.Schema{
				Type:     schema.TypeFloat,
				Optional: true,
				Description: "Cost price associated with this template when " +
					"performing market analysis.",
			},

			"vendor": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Hardware, software vendor",
			},

			// -- Attributes --

			"discovered": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
				Description: "Whether or not the template is discovered or manually " +
					"created.",
			},
		},
	}
}

// resourceTurboResource defines the resource and schema definition for the
// various template resources. This loosely translates to a hybrid of the
// `ResourceApiDTO` and `StatApiDTO`.  The resources are used in the compute,
// storage, network, and infrastructure stats provided to the template in order
// to properly perform the analysis on where to place the object.
func resourceTurboResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: fmt.Sprintf(
					"Resource statistic name. "+
						"%s `\"cpuSpeed\"`",
					autodoc.MetaExample,
				),
			},
			"units": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: fmt.Sprintf(
					"Resource statistic units. "+
						"%s `\"MHz\"`",
					autodoc.MetaExample,
				),
			},
			"value": &schema.Schema{
				Type:     schema.TypeFloat,
				Required: true,
				Description: fmt.Sprintf(
					"Resource statistic value. "+
						"%s `1024`",
					autodoc.MetaExample,
				),
			},
		},
	}
}

// -----------------------------------------------------------------------------
// Resource Helpers and Validation
// -----------------------------------------------------------------------------

// diffSuppressTemplateClassName suppress an update on the class name attribute
// of the template resource. The value for class name differs as an input
// versus the API endpoint output. As an output, it appends "Profile" to the
// class name. For example, input class_name="VirtualMachine" => output
// class_name="VirtualMachineProfile".
func diffSuppressTemplateClassName(k, old, new string, d *schema.ResourceData) bool {
	log.Tracef("diffSuppressTemplateClassName")
	log.Debugf("k: [%s] old: [%s] new: [%s]", k, old, new)
	return strings.HasPrefix(old, new)
}

// -----------------------------------------------------------------------------
// Conversion Helpers
// -----------------------------------------------------------------------------

// buildTemplate constructs a concrete TemplateApiDTO struct from the
// ResourceData reference.
func buildTemplate(d *schema.ResourceData) *api.TemplateApiDTO {
	log.Tracef("buildTemplate")

	obj := api.TemplateApiDTO{}

	obj.UUID = d.Id()

	obj.ClassName = d.Get("class_name").(string)
	obj.DisplayName = d.Get("display_name").(string)

	var attr interface{}
	var ok bool

	if attr, ok = d.GetOk("compute_resource"); ok {
		obj.ComputeResources = setToResourceApiDTO(attr.(*schema.Set))
	}
	if attr, ok = d.GetOk("infrastructure_resource"); ok {
		obj.InfrastructureResources = setToResourceApiDTO(attr.(*schema.Set))
	}
	if attr, ok = d.GetOk("network_resource"); ok {
		obj.NetworkResources = setToResourceApiDTO(attr.(*schema.Set))
	}
	if attr, ok = d.GetOk("storage_resource"); ok {
		obj.StorageResources = setToStorageResource(attr.(*schema.Set))
	}

	if attr, ok = d.GetOk("deployment_profile_id"); ok {
		obj.DeploymentProfile.UUID = attr.(string)
	}
	if attr, ok = d.GetOk("description"); ok {
		obj.Description = attr.(string)
	}
	if attr, ok = d.GetOk("price"); ok {
		obj.Price = attr.(float64)
	}
	if attr, ok = d.GetOk("vendor"); ok {
		obj.Vendor = attr.(string)
	}

	if attr, ok = d.GetOk("discovered"); ok {
		obj.Discovered = attr.(bool)
	}

	return &obj
}

// setResourceDataFromTemplate takes a TemplateApiDTO references and writes
// back the state to the ResourceData reference.
func setResourceDataFromTemplate(d *schema.ResourceData, obj *api.TemplateApiDTO) {
	log.Tracef("setResourceDataFromTemplate")

	d.SetId(obj.UUID)

	d.Set("class_name", obj.ClassName)
	d.Set("display_name", obj.DisplayName)

	d.Set("compute_resource", resourceApiDTOToSet(obj.ComputeResources))
	d.Set("infrastructure_resource", resourceApiDTOToSet(obj.InfrastructureResources))
	d.Set("network_resource", resourceApiDTOToSet(obj.NetworkResources))
	d.Set("storage_resource", resourceApiDTOToSet(obj.StorageResources))

	d.Set("deployment_profile_id", obj.DeploymentProfile.UUID)
	d.Set("description", obj.Description)
	d.Set("price", obj.Price)
	d.Set("vendor", obj.Vendor)

	d.Set("discovered", obj.Discovered)
}

// setToStorageResource converts a schema.Set reference from the Template's
// storage resource set into a concrete ResourceApiDTO representation in the API
// layer. It serves as a wrapper to the setToResourceApiDTO function, but
// adds type and other meta information to the struct.
func setToStorageResource(s *schema.Set) []api.ResourceApiDTO {
	log.Tracef("setToResourceApiDTO")
	resourceArr := setToResourceApiDTO(s)
	for idx, _ := range resourceArr {
		resourceArr[idx].Type = api.ResourceTypeDisk
	}
	log.Debugf("resourceArr: [%+v]", resourceArr)
	return resourceArr
}

// setToResourceApiDTO converts a schema.Set reference from the Template's
// compute, infrastructure, network, or storage resource set into the concrete
// ResourceApiDTO representation in the API layer. Each of the items in the set
// are StatApiDTO structs which are associated to a ResourceApiDTO through the
// Stats attribute.
func setToResourceApiDTO(s *schema.Set) []api.ResourceApiDTO {
	log.Tracef("setToResourceApiDTO")
	// type assert the underlying *schema.Set to a list. The list will be
	// []interface{}
	attrList := s.List()
	attrListLen := len(attrList)
	statList := make([]api.StatApiDTO, attrListLen)
	// Each of the entries in the list is map[string]interface{}. Iterate over
	// each of the mapstruct entries in the set and convert to a concrete struct
	// implementation
	for idx, attrIface := range attrList {
		attrMap := attrIface.(map[string]interface{})
		statList[idx] = mapstructToStatApiDTO(attrMap)
	}
	log.Debugf("[]StatApiDTO: [%+v]", statList)
	// Embed the []api.StatApiDTO models into the stats attribute of the
	// ResourceApiDTO
	return []api.ResourceApiDTO{
		api.ResourceApiDTO{
			Stats: statList,
		},
	}
}

// resourceApiDTOToSet converts the concrete ResourceApiDTO list from the API
// layer into a *schema.Set for use in the template's compute, infrastructure,
// network, storage resource set.
func resourceApiDTOToSet(s []api.ResourceApiDTO) *schema.Set {
	log.Tracef("resourceApiDTOToSet")
	log.Debugf("[]api.ResourceApiDTO: [%+v]", s)
	ifaceArr := make([]interface{}, 0)
	log.Debugf("len([]api.ResourceApiDTO): [%d]", len(s))
	// for each api.ResourceApiDTO in the resource list
	for _, res := range s {
		log.Debugf("len([]api.ResourceApiDTO.Stats): [%d]", len(res.Stats))
		// for each api.StatApiDTO in the resource's stat list
		for _, stat := range res.Stats {
			// convert the StatApiDTO into a mapstruct and append it to the
			// interface list.
			ifaceArr = append(
				ifaceArr,
				statApiDTOToMapstruct(stat),
			)
		}
	}
	log.Debugf("ifaceArr: [%+v]", ifaceArr)
	// Convert the interface list into a set by hashing the resource definition
	return schema.NewSet(
		schema.HashResource(resourceTurboResource()),
		ifaceArr,
	)
}

// mapstructToStatApiDTO converts a map[string]interface{} from an entry in
// the Turbo Template resource's compute, network, etc sets into a StatApiDTO
// struct.
func mapstructToStatApiDTO(m map[string]interface{}) api.StatApiDTO {
	log.Tracef("mapstructToStatApiDTO")

	obj := api.StatApiDTO{}
	var ok bool

	if obj.Name, ok = m["name"].(string); !ok {
		obj.Name = ""
	}
	if obj.Units, ok = m["units"].(string); !ok {
		obj.Units = ""
	}
	if obj.Value, ok = m["value"].(float64); !ok {
		obj.Value = 0
	}

	log.Debugf("StatApiDTO: [%+v]", obj)
	return obj
}

// statApiDTOToMapstruct converts a StatApiDTO into a map[string]interface{}.
// The mapstructs are used in resourceApiDTOToSet to create a *schema.Set
// and write back the state.
func statApiDTOToMapstruct(obj api.StatApiDTO) map[string]interface{} {
	log.Tracef("statApiDTOToMapstruct")
	return map[string]interface{}{
		"name":  obj.Name,
		"units": obj.Units,
		"value": obj.Value,
	}
}

// -----------------------------------------------------------------------------
// CRUD Functions
// -----------------------------------------------------------------------------

func resourceTurboTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	log.Tracef("resourceTurboTemplateCreate")

	client := meta.(*api.Client)
	obj := buildTemplate(d)

	log.Debugf("TemplateApiDTO: [%+v]", obj)

	inObj, convErr := obj.TemplateApiInputDTO()
	if convErr != nil {
		return convErr
	}

	log.Debugf("TemplateApiInputDTO: [%+v]", inObj)

	createObj, createErr := client.CreateTemplate(inObj)
	if createErr != nil {
		return createErr
	}

	log.Debugf("Created TemplateApiDTO: [%+v]", createObj)

	setResourceDataFromTemplate(d, createObj)

	return nil
}

func resourceTurboTemplateRead(d *schema.ResourceData, meta interface{}) error {
	log.Tracef("resourceTurboTemplateRead")

	client := meta.(*api.Client)
	obj := buildTemplate(d)

	log.Debugf("TemplateApiDTO: [%+v]", obj)

	readObj, readErr := client.ReadTemplate(obj.UUID)
	if readErr != nil {
		return readErr
	}

	log.Debugf("Read TemplateApiDTO: [%+v]", readObj)

	setResourceDataFromTemplate(d, readObj)

	return nil
}

func resourceTurboTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Tracef("resourceTurboTemplateUpdate")

	client := meta.(*api.Client)
	obj := buildTemplate(d)

	log.Debugf("TemplateApiDTO: [%+v]", obj)

	inObj, convErr := obj.TemplateApiInputDTO()
	if convErr != nil {
		return convErr
	}

	log.Debugf("TemplateApiInputDTO: [%+v]", inObj)

	updateObj, updateErr := client.UpdateTemplate(obj.UUID, inObj)
	if updateErr != nil {
		return updateErr
	}

	log.Debugf("Update TemplateApiDTO: [%+v]", updateObj)

	setResourceDataFromTemplate(d, updateObj)

	return nil
}

func resourceTurboTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	log.Tracef("resourceTurboTemplateDelete")

	client := meta.(*api.Client)
	obj := buildTemplate(d)

	log.Debugf("TemplateApiDTO: [%+v]", obj)

	return client.DeleteTemplate(obj.UUID)
}
