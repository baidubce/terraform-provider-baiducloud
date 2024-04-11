/*
Provide a resource to create a Peer Conn Acceptor.

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

resource "baiducloud_peer_conn_acceptor" "default" {
  peer_conn_id = "${baiducloud_peer_conn.default.id}"
  auto_accept = true
  dns_sync = true
}
```

Import

Peer Conn Acceptor instance can be imported, e.g.

```hcl
$ terraform import baiducloud_peer_conn_acceptor.default peer_conn_id
```
*/
package baiducloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/baidubce/bce-sdk-go/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudPeerConnAcceptor() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudPeerConnAcceptorCreate,
		Read:   resourceBaiduCloudPeerConnAcceptorRead,
		Update: resourceBaiduCloudPeerConnAcceptorUpdate,
		Delete: resourceBaiduCloudPeerConnAcceptorDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"peer_conn_id": {
				Type:        schema.TypeString,
				Description: "ID of the peer connection.",
				Required:    true,
				ForceNew:    true,
			},
			"auto_accept": {
				Type:          schema.TypeBool,
				Description:   "Whether to accept the peer connection request, default to false.",
				Optional:      true,
				Default:       false,
				ConflictsWith: []string{"auto_reject"},
			},
			"auto_reject": {
				Type:          schema.TypeBool,
				Description:   "Whether to reject the peer connection request, default to false.",
				Optional:      true,
				Default:       false,
				ConflictsWith: []string{"auto_accept"},
			},
			"dns_sync": {
				Type:        schema.TypeBool,
				Description: "Whether to open the switch of dns synchronization.",
				Optional:    true,
				Default:     false,
			},
			"bandwidth_in_mbps": {
				Type:        schema.TypeInt,
				Description: "Bandwidth(Mbps) of the peer connection.",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the peer connection.",
				Computed:    true,
			},
			"local_if_name": {
				Type:        schema.TypeString,
				Description: "Local interface name of the peer connection.",
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
				Computed:    true,
			},
			"peer_account_id": {
				Type:        schema.TypeString,
				Description: "Peer account ID of the peer VPC, which is required only when creating a peer connection across accounts.",
				Computed:    true,
			},
			"peer_vpc_id": {
				Type:        schema.TypeString,
				Description: "Peer VPC ID of the peer connection.",
				Computed:    true,
			},
			"peer_region": {
				Type:        schema.TypeString,
				Description: "Peer region of the peer connection.",
				Computed:    true,
			},
			"peer_if_name": {
				Type:        schema.TypeString,
				Description: "Peer interface name of the peer connection, which is allowed to be set only when the peer connection within this account.",
				Computed:    true,
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
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"payment_timing": {
							Type:        schema.TypeString,
							Description: "Payment timing of billing, which can be Prepaid or Postpaid. The default is Postpaid.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func resourceBaiduCloudPeerConnAcceptorCreate(d *schema.ResourceData, meta interface{}) error {
	action := "create peer conn of acceptor"

	peerConnID := d.Get("peer_conn_id").(string)
	d.SetId(peerConnID)

	if err := resourceBaiduCloudPeerConnAcceptorRead(d, meta); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_peer_conn_acceptor", action, BCESDKGoERROR)
	}

	if d.Id() == "" {
		err := fmt.Errorf("Peer Conn %s is not found.", peerConnID)
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_peer_conn_acceptor", action, BCESDKGoERROR)
	}

	return resourceBaiduCloudPeerConnAcceptorUpdate(d, meta)
}

func resourceBaiduCloudPeerConnAcceptorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	vpcService := VpcService{client}

	peerConnId := d.Id()
	action := "Query Peer Conn " + peerConnId

	result, state, err := vpcService.PeerConnStateRefresh(peerConnId, vpc.PEERCONN_ROLE_ACCEPTOR)()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_peer_conn_acceptor", action, BCESDKGoERROR)
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
	result, _, err = vpcService.PeerConnStateRefresh(peerConnId, vpc.PEERCONN_ROLE_INITIATOR)()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_peer_conn_acceptor", action, BCESDKGoERROR)
	}
	initiatorConn := result.(*vpc.PeerConn)
	autoAccept := true
	autoReject := false
	if initiatorConn.PeerAccountId != peerConn.PeerAccountId && peerConn.Status == vpc.PEERCONN_STATUS_CONSULTING {
		autoAccept = false
		autoReject = true
	}
	d.Set("auto_accept", autoAccept)
	d.Set("auto_reject", autoReject)

	return nil
}

func resourceBaiduCloudPeerConnAcceptorUpdate(d *schema.ResourceData, meta interface{}) error {
	peerConnID := d.Id()
	action := "Update the peer conn of acceptor " + peerConnID
	client := meta.(*connectivity.BaiduClient)
	vpcService := &VpcService{client}

	d.Partial(true)

	peerConn, err := vpcService.GetPeerConnDetail(peerConnID, vpc.PEERCONN_ROLE_ACCEPTOR)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_peer_conn_acceptor", action, BCESDKGoERROR)
	}

	autoAccept := d.Get("auto_accept").(bool)
	if autoAccept && peerConn.Status == vpc.PEERCONN_STATUS_CONSULTING {
		if _, err := client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
			return nil, vpcClient.AcceptPeerConnApply(peerConnID, buildClientToken())
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_peer_conn_acceptor", action, BCESDKGoERROR)
		}
		d.SetPartial("auto_accept")
	}

	autoReject := d.Get("auto_reject").(bool)
	if autoReject && peerConn.Status == vpc.PEERCONN_STATUS_CONSULTING {
		if _, err := client.WithVpcClient(func(vpcClient *vpc.Client) (i interface{}, e error) {
			return nil, vpcClient.RejectPeerConnApply(peerConnID, buildClientToken())
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_peer_conn_acceptor", action, BCESDKGoERROR)
		}
		d.SetPartial("auto_reject")
	}

	if d.HasChange("dns_sync") {
		if err := resourceBaiduCloudPeerConnDNSSync(d, meta, vpc.PEERCONN_ROLE_ACCEPTOR); err != nil {
			return err
		}
		d.SetPartial("dns_sync")
	}

	d.Partial(false)

	return resourceBaiduCloudPeerConnAcceptorRead(d, meta)
}

func resourceBaiduCloudPeerConnAcceptorDelete(d *schema.ResourceData, meta interface{}) error {
	action := "The peer conn will be removed from the state file, however the resources will remain."
	addDebug(action, d.Id())
	return nil
}
