/*
Use this data source to query vpn conn list.

Example Usage

```hcl
data "baiducloud_vpn_conns" "default" {
    vpn_id = "vpn-xxxxxxx"
}

output "conns" {
  value = "${data.baiducloud_vpn_conns.default.conns}"
}
```
*/
package baiducloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudVPNConns() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudVPNConnsRead,

		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:        schema.TypeString,
				Description: "Output file for saving result.",
				Optional:    true,
				ForceNew:    true,
			},
			"vpn_id": {
				Type:        schema.TypeString,
				Description: "VPN id which vpn conn belong to.",
				Required:    true,
			},
			"filter": dataSourceFiltersSchema(),
			"vpn_conns": {
				Type:        schema.TypeList,
				Description: "Result of VPN conns.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpn_id": {
							Type:        schema.TypeString,
							Description: "VPN id which vpn conn belong to.",
							Computed:    true,
						},
						"vpn_conn_id": {
							Type:        schema.TypeString,
							Description: "ID of the VPN conn.",
							Computed:    true,
						},
						"secret_key": {
							Type:        schema.TypeString,
							Description: "Shared secret key, 8 to 17 characters, English, numbers and symbols must exist at the same time, the symbols are limited to !@#$%^*()_.",
							Computed:    true,
						},
						"local_subnets": {
							Type:        schema.TypeList,
							Description: "Local network cidr list.",
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"remote_ip": {
							Type:        schema.TypeString,
							Description: "Public IP of the peer VPN gateway.",
							Computed:    true,
						},
						"local_ip": {
							Type:        schema.TypeString,
							Description: "Public IP of the VPN gateway.",
							Computed:    true,
						},
						"remote_subnets": {
							Type:        schema.TypeList,
							Description: "Peer network cidr list.",
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"description": {
							Type:        schema.TypeString,
							Description: "Description of VPN conn.",
							Optional:    true,
						},
						"vpn_conn_name": {
							Type:        schema.TypeString,
							Description: "Name of vpn conn.",
							Computed:    true,
						},
						"created_time": {
							Type:        schema.TypeString,
							Description: "Create time of VPN conn.",
							Computed:    true,
						},
						"health_status": {
							Type:        schema.TypeString,
							Description: "Health status of the vpn conn.",
							Computed:    true,
						},
						"status": {
							Type:        schema.TypeString,
							Description: "Status of the vpn conn.",
							Computed:    true,
						},
						"ike_config": {
							Type:        schema.TypeMap,
							Description: "IKE config detail.",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ike_version": {
										Type:        schema.TypeString,
										Description: "Version of IKE.",
										Optional:    true,
									},
									"ike_mode": {
										Type:        schema.TypeString,
										Description: "Negotiation mode.",
										Optional:    true,
									},
									"ike_enc_alg": {
										Type:        schema.TypeString,
										Description: "IKE Encryption Algorithm.",
										Optional:    true,
									},
									"ike_auth_alg": {
										Type:        schema.TypeString,
										Description: "IKE Authenticate Algorithm",
										Optional:    true,
									},
									"ike_pfs": {
										Type:        schema.TypeString,
										Description: "Diffie-Hellman key exchange algorithm.",
										Optional:    true,
									},
									"ike_life_time": {
										Type:        schema.TypeString,
										Description: "IKE life time.",
										Optional:    true,
									},
								},
							},
						},
						"ipsec_config": {
							Type:        schema.TypeMap,
							Description: "Ipsec config of vpn conn.",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipsec_enc_alg": {
										Type:        schema.TypeString,
										Description: "Ipsec Encryption Algorithm.",
										Optional:    true,
									},
									"ipsec_auth_alg": {
										Type:        schema.TypeString,
										Description: "Ipsec Authenticate Algorithm.",
										Optional:    true,
									},
									"ipsec_pfs": {
										Type:        schema.TypeString,
										Description: "Diffie-Hellman key exchange algorithm.",
										Optional:    true,
									},
									"ipsec_life_time": {
										Type:        schema.TypeString,
										Description: "Ipsec life time.",
										Optional:    true,
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

func dataSourceBaiduCloudVPNConnsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	service := VpnService{client: client}
	action := "Query all vpn conns"
	var vpnId string
	if v, ok := d.GetOk("vpn_id"); ok {
		vpnId = v.(string)
	}
	conns, err := service.ListVpnConns(vpnId)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_vpn_conns", action, BCESDKGoERROR)
	}
	vpnConnsResult := make([]map[string]interface{}, 0)
	for _, conn := range conns {
		vpnConnsResult = append(vpnConnsResult, service.connToMap(conn))
	}
	d.SetId(resource.UniqueId())
	if err := d.Set("vpn_conns", vpnConnsResult); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_vpn_conns", action, BCESDKGoERROR)
	}
	return nil
}
