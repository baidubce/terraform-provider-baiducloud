/*
Use this data source to query APPBLB Server Group list.

Example Usage

```hcl
data "baiducloud_appblb_server_groups" "default" {
 name = "testServerGroup"
}

output "server_groups" {
 value = "${data.baiducloud_appblb_server_groups.default.server_groups}"
}
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/appblb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudAppBLBServerGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudAppBLBServerGroupsRead,

		Schema: map[string]*schema.Schema{
			"blb_id": {
				Type:        schema.TypeString,
				Description: "ID of the LoadBalance instance to be queried",
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the Server Group to be queried",
				Optional:    true,
			},
			"exactly_match": {
				Type:        schema.TypeBool,
				Description: "Whether the name is an exact match or not, default false",
				Optional:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("name"); ok && v.(string) != "" {
						return false
					}

					return true
				},
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: "Query result output file path",
				Optional:    true,
				ForceNew:    true,
			},
			"filter": dataSourceFiltersSchema(),

			// Attributes used for result
			"server_groups": {
				Type:        schema.TypeList,
				Description: "A list of Application LoadBalance Server Group",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sg_id": {
							Type:        schema.TypeString,
							Description: "Server Group's ID",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Server Group's name",
							Computed:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "Server Group's description",
							Computed:    true,
						},
						"status": {
							Type:        schema.TypeString,
							Description: "Server Group status",
							Computed:    true,
						},
						"port_list": {
							Type:        schema.TypeList,
							Description: "Server Group backend port list",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Description: "Server Group port ID",
										Computed:    true,
									},
									"port": {
										Type:        schema.TypeInt,
										Description: "Server Group port",
										Computed:    true,
									},
									"type": {
										Type:        schema.TypeString,
										Description: "Server Group port protocol type",
										Computed:    true,
									},
									"status": {
										Type:        schema.TypeString,
										Description: "Server Group port status",
										Computed:    true,
									},
									"health_check": {
										Type:        schema.TypeString,
										Description: "Server Group port health check protocol",
										Computed:    true,
									},
									"health_check_port": {
										Type:        schema.TypeInt,
										Description: "Server Group port health check port",
										Computed:    true,
									},
									"health_check_timeout_in_second": {
										Type:        schema.TypeInt,
										Description: "Server Group health check timeout(second)",
										Computed:    true,
									},
									"health_check_interval_in_second": {
										Type:        schema.TypeInt,
										Description: "Server Group health check interval time(second)",
										Computed:    true,
									},
									"health_check_down_retry": {
										Type:        schema.TypeInt,
										Description: "Server Group health check down retry time",
										Computed:    true,
									},
									"health_check_up_retry": {
										Type:        schema.TypeInt,
										Description: "Server Group health check up retry time",
										Computed:    true,
									},
									"health_check_normal_status": {
										Type:        schema.TypeString,
										Description: "Server Group health check normal http status code",
										Computed:    true,
									},
									"health_check_url_path": {
										Type:        schema.TypeString,
										Description: "Server Group health check url path",
										Computed:    true,
									},
									"udp_health_check_string": {
										Type:        schema.TypeString,
										Description: "Server Group udp health check string",
										Computed:    true,
									},
								},
							},
						},
						"backend_server_list": {
							Type:        schema.TypeList,
							Description: "Server group bound backend server list",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:        schema.TypeString,
										Description: "Backend server instance ID",
										Computed:    true,
									},
									"weight": {
										Type:        schema.TypeInt,
										Description: "Backend server instance weight in this group",
										Computed:    true,
									},
									"private_ip": {
										Type:        schema.TypeString,
										Description: "Backend server instance bind private ip",
										Computed:    true,
									},
									"port_list": {
										Type:        schema.TypeSet,
										Description: "Backend server open port list",
										Computed:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"listener_port": {
													Type:        schema.TypeInt,
													Description: "Listener port",
													Computed:    true,
												},
												"backend_port": {
													Type:        schema.TypeInt,
													Description: "Backend open port",
													Computed:    true,
												},
												"port_type": {
													Type:        schema.TypeString,
													Description: "Port protocol type",
													Computed:    true,
												},
												"health_check_port_type": {
													Type:        schema.TypeString,
													Description: "Health check port protocol type",
													Computed:    true,
												},
												"status": {
													Type:        schema.TypeString,
													Description: "Port status, include Alive/Dead/Unknown",
													Computed:    true,
												},
												"port_id": {
													Type:        schema.TypeString,
													Description: "Port ID",
													Computed:    true,
												},
												"policy_id": {
													Type:        schema.TypeString,
													Description: "Port bind policy ID",
													Computed:    true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudAppBLBServerGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	appblbService := APPBLBService{client}

	blbId := ""
	if v, ok := d.GetOk("blb_id"); ok && v.(string) != "" {
		blbId = v.(string)
	}

	listServerGroupArgs := &appblb.DescribeAppServerGroupArgs{}
	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		listServerGroupArgs.Name = v.(string)
	}

	if v, ok := d.GetOk("exactly_match"); ok {
		listServerGroupArgs.ExactlyMatch = v.(bool)
	}

	action := "Query APPBLB " + blbId + " Server Groups " + listServerGroupArgs.Name
	serverGroupList, err := appblbService.ListAllServerGroups(blbId, listServerGroupArgs)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_appblb_server_groups", action, BCESDKGoERROR)
	}

	FilterDataSourceResult(d, &serverGroupList)

	if err := d.Set("server_groups", serverGroupList); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_appblb_server_groups", action, BCESDKGoERROR)
	}
	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), serverGroupList); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_appblb_server_groups", action, BCESDKGoERROR)
		}
	}

	return nil
}
