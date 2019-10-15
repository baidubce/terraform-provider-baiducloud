/*
Provide a resource to create an APPBLB Listener.

Example Usage

```hcl
[TCP/UDP] Listener
resource "baiducloud_appblb_listener" "default" {
  blb_id               = "lb-0d29a3f6"
  listener_port        = 124
  protocol             = "TCP"
  scheduler            = "LeastConnection"
}

[HTTP] Listener
resource "baiducloud_appblb_listener" "default" {
  blb_id        = "lb-0d29a3f6"
  listener_port = 129
  protocol      = "HTTP"
  scheduler     = "RoundRobin"
  keep_session  = true

  policies {
    description         = "acceptance test"
    app_server_group_id = "sg-11bd8054"
    backend_port        = 70
    priority            = 50

    rule_list {
      key   = "host"
      value = "baidu.com"
    }
  }
}

[HTTPS] Listener
resource "baiducloud_appblb_listener" "default" {
  blb_id               = "lb-0d29a3f6"
  listener_port        = 130
  protocol             = "HTTPS"
  scheduler            = "LeastConnection"
  keep_session         = true
  cert_ids             = ["cert-xvysj80uif1y"]
  encryption_protocols = ["sslv3", "tlsv10", "tlsv11"]
  encryption_type      = "userDefind"
}

[SSL] Listener
resource "baiducloud_appblb_listener" "default" {
  blb_id               = "lb-0d29a3f6"
  listener_port        = 131
  protocol             = "SSL"
  scheduler            = "LeastConnection"
  cert_ids             = ["cert-xvysj80uif1y"]
  encryption_protocols = ["sslv3", "tlsv10", "tlsv11"]
  encryption_type      = "userDefind"
}
```
*/
package baiducloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/appblb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudAppBlbListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudAppBlbListenerCreate,
		Read:   resourceBaiduCloudAppBlbListenerRead,
		Update: resourceBaiduCloudAppBlbListenerUpdate,
		Delete: resourceBaiduCloudAppBlbListenerDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"blb_id": {
				Type:        schema.TypeString,
				Description: "ID of the Application LoadBalance instance",
				Required:    true,
				ForceNew:    true,
			},
			"listener_port": {
				Type:         schema.TypeInt,
				Description:  "Listening port, range from 1-65535",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validatePort(),
			},
			"protocol": {
				Type:         schema.TypeString,
				Description:  "Listening protocol, support TCP/UDP/HTTP/HTTPS/SSL",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{TCP, UDP, HTTP, HTTPS, SSL}, false),
			},
			"scheduler": {
				Type:         schema.TypeString,
				Description:  "Load balancing algorithm, support RoundRobin/LeastConnection/Hash, if protocol is HTTP/HTTPS, only support RoundRobin/LeastConnection",
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"RoundRobin", "LeastConnection", "Hash"}, false),
			},
			"tcp_session_timeout": {
				Type:         schema.TypeInt,
				Description:  "TCP Listener connection session timeout time(second), default 900, support 10-4000",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(10, 4000),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.Get("protocol").(string); ok {
						return v != TCP
					}
					return true
				},
			},
			// http & https
			"keep_session": {
				Type:             schema.TypeBool,
				Description:      "KeepSession or not",
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: appBlbProtocolTCPUDPSSLSuppressFunc,
			},
			//http & https
			"keep_session_type": {
				Type:             schema.TypeString,
				Description:      "KeepSessionType option, support insert/rewrite, default insert",
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.StringInSlice([]string{"insert", "rewrite"}, false),
				DiffSuppressFunc: appBlbProtocolTCPUDPSSLSuppressFunc,
			},
			// http & https
			"keep_session_timeout": {
				Type:             schema.TypeInt,
				Description:      "KeepSession Cookie timeout time(second), support in [1, 15552000], default 3600s",
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.IntBetween(1, 15552000),
				DiffSuppressFunc: appBlbProtocolTCPUDPSSLSuppressFunc,
			},
			// http & https
			"keep_session_cookie_name": {
				Type:        schema.TypeString,
				Description: "CookieName which need to covered, useful when keep_session_type is rewrite",
				Optional:    true,
				Computed:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					protocolCheck := appBlbProtocolTCPUDPSSLSuppressFunc(k, old, new, d)
					if protocolCheck {
						return true
					}

					if v, ok := d.GetOk("keep_session"); !ok || !(v.(bool)) {
						return true
					}

					if v, ok := d.GetOk("keep_session_type"); ok {
						return v.(string) != "rewrite"
					}

					return true
				},
			},
			// http & https
			"x_forwarded_for": {
				Type:             schema.TypeBool,
				Description:      "Listener xForwardedFor, determine get client real ip or not, default false",
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: appBlbProtocolTCPUDPSSLSuppressFunc,
			},
			// http & https
			"server_timeout": {
				Type:             schema.TypeInt,
				Description:      "Backend server maximum timeout time, only support in [1, 3600] second, default 30s",
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: appBlbProtocolTCPUDPSSLSuppressFunc,
				ValidateFunc:     validation.IntBetween(1, 3600),
			},
			// http
			"redirect_port": {
				Type:        schema.TypeInt,
				Description: "Redirect HTTP request to HTTPS Listener, HTTPS Listener port set by this parameter",
				Optional:    true,
				Computed:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("protocol").(string) != HTTP
				},
			},
			// https && ssl
			"cert_ids": {
				Type:        schema.TypeSet,
				Description: "Listener bind certifications",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: appBlbProtocolTCPUDPHTTPSuppressFunc,
			},
			// https && ssl
			"ie6_compatible": {
				Type:             schema.TypeBool,
				Description:      "Listener support ie6 option, default true",
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: appBlbProtocolTCPUDPHTTPSuppressFunc,
			},
			// https && ssl
			"encryption_type": {
				Type:             schema.TypeString,
				Description:      "Listener encryption option, support [compatibleIE, incompatibleIE, userDefind]",
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.StringInSlice([]string{"compatibleIE", "incompatibleIE", "userDefind"}, false),
				DiffSuppressFunc: appBlbProtocolTCPUDPHTTPSuppressFunc,
			},
			// https && ssl
			"encryption_protocols": {
				Type:        schema.TypeSet,
				Description: "Listener encryption protocol, only useful when encryption_type is userDefind, support [sslv3, tlsv10, tlsv11, tlsv12]",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"sslv3", "tlsv10", "tlsv11", "tlsv12"}, false),
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("encryption_type"); ok {
						return v.(string) != "userDefind"
					}

					return true
				},
			},
			// https && ssl
			"dual_auth": {
				Type:             schema.TypeBool,
				Description:      "Listener open dual authorization or not, default false",
				Optional:         true,
				Computed:         false,
				DiffSuppressFunc: appBlbProtocolTCPUDPHTTPSuppressFunc,
			},
			// https && ssl
			"client_cert_ids": {
				Type:        schema.TypeSet,
				Description: "Listener import cert list, only useful when dual_auth is true",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: appBlbProtocolTCPUDPHTTPSuppressFunc,
			},
			"policies": {
				Type:        schema.TypeSet,
				Description: "Listener's policy",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "Policy's id",
							Computed:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "Policy's description",
							Optional:    true,
						},
						"port_type": {
							Type:        schema.TypeString,
							Description: "Policy bind port protocol type",
							Computed:    true,
						},
						"app_server_group_id": {
							Type:        schema.TypeString,
							Description: "Policy bind server group id",
							Required:    true,
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
							Type:         schema.TypeInt,
							Description:  "Backend port",
							Required:     true,
							ValidateFunc: validatePort(),
						},
						"priority": {
							Type:         schema.TypeInt,
							Description:  "Policy priority, support in [1, 32768]",
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 32768),
						},
						"rule_list": {
							Type:        schema.TypeSet,
							Description: "Policy rule list",
							Optional:    true,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:         schema.TypeString,
										Description:  "Rule key",
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"*", "host", "uri"}, false),
									},
									"value": {
										Type:        schema.TypeString,
										Description: "Rule value",
										Required:    true,
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

func resourceBaiduCloudAppBlbListenerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	blbId := d.Get("blb_id").(string)
	protocol := d.Get("protocol").(string)
	listenerPort := d.Get("listener_port").(int)
	action := fmt.Sprintf("Create APPBLB %s Listener [%s:%d]", blbId, protocol, listenerPort)

	listenerArgs, err := buildBaiduCloudCreateAppBlbListenerArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}

	var policyArgs *appblb.CreatePolicysArgs
	if v, ok := d.GetOk("policies"); ok {
		policyArgs, err = buildBaiduCloudCreatePolicyArgs(listenerPort, protocol, v.(*schema.Set).List())
		if err != nil {
			return WrapError(err)
		}
	}

	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithAppBLBClient(func(client *appblb.Client) (i interface{}, e error) {
			switch protocol {
			case TCP:
				return blbId, client.CreateAppTCPListener(blbId, listenerArgs.(*appblb.CreateAppTCPListenerArgs))
			case UDP:
				return blbId, client.CreateAppUDPListener(blbId, listenerArgs.(*appblb.CreateAppUDPListenerArgs))
			case HTTP:
				return blbId, client.CreateAppHTTPListener(blbId, listenerArgs.(*appblb.CreateAppHTTPListenerArgs))
			case HTTPS:
				return blbId, client.CreateAppHTTPSListener(blbId, listenerArgs.(*appblb.CreateAppHTTPSListenerArgs))
			case SSL:
				return blbId, client.CreateAppSSLListener(blbId, listenerArgs.(*appblb.CreateAppSSLListenerArgs))
			default:
				// never run here
				return blbId, fmt.Errorf("unsupport protocol")
			}
		})

		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		addDebug(action, raw)
		d.SetId(strconv.Itoa(listenerPort))

		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_appblb_listener", action, BCESDKGoERROR)
	}

	if policyArgs != nil {
		_, err = client.WithAppBLBClient(func(client *appblb.Client) (i interface{}, e error) {
			return nil, client.CreatePolicys(blbId, policyArgs)
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_appblb_listener", action, BCESDKGoERROR)
		}
	}

	return resourceBaiduCloudAppBlbListenerRead(d, meta)
}

func resourceBaiduCloudAppBlbListenerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	appblbService := APPBLBService{client}

	blbId := d.Get("blb_id").(string)
	protocol := d.Get("protocol").(string)
	listenerPort := d.Get("listener_port").(int)
	action := fmt.Sprintf("Query APPBLB %s Listener [%s:%d]", blbId, protocol, listenerPort)

	raw, err := appblbService.DescribeListener(blbId, protocol, listenerPort)
	if err != nil {
		d.SetId("")
		return WrapError(err)
	}
	addDebug(action, raw)

	switch protocol {
	case HTTP:
		listenerMeta := raw.(*appblb.AppHTTPListenerModel)
		d.Set("scheduler", listenerMeta.Scheduler)
		d.Set("keep_session", listenerMeta.KeepSession)
		d.Set("keep_session_type", listenerMeta.KeepSessionType)
		d.Set("keep_session_timeout", listenerMeta.KeepSessionTimeout)
		d.Set("keep_session_cookie_name", listenerMeta.KeepSessionCookieName)
		d.Set("x_forwarded_for", listenerMeta.XForwardedFor)
		d.Set("server_timeout", listenerMeta.ServerTimeout)
		d.Set("redirect_port", listenerMeta.RedirectPort)
		d.Set("listener_port", listenerMeta.ListenerPort)
	case HTTPS:
		listenerMeta := raw.(*appblb.AppHTTPSListenerModel)
		d.Set("scheduler", listenerMeta.Scheduler)
		d.Set("keep_session", listenerMeta.KeepSession)
		d.Set("keep_session_type", listenerMeta.KeepSessionType)
		d.Set("keep_session_timeout", listenerMeta.KeepSessionTimeout)
		d.Set("keep_session_cookie_name", listenerMeta.KeepSessionCookieName)
		d.Set("x_forwarded_for", listenerMeta.XForwardedFor)
		d.Set("server_timeout", listenerMeta.ServerTimeout)
		d.Set("cert_ids", listenerMeta.CertIds)
		d.Set("encryption_type", listenerMeta.EncryptionType)
		d.Set("encryption_protocols", listenerMeta.EncryptionProtocols)
		d.Set("dual_auth", listenerMeta.DualAuth)
		d.Set("client_cert_ids", listenerMeta.ClientCertIds)
		d.Set("listener_port", listenerMeta.ListenerPort)
	case SSL:
		listenerMeta := raw.(*appblb.AppSSLListenerModel)
		d.Set("scheduler", listenerMeta.Scheduler)
		d.Set("cert_ids", listenerMeta.CertIds)
		d.Set("encryption_type", listenerMeta.EncryptionType)
		d.Set("encryption_protocols", listenerMeta.EncryptionProtocols)
		d.Set("dual_auth", listenerMeta.DualAuth)
		d.Set("client_cert_ids", listenerMeta.ClientCertIds)
		d.Set("listener_port", listenerMeta.ListenerPort)
	case TCP:
		listenerMeta := raw.(*appblb.AppTCPListenerModel)
		d.Set("scheduler", listenerMeta.Scheduler)
		d.Set("tcp_session_timeout", listenerMeta.TcpSessionTimeout)
	case UDP:
		listenerMeta := raw.(*appblb.AppUDPListenerModel)
		d.Set("scheduler", listenerMeta.Scheduler)
	default:
		return WrapError(fmt.Errorf("unsupport listener type"))
	}
	d.SetId(strconv.Itoa(listenerPort))

	policies, err := appblbService.DescribePolicys(blbId, protocol, listenerPort)
	if err != nil {
		return WrapError(err)
	}

	if err := d.Set("policies", appblbService.FlattenAppPolicysToMap(policies)); err != nil {
		return WrapError(err)
	}

	return nil
}

func resourceBaiduCloudAppBlbListenerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	blbId := d.Get("blb_id").(string)
	protocol := d.Get("protocol").(string)
	listenerPort := d.Get("listener_port").(int)
	action := fmt.Sprintf("Update APPBLB %s Listener [%s:%d]", blbId, protocol, listenerPort)

	update, args, err := buildBaiduCloudUpdateAppBlbListenerArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}
	if update {
		_, err := client.WithAppBLBClient(func(client *appblb.Client) (i interface{}, e error) {
			switch protocol {
			case TCP:
				return nil, client.UpdateAppTCPListener(blbId, args.(*appblb.UpdateAppTCPListenerArgs))
			case UDP:
				return nil, client.UpdateAppUDPListener(blbId, args.(*appblb.UpdateAppUDPListenerArgs))
			case HTTP:
				return nil, client.UpdateAppHTTPListener(blbId, args.(*appblb.UpdateAppHTTPListenerArgs))
			case HTTPS:
				return nil, client.UpdateAppHTTPSListener(blbId, args.(*appblb.UpdateAppHTTPSListenerArgs))
			case SSL:
				return nil, client.UpdateAppSSLListener(blbId, args.(*appblb.UpdateAppSSLListenerArgs))
			default:
				return nil, fmt.Errorf("unsupport listener type")
			}
		})
		addDebug(action, nil)

		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_appblb_listener", action, BCESDKGoERROR)
		}
	}

	if d.HasChange("policies") {
		o, n := d.GetChange("policies")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		add := ns.Difference(os).List()
		remove := os.Difference(ns).List()

		createArgs, err := buildBaiduCloudCreatePolicyArgs(listenerPort, protocol, add)
		if err != nil {
			return WrapError(err)
		}

		if len(remove) > 0 {
			deleteArgs := &appblb.DeletePolicysArgs{
				Port: uint16(listenerPort),
			}

			for _, p := range remove {
				id := p.(map[string]interface{})["id"].(string)
				deleteArgs.PolicyIdList = append(deleteArgs.PolicyIdList, id)
			}

			_, err := client.WithAppBLBClient(func(client *appblb.Client) (i interface{}, e error) {
				return nil, client.DeletePolicys(blbId, deleteArgs)
			})

			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, "baiducloud_appblb_listener", action, BCESDKGoERROR)
			}
		}

		if len(add) > 0 {
			_, err = client.WithAppBLBClient(func(client *appblb.Client) (i interface{}, e error) {
				return nil, client.CreatePolicys(blbId, createArgs)
			})

			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, "baiducloud_appblb_listener", action, BCESDKGoERROR)
			}
		}
	}

	return resourceBaiduCloudAppBlbListenerRead(d, meta)
}

func resourceBaiduCloudAppBlbListenerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	blbId := d.Get("blb_id").(string)
	protocol := d.Get("protocol").(string)
	listenerPort := d.Get("listener_port").(int)
	action := fmt.Sprintf("Delete APPBLB %s Listener [%s:%d]", blbId, protocol, listenerPort)

	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err := client.WithAppBLBClient(func(client *appblb.Client) (i interface{}, e error) {
			return blbId, client.DeleteAppListeners(blbId, &appblb.DeleteAppListenersArgs{
				PortList:    []uint16{uint16(listenerPort)},
				ClientToken: buildClientToken(),
			})
		})
		addDebug(action, blbId)

		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		if IsExceptedErrors(err, ObjectNotFound) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_appblb_listener", action, BCESDKGoERROR)
	}

	return nil
}

func buildBaiduCloudCreateAppBlbListenerArgs(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	protocol := d.Get("protocol").(string)

	switch protocol {
	case TCP:
		return &appblb.CreateAppTCPListenerArgs{
			TcpSessionTimeout: d.Get("tcp_session_timeout").(int),
			ListenerPort:      uint16(d.Get("listener_port").(int)),
			Scheduler:         d.Get("scheduler").(string),
			ClientToken:       buildClientToken(),
		}, nil
	case UDP:
		return &appblb.CreateAppUDPListenerArgs{
			ListenerPort: uint16(d.Get("listener_port").(int)),
			Scheduler:    d.Get("scheduler").(string),
			ClientToken:  buildClientToken(),
		}, nil
	case HTTP:
		return buildBaiduCloudCreateAppBlbHTTPListenerArgs(d, meta)
	case HTTPS:
		return buildBaiduCloudCreateAppBlbHTTPSListenerArgs(d, meta)
	case SSL:
		return buildBaiduCloudCreateAppBlbSSLListenerArgs(d, meta)
	default:
		// never run here
		return nil, fmt.Errorf("listener only support protocol [TCP, UDP, HTTP, HTTPS, SSL], but now set: %s", protocol)
	}
}

