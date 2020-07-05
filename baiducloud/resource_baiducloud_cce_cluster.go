/*
Use this resource to get information about a CCE Cluster.

~> **NOTE:** The terminate operation of cce does NOT take effect immediately，maybe takes for several minites.

Example Usage

```hcl
resource "baiducloud_cce_cluster" "my-cluster" {
  cluster_name        = "test-cce-cluster"
  main_available_zone = "zoneA"
  container_net       = "172.16.0.0/16"
  deploy_mode		  = "BCC"
  master_config {
    instance_type     = "10"
    cpu               = 4
    memory            = 8
    image_type        = "common"
    logical_zone      = "zoneA"
    subnet_uuid       = baiducloud_subnet.defaultA.id
    security_group_id = baiducloud_security_group.defualt.id
    product_type      = "postpay"
    image_id          = data.baiducloud_images.default.images.0.id
  }
  worker_config {
    count = {
      "zoneA" : 2
    }
    instance_type = "10"
    cpu           = 1
    memory        = 2
    subnet_uuid = {
      "zoneA" : baiducloud_subnet.defaultA.id
      "zoneB" : baiducloud_subnet.defaultB.id
    }
    security_group_id = baiducloud_security_group.defualt.id
    product_type      = "postpay"
    image_id          = data.baiducloud_images.default.images.0.id
  }
}
```
*/
package baiducloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/cce"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudCCECluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudCCEClusterCreate,
		Read:   resourceBaiduCloudCCEClusterRead,
		Update: resourceBaiduCloudCCEClusterUpdate,
		Delete: resourceBaiduCloudCCEClusterDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:        schema.TypeString,
				Description: "Name of the Cluster. Support for uppercase and lowercase letters, numbers, Chinese and special characters, such as \"-\",\"_\",\"/\",\".\", the value must start with a letter, length 1-65.",
				Required:    true,
				ForceNew:    true,
			},
			"main_available_zone": {
				Type:        schema.TypeString,
				Description: "Main available zone of the cce cluster, support zoneA, zoneB, etc.",
				Optional:    true,
				ForceNew:    true,
			},
			"version": {
				Type:        schema.TypeString,
				Description: "Kubernetes version of the cce cluster.",
				Optional:    true,
				ForceNew:    true,
			},
			"container_net": {
				Type:        schema.TypeString,
				Description: "Container network type of the cce cluster.",
				Required:    true,
				ForceNew:    true,
			},
			"advanced_options": {
				Type:        schema.TypeMap,
				Description: "Advanced options configuration of the cce cluster.",
				Optional:    true,
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kube_proxy_mode": {
							Type:         schema.TypeString,
							Description:  "Mode of kube-proxy, which can only be iptables or ipvs.",
							Optional:     true,
							ForceNew:     true,
							Default:      string(cce.KubeProxyModeIpvs),
							ValidateFunc: validation.StringInSlice([]string{string(cce.KubeProxyModeIptables), string(cce.KubeProxyModeIpvs)}, false),
						},
						"cni_mode": {
							Type:         schema.TypeString,
							Description:  "Mode of the container network interface, which can only be cni or kubenet.",
							Optional:     true,
							ForceNew:     true,
							Default:      string(cce.CniModeKubenet),
							ValidateFunc: validation.StringInSlice([]string{string(cce.CniModeCni), string(cce.CniModeKubenet)}, false),
						},
						"cni_type": {
							Type:        schema.TypeString,
							Description: "Type of the container network interface, which can be VPC_ROUTE_AUTODETECT, VPC_SECONDARY_IP_VETH.",
							Optional:    true,
							ForceNew:    true,
							Default:     string(cce.CniTypeEmpty),
							ValidateFunc: validation.StringInSlice([]string{
								string(cce.CniTypeEmpty),
								string(cce.CniTypeRouteAutoDetect),
								string(cce.CniTypeSecondaryIpVeth),
							}, false),
						},
						"dns_mode": {
							Type:         schema.TypeString,
							Description:  "Mode of the dns, which can be coreDNS or kubeDNS.",
							Optional:     true,
							ForceNew:     true,
							Default:      string(cce.DNSModeCoreDNS),
							ValidateFunc: validation.StringInSlice([]string{string(cce.DNSModeCoreDNS), string(cce.DNSModeKubeDNS)}, false),
						},
						"max_pod_num": {
							Type:        schema.TypeString,
							Description: "Maximum number of pods in a node.",
							Optional:    true,
							ForceNew:    true,
							Default:     "256",
							ValidateFunc: func(i interface{}, s string) (strings []string, errors []error) {
								value := i.(string)

								if _, err := strconv.Atoi(value); err != nil {
									errors = append(errors, fmt.Errorf(
										"%q convert to int failed with error: %s",
										s, err))
								}

								return
							},
						},
					},
				},
			},
			"comment": {
				Type:        schema.TypeString,
				Description: "Comment information of the cce cluster.",
				Optional:    true,
				ForceNew:    true,
			},
			"deploy_mode": {
				Type:         schema.TypeString,
				Description:  "Deployment mode of the cce cluster, which can only be BCC.",
				Optional:     true,
				ForceNew:     true,
				Default:      string(cce.DeployModeBcc),
				ValidateFunc: validation.StringInSlice([]string{string(cce.DeployModeBcc)}, false),
			},
			"worker_config": {
				Type:        schema.TypeList,
				Description: "Worker node config of the cce cluster.",
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": {
							Type:        schema.TypeMap,
							Description: "Count of the worker node.",
							Required:    true,
							Elem: &schema.Schema{
								Type:         schema.TypeInt,
								MinItems:     1,
								ValidateFunc: validation.IntAtLeast(1),
							},
						},
						"instance_type": {
							Type:        schema.TypeString,
							Description: "Instance type of the worker node.",
							Required:    true,
							ForceNew:    true,
						},
						"gpu_card": {
							Type:        schema.TypeString,
							Description: "Gpu card of the worker node.",
							Optional:    true,
							ForceNew:    true,
						},
						"gpu_count": {
							Type:        schema.TypeInt,
							Description: "Gpu count of the worker node.",
							Optional:    true,
							ForceNew:    true,
						},
						"cpu": {
							Type:        schema.TypeInt,
							Description: "Number of cpu cores for the worker node.",
							Required:    true,
							ForceNew:    true,
						},
						"memory": {
							Type:        schema.TypeInt,
							Description: "Memory size of the worker node.",
							Required:    true,
							ForceNew:    true,
						},
						"image_type": {
							Type:        schema.TypeString,
							Description: "Image type of the worker node, which can be common, custom, gpuBccImage, gpuBccCustom, sharing.",
							Required:    true,
							ForceNew:    true,
							ValidateFunc: validation.StringInSlice([]string{
								string(cce.ImageTypeCommon),
								string(cce.ImageTypeCustom),
								string(cce.ImageTypeGpu),
								string(cce.ImageTypeGpuCustom),
								string(cce.ImageTypeSharing),
							}, false),
						},
						"subnet_uuid": {
							Type:        schema.TypeMap,
							Description: "Subnet uuid of the worker node.",
							Required:    true,
							Elem: &schema.Schema{
								Type:     schema.TypeString,
								MinItems: 1,
							},
						},
						"security_group_id": {
							Type:        schema.TypeString,
							Description: "ID of the security group.",
							Required:    true,
							ForceNew:    true,
						},
						"admin_pass": {
							Type:        schema.TypeString,
							Description: "Password of the worker node.",
							Optional:    true,
							ForceNew:    true,
						},
						"root_disk_storage_type": {
							Type:        schema.TypeString,
							Description: "System disk storage type of the worker node.",
							Optional:    true,
							ForceNew:    true,
							Default:     cce.VolumeTypePremiumSsd,
							ValidateFunc: validation.StringInSlice([]string{
								string(cce.VolumeTypeSata),
								string(cce.VolumeTypeSsd),
								string(cce.VolumeTypePremiumSsd),
							}, false),
						},
						"root_disk_size_in_gb": {
							Type:        schema.TypeInt,
							Description: "System disk size(GB) of the worker node.",
							Optional:    true,
							ForceNew:    true,
							Default:     40,
						},
						"product_type": {
							Type:         schema.TypeString,
							Description:  "Product type of the worker node, which can be postpay or prepay.",
							Optional:     true,
							ForceNew:     true,
							Default:      string(cce.ProductTypePostpay),
							ValidateFunc: validation.StringInSlice([]string{string(cce.ProductTypePostpay), string(cce.ProductTypePrepay)}, false),
						},
						"purchase_length": {
							Type:        schema.TypeInt,
							Description: "Purchase duration of the worker node.",
							Optional:    true,
							ForceNew:    true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return d.Get("product_type").(string) == "postpay"
							},
						},
						"auto_renew_time_unit": {
							Type:        schema.TypeString,
							Description: "Time unit of automatic renewal, the default value is month, It is valid only when the product_type is prepay and auto_renew is true.",
							Optional:    true,
							ForceNew:    true,
							Default:     "month",
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								if d.Get("worker_config.0.product_type").(string) != "postpay" {
									return false
								}

								if d.Get("worker_config.0.auto_renew").(bool) {
									return false
								}

								return true
							},
						},
						"auto_renew_time": {
							Type:        schema.TypeInt,
							Description: "The time length of automatic renewal. It is valid only when the product_type is prepay and auto_renew is true.",
							Optional:    true,
							ForceNew:    true,
							Default:     0,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return d.Get("worker_config.0.product_type").(string) == "postpay" || !d.Get("worker_config.0.auto_renew").(bool)
							},
						},
						"auto_renew": {
							Type:        schema.TypeBool,
							Description: "Whether the worker is automatically renewed.",
							Optional:    true,
							ForceNew:    true,
							Default:     false,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return d.Get("worker_config.0.product_type").(string) == "postpay"
							},
						},
						"image_id": {
							Type:        schema.TypeString,
							Description: "Image id of the worker node.",
							Required:    true,
							ForceNew:    true,
						},
						"cds_disks": {
							Type:        schema.TypeList,
							Description: "CDS disks of the worker node.",
							Optional:    true,
							ForceNew:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"volume_type": {
										Type:        schema.TypeString,
										Description: "Types of disk，available values: CLOUD_PREMIUM and CLOUD_SSD.",
										ForceNew:    true,
										Required:    true,
									},
									"disk_size_in_gb": {
										Type:        schema.TypeInt,
										Description: "Volume of disk in GB. Default is 0.",
										ForceNew:    true,
										Required:    true,
									},
									"snapshot_id": {
										Description: "Data disk snapshot ID.",
										Type:        schema.TypeString,
										ForceNew:    true,
										Optional:    true,
									},
								},
							},
						},
						"eip": {
							Type:        schema.TypeMap,
							Description: "Eip of the worker node.",
							Optional:    true,
							ForceNew:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bandwidth_in_mbps": {
										Type:        schema.TypeInt,
										Description: "Eip bandwidth(Mbps) of the worker node.",
										Required:    true,
										ForceNew:    true,
									},
									"sub_product_type": {
										Type:         schema.TypeString,
										Description:  "Eip product type of the worker node, which can be bandwidth or netraffic.",
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validation.StringInSlice([]string{string(cce.EipTypeBandwidth), string(cce.EipTypeNetraffic)}, false),
									},
									"eip_name": {
										Type:        schema.TypeString,
										Description: "Eip name of the worker node.",
										Optional:    true,
										ForceNew:    true,
									},
								},
							},
						},
					},
				},
			},
			"master_config": {
				Type:        schema.TypeList,
				Description: "Master config of the cce cluster.",
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"logical_zone": {
							Type:        schema.TypeString,
							Description: "Logical zone of the master node.",
							Required:    true,
							ForceNew:    true,
						},
						"instance_type": {
							Type:        schema.TypeString,
							Description: "Instance type of the master node.",
							Required:    true,
							ForceNew:    true,
						},
						"gpu_card": {
							Type:        schema.TypeString,
							Description: "Gpu card of the master node.",
							Optional:    true,
							ForceNew:    true,
						},
						"gpu_count": {
							Type:        schema.TypeInt,
							Description: "Count of gpu card.",
							Optional:    true,
							ForceNew:    true,
						},
						"cpu": {
							Type:        schema.TypeInt,
							Description: "Number of cpu cores.",
							Required:    true,
							ForceNew:    true,
						},
						"memory": {
							Type:        schema.TypeInt,
							Description: "Memory capacity(GB) of the master node.",
							Required:    true,
							ForceNew:    true,
						},
						"image_type": {
							Type:        schema.TypeString,
							Description: "Image type of the master node.",
							Required:    true,
							ForceNew:    true,
						},
						"subnet_uuid": {
							Type:        schema.TypeString,
							Description: "Subnet uuid of the master node.",
							Required:    true,
							ForceNew:    true,
						},
						"security_group_id": {
							Type:        schema.TypeString,
							Description: "ID of the security group.",
							Required:    true,
							ForceNew:    true,
						},
						"admin_pass": {
							Type:        schema.TypeString,
							Description: "Password of the worker node.",
							Optional:    true,
							ForceNew:    true,
						},
						"root_disk_storage_type": {
							Type:        schema.TypeString,
							Description: "System disk storage type of the master node.",
							Optional:    true,
							ForceNew:    true,
							Default:     string(cce.VolumeTypePremiumSsd),
							ValidateFunc: validation.StringInSlice([]string{
								string(cce.VolumeTypeSata),
								string(cce.VolumeTypeSsd),
								string(cce.VolumeTypePremiumSsd),
							}, false),
						},
						"root_disk_size_in_gb": {
							Type:        schema.TypeInt,
							Description: "System disk size(GB) of the master node.",
							Optional:    true,
							ForceNew:    true,
							Default:     40,
						},
						"product_type": {
							Type:         schema.TypeString,
							Description:  "Product type of the master node, which can be postpay or prepay.",
							Optional:     true,
							ForceNew:     true,
							Default:      string(cce.ProductTypePostpay),
							ValidateFunc: validation.StringInSlice([]string{string(cce.ProductTypePostpay), string(cce.ProductTypePrepay)}, false),
						},
						"purchase_length": {
							Type:        schema.TypeInt,
							Description: "Purchase duration of the master node.",
							Optional:    true,
							ForceNew:    true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return d.Get("master_config.0.product_type").(string) == "postpay"
							},
						},
						"auto_renew_time_unit": {
							Type:        schema.TypeString,
							Description: "Time unit of automatic renewal, the default value is month, It is valid only when the product_type is prepay and auto_renew is true.",
							Optional:    true,
							ForceNew:    true,
							Default:     "month",
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return d.Get("master_config.0.product_type").(string) == "postpay" || !d.Get("master_config.0.auto_renew").(bool)
							},
						},
						"auto_renew_time": {
							Type:        schema.TypeInt,
							Description: "The time length of automatic renewal. It is valid only when the product_type is prepay and auto_renew is true.",
							Optional:    true,
							ForceNew:    true,
							Default:     0,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return d.Get("master_config.0.product_type").(string) == "postpay" || !d.Get("master_config.0.auto_renew").(bool)
							},
						},
						"auto_renew": {
							Type:        schema.TypeBool,
							Description: "Whether the master is automatically renewed.",
							Optional:    true,
							ForceNew:    true,
							Default:     false,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return d.Get("master_config.0.product_type").(string) == "postpay"
							},
						},
						"image_id": {
							Type:        schema.TypeString,
							Description: "Image id of the master node.",
							Required:    true,
							ForceNew:    true,
						},
					},
				},
			},
			"delete_eip_cds": {
				Type:        schema.TypeBool,
				Description: "Whether to delete the eip and cds, default to true.",
				Optional:    true,
				Default:     true,
			},
			"delete_snapshots": {
				Type:        schema.TypeBool,
				Description: "Whether to delete the snapshots, default to true.",
				Optional:    true,
				Default:     true,
			},
			"create_start_time": {
				Type:        schema.TypeString,
				Description: "Create time of the cce cluster.",
				Computed:    true,
			},
			"delete_time": {
				Type:        schema.TypeString,
				Description: "Delete time of the cce cluster.",
				Computed:    true,
			},
			"instance_mode": {
				Type:        schema.TypeString,
				Description: "Instance mode of the cce cluster.",
				Computed:    true,
			},
			"has_prepay": {
				Type:        schema.TypeBool,
				Description: "Whether to include prepaid nodes.",
				Computed:    true,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Description: "VPC id of the cce cluster.",
				Computed:    true,
			},
			"vpc_uuid": {
				Type:        schema.TypeString,
				Description: "VPC uuid of the cce cluster.",
				Computed:    true,
			},
			"vpc_cidr": {
				Type:        schema.TypeString,
				Description: "VPC cidr of the cce cluster.",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Status of the cce cluster.",
				Computed:    true,
			},
			"master_zone_subnet_map": {
				Type:        schema.TypeMap,
				Description: "Availability zone of master node.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cluster_uuid": {
				Type:        schema.TypeString,
				Description: "UUID of cce cluster.",
				Computed:    true,
			},
			"region": {
				Type:        schema.TypeString,
				Description: "Region of the cce cluster.",
				Computed:    true,
			},
			"master_vm_count": {
				Type:        schema.TypeInt,
				Description: "Number of virtual machines in the master node of the cce cluster.",
				Computed:    true,
			},
			"zone_subnet_map": {
				Type:        schema.TypeMap,
				Description: "Subnet of the zone.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cluster_node_num": {
				Type:        schema.TypeInt,
				Description: "Number of nodes in the cluster.",
				Computed:    true,
			},
			"worker_instances_list": {
				Type:        schema.TypeList,
				Description: "List of the worker instances.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Description: "ID of the instance.",
							Computed:    true,
						},
						"available_zone": {
							Type:        schema.TypeString,
							Description: "Available zone of the instance.",
							Computed:    true,
						},
						"eip": {
							Type:        schema.TypeString,
							Description: "Eip of the instance.",
							Computed:    true,
						},
						"status": {
							Type:        schema.TypeString,
							Description: "Status of the instance.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func resourceBaiduCloudCCEClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	cceService := CceService{client}

	createClusterArgs, err := buildBaiduCloudCCEClusterArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}

	action := "Create CCE cluster " + createClusterArgs.ClusterName
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithCCEClient(func(client *cce.Client) (interface{}, error) {
			return client.CreateCluster(createClusterArgs)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw)
		response, _ := raw.(*cce.CreateClusterResult)
		d.SetId(response.ClusterUuid)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cce_cluster", action, BCESDKGoERROR)
	}

	stateConf := buildStateConf(
		[]string{string(cce.ClusterStatusCreating)},
		[]string{string(cce.ClusterStatusRunning)},
		d.Timeout(schema.TimeoutCreate),
		cceService.ClusterStateRefresh(d.Id(), []cce.ClusterStatus{
			cce.ClusterStatusCreateFailed,
			cce.ClusterStatusError,
			cce.ClusterStatusMasterUpgradeFailed,
			cce.ClusterStatusDeleting,
			cce.ClusterStatusDeleted,
		}),
	)
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cce_cluster", action, BCESDKGoERROR)
	}

	return resourceBaiduCloudCCEClusterRead(d, meta)
}

func resourceBaiduCloudCCEClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	clusterId := d.Id()
	action := "Get CCE Cluster " + clusterId

	raw, err := client.WithCCEClient(func(client *cce.Client) (i interface{}, e error) {
		return client.GetCluster(clusterId)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cce_cluster", action, BCESDKGoERROR)
	}

	response := raw.(*cce.GetClusterResult)
	d.Set("cluster_uuid", response.ClusterUuid)
	d.Set("cluster_name", response.ClusterName)
	d.Set("version", response.Version)
	d.Set("region", response.Region)
	d.Set("master_vm_count", response.MasterVmCount)
	d.Set("vpc_id", response.VpcId)
	d.Set("vpc_uuid", response.VpcUuid)
	d.Set("vpc_name", response.VpcName)
	d.Set("vpc_cidr", response.VpcCidr)
	d.Set("container_name", response.ContainerNet)
	d.Set("status", response.Status)
	d.Set("create_start_time", response.CreateStartTime.String())
	d.Set("delete_time", response.DeleteTime.String())
	d.Set("comment", response.Comment)
	d.Set("instance_mode", response.InstanceMode)
	d.Set("has_prepay", response.HasPrepay)
	d.Set("master_zone_subnet_map", response.MasterZoneSubnetMap)

	action = "Get CCE Cluster " + clusterId + " Node list"
	args := &cce.ListNodeArgs{ClusterUuid: clusterId}
	rawList, errList := client.WithCCEClient(func(client *cce.Client) (i interface{}, e error) {
		return client.ListNodes(args)
	})

	if errList != nil {
		return WrapErrorf(errList, DefaultErrorMsg, "baiducloud_cce_cluster", action, BCESDKGoERROR)
	}

	responseList := rawList.(*cce.ListNodeResult)
	d.Set("cluster_node_num", len(responseList.Nodes))
	instanceList := make([]map[string]string, 0)
	for _, instance := range responseList.Nodes {
		instanceList = append(instanceList, map[string]string{
			"instance_id":    instance.InstanceShortId,
			"available_zone": instance.AvailableZone,
			"eip":            instance.Eip,
			"status":         instance.Status,
		})
	}
	if err := d.Set("worker_instances_list", instanceList); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cce_cluster", action, BCESDKGoERROR)
	}
	if err := d.Set("zone_subnet_map", response.ZoneSubnetMap); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cce_cluster", action, BCESDKGoERROR)
	}

	return nil
}

func resourceBaiduCloudCCEClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	if !d.HasChange("worker_config.0.count") {
		return nil
	}

	clusterId := d.Id()
	action := "Update CCE Cluster " + clusterId
	scalingUpArgs, scalingDownArgs, err := buildScalingArgs(d)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cce_cluster", action, BCESDKGoERROR)
	}

	if scalingDownArgs != nil && len(scalingDownArgs.NodeInfo) > 0 {
		action = "Scaling down CCE Cluster " + clusterId
		err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			raw, err := client.WithCCEClient(func(cceClient *cce.Client) (interface{}, error) {
				return clusterId, cceClient.ScalingDown(scalingDownArgs)
			})
			if err != nil {
				if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, raw)
			return nil
		})
		if err != nil {
			if IsExceptedErrors(err, CceClusterNotFound) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cce_cluster", action, BCESDKGoERROR)
		}
	}

	if scalingUpArgs != nil && len(scalingUpArgs.OrderContent.Items) > 0 {
		action = "Scaling up CCE Cluster " + clusterId
		err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			raw, err := client.WithCCEClient(func(cceClient *cce.Client) (interface{}, error) {
				return cceClient.ScalingUp(scalingUpArgs)
			})
			if err != nil {
				if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, raw)
			return nil
		})
		if err != nil {
			if IsExceptedErrors(err, CceClusterNotFound) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cce_cluster", action, BCESDKGoERROR)
		}
	}

	return resourceBaiduCloudCCEClusterRead(d, meta)
}

func resourceBaiduCloudCCEClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	cceService := CceService{client}

	clusterId := d.Id()
	action := "Delete CCE Cluster " + clusterId

	args := &cce.DeleteClusterArgs{
		ClusterUuid: clusterId,
	}
	if v, ok := d.GetOk("delete_eip_cds"); ok {
		args.DeleteEipCds = v.(bool)
	}
	if v, ok := d.GetOk("delete_snapshots"); ok {
		args.DeleteSnap = v.(bool)
	}
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := client.WithCCEClient(func(cceClient *cce.Client) (interface{}, error) {
			return clusterId, cceClient.DeleteCluster(args)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw)
		return nil
	})
	if err != nil {
		if IsExceptedErrors(err, CceClusterNotFound) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cce_cluster", action, BCESDKGoERROR)
	}

	stateConf := buildStateConf(
		[]string{string(cce.ClusterStatusRunning),
			string(cce.ClusterStatusDeleting),
			string(cce.ClusterStatusCreateFailed),
			string(cce.ClusterStatusError),
			string(cce.ClusterStatusMasterUpgradeFailed),
		},
		[]string{"DELETED"},
		d.Timeout(schema.TimeoutDelete),
		cceService.ClusterStateRefresh(clusterId, []cce.ClusterStatus{
			cce.ClusterStatusMasterUpgrading,
			cce.ClusterStatusCreating,
		}),
	)
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cce_cluster", action, BCESDKGoERROR)
	}

	return nil
}

