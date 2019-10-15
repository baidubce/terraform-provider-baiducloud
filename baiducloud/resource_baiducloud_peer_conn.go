/*
Provide a resource to create a Peer Conn.

Example Usage

```hcl
resource "baiducloud_peer_conn" "default" {
  bandwidth_in_mbps = 10
  local_vpc_id = "vpc-y4p102r3mz6m"
  peer_vpc_id = "vpc-4njbqurm0uag"
  peer_region = "bj"
  billing = {
    payment_timing = "Postpaid"
  }
}
```

Import

Peer Conn instance can be imported, e.g.

```hcl
$ terraform import baiducloud_peer_conn.default peer_conn_id
```
*/
package baiducloud

import (
	"strings"
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudPeerConn() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudPeerConnCreate,
		Read:   resourceBaiduCloudPeerConnRead,
		Update: resourceBaiduCloudPeerConnUpdate,
		Delete: resourceBaiduCloudPeerConnDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"bandwidth_in_mbps": {
				Type:        schema.TypeInt,
				Description: "Bandwidth(Mbps) of the peer connection.",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the peer connection.",
				Optional:    true,
			},
			"local_if_name": {
				Type:        schema.TypeString,
				Description: "Local interface name of the peer connection.",
				Optional:    true,
				Computed:    true,
			},
			"local_if_id": {
				Type:        schema.TypeString,
				Description: "Local interface ID of the peer connection.",
				Computed:    true,
			},
			"local_vpc_id": {
				Type:        schema.TypeString,
				Description: "Local VPC ID of the peer connection.",
				Required:    true,
				ForceNew:    true,
			},
			"peer_account_id": {
				Type:        schema.TypeString,
				Description: "Peer account ID of the peer VPC, which is required only when creating a peer connection across accounts.",
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"peer_vpc_id": {
				Type:        schema.TypeString,
				Description: "Peer VPC ID of the peer connection.",
				Required:    true,
				ForceNew:    true,
			},
			"peer_region": {
				Type:        schema.TypeString,
				Description: "Peer region of the peer connection.",
				Required:    true,
				ForceNew:    true,
			},
			"peer_if_name": {
				Type:        schema.TypeString,
				Description: "Peer interface name of the peer connection, which is allowed to be set only when the peer connection within this account.",
				Optional:    true,
			},
			"dns_sync": {
				Type:        schema.TypeBool,
				Description: "Whether to open the switch of dns synchronization.",
				Optional:    true,
				Default:     false,
			},
			"role": {
				Type:        schema.TypeString,
				Description: "Role of the peer connection, which can be initiator or acceptor.",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Status of the peer connection.",
				Computed:    true,
			},
			"created_time": {
				Type:        schema.TypeString,
				Description: "Created time of the peer connection.",
				Computed:    true,
			},
			"expired_time": {
				Type:        schema.TypeString,
				Description: "Expired time of the peer connection, which will be empty when the payment_timing is Postpaid.",
				Computed:    true,
			},
			"dns_status": {
				Type:        schema.TypeString,
				Description: "DNS status of the peer connection.",
				Computed:    true,
			},
			"billing": {
				Type:        schema.TypeMap,
				Description: "Billing information of the peer connection.",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"payment_timing": {
							Type:         schema.TypeString,
							Description:  "Payment timing of the billing, which can be Prepaid or Postpaid. The default is Postpaid.",
							Required:     true,
							ForceNew:     true,
							Default:      PAYMENT_TIMING_POSTPAID,
							ValidateFunc: validatePaymentTiming(),
						},
						"reservation": {
							Type:             schema.TypeMap,
							Description:      "Reservation of the peer connection.",
							Optional:         true,
							DiffSuppressFunc: postPaidDiffSuppressFunc,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"reservation_length": {
										Type:             schema.TypeInt,
										Description:      "Reservation length that you will pay for your resource. It is valid when payment_timing is Prepaid. Valid values: [1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36].",
										Optional:         true,
										Default:          1,
										ForceNew:         true,
										DiffSuppressFunc: postPaidDiffSuppressFunc,
										ValidateFunc:     validateReservationLength(),
									},
									"reservation_time_unit": {
										Type:             schema.TypeString,
										Description:      "Reservation time unit that you will pay for your resource. It is valid when payment_timing is Prepaid. The value can only be month currently, which is also the default value.",
										Optional:         true,
										Default:          "month",
										ValidateFunc:     validateReservationUnit(),
										DiffSuppressFunc: postPaidDiffSuppressFunc,
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

func resourceBaiduCloudPeerConnCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	vpcService := VpcService{client}

	args := buildBaiduCloudPeerConnArgs(d)
	action := "Create Peer Conn"

	if err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
			return vpcClient.CreatePeerConn(args)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw)
		result, _ := raw.(*vpc.CreatePeerConnResult)
		d.SetId(result.PeerConnId)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_peer_conn", action, BCESDKGoERROR)
	}

	stateConf := buildStateConf(
		[]string{string(vpc.PEERCONN_STATUS_CREATING)},
		[]string{string(vpc.PEERCONN_STATUS_ACTIVE), string(vpc.PEERCONN_STATUS_CONSULTING)},
		d.Timeout(schema.TimeoutCreate),
		vpcService.PeerConnStateRefresh(d.Id(), vpc.PEERCONN_ROLE_INITIATOR))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_peer_conn", action, BCESDKGoERROR)
	}

	dnsSync := d.Get("dns_sync").(bool)
	if dnsSync {
		if err := vpcService.OpenPeerConnDNSSync(d, d.Id(), vpc.PEERCONN_ROLE_INITIATOR); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_peer_conn", action, BCESDKGoERROR)
		}
	}

	return resourceBaiduCloudPeerConnRead(d, meta)
}

func resourceBaiduCloudPeerConnRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	vpcService := VpcService{client}

	peerConnId := d.Id()
	action := "Query Peer Conn " + peerConnId

	result, state, err := vpcService.PeerConnStateRefresh(peerConnId, vpc.PEERCONN_ROLE_INITIATOR)()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_peer_conn", action, BCESDKGoERROR)
	}
	addDebug(action, result)

	status := map[string]bool{
		string(vpc.PEERCONN_STATUS_DELETED):        true,
		string(vpc.PEERCONN_STATUS_DELETING):       true,
		string(vpc.PEERCONN_STATUS_EXPIRED):        true,
		string(vpc.PEERCONN_STATUS_ERROR):          true,
		string(vpc.PEERCONN_STATUS_CONSULT_FAILED): true,
		string(vpc.PEERCONN_STATUS_DOWN):           true,
	}
	if _, ok := status[strings.ToLower(state)]; ok || result == nil {
		d.SetId("")
		return nil
	}

	peerConn := result.(*vpc.PeerConn)
	setAttributeForPeerConn(d, peerConn)

	return nil
}

func resourceBaiduCloudPeerConnUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	peerConnId := d.Id()
	action := "Update Peer Conn " + peerConnId

	d.Partial(true)
	update := false
	args := &vpc.UpdatePeerConnArgs{
		LocalIfId: d.Get("local_if_id").(string),
	}
	if d.HasChange("local_if_name") {
		update = true
		args.LocalIfName = d.Get("local_if_name").(string)
	}
	if d.HasChange("description") {
		update = true
		args.Description = d.Get("description").(string)
	}
	if update {
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
			return nil, vpcClient.UpdatePeerConn(peerConnId, args)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_peer_conn", action, BCESDKGoERROR)
		}
		d.SetPartial("local_if_name")
		d.SetPartial("description")
	}

	if d.HasChange("bandwidth_in_mbps") {
		newBandwidthInMbps := d.Get("bandwidth_in_mbps").(int)
		args := &vpc.ResizePeerConnArgs{
			NewBandwidthInMbps: newBandwidthInMbps,
		}
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
			return nil, vpcClient.ResizePeerConn(peerConnId, args)
		})
		addDebug(action, err)
		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_peer_conn", action, BCESDKGoERROR)
		}
		d.SetPartial("bandwidth_in_mbps")
	}

	if d.HasChange("dns_sync") {
		if err := resourceBaiduCloudPeerConnDNSSync(d, meta, vpc.PEERCONN_ROLE_INITIATOR); err != nil {
			return err
		}
		d.SetPartial("dns_sync")
	}

	d.Partial(false)

	return resourceBaiduCloudPeerConnRead(d, meta)
}

func resourceBaiduCloudPeerConnDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	vpcService := VpcService{client}

	peerConnId := d.Id()
	action := "Delete Peer Conn " + peerConnId

	_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
		return nil, vpcClient.DeletePeerConn(peerConnId, buildClientToken())
	})
	if err != nil {
		if NotFoundError(err) || IsExceptedErrors(err, PeerConnNotFound) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_peer_conn", action, BCESDKGoERROR)
	}
	addDebug(action, nil)

	stateConf := buildStateConf(
		[]string{string(vpc.PEERCONN_STATUS_DELETING), string(vpc.PEERCONN_STATUS_CONSULTING),
			string(vpc.PEERCONN_STATUS_CONSULT_FAILED), string(vpc.PEERCONN_STATUS_ACTIVE)},
		[]string{string(vpc.PEERCONN_STATUS_DELETED)},
		d.Timeout(schema.TimeoutDelete),
		vpcService.PeerConnStateRefresh(peerConnId, vpc.PEERCONN_ROLE_INITIATOR))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_peer_conn", action, BCESDKGoERROR)
	}

	return nil
}

func buildBaiduCloudPeerConnArgs(d *schema.ResourceData) *vpc.CreatePeerConnArgs {
	args := &vpc.CreatePeerConnArgs{
		Billing:     &vpc.Billing{},
		ClientToken: buildClientToken(),
	}

	if v := d.Get("bandwidth_in_mbps").(int); v != 0 {
		args.BandwidthInMbps = v
	}
	if v := d.Get("description").(string); v != "" {
		args.Description = v
	}
	if v := d.Get("local_if_name").(string); v != "" {
		args.LocalIfName = v
	}
	if v := d.Get("local_vpc_id").(string); v != "" {
		args.LocalVpcId = v
	}
	if v := d.Get("peer_account_id").(string); v != "" {
		args.PeerAccountId = v
	}
	if v := d.Get("peer_vpc_id").(string); v != "" {
		args.PeerVpcId = v
	}
	if v := d.Get("peer_region").(string); v != "" {
		args.PeerRegion = v
	}
	if v := d.Get("peer_if_name").(string); v != "" {
		args.PeerIfName = v
	}
	if v, ok := d.GetOk("billing"); ok {
		billing := v.(map[string]interface{})
		if p, ok := billing["payment_timing"]; ok {
			paymentTiming := vpc.PaymentTimingType(p.(string))
			args.Billing.PaymentTiming = paymentTiming
		}
		if args.Billing.PaymentTiming == PAYMENT_TIMING_PREPAID {
			if r, ok := billing["reservation"]; ok {
				args.Billing.Reservation = &vpc.Reservation{}
				reservation := r.(map[string]interface{})
				if reservationLength, ok := reservation["reservation_length"]; ok {
					args.Billing.Reservation.ReservationLength = reservationLength.(int)
				}
				if reservationTimeUnit, ok := reservation["reservation_time_unit"]; ok {
					args.Billing.Reservation.ReservationTimeUnit = reservationTimeUnit.(string)
				}
			}
		}
	}

	return args
}

func resourceBaiduCloudPeerConnDNSSync(d *schema.ResourceData, meta interface{}, role vpc.PeerConnRoleType) error {
	peerConnId := d.Id()
	action := "Update peer conn DNS sync " + peerConnId

	client := meta.(*connectivity.BaiduClient)
	vpcService := &VpcService{client}

	dnsSync := d.Get("dns_sync").(bool)

	var err error
	if dnsSync {
		// open DNS sync
		err = vpcService.OpenPeerConnDNSSync(d, peerConnId, role)
	} else {
		// close DNS sync
		err = vpcService.ClosePeerConnDNSSync(d, peerConnId, role)
	}

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_peer_conn", action, BCESDKGoERROR)
	}

	return nil
}

func setAttributeForPeerConn(d *schema.ResourceData, peerConn *vpc.PeerConn) {
	d.Set("role", peerConn.Role)
	d.Set("status", peerConn.Status)
	d.Set("bandwidth_in_mbps", peerConn.BandwidthInMbps)
	d.Set("description", peerConn.Description)
	d.Set("local_if_id", peerConn.LocalIfId)
	d.Set("local_if_name", peerConn.LocalIfName)
	d.Set("local_vpc_id", peerConn.LocalVpcId)
	d.Set("peer_account_id", peerConn.PeerAccountId)
	d.Set("peer_vpc_id", peerConn.PeerVpcId)
	d.Set("peer_region", peerConn.PeerRegion)
	d.Set("created_time", peerConn.CreatedTime)
	d.Set("expired_time", peerConn.ExpiredTime)
	d.Set("dns_status", peerConn.DnsStatus)

	billingMap := map[string]interface{}{"payment_timing": peerConn.PaymentTiming}
	d.Set("billing", billingMap)
}