func buildBaiduCloudCreateAppBlbHTTPListenerArgs(d *schema.ResourceData, meta interface{}) (*appblb.CreateAppHTTPListenerArgs, error) {
	result := &appblb.CreateAppHTTPListenerArgs{
		ClientToken:  buildClientToken(),
		ListenerPort: uint16(d.Get("listener_port").(int)),
		Scheduler:    d.Get("scheduler").(string),
	}
	if result.Scheduler != "RoundRobin" && result.Scheduler != "LeastConnection" {
		return nil, fmt.Errorf("HTTP Listener scheduler only support [RoundRobin, LeastConnection], but you set: %s", result.Scheduler)
	}

	if v, ok := d.GetOk("keep_session"); ok {
		result.KeepSession = v.(bool)
	}

	if v, ok := d.GetOk("keep_session_type"); ok {
		result.KeepSessionType = v.(string)
	}

	if v, ok := d.GetOk("keep_session_timeout"); ok {
		result.KeepSessionTimeout = v.(int)
	}

	if v, ok := d.GetOk("keep_session_cookie_name"); ok {
		result.KeepSessionCookieName = v.(string)
	}

	if v, ok := d.GetOk("x_forwarded_for"); ok {
		result.XForwardedFor = v.(bool)
	}

	if v, ok := d.GetOk("server_timeout"); ok {
		result.ServerTimeout = v.(int)
	}

	if v, ok := d.GetOk("redirect_port"); ok {
		result.RedirectPort = uint16(v.(int))
	}

	return result, nil
}

func buildBaiduCloudCreateAppBlbHTTPSListenerArgs(d *schema.ResourceData, meta interface{}) (*appblb.CreateAppHTTPSListenerArgs, error) {
	result := &appblb.CreateAppHTTPSListenerArgs{
		ClientToken:  buildClientToken(),
		ListenerPort: uint16(d.Get("listener_port").(int)),
		Scheduler:    d.Get("scheduler").(string),
	}
	if result.Scheduler != "RoundRobin" && result.Scheduler != "LeastConnection" {
		return nil, fmt.Errorf("HTTPS Listener scheduler only support [RoundRobin, LeastConnection], but you set: %s", result.Scheduler)
	}

	if v, ok := d.GetOk("keep_session"); ok {
		result.KeepSession = v.(bool)
	}

	if v, ok := d.GetOk("keep_session_type"); ok {
		result.KeepSessionType = v.(string)
	}

	if v, ok := d.GetOk("keep_session_timeout"); ok {
		result.KeepSessionTimeout = v.(int)
	}

	if v, ok := d.GetOk("keep_session_cookie_name"); ok {
		result.KeepSessionCookieName = v.(string)
	}

	if v, ok := d.GetOk("x_forwarded_for"); ok {
		result.XForwardedFor = v.(bool)
	}

	if v, ok := d.GetOk("server_timeout"); ok {
		result.ServerTimeout = v.(int)
	}

	if v, ok := d.GetOk("cert_ids"); ok {
		for _, id := range v.(*schema.Set).List() {
			result.CertIds = append(result.CertIds, id.(string))
		}
	}
	if len(result.CertIds) <= 0 {
		return nil, fmt.Errorf("HTTPS Listener require cert, but not set")
	}

	if v, ok := d.GetOk("ie6_compatible"); ok {
		result.Ie6Compatible = v.(bool)
	}

	if v, ok := d.GetOk("encryption_type"); ok {
		result.EncryptionType = v.(string)
	}

	if v, ok := d.GetOk("encryption_protocols"); ok {
		for _, p := range v.(*schema.Set).List() {
			result.EncryptionProtocols = append(result.EncryptionProtocols, p.(string))
		}
	}

	if v, ok := d.GetOk("dual_auth"); ok {
		result.DualAuth = v.(bool)
	}

	if v, ok := d.GetOk("client_cert_ids"); ok {
		for _, id := range v.(*schema.Set).List() {
			result.ClientCertIds = append(result.ClientCertIds, id.(string))
		}
	}

	return result, nil
}