func buildBaiduCloudCCEClusterArgs(d *schema.ResourceData, meta interface{}) (*cce.CreateClusterArgs, error) {
	request := &cce.CreateClusterArgs{
		ClusterName:       d.Get("cluster_name").(string),
		Version:           d.Get("version").(string),
		MainAvailableZone: d.Get("main_available_zone").(string),
		ContainerNet:      d.Get("container_net").(string),
		DeployMode:        cce.DeployMode(d.Get("deploy_mode").(string)),
	}

	zoneCountMap := map[string]interface{}{}
	if value, ok := d.GetOk("worker_config"); ok {
		workConfig := value.([]interface{})

		// get zone count config
		configMap := workConfig[0].(map[string]interface{})
		zoneCount, ok := configMap["count"]
		if !ok {
			return nil, fmt.Errorf("worker_config.count required at leaset one")
		}
		zoneCountMap = zoneCount.(map[string]interface{})

		workerConfig := &cce.BaseCreateOrderRequestVo{}

		// get bcc cds config
		instanceCount := 0
		subnetMap := configMap["subnet_uuid"].(map[string]interface{})
		for zone, count := range zoneCountMap {
			subnetUuid, has := subnetMap[zone]
			if !has {
				return nil, fmt.Errorf("please set worker_config subnet_uuid in zone: %s", zone)
			}
			instanceCount += count.(int)
			items, err := buildBaiduCloudCCENodeBCCCDSArgs(configMap, count.(int), zone, subnetUuid.(string))
			if err != nil {
				return nil, fmt.Errorf("get worker bcc cds config args failed with error: %v", err)
			}

			workerConfig.Items = append(workerConfig.Items, items...)
		}

		// get eip config
		if _, ok := configMap["eip"]; ok {
			item, err := buildBaiduCloudCCENodeEipArgs(configMap, instanceCount)
			if err != nil {
				return nil, fmt.Errorf("get worker eip config args failed with error: %v", err)
			}

			if item != nil {
				workerConfig.Items = append(workerConfig.Items, *item)
			}
		}

		request.OrderContent = workerConfig
	}

	if value, ok := d.GetOk("master_config"); ok {
		masterConfig := value.([]interface{})

		configMap := masterConfig[0].(map[string]interface{})
		itemConfig := &cce.BaseCreateOrderRequestVo{}

		masterZone := configMap["logical_zone"].(string)
		subnetUuid := configMap["subnet_uuid"].(string)
		// get bcc cds
		items, err := buildBaiduCloudCCENodeBCCCDSArgs(configMap, 3, masterZone, subnetUuid)
		if err != nil {
			return nil, fmt.Errorf("get master config args failed with error: %v", err)
		}
		itemConfig.Items = append(itemConfig.Items, items...)

		request.MasterOrderContent = itemConfig
		request.MasterExposed = true
	}

	for _, item := range request.OrderContent.Items {
		if value, ok := item.Config.(*cce.CdsConfig); !ok {
			continue
		} else {
			request.CdsPreMountInfo = &cce.CdsPreMountInfo{
				CdsConfig: []cce.DiskSizeConfig{
					{
						Size:       value.CdsDiskSize[0].Size,
						VolumeType: value.CdsDiskSize[0].VolumeType,
					},
				},
				MountPath: "/data",
			}
			break
		}
	}

	if value, ok := d.GetOk("advanced_options"); ok {
		adVancedOptionsMap := value.(map[string]interface{})
		request.AdvancedOptions = &cce.AdvancedOptions{}

		if v, ok := adVancedOptionsMap["kube_proxy_mode"]; ok {
			request.AdvancedOptions.KubeProxyMode = cce.KubeProxyMode(v.(string))
		}

		if v, ok := adVancedOptionsMap["cni_mode"]; ok {
			request.AdvancedOptions.CniMode = cce.CniMode(v.(string))
		}

		if v, ok := adVancedOptionsMap["cni_type"]; ok {
			request.AdvancedOptions.CniType = cce.CniType(v.(string))
		}

		if v, ok := adVancedOptionsMap["dns_mode"]; ok {
			request.AdvancedOptions.DnsMode = cce.DNSMode(v.(string))
		}

		if v, ok := adVancedOptionsMap["max_pod_num"]; ok {
			request.AdvancedOptions.MaxPodNum, _ = strconv.Atoi(v.(string))
		}
	}

	return request, nil
}

