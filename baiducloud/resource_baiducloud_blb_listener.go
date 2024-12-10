/*
Provide a resource to create an BLB Listener.

Example Usage

```hcl
[TCP/UDP] Listener
resource "baiducloud_blb_listener" "default" {
  blb_id               = "lb-0d29a3f6"
  listener_port        = 124
  protocol             = "TCP"
  scheduler            = "LeastConnection"
}

[HTTP] Listener
resource "baiducloud_blb_listener" "default" {
  blb_id        = "lb-0d29a3f6"
  listener_port = 129
  protocol      = "HTTP"
  scheduler     = "RoundRobin"

}

[HTTPS] Listener
resource "baiducloud_blb_listener" "default" {
  blb_id               = "lb-0d29a3f6"
  listener_port        = 130
  protocol             = "HTTPS"
  scheduler            = "LeastConnection"
  keep_session         = true
  cert_ids             = ["cert-xvysj8xxx"]
  encryption_protocols = ["sslv3", "tlsv10", "tlsv11"]
  encryption_type      = "userDefind"
}

[SSL] Listener
resource "baiducloud_blb_listener" "default" {
  blb_id               = "lb-0d29a3f6"
  listener_port        = 131
  protocol             = "SSL"
  scheduler            = "LeastConnection"
  cert_ids             = ["cert-xvysjxxxx"]
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
	"github.com/baidubce/bce-sdk-go/services/blb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudBlbListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudBlbListenerCreate,
		Read:   resourceBaiduCloudBlbListenerRead,
		Update: resourceBaiduCloudBlbListenerUpdate,
		Delete: resourceBaiduCloudBlbListenerDelete,

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
			"backend_port": {
				Type:         schema.TypeInt,
				Description:  "backend port, range from 1-65535",
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
				ForceNew:     true,
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
			"udp_session_timeout": {
				Type:         schema.TypeInt,
				Description:  "UDP Listener connection session timeout time(second), default 900, support 10-4000",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(10, 4000),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.Get("protocol").(string); ok {
						return v != UDP
					}
					return true
				},
			},
			"health_check_timeout_in_second": {
				Type:        schema.TypeInt,
				Description: "health check timeout in second",
				Optional:    true,
				Computed:    true,
			},
			"health_check_interval": {
				Type:        schema.TypeInt,
				Description: "health check interval",
				Optional:    true,
				Computed:    true,
			},
			"healthy_threshold": {
				Type:        schema.TypeInt,
				Description: "healthy threshold",
				Optional:    true,
				Computed:    true,
			},
			"unhealthy_threshold": {
				Type:        schema.TypeInt,
				Description: "unhealthy threshold",
				Optional:    true,
				Computed:    true,
			},
			"get_blb_ip": {
				Type:        schema.TypeBool,
				Description: "get blb ip or not",
				Optional:    true,
				Computed:    true,
			},
			// UCP
			"health_check_string": {
				Type:        schema.TypeString,
				Description: "health check string, This parameter is mandatory when the listening protocol is UDP",
				Optional:    true,
				Computed:    true,
			},
			// SSL HTTPS
			"applied_ciphers": {
				Type:        schema.TypeString,
				Description: "applied ciphers",
				Optional:    true,
				Computed:    true,
			},
			// http https
			"keep_session_duration": {
				Type:        schema.TypeInt,
				Description: "keep session duration",
				Optional:    true,
				Computed:    true,
			},
			// http https
			"health_check_type": {
				Type:        schema.TypeString,
				Description: "health check type",
				Optional:    true,
				Computed:    true,
			},
			// http https
			"health_check_port": {
				Type:        schema.TypeInt,
				Description: "health check port",
				Optional:    true,
				Computed:    true,
			},
			// http https
			"health_check_uri": {
				Type:        schema.TypeString,
				Description: "health check uri",
				Optional:    true,
				Computed:    true,
			},
			// http https
			"health_check_normal_status": {
				Type:        schema.TypeString,
				Description: "health check normal status",
				Optional:    true,
				Computed:    true,
			},
			// http & https
			"keep_session": {
				Type:             schema.TypeBool,
				Description:      "KeepSession or not",
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: blbProtocolTCPUDPSSLSuppressFunc,
			},
			//http & https
			"keep_session_type": {
				Type:             schema.TypeString,
				Description:      "KeepSessionType option, support insert/rewrite, default insert",
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.StringInSlice([]string{"insert", "rewrite"}, false),
				DiffSuppressFunc: blbProtocolTCPUDPSSLSuppressFunc,
			},
			// http & https
			"keep_session_cookie_name": {
				Type:        schema.TypeString,
				Description: "CookieName which need to covered, useful when keep_session_type is rewrite",
				Optional:    true,
				Computed:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					protocolCheck := blbProtocolTCPUDPSSLSuppressFunc(k, old, new, d)
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
				DiffSuppressFunc: blbProtocolTCPUDPSSLSuppressFunc,
			},
			// http & https
			"server_timeout": {
				Type:             schema.TypeInt,
				Description:      "Backend server maximum timeout time, only support in [1, 3600] second, default 30s",
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: blbProtocolTCPUDPSSLSuppressFunc,
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
				DiffSuppressFunc: blbProtocolTCPUDPHTTPSuppressFunc,
			},
			// https && ssl
			"ie6_compatible": {
				Type:             schema.TypeBool,
				Description:      "Listener support ie6 option, default true",
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: blbProtocolTCPUDPHTTPSuppressFunc,
			},
			// https && ssl
			"encryption_type": {
				Type:             schema.TypeString,
				Description:      "Listener encryption option, support [compatibleIE, incompatibleIE, userDefind]",
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.StringInSlice([]string{"compatibleIE", "incompatibleIE", "userDefind"}, false),
				DiffSuppressFunc: blbProtocolTCPUDPHTTPSuppressFunc,
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
				DiffSuppressFunc: blbProtocolTCPUDPHTTPSuppressFunc,
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
				DiffSuppressFunc: blbProtocolTCPUDPHTTPSuppressFunc,
			},
		},
	}
}

func resourceBaiduCloudBlbListenerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	blbId := d.Get("blb_id").(string)
	protocol := d.Get("protocol").(string)
	listenerPort := d.Get("listener_port").(int)
	action := fmt.Sprintf("Create BLB %s Listener [%s:%d]", blbId, protocol, listenerPort)

	listenerArgs, err := buildBaiduCloudCreateBlbListenerArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}

	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithBLBClient(func(client *blb.Client) (i interface{}, e error) {
			switch protocol {
			case TCP:
				return blbId, client.CreateTCPListener(blbId, listenerArgs.(*blb.CreateTCPListenerArgs))
			case UDP:
				return blbId, client.CreateUDPListener(blbId, listenerArgs.(*blb.CreateUDPListenerArgs))
			case HTTP:
				return blbId, client.CreateHTTPListener(blbId, listenerArgs.(*blb.CreateHTTPListenerArgs))
			case HTTPS:
				return blbId, client.CreateHTTPSListener(blbId, listenerArgs.(*blb.CreateHTTPSListenerArgs))
			case SSL:
				return blbId, client.CreateSSLListener(blbId, listenerArgs.(*blb.CreateSSLListenerArgs))
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
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_blb_listener", action, BCESDKGoERROR)
	}

	return resourceBaiduCloudBlbListenerRead(d, meta)
}

func resourceBaiduCloudBlbListenerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	blbService := BLBService{client}

	blbId := d.Get("blb_id").(string)
	protocol := d.Get("protocol").(string)
	listenerPort := d.Get("listener_port").(int)
	action := fmt.Sprintf("Query BLB %s Listener [%s:%d]", blbId, protocol, listenerPort)

	raw, err := blbService.DescribeListener(blbId, protocol, listenerPort)
	if err != nil {
		d.SetId("")
		return WrapError(err)
	}
	addDebug(action, raw)

	switch protocol {
	case HTTP:
		listenerMeta := raw.(*blb.HTTPListenerModel)
		d.Set("listener_port", listenerMeta.ListenerPort)
		d.Set("backend_port", listenerMeta.BackendPort)
		d.Set("scheduler", listenerMeta.Scheduler)
		d.Set("keep_session", listenerMeta.KeepSession)
		d.Set("keep_session_type", listenerMeta.KeepSessionType)
		d.Set("keep_session_duration", listenerMeta.KeepSessionDuration)
		d.Set("keep_session_cookie_name", listenerMeta.KeepSessionCookieName)
		d.Set("x_forwarded_for", listenerMeta.XForwardedFor)
		d.Set("health_check_type", listenerMeta.HealthCheckType)
		d.Set("health_check_port", listenerMeta.HealthCheckPort)
		d.Set("health_check_uri", listenerMeta.HealthCheckURI)
		d.Set("health_check_timeout_in_second", listenerMeta.HealthCheckTimeoutInSecond)
		d.Set("health_check_interval", listenerMeta.HealthCheckInterval)
		d.Set("unhealthy_threshold", listenerMeta.UnhealthyThreshold)
		d.Set("healthy_threshold", listenerMeta.HealthyThreshold)
		d.Set("get_blb_ip", listenerMeta.GetBlbIp)
		d.Set("health_check_normal_status", listenerMeta.HealthCheckNormalStatus)
		d.Set("server_timeout", listenerMeta.ServerTimeout)
		d.Set("redirect_port", listenerMeta.RedirectPort)
	case HTTPS:
		listenerMeta := raw.(*blb.HTTPSListenerModel)
		d.Set("listener_port", listenerMeta.ListenerPort)
		d.Set("backend_port", listenerMeta.BackendPort)
		d.Set("scheduler", listenerMeta.Scheduler)
		d.Set("keep_session", listenerMeta.KeepSession)
		d.Set("keep_session_type", listenerMeta.KeepSessionType)
		d.Set("keep_session_duration", listenerMeta.KeepSessionDuration)
		d.Set("keep_session_cookie_name", listenerMeta.KeepSessionCookieName)
		d.Set("x_forwarded_for", listenerMeta.XForwardedFor)
		d.Set("health_check_type", listenerMeta.HealthCheckType)
		d.Set("health_check_port", listenerMeta.HealthCheckPort)
		d.Set("health_check_uri", listenerMeta.HealthCheckURI)
		d.Set("health_check_timeout_in_second", listenerMeta.HealthCheckTimeoutInSecond)
		d.Set("health_check_interval", listenerMeta.HealthCheckInterval)
		d.Set("unhealthy_threshold", listenerMeta.UnhealthyThreshold)
		d.Set("healthy_threshold", listenerMeta.HealthyThreshold)
		d.Set("get_blb_ip", listenerMeta.GetBlbIp)
		d.Set("health_check_normal_status", listenerMeta.HealthCheckNormalStatus)
		d.Set("server_timeout", listenerMeta.ServerTimeout)
		d.Set("cert_ids", listenerMeta.CertIds)
		d.Set("dual_auth", listenerMeta.DualAuth)
		d.Set("client_cert_ids", listenerMeta.ClientCertIds)
		d.Set("encryption_type", listenerMeta.EncryptionType)
		d.Set("encryption_protocols", listenerMeta.EncryptionProtocols)
		d.Set("applied_ciphers", listenerMeta.AppliedCiphers)
	case SSL:
		listenerMeta := raw.(*blb.SSLListenerModel)
		d.Set("listener_port", listenerMeta.ListenerPort)
		d.Set("backend_port", listenerMeta.BackendPort)
		d.Set("scheduler", listenerMeta.Scheduler)
		d.Set("health_check_timeout_in_second", listenerMeta.HealthCheckTimeoutInSecond)
		d.Set("health_check_interval", listenerMeta.HealthCheckInterval)
		d.Set("unhealthy_threshold", listenerMeta.UnhealthyThreshold)
		d.Set("healthy_threshold", listenerMeta.HealthyThreshold)
		d.Set("get_blb_ip", listenerMeta.GetBlbIp)
		d.Set("server_timeout", listenerMeta.ServerTimeout)
		d.Set("cert_ids", listenerMeta.CertIds)
		d.Set("dual_auth", listenerMeta.DualAuth)
		d.Set("client_cert_ids", listenerMeta.ClientCertIds)
		d.Set("encryption_type", listenerMeta.EncryptionType)
		d.Set("encryption_protocols", listenerMeta.EncryptionProtocols)
		d.Set("applied_ciphers", listenerMeta.AppliedCiphers)
	case TCP:
		listenerMeta := raw.(*blb.TCPListenerModel)
		d.Set("listener_port", listenerMeta.ListenerPort)
		d.Set("backend_port", listenerMeta.BackendPort)
		d.Set("scheduler", listenerMeta.Scheduler)
		d.Set("tcp_session_timeout", listenerMeta.TcpSessionTimeout)
		d.Set("health_check_timeout_in_second", listenerMeta.HealthCheckTimeoutInSecond)
		d.Set("health_check_interval", listenerMeta.HealthCheckInterval)
		d.Set("unhealthy_threshold", listenerMeta.UnhealthyThreshold)
		d.Set("healthy_threshold", listenerMeta.HealthyThreshold)
		d.Set("get_blb_ip", listenerMeta.GetBlbIp)
	case UDP:
		listenerMeta := raw.(*blb.UDPListenerModel)
		d.Set("listener_port", listenerMeta.ListenerPort)
		d.Set("backend_port", listenerMeta.BackendPort)
		d.Set("scheduler", listenerMeta.Scheduler)
		d.Set("udp_session_timeout", listenerMeta.UdpSessionTimeout)
		d.Set("health_check_timeout_in_second", listenerMeta.HealthCheckTimeoutInSecond)
		d.Set("health_check_interval", listenerMeta.HealthCheckInterval)
		d.Set("unhealthy_threshold", listenerMeta.UnhealthyThreshold)
		d.Set("healthy_threshold", listenerMeta.HealthyThreshold)
		d.Set("health_check_string", listenerMeta.HealthCheckString)
		d.Set("get_blb_ip", listenerMeta.GetBlbIp)
	default:
		return WrapError(fmt.Errorf("unsupport listener type"))
	}
	d.SetId(strconv.Itoa(listenerPort))

	return nil
}

func resourceBaiduCloudBlbListenerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	blbId := d.Get("blb_id").(string)
	protocol := d.Get("protocol").(string)
	listenerPort := d.Get("listener_port").(int)
	action := fmt.Sprintf("Update BLB %s Listener [%s:%d]", blbId, protocol, listenerPort)

	update, args, err := buildBaiduCloudUpdateBlbListenerArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}
	if update {
		_, err := client.WithBLBClient(func(client *blb.Client) (i interface{}, e error) {
			switch protocol {
			case TCP:
				return nil, client.UpdateTCPListener(blbId, args.(*blb.UpdateTCPListenerArgs))
			case UDP:
				return nil, client.UpdateUDPListener(blbId, args.(*blb.UpdateUDPListenerArgs))
			case HTTP:
				return nil, client.UpdateHTTPListener(blbId, args.(*blb.UpdateHTTPListenerArgs))
			case HTTPS:
				return nil, client.UpdateHTTPSListener(blbId, args.(*blb.UpdateHTTPSListenerArgs))
			case SSL:
				return nil, client.UpdateSSLListener(blbId, args.(*blb.UpdateSSLListenerArgs))
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
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_blb_listener", action, BCESDKGoERROR)
		}
	}

	return resourceBaiduCloudBlbListenerRead(d, meta)
}

func resourceBaiduCloudBlbListenerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	blbId := d.Get("blb_id").(string)
	protocol := d.Get("protocol").(string)
	listenerPort := d.Get("listener_port").(int)
	action := fmt.Sprintf("Delete BLB %s Listener [%s:%d]", blbId, protocol, listenerPort)

	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err := client.WithBLBClient(func(client *blb.Client) (i interface{}, e error) {
			return blbId, client.DeleteListeners(blbId, &blb.DeleteListenersArgs{
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
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_blb_listener", action, BCESDKGoERROR)
	}

	return nil
}

func buildBaiduCloudCreateBlbListenerArgs(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	protocol := d.Get("protocol").(string)

	switch protocol {
	case TCP:
		return buildBaiduCloudCreateBlbTCPListenerArgs(d, meta)
	case UDP:
		return buildBaiduCloudCreateBlbUDPListenerArgs(d, meta)
	case HTTP:
		return buildBaiduCloudCreateBlbHTTPListenerArgs(d, meta)
	case HTTPS:
		return buildBaiduCloudCreateBlbHTTPSListenerArgs(d, meta)
	case SSL:
		return buildBaiduCloudCreateBlbSSLListenerArgs(d, meta)
	default:
		// never run here
		return nil, fmt.Errorf("listener only support protocol [TCP, UDP, HTTP, HTTPS, SSL], but now set: %s", protocol)
	}
}

func buildBaiduCloudCreateBlbTCPListenerArgs(d *schema.ResourceData, meta interface{}) (*blb.CreateTCPListenerArgs, error) {
	result := &blb.CreateTCPListenerArgs{
		ClientToken:  buildClientToken(),
		ListenerPort: uint16(d.Get("listener_port").(int)),
		BackendPort:  uint16(d.Get("backend_port").(int)),
		Scheduler:    d.Get("scheduler").(string),
	}

	if v, ok := d.GetOk("tcp_session_timeout"); ok {
		result.TcpSessionTimeout = v.(int)
	}

	if v, ok := d.GetOk("health_check_timeout_in_second"); ok {
		result.HealthCheckTimeoutInSecond = v.(int)
	}

	if v, ok := d.GetOk("health_check_interval"); ok {
		result.HealthCheckInterval = v.(int)
	}

	if v, ok := d.GetOk("unhealthy_threshold"); ok {
		result.UnhealthyThreshold = v.(int)
	}

	if v, ok := d.GetOk("healthy_threshold"); ok {
		result.HealthyThreshold = v.(int)
	}

	return result, nil
}

func buildBaiduCloudCreateBlbUDPListenerArgs(d *schema.ResourceData, meta interface{}) (*blb.CreateUDPListenerArgs, error) {
	result := &blb.CreateUDPListenerArgs{
		ClientToken:  buildClientToken(),
		ListenerPort: uint16(d.Get("listener_port").(int)),
		BackendPort:  uint16(d.Get("backend_port").(int)),
		Scheduler:    d.Get("scheduler").(string),
	}

	if v, ok := d.GetOk("udp_session_timeout"); ok {
		result.UdpSessionTimeout = v.(int)
	}

	if v, ok := d.GetOk("health_check_timeout_in_second"); ok {
		result.HealthCheckTimeoutInSecond = v.(int)
	}

	if v, ok := d.GetOk("health_check_interval"); ok {
		result.HealthCheckInterval = v.(int)
	}

	if v, ok := d.GetOk("unhealthy_threshold"); ok {
		result.UnhealthyThreshold = v.(int)
	}

	if v, ok := d.GetOk("healthy_threshold"); ok {
		result.HealthyThreshold = v.(int)
	}

	if v, ok := d.GetOk("health_check_string"); ok {
		result.HealthCheckString = v.(string)
	}

	return result, nil
}

func buildBaiduCloudCreateBlbHTTPListenerArgs(d *schema.ResourceData, meta interface{}) (*blb.CreateHTTPListenerArgs, error) {
	result := &blb.CreateHTTPListenerArgs{
		ClientToken:  buildClientToken(),
		ListenerPort: uint16(d.Get("listener_port").(int)),
		BackendPort:  uint16(d.Get("backend_port").(int)),
		Scheduler:    d.Get("scheduler").(string),
	}

	if result.Scheduler != "RoundRobin" && result.Scheduler != "LeastConnection" {
		return nil, fmt.Errorf("HTTP Listener scheduler only support [RoundRobin, LeastConnection], but you set: %s", result.Scheduler)
	}

	if v, ok := d.GetOkExists("keep_session"); ok {
		if boolValue, ok := v.(bool); ok {
			result.KeepSession = &boolValue
		}
	}

	if v, ok := d.GetOk("keep_session_type"); ok {
		result.KeepSessionType = v.(string)
	}

	if v, ok := d.GetOk("keep_session_duration"); ok {
		result.KeepSessionDuration = v.(int)
	}

	if v, ok := d.GetOk("keep_session_cookie_name"); ok {
		result.KeepSessionCookieName = v.(string)
	}

	if v, ok := d.GetOkExists("x_forwarded_for"); ok {
		if boolValue, ok := v.(bool); ok {
			result.XForwardedFor = &boolValue
		}
	}

	if v, ok := d.GetOk("health_check_type"); ok {
		result.HealthCheckType = v.(string)
	}

	if v, ok := d.GetOk("health_check_port"); ok {
		result.HealthCheckPort = v.(uint16)
	}

	if v, ok := d.GetOk("health_check_uri"); ok {
		result.HealthCheckURI = v.(string)
	}

	if v, ok := d.GetOk("health_check_timeout_in_second"); ok {
		result.HealthCheckTimeoutInSecond = v.(int)
	}

	if v, ok := d.GetOk("health_check_interval"); ok {
		result.HealthCheckInterval = v.(int)
	}

	if v, ok := d.GetOk("unhealthy_threshold"); ok {
		result.UnhealthyThreshold = v.(int)
	}

	if v, ok := d.GetOk("healthy_threshold"); ok {
		result.HealthyThreshold = v.(int)
	}

	if v, ok := d.GetOk("health_check_normal_status"); ok {
		result.HealthCheckNormalStatus = v.(string)
	}

	if v, ok := d.GetOk("server_timeout"); ok {
		result.ServerTimeout = v.(int)
	}

	if v, ok := d.GetOk("redirect_port"); ok {
		result.RedirectPort = uint16(v.(int))
	}

	return result, nil
}

func buildBaiduCloudCreateBlbHTTPSListenerArgs(d *schema.ResourceData, meta interface{}) (*blb.CreateHTTPSListenerArgs, error) {
	result := &blb.CreateHTTPSListenerArgs{
		ClientToken:  buildClientToken(),
		ListenerPort: uint16(d.Get("listener_port").(int)),
		BackendPort:  uint16(d.Get("backend_port").(int)),
		Scheduler:    d.Get("scheduler").(string),
	}

	if result.Scheduler != "RoundRobin" && result.Scheduler != "LeastConnection" {
		return nil, fmt.Errorf("HTTP Listener scheduler only support [RoundRobin, LeastConnection], but you set: %s", result.Scheduler)
	}

	if v, ok := d.GetOkExists("keep_session"); ok {
		if boolValue, ok := v.(bool); ok {
			result.KeepSession = &boolValue
		}
	}

	if v, ok := d.GetOk("keep_session_type"); ok {
		result.KeepSessionType = v.(string)
	}

	if v, ok := d.GetOk("keep_session_duration"); ok {
		result.KeepSessionDuration = v.(int)
	}

	if v, ok := d.GetOk("keep_session_cookie_name"); ok {
		result.KeepSessionCookieName = v.(string)
	}

	if v, ok := d.GetOkExists("x_forwarded_for"); ok {
		if boolValue, ok := v.(bool); ok {
			result.XForwardedFor = &boolValue
		}
	}

	if v, ok := d.GetOk("health_check_type"); ok {
		result.HealthCheckType = v.(string)
	}

	if v, ok := d.GetOk("health_check_port"); ok {
		result.HealthCheckPort = v.(uint16)
	}

	if v, ok := d.GetOk("health_check_uri"); ok {
		result.HealthCheckURI = v.(string)
	}

	if v, ok := d.GetOk("health_check_timeout_in_second"); ok {
		result.HealthCheckTimeoutInSecond = v.(int)
	}

	if v, ok := d.GetOk("health_check_interval"); ok {
		result.HealthCheckInterval = v.(int)
	}

	if v, ok := d.GetOk("unhealthy_threshold"); ok {
		result.UnhealthyThreshold = v.(int)
	}

	if v, ok := d.GetOk("healthy_threshold"); ok {
		result.HealthyThreshold = v.(int)
	}

	if v, ok := d.GetOk("health_check_normal_status"); ok {
		result.HealthCheckNormalStatus = v.(string)
	}

	if v, ok := d.GetOk("server_timeout"); ok {
		result.ServerTimeout = v.(int)
	}

	if v, ok := d.GetOk("redirect_port"); ok {
		result.RedirectPort = uint16(v.(int))
	}

	if v, ok := d.GetOk("cert_ids"); ok {
		for _, id := range v.(*schema.Set).List() {
			result.CertIds = append(result.CertIds, id.(string))
		}
	}
	if len(result.CertIds) <= 0 {
		return nil, fmt.Errorf("HTTPS Listener require cert, but not set")
	}

	if v, ok := d.GetOk("encryption_type"); ok {
		result.EncryptionType = v.(string)
	}

	if v, ok := d.GetOk("encryption_protocols"); ok {
		for _, p := range v.(*schema.Set).List() {
			result.EncryptionProtocols = append(result.EncryptionProtocols, p.(string))
		}
	}

	if v, ok := d.GetOkExists("dual_auth"); ok {
		if boolValue, ok := v.(bool); ok {
			result.DualAuth = &boolValue
		}
	}

	if v, ok := d.GetOk("client_cert_ids"); ok {
		for _, id := range v.(*schema.Set).List() {
			result.ClientCertIds = append(result.ClientCertIds, id.(string))
		}
	}

	if v, ok := d.GetOk("applied_ciphers"); ok {
		result.AppliedCiphers = v.(string)
	}

	return result, nil
}

func buildBaiduCloudCreateBlbSSLListenerArgs(d *schema.ResourceData, meta interface{}) (*blb.CreateSSLListenerArgs, error) {
	result := &blb.CreateSSLListenerArgs{
		ClientToken:  buildClientToken(),
		ListenerPort: uint16(d.Get("listener_port").(int)),
		BackendPort:  uint16(d.Get("backend_port").(int)),
		Scheduler:    d.Get("scheduler").(string),
	}

	if v, ok := d.GetOk("health_check_timeout_in_second"); ok {
		result.HealthCheckTimeoutInSecond = v.(int)
	}

	if v, ok := d.GetOk("health_check_interval"); ok {
		result.HealthCheckInterval = v.(int)
	}

	if v, ok := d.GetOk("unhealthy_threshold"); ok {
		result.UnhealthyThreshold = v.(int)
	}

	if v, ok := d.GetOk("healthy_threshold"); ok {
		result.HealthyThreshold = v.(int)
	}

	if v, ok := d.GetOk("cert_ids"); ok {
		for _, id := range v.(*schema.Set).List() {
			result.CertIds = append(result.CertIds, id.(string))
		}
	}
	if len(result.CertIds) <= 0 {
		return nil, fmt.Errorf("HTTPS Listener require cert, but not set")
	}

	if v, ok := d.GetOk("encryption_type"); ok {
		result.EncryptionType = v.(string)
	}

	if v, ok := d.GetOk("encryption_protocols"); ok {
		for _, p := range v.(*schema.Set).List() {
			result.EncryptionProtocols = append(result.EncryptionProtocols, p.(string))
		}
	}

	if v, ok := d.GetOkExists("dual_auth"); ok {
		if boolValue, ok := v.(bool); ok {
			result.DualAuth = &boolValue
		}
	}

	if v, ok := d.GetOk("client_cert_ids"); ok {
		for _, id := range v.(*schema.Set).List() {
			result.ClientCertIds = append(result.ClientCertIds, id.(string))
		}
	}

	if v, ok := d.GetOk("applied_ciphers"); ok {
		result.AppliedCiphers = v.(string)
	}

	return result, nil
}

func buildBaiduCloudUpdateBlbListenerArgs(d *schema.ResourceData, meta interface{}) (bool, interface{}, error) {
	protocol := d.Get("protocol").(string)

	switch protocol {
	case TCP:
		return buildBaiduCloudUpdateBlbTCPListenerArgs(d, meta)
	case UDP:
		return buildBaiduCloudUpdateBlbUDPListenerArgs(d, meta)
	case HTTP:
		return buildBaiduCloudUpdateBlbHTTPListenerArgs(d, meta)
	case HTTPS:
		return buildBaiduCloudUpdateBlbHTTPSListenerArgs(d, meta)
	case SSL:
		return buildBaiduCloudUpdateBlbSSLListenerArgs(d, meta)
	default:
		// never run here
		return false, nil, fmt.Errorf("listener only support protocol [TCP, UDP, HTTP, HTTPS, SSL], but now set: %s", protocol)
	}
}

func buildBaiduCloudUpdateBlbTCPListenerArgs(d *schema.ResourceData, meta interface{}) (bool, *blb.UpdateTCPListenerArgs, error) {
	update := false
	result := &blb.UpdateTCPListenerArgs{
		ClientToken:  buildClientToken(),
		ListenerPort: uint16(d.Get("listener_port").(int)),
		Scheduler:    d.Get("scheduler").(string),
	}

	update = d.HasChange("scheduler")

	if v, ok := d.GetOk("backend_port"); ok {
		if !update {
			update = d.HasChange("backend_port")
		}

		result.BackendPort = v.(uint16)
	}

	if v, ok := d.GetOk("tcp_session_timeout"); ok {
		if !update {
			update = d.HasChange("tcp_session_timeout")
		}

		result.TcpSessionTimeout = v.(int)
	}

	if v, ok := d.GetOk("health_check_timeout_in_second"); ok {
		if !update {
			update = d.HasChange("health_check_timeout_in_second")
		}

		result.HealthCheckTimeoutInSecond = v.(int)
	}

	if v, ok := d.GetOk("health_check_interval"); ok {
		if !update {
			update = d.HasChange("health_check_interval")
		}

		result.HealthCheckInterval = v.(int)
	}

	if v, ok := d.GetOk("unhealthy_threshold"); ok {
		if !update {
			update = d.HasChange("unhealthy_threshold")
		}

		result.UnhealthyThreshold = v.(int)
	}

	if v, ok := d.GetOk("healthy_threshold"); ok {
		if !update {
			update = d.HasChange("healthy_threshold")
		}

		result.HealthyThreshold = v.(int)
	}

	return update, result, nil
}

func buildBaiduCloudUpdateBlbUDPListenerArgs(d *schema.ResourceData, meta interface{}) (bool, *blb.UpdateUDPListenerArgs, error) {
	update := false
	result := &blb.UpdateUDPListenerArgs{
		ClientToken:  buildClientToken(),
		ListenerPort: uint16(d.Get("listener_port").(int)),
		Scheduler:    d.Get("scheduler").(string),
	}

	update = d.HasChange("scheduler")

	if v, ok := d.GetOk("backend_port"); ok {
		if !update {
			update = d.HasChange("backend_port")
		}

		result.BackendPort = v.(uint16)
	}

	if v, ok := d.GetOk("udp_session_timeout"); ok {
		if !update {
			update = d.HasChange("udp_session_timeout")
		}

		result.UdpSessionTimeout = v.(int)
	}

	if v, ok := d.GetOk("health_check_timeout_in_second"); ok {
		if !update {
			update = d.HasChange("health_check_timeout_in_second")
		}

		result.HealthCheckTimeoutInSecond = v.(int)
	}

	if v, ok := d.GetOk("health_check_interval"); ok {
		if !update {
			update = d.HasChange("health_check_interval")
		}

		result.HealthCheckInterval = v.(int)
	}

	if v, ok := d.GetOk("unhealthy_threshold"); ok {
		if !update {
			update = d.HasChange("unhealthy_threshold")
		}

		result.UnhealthyThreshold = v.(int)
	}

	if v, ok := d.GetOk("healthy_threshold"); ok {
		if !update {
			update = d.HasChange("healthy_threshold")
		}

		result.HealthyThreshold = v.(int)
	}

	if v, ok := d.GetOk("health_check_string"); ok {
		if !update {
			update = d.HasChange("health_check_string")
		}

		result.HealthCheckString = v.(string)
	}

	return update, result, nil
}

func buildBaiduCloudUpdateBlbHTTPListenerArgs(d *schema.ResourceData, meta interface{}) (bool, *blb.UpdateHTTPListenerArgs, error) {
	update := false
	result := &blb.UpdateHTTPListenerArgs{
		ClientToken:  buildClientToken(),
		ListenerPort: uint16(d.Get("listener_port").(int)),
		Scheduler:    d.Get("scheduler").(string),
	}

	update = d.HasChange("scheduler")
	if result.Scheduler != "RoundRobin" && result.Scheduler != "LeastConnection" {
		return false, nil, fmt.Errorf("HTTP Listener scheduler only support [RoundRobin, LeastConnection], but you set: %s", result.Scheduler)
	}

	if v, ok := d.GetOk("backend_port"); ok {
		if !update {
			update = d.HasChange("backend_port")
		}

		result.BackendPort = v.(uint16)
	}

	if v, ok := d.GetOkExists("keep_session"); ok {
		if !update {
			update = d.HasChange("keep_session")
		}

		if boolValue, ok := v.(bool); ok {
			result.KeepSession = &boolValue
		}
	}

	if v, ok := d.GetOk("keep_session_type"); ok {
		if !update {
			update = d.HasChange("keep_session_type")
		}

		result.KeepSessionType = v.(string)
	}

	if v, ok := d.GetOk("keep_session_duration"); ok {
		if !update {
			update = d.HasChange("keep_session_duration")
		}

		result.KeepSessionDuration = v.(int)
	}

	if v, ok := d.GetOk("keep_session_cookie_name"); ok {
		if !update {
			update = d.HasChange("keep_session_cookie_name")
		}

		result.KeepSessionCookieName = v.(string)
	}

	if v, ok := d.GetOkExists("x_forwarded_for"); ok {
		if !update {
			update = d.HasChange("x_forwarded_for")
		}

		if boolValue, ok := v.(bool); ok {
			result.XForwardedFor = &boolValue
		}
	}

	if v, ok := d.GetOk("health_check_type"); ok {
		if !update {
			update = d.HasChange("health_check_type")
		}

		result.HealthCheckType = v.(string)
	}

	if v, ok := d.GetOk("health_check_port"); ok {
		if !update {
			update = d.HasChange("health_check_port")
		}

		result.HealthCheckPort = v.(uint16)
	}

	if v, ok := d.GetOk("health_check_uri"); ok {
		if !update {
			update = d.HasChange("health_check_uri")
		}

		result.HealthCheckURI = v.(string)
	}

	if v, ok := d.GetOk("health_check_timeout_in_second"); ok {
		if !update {
			update = d.HasChange("health_check_timeout_in_second")
		}

		result.HealthCheckTimeoutInSecond = v.(int)
	}

	if v, ok := d.GetOk("health_check_interval"); ok {
		if !update {
			update = d.HasChange("health_check_interval")
		}

		result.HealthCheckInterval = v.(int)
	}

	if v, ok := d.GetOk("unhealthy_threshold"); ok {
		if !update {
			update = d.HasChange("unhealthy_threshold")
		}

		result.UnhealthyThreshold = v.(int)
	}

	if v, ok := d.GetOk("healthy_threshold"); ok {
		if !update {
			update = d.HasChange("healthy_threshold")
		}

		result.HealthyThreshold = v.(int)
	}

	if v, ok := d.GetOk("health_check_normal_status"); ok {
		if !update {
			update = d.HasChange("health_check_normal_status")
		}

		result.HealthCheckNormalStatus = v.(string)
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

func buildBaiduCloudUpdateBlbHTTPSListenerArgs(d *schema.ResourceData, meta interface{}) (bool, *blb.UpdateHTTPSListenerArgs, error) {
	update := false
	result := &blb.UpdateHTTPSListenerArgs{
		ClientToken:  buildClientToken(),
		ListenerPort: uint16(d.Get("listener_port").(int)),
		Scheduler:    d.Get("scheduler").(string),
	}

	update = d.HasChange("scheduler")
	if result.Scheduler != "RoundRobin" && result.Scheduler != "LeastConnection" {
		return false, nil, fmt.Errorf("HTTPS Listener scheduler only support [RoundRobin, LeastConnection], but you set: %s", result.Scheduler)
	}

	if v, ok := d.GetOk("backend_port"); ok {
		if !update {
			update = d.HasChange("backend_port")
		}

		result.BackendPort = v.(uint16)
	}

	if v, ok := d.GetOkExists("keep_session"); ok {
		if !update {
			update = d.HasChange("keep_session")
		}

		if boolValue, ok := v.(bool); ok {
			result.KeepSession = &boolValue
		}
	}

	if v, ok := d.GetOk("keep_session_type"); ok {
		if !update {
			update = d.HasChange("keep_session_type")
		}

		result.KeepSessionType = v.(string)
	}

	if v, ok := d.GetOk("keep_session_duration"); ok {
		if !update {
			update = d.HasChange("keep_session_duration")
		}

		result.KeepSessionDuration = v.(int)
	}

	if v, ok := d.GetOk("keep_session_cookie_name"); ok {
		if !update {
			update = d.HasChange("keep_session_cookie_name")
		}

		result.KeepSessionCookieName = v.(string)
	}

	if v, ok := d.GetOkExists("x_forwarded_for"); ok {
		if !update {
			update = d.HasChange("x_forwarded_for")
		}

		if boolValue, ok := v.(bool); ok {
			result.XForwardedFor = &boolValue
		}
	}

	if v, ok := d.GetOk("health_check_type"); ok {
		if !update {
			update = d.HasChange("health_check_type")
		}

		result.HealthCheckType = v.(string)
	}

	if v, ok := d.GetOk("health_check_port"); ok {
		if !update {
			update = d.HasChange("health_check_port")
		}

		result.HealthCheckPort = v.(uint16)
	}

	if v, ok := d.GetOk("health_check_uri"); ok {
		if !update {
			update = d.HasChange("health_check_uri")
		}

		result.HealthCheckURI = v.(string)
	}

	if v, ok := d.GetOk("health_check_timeout_in_second"); ok {
		if !update {
			update = d.HasChange("health_check_timeout_in_second")
		}

		result.HealthCheckTimeoutInSecond = v.(int)
	}

	if v, ok := d.GetOk("health_check_interval"); ok {
		if !update {
			update = d.HasChange("health_check_interval")
		}

		result.HealthCheckInterval = v.(int)
	}

	if v, ok := d.GetOk("unhealthy_threshold"); ok {
		if !update {
			update = d.HasChange("unhealthy_threshold")
		}

		result.UnhealthyThreshold = v.(int)
	}

	if v, ok := d.GetOk("healthy_threshold"); ok {
		if !update {
			update = d.HasChange("healthy_threshold")
		}

		result.HealthyThreshold = v.(int)
	}

	if v, ok := d.GetOk("health_check_normal_status"); ok {
		if !update {
			update = d.HasChange("health_check_normal_status")
		}

		result.HealthCheckNormalStatus = v.(string)
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

	if v, ok := d.GetOk("applied_ciphers"); ok {
		if !update {
			update = d.HasChange("applied_ciphers")
		}

		result.AppliedCiphers = v.(string)
	}

	return update, result, nil
}

func buildBaiduCloudUpdateBlbSSLListenerArgs(d *schema.ResourceData, meta interface{}) (bool, *blb.UpdateSSLListenerArgs, error) {
	update := false
	result := &blb.UpdateSSLListenerArgs{
		ClientToken:  buildClientToken(),
		ListenerPort: uint16(d.Get("listener_port").(int)),
		Scheduler:    d.Get("scheduler").(string),
	}

	update = d.HasChange("scheduler")
	if v, ok := d.GetOk("backend_port"); ok {
		if !update {
			update = d.HasChange("backend_port")
		}

		result.BackendPort = v.(uint16)
	}

	if v, ok := d.GetOk("health_check_timeout_in_second"); ok {
		if !update {
			update = d.HasChange("health_check_timeout_in_second")
		}

		result.HealthCheckTimeoutInSecond = v.(int)
	}

	if v, ok := d.GetOk("health_check_interval"); ok {
		if !update {
			update = d.HasChange("health_check_interval")
		}

		result.HealthCheckInterval = v.(int)
	}

	if v, ok := d.GetOk("unhealthy_threshold"); ok {
		if !update {
			update = d.HasChange("unhealthy_threshold")
		}

		result.UnhealthyThreshold = v.(int)
	}

	if v, ok := d.GetOk("healthy_threshold"); ok {
		if !update {
			update = d.HasChange("healthy_threshold")
		}

		result.HealthyThreshold = v.(int)
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

	if v, ok := d.GetOk("applied_ciphers"); ok {
		if !update {
			update = d.HasChange("applied_ciphers")
		}

		result.AppliedCiphers = v.(string)
	}

	if v, ok := d.GetOkExists("dual_auth"); ok {
		if !update {
			update = d.HasChange("dual_auth")
		}

		if boolValue, ok := v.(bool); ok {
			result.DualAuth = &boolValue
		}
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