func buildBaiduCloudCreateAppBlbSSLListenerArgs(d *schema.ResourceData, meta interface{}) (*appblb.CreateAppSSLListenerArgs, error) {
	result := &appblb.CreateAppSSLListenerArgs{
		ClientToken:  buildClientToken(),
		ListenerPort: uint16(d.Get("listener_port").(int)),
		Scheduler:    d.Get("scheduler").(string),
	}

	if v, ok := d.GetOk("cert_ids"); ok {
		for _, id := range v.(*schema.Set).List() {
			result.CertIds = append(result.CertIds, id.(string))
		}
	}
	if len(result.CertIds) <= 0 {
		return nil, fmt.Errorf("SSL Listener require cert, but not set")
	}

	if v, ok := d.GetOk("ie6_compatible"); ok {
		result.Ie6Compatible = v.(bool)
	}

	if v, ok := d.GetOk("encryption_type"); ok {
		result.EncryptionType = v.(string)
	}

	if v, ok := d.GetOk("encryption_protocols"); ok {
		for _, p := range v.(*schema.Set).List() {
			result.EncryptionProtocols = append(result.EncryptionProtocols, p.(string))
		}
	}

	if v, ok := d.GetOk("dual_auth"); ok {
		result.DualAuth = v.(bool)
	}

	if v, ok := d.GetOk("client_cert_ids"); ok {
		for _, id := range v.(*schema.Set).List() {
			result.ClientCertIds = append(result.ClientCertIds, id.(string))
		}
	}

	return result, nil
}

func buildBaiduCloudUpdateAppBlbListenerArgs(d *schema.ResourceData, meta interface{}) (bool, interface{}, error) {
	protocol := d.Get("protocol").(string)

	switch protocol {
	case TCP:
		if d.HasChange("scheduler") || d.HasChange("tcp_session_timeout") {
			return true, &appblb.UpdateAppTCPListenerArgs{
				UpdateAppListenerArgs: appblb.UpdateAppListenerArgs{
					TcpSessionTimeout: d.Get("tcp_session_timeout").(int),
					ListenerPort:      uint16(d.Get("listener_port").(int)),
					Scheduler:         d.Get("scheduler").(string),
					ClientToken:       buildClientToken(),
				},
			}, nil
		}
		return false, nil, nil
	case UDP:
		if d.HasChange("scheduler") {
			return true, &appblb.UpdateAppUDPListenerArgs{
				UpdateAppListenerArgs: appblb.UpdateAppListenerArgs{
					ListenerPort: uint16(d.Get("listener_port").(int)),
					Scheduler:    d.Get("scheduler").(string),
					ClientToken:  buildClientToken(),
				},
			}, nil
		}
		return false, nil, nil
	case HTTP:
		return buildBaiduCloudUpdateAppBlbHTTPListenerArgs(d, meta)
	case HTTPS:
		return buildBaiduCloudUpdateAppBlbHTTPSListenerArgs(d, meta)
	case SSL:
		return buildBaiduCloudUpdateAppBlbSSLListenerArgs(d, meta)
	default:
		// never run here
		return false, nil, fmt.Errorf("listener only support protocol [TCP, UDP, HTTP, HTTPS, SSL], but now set: %s", protocol)
	}
}