func buildBaiduCloudCCENodeBCCCDSArgs(config map[string]interface{}, purchaseNum int, zone, subnetUuid string) ([]cce.Item, error) {
	bccConfig := cce.BccConfig{
		ProductType:     cce.ProductType(config["product_type"].(string)),
		InstanceType:    cce.InstanceType(config["instance_type"].(string)),
		Cpu:             config["cpu"].(int),
		Memory:          config["memory"].(int),
		ImageType:       cce.ImageType(config["image_type"].(string)),
		ImageId:         config["image_id"].(string),
		SecurityGroupId: config["security_group_id"].(string),
		PurchaseNum:     purchaseNum,
		ServiceType:     cce.ServiceTypeBCC,
		LogicalZone:     zone,
		SubnetUuid:      subnetUuid,
	}

	if bccConfig.ProductType == cce.ProductTypePrepay {
		value, ok := config["purchase_length"]
		if !ok {
			return nil, fmt.Errorf("purchase_length is needed if purchase_type is prepay")
		}
		bccConfig.PurchaseLength = value.(int)

		if v, ok := config["auto_renew"]; ok && v.(bool) {
			if value, ok := config["auto_renew_time_unit"]; ok {
				bccConfig.AutoRenewTimeUnit = value.(string)
			}
			if value, ok := config["auto_renew_time"]; ok {
				bccConfig.AutoRenewTime = value.(int)
			}
		}
	}

	if value, ok := config["name"]; ok {
		bccConfig.Name = value.(string)
	}

	if value, ok := config["gpu_card"]; ok {
		bccConfig.GpuCard = value.(string)
	}

	if value, ok := config["gpu_count"]; ok {
		bccConfig.GpuCount = value.(int)
	}

	if value, ok := config["admin_pass"]; ok {
		bccConfig.AdminPass = value.(string)
	}

	if value, ok := config["root_disk_size_in_gb"]; ok {
		bccConfig.RootDiskSizeInGb = value.(int)
	}

	if value, ok := config["root_disk_storage_type"]; ok {
		bccConfig.RootDiskStorageType = cce.VolumeType(value.(string))
	}

	result := []cce.Item{
		{Config: bccConfig},
	}

	if value, ok := config["cds_disks"]; ok && len(value.([]interface{})) > 0 {
		cdsConfig := &cce.CdsConfig{
			PurchaseNum: purchaseNum,
			LogicalZone: bccConfig.LogicalZone,
			ProductType: bccConfig.ProductType,
			ServiceType: cce.ServiceTypeCDS,
		}

		if cdsConfig.ProductType == cce.ProductTypePrepay {
			cdsConfig.PurchaseLength = bccConfig.PurchaseLength
			if bccConfig.AutoRenew {
				cdsConfig.AutoRenewTime = bccConfig.AutoRenewTime
				cdsConfig.AutoRenewTimeUnit = bccConfig.AutoRenewTimeUnit
			}
		}

		cdsConfig.CdsDiskSize = []cce.DiskSizeConfig{}
		cdsList := value.([]interface{})
		for _, cds := range cdsList {
			cdsMap := cds.(map[string]interface{})
			cdsSize := cce.DiskSizeConfig{
				Size:       strconv.Itoa(cdsMap["disk_size_in_gb"].(int)),
				VolumeType: cce.VolumeType(cdsMap["volume_type"].(string)),
			}

			if snp, ok := cdsMap["snapshot_id"]; ok {
				cdsSize.SnapshotId = snp.(string)
			}

			cdsConfig.CdsDiskSize = append(cdsConfig.CdsDiskSize, cdsSize)
		}

		result = append(result, cce.Item{Config: cdsConfig})
	}

	return result, nil
}

