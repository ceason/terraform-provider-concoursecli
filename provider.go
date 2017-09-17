package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

type providerConfig struct {
	fly *flyCli
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"target": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ConfigureFunc: configureProvider,
		ResourcesMap: map[string]*schema.Resource{
			"concoursecli_pipeline": resourcePipeline(),
		},
	}
}


func configureProvider(d *schema.ResourceData) (interface{}, error) {
	// todo: validate target exists
	cfg := &providerConfig{
		fly: &flyCli{
			target: d.Get("target").(string),
		},
	}
	return cfg, nil
}
