/*
Provide a resource to create a VPN conn.

Example Usage

```hcl
resource "baiducloud_vpn_conn" "default" {
  vpn_id = baiducloud_vpn_gateway.default.id
  secret_key = "ddd22@www"
  local_subnets = ["192.168.0.0/20"]
  remote_ip = "11.11.11.133"
  remote_subnets = ["192.168.100.0/24"]
  description = "just for test"
  vpn_conn_name = "vpn-conn"
  ike_config {
    ike_version = "v1"
    ike_mode = "main"
    ike_enc_alg = "aes"
    ike_auth_alg = "sha1"
    ike_pfs = "group2"
    ike_life_time = 300
  }
  ipsec_config {
    ipsec_enc_alg = "aes"
    ipsec_auth_alg = "sha1"
    ipsec_pfs = "group2"
    ipsec_life_time = 200
  }
}
```

Import

VPN conn can be imported, e.g.

```hcl
$ terraform import baiducloud_vpn_conn.default vpn_conn_id
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/vpn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
	"log"
	"time"
)

func resourceBaiduCloudVpnConn() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudVpnConnCreate,
		Read:   resourceBaiduCloudVpnConnRead,
		Update: resourceBaiduCloudVpnConnUpdate,
		Delete: resourceBaiduCloudVpnConnDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"vpn_id": {
				Type:        schema.TypeString,
				Description: "VPN id which vpn conn belong to.",
				Required:    true,
			},
			"secret_key": {
				Type:        schema.TypeString,
				Description: "Shared secret key, 8 to 17 characters, English, numbers and symbols must exist at the same time, the symbols are limited to !@#$%^*()_.",
				Required:    true,
			},
			"local_subnets": {
				Type:        schema.TypeList,
				Description: "Local network cidr list.",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"remote_ip": {
				Type:        schema.TypeString,
				Description: "Public IP of the peer VPN gateway.",
				Required:    true,
			},
			"local_ip": {
				Type:        schema.TypeString,
				Description: "Public IP of the VPN gateway.",
				Computed:    true,
			},
			"remote_subnets": {
				Type:        schema.TypeList,
				Description: "Peer network cidr list.",
				Required:    true,
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
				Optional:    true,
				Computed:    true,
			},
			"ike_config": {
				Type:        schema.TypeSet,
				Description: "IKE config.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ike_version": {
							Type:         schema.TypeString,
							Description:  "Version of IKE.",
							ValidateFunc: validation.StringInSlice([]string{"v1", "v2"}, false),
							Required:     true,
						},
						"ike_mode": {
							Type:         schema.TypeString,
							Description:  "Negotiation mode.",
							ValidateFunc: validation.StringInSlice([]string{"main", "aggressive"}, false),
							Required:     true,
						},
						"ike_enc_alg": {
							Type:         schema.TypeString,
							Description:  "IKE Encryption Algorithm.",
							ValidateFunc: validation.StringInSlice([]string{"aes", "aes192", "aes256", "3des"}, false),
							Required:     true,
						},
						"ike_auth_alg": {
							Type:         schema.TypeString,
							Description:  "IKE Authenticate Algorithm",
							ValidateFunc: validation.StringInSlice([]string{"sha1", "md5", "sha2_256", "sha2_384", "sha2_512"}, false),
							Required:     true,
						},
						"ike_pfs": {
							Type:         schema.TypeString,
							Description:  "Diffie-Hellman key exchange algorithm.",
							ValidateFunc: validation.StringInSlice([]string{"group2", "group5", "group14", "group24"}, false),
							Required:     true,
						},
						"ike_life_time": {
							Type:        schema.TypeInt,
							Description: "IKE life time.",
							Required:    true,
						},
					},
				},
			},
			"ipsec_config": {
				Type:        schema.TypeSet,
				Description: "Ipsec config details.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipsec_enc_alg": {
							Type:         schema.TypeString,
							Description:  "Ipsec Encryption Algorithm.",
							ValidateFunc: validation.StringInSlice([]string{"aes", "aes192", "aes256", "3des"}, false),
							Required:     true,
						},
						"ipsec_auth_alg": {
							Type:         schema.TypeString,
							Description:  "Ipsec Authenticate Algorithm.",
							ValidateFunc: validation.StringInSlice([]string{"sha1", "md5", "sha2_256", "sha2_384", "sha2_512"}, false),
							Required:     true,
						},
						"ipsec_pfs": {
							Type:         schema.TypeString,
							Description:  "Diffie-Hellman key exchange algorithm.",
							ValidateFunc: validation.StringInSlice([]string{"group2", "group5", "group14", "group24", "disabled"}, false),
							Required:     true,
						},
						"ipsec_life_time": {
							Type:        schema.TypeInt,
							Description: "Ipsec life time.",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func resourceBaiduCloudVpnConnCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	createVpnGatewayConnArgs, err := buildBaiduCloudVpnGatewayConnCreateArgs(d, meta)

	action := "Create VPN Gateway Conn" + createVpnGatewayConnArgs.VpnConnName
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_vpn_conn", action, BCESDKGoERROR)
	}
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithVPNClient(func(vpnClient *vpn.Client) (i interface{}, e error) {
			return vpnClient.CreateVpnConn(createVpnGatewayConnArgs)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw)
		result, _ := raw.(*vpn.CreateVpnConnResult)
		d.SetId(result.VpnConnId)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_vpn_conn", action, BCESDKGoERROR)
	}
	return resourceBaiduCloudVpnConnRead(d, meta)
}

func resourceBaiduCloudVpnConnRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	action := "Query VPN Gateway Conn"
	var vpnId string
	if v, ok := d.GetOk("vpn_id"); ok {
		vpnId = v.(string)
	}
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithVPNClient(func(vpnClient *vpn.Client) (i interface{}, e error) {
			return vpnClient.ListVpnConn(vpnId)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw)
		connList, _ := raw.(*vpn.ListVpnConnResult)
		for _, conn := range connList.VpnConns {
			if conn.VpnConnId == d.Id() {
				resultToSchema(d, conn)
			}
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_vpn_conn", action, BCESDKGoERROR)
	}
	return nil
}

func resourceBaiduCloudVpnConnUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	updateVPNConnArgs, err := buildBaiduCloudVpnGatewayConnUpdateArgs(d, meta)
	action := "Update VPN Gateway Conn" + updateVPNConnArgs.UpdateVpnconn.VpnConnName
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_vpn_conn", action, BCESDKGoERROR)
	}
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		_, err := client.WithVPNClient(func(vpnClient *vpn.Client) (i interface{}, e error) {
			log.Print(updateVPNConnArgs.UpdateVpnconn.RemoteIp)
			log.Print(updateVPNConnArgs.UpdateVpnconn.RemoteSubnets)
			log.Print(updateVPNConnArgs.UpdateVpnconn.LocalSubnets)
			return nil, vpnClient.UpdateVpnConn(updateVPNConnArgs)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_vpn_conn", action, BCESDKGoERROR)
	}
	return resourceBaiduCloudVpnConnRead(d, meta)
}

func resourceBaiduCloudVpnConnDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	action := "Delete VPN Gateway Conn"

	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err := client.WithVPNClient(func(vpnClient *vpn.Client) (i interface{}, e error) {
			return nil, vpnClient.DeleteVpnConn(d.Id(), buildClientToken())
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_vpn_conn", action, BCESDKGoERROR)
	}
	return nil
}

func buildBaiduCloudVpnGatewayConnCreateArgs(d *schema.ResourceData, meta interface{}) (*vpn.CreateVpnConnArgs, error) {
	res := &vpn.CreateVpnConnArgs{
		ClientToken: buildClientToken(),
	}
	if vpnID, ok := d.GetOk("vpn_id"); ok {
		res.VpnId = vpnID.(string)
	}
	if desc, ok := d.GetOk("description"); ok {
		res.Description = desc.(string)
	}
	if name, ok := d.GetOk("vpn_conn_name"); ok {
		res.VpnConnName = name.(string)
	}
	if secKey, ok := d.GetOk("secret_key"); ok {
		res.SecretKey = secKey.(string)
	}
	if remoteSubnets, ok := d.GetOk("remote_subnets"); ok {
		res.RemoteSubnets = interfaceToStringSlice(remoteSubnets)
	}
	if localSubnets, ok := d.GetOk("local_subnets"); ok {
		res.LocalSubnets = interfaceToStringSlice(localSubnets)
	}
	if remoteIp, ok := d.GetOk("remote_ip"); ok {
		res.RemoteIp = remoteIp.(string)
	}
	if r, ok := d.GetOk("ike_config"); ok {
		createIkeConfig := &vpn.CreateIkeConfig{}
		ikeConfig := r.(*schema.Set).List()[0].(map[string]interface{})
		if ikeVersion, ok := ikeConfig["ike_version"]; ok {
			createIkeConfig.IkeVersion = ikeVersion.(string)
		}
		if ikeMode, ok := ikeConfig["ike_mode"]; ok {
			createIkeConfig.IkeMode = ikeMode.(string)
		}
		if ikeEncAlg, ok := ikeConfig["ike_enc_alg"]; ok {
			createIkeConfig.IkeEncAlg = ikeEncAlg.(string)
		}
		if ikePfs, ok := ikeConfig["ike_pfs"]; ok {
			createIkeConfig.IkePfs = ikePfs.(string)
		}
		if ikeAuthAlg, ok := ikeConfig["ike_auth_alg"]; ok {
			createIkeConfig.IkeAuthAlg = ikeAuthAlg.(string)
		}
		if ikeLifeTime, ok := ikeConfig["ike_life_time"]; ok {
			createIkeConfig.IkeLifeTime, _ = ikeLifeTime.(int)
		}
		res.CreateIkeConfig = createIkeConfig
	}
	if r, ok := d.GetOk("ipsec_config"); ok {
		ipsecConfig := r.(*schema.Set).List()[0].(map[string]interface{})
		createIpsecConfig := &vpn.CreateIpsecConfig{}
		// ipsec_enc_alg
		if ipsecEncAlg, ok := ipsecConfig["ipsec_enc_alg"]; ok {
			createIpsecConfig.IpsecEncAlg = ipsecEncAlg.(string)
		}
		// ipsec_auth_alg
		if ipsecAuthAlg, ok := ipsecConfig["ipsec_auth_alg"]; ok {
			createIpsecConfig.IpsecAuthAlg = ipsecAuthAlg.(string)
		}
		// ipsec_pfs
		if ipsecPfs, ok := ipsecConfig["ipsec_pfs"]; ok {
			createIpsecConfig.IpsecPfs = ipsecPfs.(string)
		}
		// ipsec_life_time
		log.Print(ipsecConfig["ipsec_life_time"])
		if ipsecLifeTime, ok := ipsecConfig["ipsec_life_time"]; ok {
			createIpsecConfig.IpsecLifetime, _ = ipsecLifeTime.(int)
		}
		res.CreateIpsecConfig = createIpsecConfig
	}
	return res, nil
}

func buildBaiduCloudVpnGatewayConnUpdateArgs(d *schema.ResourceData, meta interface{}) (*vpn.UpdateVpnConnArgs, error) {
	res := &vpn.UpdateVpnConnArgs{}
	res.VpnConnId = d.Id()
	var err error
	res.UpdateVpnconn, err = buildBaiduCloudVpnGatewayConnCreateArgs(d, meta)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func resultToSchema(d *schema.ResourceData, conn vpn.VpnConn) {
	d.Set("vpn_conn_name", conn.VpnConnName)
	d.Set("status", conn.Status)
	d.Set("secret_key", conn.SecretKey)
	d.Set("remote_subnets", conn.RemoteSubnets)
	d.Set("remote_ip", conn.RemoteIp)
	d.Set("local_subnets", conn.LocalSubnets)
	d.Set("local_ip", conn.LocalIp)
	d.Set("description", conn.Description)
	ikeMap := make(map[string]interface{}, 0)
	ikeMap["ike_pfs"] = conn.IkeConfig.IkePfs
	ikeMap["ike_auth_alg"] = conn.IkeConfig.IkeAuthAlg
	ikeMap["ike_enc_alg"] = conn.IkeConfig.IkeEncAlg
	ikeLifeTime := conn.IkeConfig.IkeLifeTime
	ikeMap["ike_life_time"] = ikeLifeTime[0 : len(ikeLifeTime)-1]
	ikeMap["ike_mode"] = conn.IkeConfig.IkeMode
	ikeMap["ike_version"] = conn.IkeConfig.IkeVersion
	d.Set("ike_config", append(make([]map[string]interface{}, 0), ikeMap))
	ipsecMap := make(map[string]interface{})
	ipsecMap["ipsec_auth_alg"] = conn.IpsecConfig.IpsecAuthAlg
	ipsecMap["ipsec_enc_alg"] = conn.IpsecConfig.IpsecEncAlg
	ipsecLifetime := conn.IpsecConfig.IpsecLifetime
	ipsecMap["ipsec_life_time"] = ipsecLifetime[0 : len(ipsecLifetime)-1]
	ipsecMap["ipsec_pfs"] = conn.IpsecConfig.IpsecPfs
	d.Set("ipsec_config", append(make([]map[string]interface{}, 0), ipsecMap))
}

func interfaceToStringSlice(remoteSubnets interface{}) []string {
	var subnets []string
	for _, val := range remoteSubnets.([]interface{}) {
		subnets = append(subnets, val.(string))
	}
	return subnets
}