func buildBaiduCloudCCENodeEipArgs(config map[string]interface{}, purchaseNum int) (*cce.Item, error) {
	if purchaseNum < 1 {
		return nil, nil
	}
	eip, ok := config["eip"]
	if !ok {
		return nil, nil
	}
	if eip == nil {
		return nil, nil
	}

	eipMap := eip.(map[string]interface{})
	if len(eipMap) == 0 {
		return nil, nil
	}
	bandwidth, _ := strconv.Atoi(eipMap["bandwidth_in_mbps"].(string))
	eipConfig := &cce.EipConfig{
		ServiceType:     cce.ServiceTypeEIP,
		PurchaseNum:     purchaseNum,
		ProductType:     cce.ProductType(config["product_type"].(string)),
		BandwidthInMbps: bandwidth,
		SubProductType:  cce.EipType(eipMap["sub_product_type"].(string)),
	}

	if eipConfig.ProductType == cce.ProductTypePrepay {
		value, ok := config["purchase_length"]
		if !ok {
			return nil, fmt.Errorf("purchase_length is needed if purchase_type is prepay")
		}
		eipConfig.PurchaseLength = value.(int)

		if v, ok := config["auto_renew"]; ok && v.(bool) {
			if value, ok := config["auto_renew_time_unit"]; ok {
				eipConfig.AutoRenewTimeUnit = value.(string)
			}
			if value, ok := config["auto_renew_time"]; ok {
				eipConfig.AutoRenewTime = value.(int)
			}
		}
	}

	if name, ok := eipMap["name"]; ok {
		eipConfig.Name = name.(string)
	}

	return &cce.Item{Config: eipConfig}, nil
}

