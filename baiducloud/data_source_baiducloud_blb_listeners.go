/*
Use this data source to query BLB Listener list.

Example Usage

```hcl
data "baiducloud_blb_listeners" "default" {
 blb_id = "lb-0d29axxx6"
}

output "listeners" {
 value = "${data.baiducloud_blb_listeners.default.listeners}"
}
```
*/
package baiducloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudBLBListeners() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudBLBListenersRead,

		Schema: map[string]*schema.Schema{
			"protocol": {
				Type:        schema.TypeString,
				Description: "Protocol of the Listener to be queried",
				Optional:    true,
				ForceNew:    true,
			},
			"blb_id": {
				Type:        schema.TypeString,
				Description: "ID of the LoadBalance instance",
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
			"filter": dataSourceFiltersSchema(),

			// Attributes used for result
			"listeners": {
				Type:        schema.TypeList,
				Description: "A list of LoadBalance Listener",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"listener_port": {
							Type:        schema.TypeInt,
							Description: "Listener bind port",
							Computed:    true,
						},
						"backend_port": {
							Type:        schema.TypeInt,
							Description: "backend port, range from 1-65535",
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
						"udp_session_timeout": {
							Type:        schema.TypeInt,
							Description: "UDP Listener connection session timeout time(second), default 900, support 10-4000",
							Computed:    true,
						},
						"health_check_timeout_in_second": {
							Type:        schema.TypeInt,
							Description: "health check timeout in second",
							Computed:    true,
						},
						"health_check_interval": {
							Type:        schema.TypeInt,
							Description: "health check interval",
							Computed:    true,
						},
						"healthy_threshold": {
							Type:        schema.TypeInt,
							Description: "healthy threshold",
							Computed:    true,
						},
						"unhealthy_threshold": {
							Type:        schema.TypeInt,
							Description: "unhealthy threshold",
							Computed:    true,
						},
						"get_blb_ip": {
							Type:        schema.TypeBool,
							Description: "get blb ip or not",
							Computed:    true,
						},
						// UCP
						"health_check_string": {
							Type:        schema.TypeString,
							Description: "health check string, This parameter is mandatory when the listening protocol is UDP",
							Computed:    true,
						},
						// SSL HTTPS
						"applied_ciphers": {
							Type:        schema.TypeString,
							Description: "applied ciphers",
							Computed:    true,
						},
						// http https
						"keep_session_duration": {
							Type:        schema.TypeInt,
							Description: "keep session duration",
							Computed:    true,
						},
						// http https
						"health_check_type": {
							Type:        schema.TypeString,
							Description: "health check type",
							Computed:    true,
						},
						// http https
						"health_check_port": {
							Type:        schema.TypeInt,
							Description: "health check port",
							Computed:    true,
						},
						// http https
						"health_check_uri": {
							Type:        schema.TypeString,
							Description: "health check uri",
							Computed:    true,
						},
						// http https
						"health_check_normal_status": {
							Type:        schema.TypeString,
							Description: "health check normal status",
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
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudBLBListenersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	blbService := BLBService{client}

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

	action := "Query BLB " + blbId + "_" + protocol
	listeners, err := blbService.ListAllListeners(blbId, protocol, listenerPort)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_blb_listeners", action, BCESDKGoERROR)
	}

	FilterDataSourceResult(d, &listeners)

	if err := d.Set("listeners", listeners); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_blb_listeners", action, BCESDKGoERROR)
	}

	d.SetId(resource.UniqueId())

	if v, ok := d.GetOk("output_file"); ok && v.(string) != "" {
		if err := writeToFile(v.(string), listeners); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_blb_listeners", action, BCESDKGoERROR)
		}
	}

	return nil
}