func buildBaiduCloudUpdateAppBlbHTTPListenerArgs(d *schema.ResourceData, meta interface{}) (bool, *appblb.UpdateAppHTTPListenerArgs, error) {
	update := false
	result := &appblb.UpdateAppHTTPListenerArgs{
		ClientToken:  buildClientToken(),
		ListenerPort: uint16(d.Get("listener_port").(int)),
		Scheduler:    d.Get("scheduler").(string),
	}

	update = d.HasChange("scheduler")
	if result.Scheduler != "RoundRobin" && result.Scheduler != "LeastConnection" {
		return false, nil, fmt.Errorf("HTTP Listener scheduler only support [RoundRobin, LeastConnection], but you set: %s", result.Scheduler)
	}

	if v, ok := d.GetOk("keep_session"); ok {
		if !update {
			update = d.HasChange("keep_session")
		}

		result.KeepSession = v.(bool)
	}

	if v, ok := d.GetOk("keep_session_type"); ok {
		if !update {
			update = d.HasChange("keep_session_type")
		}

		result.KeepSessionType = v.(string)
	}

	if v, ok := d.GetOk("keep_session_timeout"); ok {
		if !update {
			update = d.HasChange("keep_session_timeout")
		}

		result.KeepSessionTimeout = v.(int)
	}

	if v, ok := d.GetOk("keep_session_cookie_name"); ok {
		if !update {
			update = d.HasChange("keep_session_cookie_name")
		}

		result.KeepSessionCookieName = v.(string)
	}

	if v, ok := d.GetOk("x_forwarded_for"); ok {
		if !update {
			update = d.HasChange("x_forwarded_for")
		}

		result.XForwardedFor = v.(bool)
	}

	if v, ok := d.GetOk("server_timeout"); ok {
		if !update {
			update = d.HasChange("server_timeout")
		}

		result.ServerTimeout = v.(int)
	}

	if v, ok := d.GetOk("redirect_port"); ok {
		if !update {
			update = d.HasChange("redirect_port")
		}

		result.RedirectPort = uint16(v.(int))
	}

	return update, result, nil
}

func buildBaiduCloudUpdateAppBlbHTTPSListenerArgs(d *schema.ResourceData, meta interface{}) (bool, *appblb.UpdateAppHTTPSListenerArgs, error) {
	update := false
	result := &appblb.UpdateAppHTTPSListenerArgs{
		ClientToken:  buildClientToken(),
		ListenerPort: uint16(d.Get("listener_port").(int)),
		Scheduler:    d.Get("scheduler").(string),
	}

	update = d.HasChange("scheduler")
	if result.Scheduler != "RoundRobin" && result.Scheduler != "LeastConnection" {
		return false, nil, fmt.Errorf("HTTPS Listener scheduler only support [RoundRobin, LeastConnection], but you set: %s", result.Scheduler)
	}

	if v, ok := d.GetOk("keep_session"); ok {
		if !update {
			update = d.HasChange("keep_session")
		}

		result.KeepSession = v.(bool)
	}

	if v, ok := d.GetOk("keep_session_type"); ok {
		if !update {
			update = d.HasChange("keep_session_type")
		}

		result.KeepSessionType = v.(string)
	}

	if v, ok := d.GetOk("keep_session_timeout"); ok {
		if !update {
			update = d.HasChange("keep_session_timeout")
		}

		result.KeepSessionTimeout = v.(int)
	}

	if v, ok := d.GetOk("keep_session_cookie_name"); ok {
		if !update {
			update = d.HasChange("keep_session_cookie_name")
		}

		result.KeepSessionCookieName = v.(string)
	}

	if v, ok := d.GetOk("x_forwarded_for"); ok {
		if !update {
			update = d.HasChange("x_forwarded_for")
		}

		result.XForwardedFor = v.(bool)
	}

	if v, ok := d.GetOk("server_timeout"); ok {
		if !update {
			update = d.HasChange("server_timeout")
		}

		result.ServerTimeout = v.(int)
	}

	if v, ok := d.GetOk("cert_ids"); ok {
		if !update {
			update = d.HasChange("cert_ids")
		}
		for _, id := range v.(*schema.Set).List() {
			result.CertIds = append(result.CertIds, id.(string))
		}
	}
	if len(result.CertIds) <= 0 {
		return false, nil, fmt.Errorf("HTTPS Listener require cert, but not set")
	}

	if v, ok := d.GetOk("ie6_compatible"); ok {
		if !update {
			update = d.HasChange("ie6_compatible")
		}

		result.Ie6Compatible = v.(bool)
	}

	if v, ok := d.GetOk("encryption_type"); ok {
		if !update {
			update = d.HasChange("encryption_type")
		}

		result.EncryptionType = v.(string)
	}

	if v, ok := d.GetOk("encryption_protocols"); ok {
		if !update {
			update = d.HasChange("encryption_protocols")
		}

		for _, p := range v.(*schema.Set).List() {
			result.EncryptionProtocols = append(result.EncryptionProtocols, p.(string))
		}
	}

	if v, ok := d.GetOk("dual_auth"); ok {
		if !update {
			update = d.HasChange("dual_auth")
		}

		result.DualAuth = v.(bool)
	}

	if v, ok := d.GetOk("client_cert_ids"); ok {
		if !update {
			update = d.HasChange("client_cert_ids")
		}

		for _, id := range v.(*schema.Set).List() {
			result.ClientCertIds = append(result.ClientCertIds, id.(string))
		}
	}

	return update, result, nil
}

