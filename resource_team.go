package main

import "github.com/hashicorp/terraform/helper/schema"

func resourceTeam() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Create: resourceTeamCreate,
		Read:   resourceTeamRead,
		Update: resourceTeamUpdate,
		Delete: resourceTeamDelete,
	}
}

func resourceTeamCreate(d *schema.ResourceData, provider interface{}) error {
	panic("not implemented")
}

func resourceTeamDelete(d *schema.ResourceData, provider interface{}) error {
	panic("not implemented")
}

func resourceTeamRead(d *schema.ResourceData, provider interface{}) error {
	panic("not implemented")
}

func resourceTeamUpdate(d *schema.ResourceData, provider interface{}) error {
	panic("not implemented")
}