func buildScalingArgs(d *schema.ResourceData) (*cce.ScalingUpArgs, *cce.ScalingDownArgs, error) {
	clusterId := d.Id()
	o, n := d.GetChange("worker_config.0.count")
	oMap := o.(map[string]interface{})
	nMap := n.(map[string]interface{})
	scaleZoneCount := map[string]int{}

	for zone, nCount := range nMap {
		if oCount, ok := oMap[zone]; ok {
			scaleZoneCount[zone] = nCount.(int) - oCount.(int)
		} else {
			scaleZoneCount[zone] = nCount.(int)
		}
	}

	for zone, oCount := range oMap {
		if _, ok := nMap[zone]; !ok {
			scaleZoneCount[zone] = -(oCount.(int))
		}
	}

	// check subnetUuid
	subnetMap := map[string]interface{}{}
	if subnets, ok := d.GetOk("worker_config.0.subnet_uuid"); ok {
		subnetMap = subnets.(map[string]interface{})

		for zone, diffCount := range scaleZoneCount {
			if diffCount > 0 {
				// need scale up
				if subnet, ok := subnetMap[zone]; !ok || len(subnet.(string)) == 0 {
					return nil, nil, fmt.Errorf("please set subnet_uuid for worker node in zone: %s", zone)
				}
			}
		}
	}

	instanceList := d.Get("worker_instances_list").([]interface{})
	instanceListMap := map[string][]string{}
	for _, instance := range instanceList {
		instanceMap := instance.(map[string]interface{})
		zone := instanceMap["available_zone"].(string)
		instanceListMap[zone] = append(instanceListMap[zone], instanceMap["instance_id"].(string))
	}

	var configMap map[string]interface{}
	if value, ok := d.GetOk("worker_config"); ok {
		workConfig := value.([]interface{})

		// get zone count config
		configMap = workConfig[0].(map[string]interface{})
	} else {
		return nil, nil, fmt.Errorf("worker_config is required")
	}
	scaleUpArgs := &cce.ScalingUpArgs{
		ClusterUuid:  clusterId,
		OrderContent: &cce.BaseCreateOrderRequestVo{},
	}
	scaleDownArgs := &cce.ScalingDownArgs{ClusterUuid: clusterId}
	if v, ok := d.GetOk("delete_eip_cds"); ok {
		scaleDownArgs.DeleteEipCds = v.(bool)
	}
	if v, ok := d.GetOk("delete_snapshots"); ok {
		scaleDownArgs.DeleteSnap = v.(bool)
	}

	scaleUpEipCount := 0
	for zone, diffCount := range scaleZoneCount {
		if diffCount == 0 {
			continue
		}

		if diffCount < 0 {
			// scaling down
			instanceList := instanceListMap[zone]
			diffCount = -diffCount
			if diffCount > len(instanceList) {
				diffCount = len(instanceList)
			}

			for _, instance := range instanceList[:diffCount] {
				scaleDownArgs.NodeInfo = append(scaleDownArgs.NodeInfo, cce.NodeInfo{InstanceId: instance})
			}
		} else {
			// scaling up
			items, err := buildBaiduCloudCCENodeBCCCDSArgs(configMap, diffCount, zone, subnetMap[zone].(string))
			if err != nil {
				return nil, nil, fmt.Errorf("build scaling up order args with error: %s", zone)
			}
			scaleUpArgs.OrderContent.Items = append(scaleUpArgs.OrderContent.Items, items...)

			scaleUpEipCount += diffCount
		}
	}

	if scaleUpEipCount > 0 {
		item, err := buildBaiduCloudCCENodeEipArgs(configMap, scaleUpEipCount)
		if err != nil {
			return nil, nil, fmt.Errorf("build scaling up eip order args with error: %s", err)
		}

		if item != nil {
			scaleUpArgs.OrderContent.Items = append(scaleUpArgs.OrderContent.Items, *item)
		}
	}

	return scaleUpArgs, scaleDownArgs, nil
}
