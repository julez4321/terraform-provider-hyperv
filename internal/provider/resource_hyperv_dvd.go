package provider

import (
	"context"
	"log"
	"path"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/taliesins/terraform-provider-hyperv/api"
)

const (
	ReadDvdTimeout   = 1 * time.Minute
	CreateDvdTimeout = 5 * time.Minute
	UpdateDvdTimeout = 5 * time.Minute
	DeleteDvdTimeout = 1 * time.Minute
)

func resourceHyperVDvd() *schema.Resource {
	return &schema.Resource{
		Description: "This Hyper-V resource allows you to manage VHDs.",
		Timeouts: &schema.ResourceTimeout{
			Read:   schema.DefaultTimeout(ReadDvdTimeout),
			Create: schema.DefaultTimeout(CreateDvdTimeout),
			//Update: schema.DefaultTimeout(UpdateDvdTimeout),
			Delete: schema.DefaultTimeout(DeleteDvdTimeout),
		},
		CreateContext: resourceHyperVDvdCreate,
		ReadContext:   resourceHyperVDvdRead,
		//UpdateContext: resourceHyperVVhdUpdate,
		DeleteContext: resourceHyperVDvdDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"path": {
				ForceNew: true,
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					extension := path.Ext(newValue)
					computedPath := strings.TrimSuffix(newValue, extension)

					// Ignore differencing
					if strings.HasPrefix(strings.ToLower(oldValue), strings.ToLower(computedPath)) && strings.HasSuffix(strings.ToLower(oldValue), strings.ToLower(extension)) {
						return true
					}

					if strings.EqualFold(oldValue, newValue) {
						return true
					}

					return false
				},
				Description: "Path to the new iso that is being created or being copied to. If a filename or relative path is specified, the new virtual hard disk path is calculated relative to the current working directory. Depending on the source selected, the path will be used to determine where to copy source vhd/vhdx/vhds file to.",
			},
			"ip": {
				ForceNew:    true,
				Type:        schema.TypeString,
				Required:    true,
				Description: "This field is mutually exclusive with the fields `source_vm`, `parent_path`, `source_disk`. This value can be a url or a path (including wildcards). Box, Zip and 7z files will automatically be expanded. The destination folder will be the directory portion of the path. If expanded files have a folder called `Virtual Machines`, then the `Virtual Machines` folder will be used instead of the entire archive contents. ",
			},
			"exists": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Does dvd exist.",
			},
		},
	}
}

func resourceHyperVDvdCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO][hyperv][create] creating hyperv dvd: %#v", d)
	c := meta.(api.Client)

	path := (d.Get("path")).(string)
	ip := (d.Get("ip")).(string)

	err := c.CreateDvd(ctx, path, ip)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(path)
	log.Printf("[INFO][hyperv][create] created hyperv vhd: %#v", d)

	return nil
}

func resourceHyperVDvdRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO][hyperv][read] reading hyperv vhd: %#v", d)
	c := meta.(api.Client)

	path := d.Id()
	ip := (d.Get("ip")).(string)

	dvd, err := c.GetDvd(ctx, path, ip)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO][hyperv][read] retrieved dvd: %+v", dvd)

	if err := d.Set("path", dvd.Path); err != nil {
		return diag.FromErr(err)
	}

	if dvd.Path == "" {
		log.Printf("[INFO][hyperv][read] unable to retrieved dvd: %+v", path)
		if err := d.Set("exists", false); err != nil {
			return diag.FromErr(err)
		}
	} else {
		log.Printf("[INFO][hyperv][read] retrieved dvd: %+v", path)
		if err := d.Set("exists", true); err != nil {
			return diag.FromErr(err)
		}
	}

	log.Printf("[INFO][hyperv][read] read hyperv dvd: %#v", d)

	return nil
}

func resourceHyperVDvdDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO][hyperv][delete] deleting hyperv vhd: %#v", d)

	c := meta.(api.Client)

	path := d.Id()

	err := c.DeleteDvd(ctx, path)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO][hyperv][delete] deleted hyperv vhd: %#v", d)
	return nil
}
