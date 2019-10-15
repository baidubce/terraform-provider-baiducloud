/*
Use this data source to query APPBLB Listener list.

Example Usage

```hcl
data "baiducloud_appblb_listeners" "default" {
 blb_id = "lb-0d29a3f6"
}

output "listeners" {
 value = "${data.baiducloud_appblb_listeners.default.listeners}"
}
```
*/
package baiducloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudAppBLBListeners() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudAppBLBListenersRead,

		Schema: map[string]*schema.Schema{
			"protocol": {
				Type:        schema.TypeString,
				Description: "Protocol of the Listener to be queried",
				Optional:    true,
				ForceNew:    true,
			},
			"blb_id": {
				Type:        schema.TypeString,
				Description: "ID of the Application LoadBalance instance",
				Required:    true,
				ForceNew:    true,
			},
			"listener_port": {
				Type:         schema.TypeInt,
				Description:  "The port of the Listener to be queried",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validatePort(),
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: "Query result output file path",
				Optional:    true,
				ForceNew:    true,
			},

			// Attributes used for result
			"listeners": {
				Type:        schema.TypeList,
				Description: "A list of Application LoadBalance Listener",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"listener_port": {
							Type:        schema.TypeInt,
							Description: "Listener bind port",
							Computed:    true,
						},
						"protocol": {
							Type:        schema.TypeString,
							Description: "Listener protocol",
							Computed:    true,
						},
						"scheduler": {
							Type:        schema.TypeString,
							Description: "Load balancing algorithm",
							Computed:    true,
						},
						"tcp_session_timeout": {
							Type:        schema.TypeInt,
							Description: "TCP Listener connetion session timeout time",
							Computed:    true,
						},
						// http & https
						"keep_session": {
							Type:        schema.TypeBool,
							Description: "Listener keepSession or not",
							Computed:    true,
						},
						//http & https
						"keep_session_type": {
							Type:        schema.TypeString,
							Description: "Listener keepSessionType option",
							Computed:    true,
						},
						// http & https
						"keep_session_timeout": {
							Type:        schema.TypeInt,
							Description: "Listener keepSessionTimeout value",
							Computed:    true,
						},
						// http & https
						"keep_session_cookie_name": {
							Type:        schema.TypeString,
							Description: "Listener keepSeesionCookieName",
							Computed:    true,
						},
						// http & https
						"x_forwarded_for": {
							Type:        schema.TypeBool,
							Description: "Listener xForwardedFor, determine get client real ip or not, default false",
							Computed:    true,
						},
						// http & https
						"server_timeout": {
							Type:        schema.TypeInt,
							Description: "Backend server maximum timeout time, only support in [1, 3600] second, default 30s",
							Computed:    true,
						},
						// http
						"redirect_port": {
							Type:        schema.TypeInt,
							Description: "Listener redirect request to https listener port",
							Computed:    true,
						},
						// https && ssl
						"cert_ids": {
							Type:        schema.TypeList,
							Description: "Listener bind certifications",
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						// https && ssl
						"ie6_compatible": {
							Type:        schema.TypeBool,
							Description: "Listener support ie6 option, default true",
							Computed:    true,
						},
						// https && ssl
						"encryption_type": {
							Type:        schema.TypeString,
							Description: "Listener encryption option",
							Computed:    true,
						},
						// https && ssl
						"encryption_protocols": {
							Type:        schema.TypeList,
							Description: "Listener encryption protocol",
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						// https && ssl
						"dual_auth": {
							Type:        schema.TypeBool,
							Description: "Listener open dual authorization or not, default false",
							Computed:    true,
						},
						// https && ssl
						"client_cert_ids": {
							Type:        schema.TypeList,
							Description: "Listener import cert list, only useful when dual_auth is true",
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"policys": {
							Type:        schema.TypeSet,
							Description: "Listener's policy",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Description: "Policy's ID",
										Computed:    true,
									},
									"description": {
										Type:        schema.TypeString,
										Description: "Policy's description",
										Computed:    true,
									},
									"port_type": {
										Type:        schema.TypeString,
										Description: "Policy bind port protocol",
										Computed:    true,
									},
									"app_server_group_id": {
										Type:        schema.TypeString,
										Description: "Policy bind server group ID",
										Computed:    true,
									},
									"app_server_group_name": {
										Type:        schema.TypeString,
										Description: "Policy bind server group name",
										Computed:    true,
									},
									"frontend_port": {
										Type:        schema.TypeInt,
										Description: "Frontend port",
										Computed:    true,
									},
									"backend_port": {
										Type:        schema.TypeInt,
										Description: "Backend port",
										Computed:    true,
									},
									"priority": {
										Type:        schema.TypeInt,
										Description: "Policy priority",
										Computed:    true,
									},
									"rule_list": {
										Type:        schema.TypeList,
										Description: "Policy rule list",
										Computed:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Description: "Rule key",
													Computed:    true,
												},
												"value": {
													Type:        schema.TypeString,
													Description: "Rule value",
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

func dataSourceBaiduCloudAppBLBListenersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	appblbService := APPBLBService{client}

	blbId := ""
	if v, ok := d.GetOk("blb_id"); ok && v.(string) != "" {
		blbId = v.(string)
	}

	protocol := ""
	if v, ok := d.GetOk("protocol"); ok && v.(string) != "" {
		protocol = v.(string)
	}

	listenerPort := 0
	if v, ok := d.GetOk("listener_port"); ok {
		listenerPort = v.(int)
	}

	action := "Query APPBLB " + blbId + "_" + protocol
	listeners, err := appblbService.ListAllListeners(blbId, protocol, listenerPort)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_appblb_listeners", action, BCESDKGoERROR)
	}
	if err := d.Set("listeners", listeners); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_appblb_listeners", action, BCESDKGoERROR)
	}

	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), listeners); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_appblb_listeners", action, BCESDKGoERROR)
		}
	}

	return nil
}
