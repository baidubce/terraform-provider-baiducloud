package bec

import (
	"fmt"
	"log"

	"github.com/baidubce/bce-sdk-go/services/bec"
	"github.com/baidubce/bce-sdk-go/services/bec/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func DataSourceVMInstances() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to query BEC VM instance list. \n\n",

		Read: dataSourceVMInstancesRead,

		Schema: map[string]*schema.Schema{
			"keyword_type": {
				Type:         schema.TypeString,
				Description:  "Filter type. Valid values: `instanceId`, `serviceId`, `instanceName`, `instanceIp`, `securityGroupId`, `deploysetId`",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"instanceId", "serviceId", "instanceName", "instanceIp", "securityGroupId", "deploysetId"}, false),
			},
			"keyword": {
				Type:        schema.TypeString,
				Description: "Filter keyword.",
				Optional:    true,
			},
			"vm_instances": {
				Type:        schema.TypeList,
				Description: "Filtered VM instance list",
				Computed:    true,
				Elem:        vmDetailSchema(),
			},
		},
	}
}

func vmDetailSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:        schema.TypeString,
				Description: "ID of the vm instance group.",
				Computed:    true,
			},
			"vm_id": {
				Type:        schema.TypeString,
				Description: "ID of the vm instance.",
				Computed:    true,
			},
			"vm_name": {
				Type:        schema.TypeString,
				Description: "Name of the vm instance.",
				Computed:    true,
			},
			"host_name": {
				Type:        schema.TypeString,
				Description: "Host name of the vm instance.",
				Computed:    true,
			},
			"region_id": {
				Type:        schema.TypeString,
				Description: "Node ID.",
				Computed:    true,
			},
			"spec": {
				Type:        schema.TypeString,
				Description: "Specification family of the vm instance.",
				Computed:    true,
			},
			"cpu": {
				Type:        schema.TypeInt,
				Description: "CPU core count of the vm instance.",
				Computed:    true,
			},
			"memory": {
				Type:        schema.TypeInt,
				Description: "Memory size (GB) of the vm instance.",
				Computed:    true,
			},
			"image_type": {
				Type:        schema.TypeString,
				Description: "Possible values: `bec`(public image or bec custom image), `bcc`(bcc custom image)",
				Computed:    true,
			},
			"image_id": {
				Type:        schema.TypeString,
				Description: "ID of the image.",
				Computed:    true,
			},
			"system_volume": {
				Type:        schema.TypeList,
				Description: "System volume config of the vm instance.",
				Computed:    true,
				Elem:        VolumeConfigReadSchema(),
			},
			"data_volume": {
				Type:        schema.TypeList,
				Description: "Data volume config of the vm instance.",
				Computed:    true,
				Elem:        VolumeConfigReadSchema(),
			},
			"need_public_ip": {
				Type:        schema.TypeBool,
				Description: "Whether public network is enabled.",
				Computed:    true,
			},
			"need_ipv6_public_ip": {
				Type:        schema.TypeBool,
				Description: "Whether IPv6 public network is enabled.",
				Computed:    true,
			},
			"bandwidth": {
				Type:        schema.TypeInt,
				Description: "Public network bandwidth size (Mbps).",
				Computed:    true,
			},
			"dns_config": {
				Type:        schema.TypeList,
				Description: "DNS config of the vm instance.",
				Computed:    true,
				Elem:        DNSConfigReadSchema(),
			},
			"status": {
				Type: schema.TypeString,
				Description: "Status of the vm instance. Possible values: `CREATING`, `RUNNING`, `STOPPING`, `STOPPED`, `RESTARTING`, " +
					"`REINSTALLING`, `STARTING`, `IMAGING`, `FAILED`, `UNKNOWN`",
				Computed: true,
			},
			"internal_ip": {
				Type:        schema.TypeString,
				Description: "Local network IPv4 address of the vm instance.",
				Computed:    true,
			},
			"public_ip": {
				Type:        schema.TypeString,
				Description: "Public network IPv4 address of the vm instance.",
				Computed:    true,
			},
			"ipv6_public_ip": {
				Type:        schema.TypeString,
				Description: "Public network IPv6 address of the vm instance.",
				Computed:    true,
			},
			"create_time": {
				Type:        schema.TypeString,
				Description: "Creation time of the vm instance.",
				Computed:    true,
			},
		},
	}
}

func VolumeConfigReadSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the disk.",
				Computed:    true,
			},
			"size_in_gb": {
				Type:        schema.TypeInt,
				Description: "Size (GB) of the disk.",
				Computed:    true,
			},
			"volume_type": {
				Type:        schema.TypeString,
				Description: "Type of the disk. Possible values: `NVME`(SSD), `SATA`(HDD).",
				Computed:    true,
			},
			"pvc_name": {
				Type:        schema.TypeString,
				Description: "PVC name of the disk.",
				Computed:    true,
			},
		},
	}
}

func DNSConfigReadSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"dns_type": {
				Type: schema.TypeString,
				Description: "DNS type. Valid values: `NONE`(no dns config), `DEFAULT`(114.114.114.114 for domestic nodes, 8.8.8.8 for overseas nodes), " +
					"`LOCAL`(local dns of node), `CUSTOMIZE`.",
				Computed: true,
			},
			"dns_address": {
				Type:        schema.TypeList,
				Description: "Custom DNS address.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func NetworkConfigListReadSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"node_type": {
				Type: schema.TypeString,
				Description: "Node type,  Valid values: `SINGLE` (single network) , `TRIPLE` (triple network)." +
					" SINGLE, TRIPLE, representing single line and triple line nodes respectively",
				Optional: true,
				Computed: true,
			},
			"networks": {
				Type: schema.TypeList,
				Description: "Network card information, including Networks;" +
					"Used to set the corresponding network card name." +
					"Note: The classic network order is fixed as inside+outside and cannot be customized; " +
					"VPC network can customize the order of network cards",
				Optional: true,
				Computed: true,
				Elem:     NetworkListReadSchema(),
			},
		},
	}
}

func NetworkListReadSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"net_type": {
				Type: schema.TypeString,
				Description: "Network type. Network card type (INTERNAL_IP represents internal network;" +
					" Single line nodes use PULIC_SP to represent single line public network cards; " +
					"Three line nodes are represented by TRIPLE_CT, TRIPLE-UN, and TRIPLE_CM " +
					"respectively for telecommunications, China Unicom, and China Mobile on the three line public network",
				Optional: true,
				Computed: true,
			},
			"net_name": {
				Type: schema.TypeString,
				Description: "Network name. Network card name, the names of internal and external network cards " +
					"cannot be duplicated; Supports uppercase and lowercase letters, numbers," +
					" \"- _\", must start with a letter, and have a length of 3-16 characters.",
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceVMInstancesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	result := []api.VmInstanceDetailsVo{}
	pageNo := 0
	for {
		raw, err := conn.WithBECClient(func(client *bec.Client) (interface{}, error) {
			return client.GetVmInstanceList(buildListArgs(d, pageNo))
		})
		log.Printf("[DEBUG] Read VM Instances result: %+v", raw)
		if err != nil {
			return fmt.Errorf("error reading vm instance list: %w", err)
		}
		response := raw.(*api.LogicPageVmInstanceResult)
		result = append(result, response.Result...)

		if len(result) < response.TotalCount {
			pageNo += 1
		} else {
			break
		}
	}

	if err := d.Set("vm_instances", flattenVMInstances(result)); err != nil {
		return fmt.Errorf("error setting vm_instances: %w", err)
	}
	d.SetId(resource.UniqueId())
	return nil
}

func buildListArgs(d *schema.ResourceData, pageNo int) *api.ListRequest {
	return &api.ListRequest{
		KeywordType: d.Get("keyword_type").(string),
		Keyword:     d.Get("keyword").(string),
		PageNo:      pageNo,
		PageSize:    1000,
	}
}
