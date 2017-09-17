package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/satori/go.uuid"
)

func resourcePipeline() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			"team": {
				Type:        schema.TypeString,
				Computed:    true, // hardcoded on creation for now until we figure out team resource impl
				Description: "",
			},
			"paused": {
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
				Description: "",
			},
			"public": {
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
				Description: "Makes the pipeline publicly viewable (ie does not require logging in)",
			},
			"yaml": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
		},
		Create: resourcePipelineCreate,
		Read:   resourcePipelineRead,
		Update: resourcePipelineUpdate,
		Delete: resourcePipelineDelete,
	}
}

func resourcePipelineCreate(d *schema.ResourceData, provider interface{}) error {
	p := provider.(*providerConfig)
	name := d.Get("name").(string)
	paused := d.Get("paused").(bool)
	public := d.Get("public").(bool)
	yaml := d.Get("yaml").(string)

	// hardcoded on creation for now until we figure out team resource impl
	d.Set("team", "main")
	// todo: how should we treat the ID? what makes a good one etc?
	d.SetId(uuid.NewV4().String())

	d.Partial(true)

	// first set pipeline
	err := p.fly.SetPipeline(name, yaml)
	if err != nil {
		return err
	}
	d.SetPartial("yaml")

	// then set whether it's paused or not
	if paused {
		err = p.fly.PausePipeline(name)
	} else {
		err = p.fly.UnpausePipeline(name)
	}
	if err != nil {
		return err
	}
	d.SetPartial("paused")

	// set public visibility
	if public {
		err = p.fly.ExposePipeline(name)
	} else {
		err = p.fly.HidePipeline(name)
	}
	if err != nil {
		return err
	}
	d.SetPartial("public")

	d.Partial(false)
	return nil
}

func resourcePipelineDelete(d *schema.ResourceData, provider interface{}) error {
	p := provider.(*providerConfig)
	return p.fly.DestroyPipeline(d.Get("name").(string))
}

func resourcePipelineRead(d *schema.ResourceData, provider interface{}) error {
	name := d.Get("name").(string)
	team := d.Get("team").(string)
	p := provider.(*providerConfig)
	pipelines, err := p.fly.Pipelines()
	if err != nil {
		return err
	}
	// find the pipeline in the list output
	// todo: improve efficiency; this pulls down the whole list of pipelines for every pipeline we read. at some point, that will have to be bad for performance
	var pipeline *pipeline
	for _, pl := range pipelines {
		if pl.name == name && pl.team == team {
			pipeline = pl
			break
		}
	}
	// update le object state
	d.Set("paused", pipeline.paused)
	d.Set("public", pipeline.public)

	// todo: we should also set the yaml, but first we need a way of dealing with the defaults the server sets.
	// ^ would it work if we read the pipeline back from the server after every time we `set-pipeline`, and
	//   put that value in some other field? when that field HasUpdate() we would `set-pipeline` the yaml field again..

	return nil
}

func resourcePipelineUpdate(d *schema.ResourceData, provider interface{}) error {
	p := provider.(*providerConfig)
	name := d.Get("name").(string)
	paused := d.Get("paused").(bool)
	public := d.Get("public").(bool)
	yaml := d.Get("yaml").(string)

	d.Partial(true)
	if d.HasChange("name") {
		oldName, newName := d.GetChange("name")
		err := p.fly.RenamePipeline(oldName.(string), newName.(string))
		if err != nil {
			return err
		}
		d.SetPartial("name")
	}
	if d.HasChange("paused") {
		var err error
		if paused {
			err = p.fly.PausePipeline(name)
		} else {
			err = p.fly.UnpausePipeline(name)
		}
		if err != nil {
			return err
		}
		d.SetPartial("paused")
	}
	if d.HasChange("public") {
		var err error
		if public {
			err = p.fly.ExposePipeline(name)
		} else {
			err = p.fly.HidePipeline(name)
		}
		if err != nil {
			return err
		}
		d.SetPartial("public")
	}
	if d.HasChange("yaml") {
		err := p.fly.SetPipeline(name, yaml)
		if err != nil {
			return err
		}
		d.SetPartial("yaml")
	}
	d.Partial(false)
	return nil
}