func buildBaiduCloudUpdateAppBlbSSLListenerArgs(d *schema.ResourceData, meta interface{}) (bool, *appblb.UpdateAppSSLListenerArgs, error) {
	update := false
	result := &appblb.UpdateAppSSLListenerArgs{
		ClientToken:  buildClientToken(),
		ListenerPort: uint16(d.Get("listener_port").(int)),
		Scheduler:    d.Get("scheduler").(string),
	}

	update = d.HasChange("scheduler")
	if v, ok := d.GetOk("cert_ids"); ok {
		if !update {
			update = d.HasChange("cert_ids")
		}

		for _, id := range v.(*schema.Set).List() {
			result.CertIds = append(result.CertIds, id.(string))
		}
	}
	if len(result.CertIds) <= 0 {
		return false, nil, fmt.Errorf("SSL Listener require cert, but not set")
	}

	if v, ok := d.GetOk("ie6_compatible"); ok {
		if !update {
			update = d.HasChange("ie6_compatible")
		}

		result.Ie6Compatible = v.(bool)
	}

	if v, ok := d.GetOk("encryption_type"); ok {
		if !update {
			update = d.HasChange("encryption_type")
		}

		result.EncryptionType = v.(string)
	}

	if v, ok := d.GetOk("encryption_protocols"); ok {
		if !update {
			update = d.HasChange("encryption_protocols")
		}

		for _, p := range v.(*schema.Set).List() {
			result.EncryptionProtocols = append(result.EncryptionProtocols, p.(string))
		}
	}

	if v, ok := d.GetOk("dual_auth"); ok {
		if !update {
			update = d.HasChange("dual_auth")
		}

		result.DualAuth = v.(bool)
	}

	if v, ok := d.GetOk("client_cert_ids"); ok {
		if !update {
			update = d.HasChange("client_cert_ids")
		}

		for _, id := range v.(*schema.Set).List() {
			result.ClientCertIds = append(result.ClientCertIds, id.(string))
		}
	}

	return update, result, nil
}

func buildBaiduCloudCreatePolicyArgs(listenerPort int, protocol string, policys []interface{}) (*appblb.CreatePolicysArgs, error) {
	result := &appblb.CreatePolicysArgs{
		ListenerPort: uint16(listenerPort),
		ClientToken:  buildClientToken(),
	}

	if len(policys) == 0 {
		return nil, nil
	}

	transportProtocol := stringInSlice(TransportProtocol, protocol)
	if transportProtocol && len(policys) > 1 {
		return nil, fmt.Errorf("%s Listener only support one policy, but now is %d", protocol, len(policys))
	}

	for _, p := range policys {
		pMap := p.(map[string]interface{})

		policy := &appblb.AppPolicy{
			AppServerGroupId: pMap["app_server_group_id"].(string),
			BackendPort:      uint16(pMap["backend_port"].(int)),
			Priority:         pMap["priority"].(int),
		}

		if v, ok := pMap["description"]; ok && v.(string) != "" {
			policy.Description = v.(string)
		}

		if r, ok := pMap["rule_list"]; ok && len(r.(*schema.Set).List()) > 0 {
			rs := r.(*schema.Set).List()

			if transportProtocol && len(rs) > 1 {
				return nil, fmt.Errorf("%s Listener only support one policy rule, but now is %d", protocol, len(rs))
			}

			for _, r := range rs {
				rMap := r.(map[string]interface{})

				rule := appblb.AppRule{}
				rule.Key = rMap["key"].(string)
				rule.Value = rMap["value"].(string)

				if transportProtocol && (rule.Key != "*" || rule.Value != "*") {
					return nil, fmt.Errorf("%s Listener only support one policy rule [key: *, value: *], but now is [key: %s, value: %s]", protocol, rule.Key, rule.Value)
				}

				policy.RuleList = append(policy.RuleList, rule)
			}
		}

		result.AppPolicyVos = append(result.AppPolicyVos, *policy)
	}

	return result, nil
}
