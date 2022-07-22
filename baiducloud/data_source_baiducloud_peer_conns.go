/*
Use this data source to query Peer Conn list.

Example Usage

```hcl
data "baiducloud_peer_conns" "default" {
  vpc_id = "vpc-y4p102r3mz6m"
}

output "peer_conns" {
  value = "${data.baiducloud_peer_conns.default.peer_conns}"
}
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudPeerConns() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudPeerConnsRead,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Description: "VPC ID where the peer connections located.",
				Optional:    true,
			},
			"peer_conn_id": {
				Type:        schema.TypeString,
				Description: "ID of the peer connection to retrieve.",
				Optional:    true,
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: "Output file for saving result.",
				Optional:    true,
				ForceNew:    true,
			},
			"filter": dataSourceFiltersSchema(),

			// Attributes used for result
			"peer_conns": {
				Type:        schema.TypeList,
				Description: "The list of the peer connections.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"peer_conn_id": {
							Type:        schema.TypeString,
							Description: "ID of the peer connection.",
							Computed:    true,
						},
						"role": {
							Type:        schema.TypeString,
							Description: "Role of the peer connection.",
							Computed:    true,
						},
						"status": {
							Type:        schema.TypeString,
							Description: "Status of the peer connection.",
							Computed:    true,
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
						"local_if_id": {
							Type:        schema.TypeString,
							Description: "Local interface ID of the peer connection.",
							Computed:    true,
						},
						"local_if_name": {
							Type:        schema.TypeString,
							Description: "Local interface name of the peer connection.",
							Computed:    true,
						},
						"local_vpc_id": {
							Type:        schema.TypeString,
							Description: "Local VPC ID of the peer connection.",
							Computed:    true,
						},
						"local_region": {
							Type:        schema.TypeString,
							Description: "Local region of the peer connection.",
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
						"peer_account_id": {
							Type:        schema.TypeString,
							Description: "Peer account ID of the peer connection.",
							Computed:    true,
						},
						"dns_status": {
							Type:        schema.TypeString,
							Description: "DNS status of the peer connection.",
							Computed:    true,
						},
						"payment_timing": {
							Type:        schema.TypeString,
							Description: "Payment timing of the peer connection.",
							Computed:    true,
						},
						"created_time": {
							Type:        schema.TypeString,
							Description: "Created time of the peer connection.",
							Computed:    true,
						},
						"expired_time": {
							Type:        schema.TypeString,
							Description: "Expired time of the peer connection.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudPeerConnsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	vpcService := &VpcService{client}

	var (
		vpcID      string
		peerConnID string
		outputFile string
	)
	if v := d.Get("vpc_id").(string); v != "" {
		vpcID = v
	}
	if v := d.Get("peer_conn_id").(string); v != "" {
		peerConnID = v
	}
	if v := d.Get("output_file").(string); v != "" {
		outputFile = v
	}

	action := "Query Peer Conns " + vpcID + "_" + peerConnID

	pcsResult := make([]map[string]interface{}, 0)

	if peerConnID != "" {
		pc, err := vpcService.GetPeerConnDetail(peerConnID, vpc.PEERCONN_ROLE_INITIATOR)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_peer_conns", action, BCESDKGoERROR)
		}

		pcMap := flattenPeerConnToMap(pc)
		pcsResult = append(pcsResult, pcMap)
	} else {
		pcs, err := vpcService.ListAllPeerConns(vpcID)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_peer_conns", action, BCESDKGoERROR)
		}

		for _, pc := range pcs {
			pcMap := flattenPeerConnToMap(&pc)
			pcsResult = append(pcsResult, pcMap)
		}
	}

	FilterDataSourceResult(d, &pcsResult)
	d.Set("peer_conns", pcsResult)

	d.SetId(resource.UniqueId())

	if outputFile != "" {
		if err := writeToFile(outputFile, pcsResult); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_peer_conns", action, BCESDKGoERROR)
		}
	}

	return nil
}

func flattenPeerConnToMap(pc *vpc.PeerConn) map[string]interface{} {
	pcMap := make(map[string]interface{})

	pcMap["peer_conn_id"] = pc.PeerConnId
	pcMap["role"] = string(pc.Role)
	pcMap["status"] = string(pc.Status)
	pcMap["bandwidth_in_mbps"] = pc.BandwidthInMbps
	pcMap["description"] = pc.Description
	pcMap["local_if_id"] = pc.LocalIfId
	pcMap["local_if_name"] = pc.LocalIfName
	pcMap["local_vpc_id"] = pc.LocalVpcId
	pcMap["local_region"] = pc.LocalRegion
	pcMap["peer_vpc_id"] = pc.PeerVpcId
	pcMap["peer_region"] = pc.PeerRegion
	pcMap["peer_account_id"] = pc.PeerAccountId
	pcMap["dns_status"] = pc.DnsStatus
	pcMap["payment_timing"] = pc.PaymentTiming
	pcMap["created_time"] = pc.CreatedTime
	pcMap["expired_time"] = pc.ExpiredTime

	return pcMap
}
