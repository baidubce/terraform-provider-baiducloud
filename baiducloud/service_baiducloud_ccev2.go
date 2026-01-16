package baiducloud

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"
	"time"

	bccapi "github.com/baidubce/bce-sdk-go/services/bcc/api"
	ccev2 "github.com/baidubce/bce-sdk-go/services/cce/v2"
	ccev2types "github.com/baidubce/bce-sdk-go/services/cce/v2/types"
	"github.com/baidubce/bce-sdk-go/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

type Ccev2Service struct {
	client *connectivity.BaiduClient
}

func (s *Ccev2Service) ClusterStateRefreshCCEv2(clusterId string, failState []ccev2types.ClusterPhase) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		action := "Query CCEv2 Cluster " + clusterId
		raw, err := s.client.WithCCEv2Client(func(ccev2Client *ccev2.Client) (i interface{}, e error) {
			return ccev2Client.GetCluster(clusterId)
		})
		addDebug(action, raw)
		if err != nil {
			if NotFoundError(err) {
				return 0, string(ccev2types.ClusterPhaseDeleted), nil
			}
			return nil, "", WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_cluster", action, BCESDKGoERROR)
		}
		if raw == nil {
			return nil, "", nil
		}

		result := raw.(*ccev2.GetClusterResponse)
		for _, statue := range failState {
			if result.Cluster.Status.ClusterPhase == statue {
				return result, string(result.Cluster.Status.ClusterPhase), WrapError(Error(GetFailTargetStatus, result.Cluster.Status.ClusterPhase))
			}
		}

		addDebug(action, raw)
		return result, string(result.Cluster.Status.ClusterPhase), nil
	}
}

func (s *Ccev2Service) InstanceEventStepsStateRefresh(instanceId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		action := "Query CCE Instance Event Steps: " + instanceId
		raw, err := s.client.WithCCEv2Client(func(cceV2Client *ccev2.Client) (i interface{}, e error) {
			return cceV2Client.GetInstanceEventSteps(instanceId)
		})

		addDebug(action, raw)
		if err != nil {
			return nil, "", WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2", action, BCESDKGoERROR)
		}
		if raw == nil {
			return nil, "", nil
		}
		result, _ := raw.(*ccev2.GetEventStepsResponse)

		return result, result.Status, nil
	}
}

func (s *Ccev2Service) waitForInstancesOperation(pending []string, target []string, timeout time.Duration, instanceIds []string) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []string

	for _, instanceId := range instanceIds {
		wg.Add(1)

		go func(id string) {
			defer wg.Done()

			stateConf := buildStateConf(pending, target, timeout, s.InstanceEventStepsStateRefresh(id))
			if result, err := stateConf.WaitForState(); err != nil {
				eventStepsResp, _ := result.(*ccev2.GetEventStepsResponse)

				var messages []string
				for _, step := range eventStepsResp.Steps {
					if step.StepStatus == "failed" {
						errorInfo := step.StepInfo.ErrorInfo
						message := fmt.Sprintf("[%s]%s [message]: %s [code]: %s [traceId]: %s", step.StepName, errorInfo.Suggestion, errorInfo.Message,
							errorInfo.Code, errorInfo.TraceID)
						messages = append(messages, message)
					}
				}

				errWithMessage := WrapErrorf(err, "instance [%s] operation failed, reason: %s", id, strings.Join(messages, "; "))

				mu.Lock()
				errs = append(errs, errWithMessage.Error())
				mu.Unlock()
			}
		}(instanceId)
	}

	wg.Wait()
	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, "; "))
	}
	return nil
}

func resourceCCEv2ClusterSpec() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:        schema.TypeString,
				Description: "Cluster Name",
				Optional:    true,
			},
			"cluster_type": {
				Type:         schema.TypeString,
				Description:  "Cluster Type. Available Value: [normal].",
				Optional:     true,
				ValidateFunc: validation.StringInSlice(ClusterTypePermitted, false),
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Cluster Description",
				Optional:    true,
			},
			"k8s_version": {
				Type: schema.TypeString,
				Description: "Kubernetes Version. Available Value: [1.18.9, 1.20.8, 1.21.14, " +
					"1.22.5, 1.24.4, 1.26.9, 1.28.8, 1.30.1].",
				Optional:     true,
				ValidateFunc: validation.StringInSlice(K8SVersionPermitted, false),
			},
			"runtime_type": {
				Type:         schema.TypeString,
				Description:  "Container Runtime Type. Available Values: [docker, containerd].",
				Optional:     true,
				ValidateFunc: validation.StringInSlice(RuntimeTypePermitted, false),
			},
			"runtime_version": {
				Type:        schema.TypeString,
				Description: "Container Runtime Version",
				Optional:    true,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Description: "VPC ID",
				Optional:    true,
			},
			"vpc_cidr": {
				Type:        schema.TypeString,
				Description: "VPC CIDR",
				Optional:    true,
			},
			"vpc_cidr_ipv6": {
				Type:        schema.TypeString,
				Description: "VPC CIDR IPv6",
				Optional:    true,
			},
			"plugins": {
				Type:        schema.TypeList,
				Description: "Plugin List",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cluster_delete_option": {
				Type:        schema.TypeList,
				Description: "Cluster Delete Option",
				Optional:    true,
				MaxItems:    1,
				Elem:        resourceCCEv2ClusterDeleteOption(),
			},
			"container_network_config": {
				Type:        schema.TypeList,
				Description: "Container Network Config",
				Optional:    true,
				MaxItems:    1,
				Elem:        resourceCCEv2ContainerNetworkConfig(),
			},
			"master_config": {
				Type:        schema.TypeList,
				Description: "Cluster Master Config",
				Optional:    true,
				MaxItems:    1,
				Elem:        resourceCCEv2MasterConfig(),
			},
			"k8s_custom_config": {
				Type:        schema.TypeList,
				Description: "Cluster k8s custom config",
				Optional:    true,
				MaxItems:    1,
				Elem:        resourceCCEv2K8SCustomConfig(),
			},
		},
	}
}

func resourceCCEv2ClusterStatus() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cluster_blb": {
				Type:        schema.TypeList,
				Description: "Cluster BLB",
				Computed:    true,
				MaxItems:    1,
				Elem:        resourceCCEv2BLB(),
			},
			"cluster_phase": {
				Type:        schema.TypeString,
				Description: "Cluster Phase",
				Computed:    true,
			},
			"node_num": {
				Type:        schema.TypeInt,
				Description: "Cluster Node Number",
				Computed:    true,
			},
		},
	}
}

func resourceCCEv2BLB() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "BLB ID",
				Optional:    true,
			},
			"vpc_ip": {
				Type:        schema.TypeString,
				Description: "VPC IP",
				Optional:    true,
			},
			"eip": {
				Type:        schema.TypeString,
				Description: "EIP",
				Optional:    true,
			},
		},
	}
}

func resourceCCEv2InstanceSpec() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cce_instance_id": {
				Type:        schema.TypeString,
				Description: "Instance ID",
				Optional:    true,
				Computed:    true,
			},
			"instance_name": {
				Type:        schema.TypeString,
				Description: "Instance Name",
				Optional:    true,
				Computed:    true,
			},
			"runtime_type": {
				Type:         schema.TypeString,
				Description:  "Container Runtime Type. Available Value: [docker].",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(RuntimeTypePermitted, false),
			},
			"runtime_version": {
				Type:        schema.TypeString,
				Description: "Container Runtime Version",
				Optional:    true,
				Computed:    true,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "Cluster ID of this Instance",
				Optional:    true,
				Computed:    true,
			},
			"cluster_role": {
				Type:         schema.TypeString,
				Description:  "Cluster Role of Instance, Master or Nodes. Available Value: [master, node].",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(ClusterRolePermitted, false),
			},
			"instance_group_id": {
				Type:        schema.TypeString,
				Description: "Instance Group ID of this Instance",
				Optional:    true,
				Computed:    true,
			},
			"instance_group_name": {
				Type:        schema.TypeString,
				Description: "Name of Instance Group",
				Optional:    true,
				Computed:    true,
			},
			"master_type": {
				Type:         schema.TypeString,
				Description:  "Master Type. Available Value: [managed, custom, serverless].",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(MasterTypePermitted, false),
			},
			"existed": {
				Type:        schema.TypeBool,
				Description: "Is the instance existed",
				Optional:    true,
				Computed:    true,
			},
			"existed_option": {
				Type:        schema.TypeList,
				Description: "Existed Instance Option",
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        resourceCCEv2ExistedOption(),
			},
			"machine_type": {
				Type:         schema.TypeString,
				Description:  "Machine Type. Available Values: [BCC, BBC, EBC, HPAS].",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(MachineTypePermitted, false),
			},
			"instance_type": {
				Type:         schema.TypeString,
				Description:  "Instance Type. Available Values: [N1, N2, N3, N4, N5, C1, C2, S1, G1, F1, HPAS].",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(BCCInstanceTypePermitted, false),
			},
			"bbc_option": {
				Type:        schema.TypeList,
				Description: "BBC Option",
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        resourceCCEv2BBCOption(),
			},
			"hpas_option": {
				Type:        schema.TypeList,
				Description: "HPAS Option",
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        resourceCCEv2HPASOption(),
			},
			"ehc_cluster_id": {
				Type:        schema.TypeString,
				Description: "EHC Cluster ID for instances",
				Optional:    true,
				Computed:    true,
			},
			"vpc_config": {
				Type:        schema.TypeList,
				Description: "VPC Config",
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        resourceCCEv2VPCConfig(),
			},
			"instance_resource": {
				Type:        schema.TypeList,
				Description: "Instance Resource Config",
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        resourceCCEv2InstanceResource(),
			},
			"image_id": {
				Type:        schema.TypeString,
				Description: "Image ID",
				Optional:    true,
				Computed:    true,
			},
			"instance_os": {
				Type:        schema.TypeList,
				Description: "OS Config of the instance",
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        resourceCCEv2InstanceOS(),
			},
			"need_eip": {
				Type:        schema.TypeBool,
				Description: "Whether the instance need a EIP",
				Optional:    true,
				Computed:    true,
			},
			"eip_option": {
				Type:        schema.TypeList,
				Description: "EIP Option",
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        resourceCCEv2EIPOption(),
			},
			"admin_password": {
				Type:        schema.TypeString,
				Description: "Admin Password",
				Optional:    true,
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return true
				},
			},
			"ssh_key_id": {
				Type:        schema.TypeString,
				Description: "SSH Key ID",
				Optional:    true,
				Computed:    true,
			},
			"instance_charging_type": {
				Type:         schema.TypeString,
				Description:  "Instance charging type. Available Value: [Prepaid, Postpaid, bidding].",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(PaymentTimingTypePermitted, false),
			},
			"instance_precharging_option": {
				Type:        schema.TypeList,
				Description: "Instance Pre-charging Option",
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        resourceCCEv2InstancePrechargingOption(),
			},
			"delete_option": {
				Type:        schema.TypeList,
				Description: "Delete Option",
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        resourceCCEv2DeleteOption(),
			},
			"deploy_custom_config": {
				Type:        schema.TypeList,
				Description: "Deploy Custom Option",
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        resourceCCEv2DeployCustomConfig(),
			},
			"tag_list": {
				Type:        schema.TypeList,
				Description: "Tag List",
				Optional:    true,
				Computed:    true,
				Elem:        resourceCCEv2Tag(),
			},
			"labels": {
				Type:        schema.TypeMap,
				Description: "Labels List",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"instance_taints": {
				Type:        schema.TypeList,
				Description: "Taint List",
				Optional:    true,
				Elem:        resourceCCEv2Taint(),
			},
			"cce_instance_priority": {
				Type:        schema.TypeInt,
				Description: "Priority of this instance.",
				Computed:    true,
			},
		},
	}
}

func resourceCCEv2InstanceStatus() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"machine": {
				Type:        schema.TypeList,
				Description: "Machine info",
				Computed:    true,
				MaxItems:    1,
				Elem:        resourceCCEv2Machine(),
			},
			"instance_phase": {
				Type:        schema.TypeString,
				Description: "Instance Phase",
				Computed:    true,
			},
			"machine_status": {
				Type:        schema.TypeString,
				Description: "Machine status",
				Computed:    true,
			},
		},
	}
}

func resourceCCEv2Machine() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Description: "Instance ID",
				Computed:    true,
			},
			"order_id": {
				Type:        schema.TypeString,
				Description: "Order ID",
				Computed:    true,
			},
			"mount_list": {
				Type:        schema.TypeList,
				Description: "Mount List of Machine",
				Computed:    true,
				Elem:        resourceCCEv2MountConfig(),
			},
			"vpc_ip": {
				Type:        schema.TypeString,
				Description: "VPC IP",
				Computed:    true,
			},
			"vpc_ip_ipv6": {
				Type:        schema.TypeString,
				Description: "VPC IPv6",
				Computed:    true,
			},
			"eip": {
				Type:        schema.TypeString,
				Description: "EIP",
				Computed:    true,
			},
		},
	}
}

func resourceCCEv2MountConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Description: "Mount Path",
				Optional:    true,
			},
			"cds_id": {
				Type:        schema.TypeString,
				Description: "CDS ID",
				Optional:    true,
			},
			"device": {
				Type:        schema.TypeString,
				Description: "Device Path",
				Optional:    true,
			},
			"cds_size": {
				Type:        schema.TypeInt,
				Description: "CDS Size",
				Optional:    true,
			},
			"storage_type": {
				Type:         schema.TypeString,
				Description:  "Storage type. Available Value: [std1, hp1, cloud_hp1, local, sata, ssd, hdd].",
				Optional:     true,
				ValidateFunc: validation.StringInSlice(StorageTypePermitted, false),
			},
		},
	}
}

func resourceCCEv2K8SCustomConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"master_feature_gates": {
				Type:        schema.TypeMap,
				Description: "custom master Feature Gates",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeBool,
				},
			},
			"node_feature_gates": {
				Type:        schema.TypeMap,
				Description: "custom node Feature Gates",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeBool,
				},
			},
			"admission_plugins": {
				Type:        schema.TypeList,
				Description: "custom Admission Plugins",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"pause_image": {
				Type:        schema.TypeString,
				Description: "custom PauseImage",
				Optional:    true,
			},
			"kube_api_qps": {
				Type:        schema.TypeInt,
				Description: "custom Kube API QPS",
				Optional:    true,
			},
			"kube_api_burst": {
				Type:        schema.TypeInt,
				Description: "custom Kube API Burst",
				Optional:    true,
			},
			"scheduler_predicated": {
				Type:        schema.TypeList,
				Description: "custom Scheduler Predicates",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"scheduler_priorities": {
				Type:        schema.TypeMap,
				Description: "custom SchedulerPriorities",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"etcd_data_path": {
				Type:        schema.TypeString,
				Description: "etcd data directory",
				Optional:    true,
			},
		},
	}
}

func resourceCCEv2MasterConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"master_type": {
				Type:         schema.TypeString,
				Description:  "Master Type. Available Value: [managed, custom, serverless].",
				Optional:     true,
				ValidateFunc: validation.StringInSlice(MasterTypePermitted, false),
			},
			"cluster_ha": {
				Type:         schema.TypeInt,
				Description:  "Number of master nodes. Available Value: [1, 3, 5, 2(for serverless)].",
				Optional:     true,
				ValidateFunc: validation.IntInSlice(ClusterHAPermitted),
			},
			"exposed_public": {
				Type:        schema.TypeBool,
				Description: "Whether exposed to public network",
				Optional:    true,
			},
			"cluster_blb_vpc_subnet_id": {
				Type:        schema.TypeString,
				Description: "Cluster BLB VPC Subnet ID",
				Optional:    true,
			},
			"managed_cluster_master_option": {
				Type:        schema.TypeList,
				Description: "Managed cluster master option",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"master_vpc_subnet_zone": {
							Type:         schema.TypeString,
							Description:  "Master VPC Subnet Zone. Available Value: [zoneA, zoneB, zoneC, zoneD, zoneE, zoneF].",
							Optional:     true,
							ValidateFunc: validation.StringInSlice(AvailableZonePermitted, false),
						},
					},
				},
			},
		},
	}
}

func resourceCCEv2ContainerNetworkConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"mode": {
				Type:         schema.TypeString,
				Description:  "Network Mode. Available Value: [kubenet, vpc-cni, vpc-route-veth, vpc-route-ipvlan, vpc-route-auto-detect, vpc-secondary-ip-veth, vpc-secondary-ip-ipvlan, vpc-secondary-ip-auto-detect].",
				Optional:     true,
				ValidateFunc: validation.StringInSlice(ContainerNetworkModePermitted, false),
			},
			"eni_vpc_subnet_ids": {
				Type:        schema.TypeList,
				Description: "ENI VPC Subnet ID",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_and_id": {
							Type:        schema.TypeMap,
							Description: "Available Zone and ENI ID",
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"eni_security_group_id": {
				Type:        schema.TypeString,
				Description: "ENI Security Group ID",
				Optional:    true,
			},
			"ip_version": {
				Type:         schema.TypeString,
				Description:  "IP Version. Available Value: [ipv4, ipv6, dualStack].",
				Optional:     true,
				ValidateFunc: validation.StringInSlice(ContainerNetworkIPTypePermitted, false),
			},
			"lb_service_vpc_subnet_id": {
				Type:        schema.TypeString,
				Description: "LB Service VPC Sunnet ID",
				Optional:    true,
			},
			"node_port_range_min": {
				Type:        schema.TypeInt,
				Description: "Node Port Service Port Range Min",
				Optional:    true,
			},
			"node_port_range_max": {
				Type:        schema.TypeInt,
				Description: "Node Port Service Port Range Max",
				Optional:    true,
			},
			"cluster_pod_cidr": {
				Type:        schema.TypeString,
				Description: "Cluster Pod IP CIDR",
				Optional:    true,
			},
			"cluster_pod_cidr_ipv6": {
				Type:        schema.TypeString,
				Description: "Cluster Pod IP CIDR IPv6",
				Optional:    true,
			},
			"cluster_ip_service_cidr": {
				Type:        schema.TypeString,
				Description: "Cluster Service ClusterIP CIDR ",
				Optional:    true,
			},
			"cluster_ip_service_cidr_ipv6": {
				Type:        schema.TypeString,
				Description: "Cluster Service ClusterIP CIDR IPv6",
				Optional:    true,
			},
			"max_pods_per_node": {
				Type:        schema.TypeInt,
				Description: "Max pod number in one node ",
				Optional:    true,
			},
			"kube_proxy_mode": {
				Type:         schema.TypeString,
				Description:  "KubeProxy Mode. Available Value: [iptables, ipvs].",
				Optional:     true,
				ValidateFunc: validation.StringInSlice(KubeProxyModePermitted, false),
			},
		},
	}
}

func resourceCCEv2ExistedOption() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"existed_instance_id": {
				Type:        schema.TypeString,
				Description: "Existed Instance ID",
				Optional:    true,
				Computed:    true,
			},
			"rebuild": {
				Type:        schema.TypeBool,
				Description: "Whether re-install OS",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func resourceCCEv2BBCOption() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"reserve_data": {
				Type:        schema.TypeBool,
				Description: "Whether reserve data",
				Optional:    true,
				Computed:    true,
			},
			"raid_id": {
				Type:        schema.TypeString,
				Description: "Disk Raid ID",
				Optional:    true,
				Computed:    true,
			},
			"sys_disk_size": {
				Type:        schema.TypeInt,
				Description: "System Disk Size",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func resourceCCEv2HPASOption() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"app_type": {
				Type:        schema.TypeString,
				Description: "Application type of the HPAS instance. e.g., `llama2_7B_train`.",
				Required:    true,
			},
			"app_performance_level": {
				Type:        schema.TypeString,
				Description: "Performance level of the application. e.g., `10k`.",
				Required:    true,
			},
		},
	}
}

func resourceCCEv2VPCConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Description: "VPC ID",
				Optional:    true,
				Computed:    true,
			},
			"vpc_subnet_id": {
				Type:        schema.TypeString,
				Description: "VPC Subnet ID",
				Optional:    true,
				Computed:    true,
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Description: "Security Group ID",
				Optional:    true,
				Computed:    true,
			},
			"security_group_type": {
				Type:         schema.TypeString,
				Description:  "Security Group type. Available Values: [normal, enterprise]. Default: `normal`",
				Optional:     true,
				Default:      "normal",
				ValidateFunc: validation.StringInSlice([]string{"normal", "enterprise"}, false),
			},
			"vpc_subnet_type": {
				Type:         schema.TypeString,
				Description:  "VPC Subnet type. Available Value: [BCC, BCC_NAT, BBC].",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(VPCSubnetTypePermitted, false),
			},
			"vpc_subnet_cidr": {
				Type:        schema.TypeString,
				Description: "VPC Subnet CIDR",
				Optional:    true,
				Computed:    true,
			},
			"vpc_subnet_cidr_ipv6": {
				Type:        schema.TypeString,
				Description: "VPC Sunbet CIDR IPv6",
				Optional:    true,
				Computed:    true,
			},
			"available_zone": {
				Type:         schema.TypeString,
				Description:  "Available Zone. Available Value: [zoneA, zoneB, zoneC, zoneD, zoneE, zoneF].",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(AvailableZonePermitted, false),
			},
		},
	}
}

func resourceCCEv2Instance() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_spec": {
				Type:        schema.TypeList,
				Description: "Instance specification",
				Computed:    true,
				MaxItems:    1,
				Elem:        resourceCCEv2InstanceSpec(),
			},
			"instance_status": {
				Type:        schema.TypeList,
				Description: "Instance status",
				Computed:    true,
				MaxItems:    1,
				Elem:        resourceCCEv2InstanceStatus(),
			},
			"created_at": {
				Type:        schema.TypeString,
				Description: "Instance create time",
				Computed:    true,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Description: "Instance update time",
				Computed:    true,
			},
		},
	}
}

func resourceCCEv2InstanceResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"machine_spec": {
				Type:        schema.TypeString,
				Description: "Machine specification for instances, e.g., 'llama_7B_train/10k'",
				Optional:    true,
				Computed:    true,
			},
			"cpu": {
				Type:        schema.TypeInt,
				Description: "CPU cores",
				Optional:    true,
				Computed:    true,
			},
			"mem": {
				Type:        schema.TypeInt,
				Description: "memory GB",
				Optional:    true,
				Computed:    true,
			},
			"node_cpu_quota": {
				Type:        schema.TypeInt,
				Description: "Node cpu quota",
				Optional:    true,
				Computed:    true,
			},
			"node_mem_quota": {
				Type:        schema.TypeInt,
				Description: "Node memory quota",
				Optional:    true,
				Computed:    true,
			},
			"root_disk_type": {
				Type:        schema.TypeString,
				Description: "Root disk type. Available Value: [std1, hp1, cloud_hp1, local, sata, ssd, hdd].",
				Optional:    true,
				Computed:    true,
			},
			"root_disk_size": {
				Type:        schema.TypeInt,
				Description: "Root disk size",
				Optional:    true,
				Computed:    true,
			},
			"local_disk_size": {
				Type:        schema.TypeInt,
				Description: "Local disk size",
				Optional:    true,
				Computed:    true,
			},
			"cds_list": {
				Type:        schema.TypeList,
				Description: "CDS List",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:        schema.TypeString,
							Description: "CDS path",
							Optional:    true,
							Computed:    true,
						},
						"storage_type": {
							Type:        schema.TypeString,
							Description: "Storage Type. Available Value: [std1, hp1, cloud_hp1, local, sata, ssd, hdd].",
							Optional:    true,
							Computed:    true,
						},
						"cds_size": {
							Type:        schema.TypeInt,
							Description: "CDS Size",
							Optional:    true,
							Computed:    true,
						},
						"snapshot_id": {
							Type:        schema.TypeString,
							Description: "Snap shot ID",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"ephemeral_disk_list": {
				Type:        schema.TypeList,
				Description: "Ephemeral Disk List for instances",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_type": {
							Type:        schema.TypeString,
							Description: "Storage Type. Available Value: [local_nvme, local_ssd].",
							Required:    true,
						},
						"size_in_gb": {
							Type:        schema.TypeInt,
							Description: "Disk size in GB",
							Required:    true,
						},
						"disk_path": {
							Type:        schema.TypeString,
							Description: "Custom disk mount path for local disks",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"gpu_type": {
				Type:         schema.TypeString,
				Description:  "GPU Type. Available Value: [V100-32, V100-16, P40, P4, K40, DLCard].",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(GPUTypePermitted, false),
			},
			"gpu_count": {
				Type:        schema.TypeInt,
				Description: "GPU Number",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func resourceCCEv2InstanceOS() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"image_type": {
				Type:         schema.TypeString,
				Description:  "Image type. Available Value: [Integration, System, All, Custom, Sharing, GpuBccSystem, GpuBccCustom, BbcSystem, BbcCustom].",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(ImageTypePermitted, false),
			},
			"image_name": {
				Type:        schema.TypeString,
				Description: "Image Name",
				Optional:    true,
				Computed:    true,
			},
			"os_type": {
				Type:         schema.TypeString,
				Description:  "OS type. Available Value: [linux, windows].",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(OSTypePermitted, false),
			},
			"os_name": {
				Type:         schema.TypeString,
				Description:  "OS name. Available Value: [CentOS, Ubuntu, Windows Server, Debian, opensuse].",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(OSNamePermitted, false),
			},
			"os_version": {
				Type:        schema.TypeString,
				Description: "OS version",
				Optional:    true,
				Computed:    true,
			},
			"os_arch": {
				Type:        schema.TypeString,
				Description: "OS arch",
				Optional:    true,
				Computed:    true,
			},
			"os_build": {
				Type:        schema.TypeString,
				Description: "OS Build Time",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func resourceCCEv2EIPOption() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"eip_name": {
				Type:        schema.TypeString,
				Description: "EIP Name",
				Optional:    true,
				Computed:    true,
			},
			"eip_charging_type": {
				Type:         schema.TypeString,
				Description:  "EIP Charging Type. Available Value: [ByTraffic, ByBandwidth].",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(EIPBillingMethodPermitted, false),
			},
			"eip_bandwidth": {
				Type:        schema.TypeInt,
				Description: "EIP Bandwidth",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func resourceCCEv2Taint() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Description: "Taint Key",
				Optional:    true,
				Computed:    true,
			},
			"value": {
				Type:        schema.TypeString,
				Description: "Taint Value",
				Optional:    true,
				Computed:    true,
			},
			"effect": {
				Type:         schema.TypeString,
				Description:  "Taint Effect. Available Value: [NoSchedule, PreferNoSchedule, NoExecute].",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(TaintEffectPermitted, false),
			},
			"time_added": {
				Type:        schema.TypeString,
				Description: "Taint Added Time. Format RFC3339",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func resourceCCEv2Tag() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"tag_key": {
				Type:        schema.TypeString,
				Description: "Tag Key",
				Optional:    true,
				Computed:    true,
			},
			"tag_value": {
				Type:        schema.TypeString,
				Description: "Tag Value",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func resourceCCEv2DeleteOption() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"move_out": {
				Type:        schema.TypeBool,
				Description: "Whether move out the instance",
				Optional:    true,
				Computed:    true,
			},
			"delete_resource": {
				Type:        schema.TypeBool,
				Description: "Whether delete resources",
				Optional:    true,
				Computed:    true,
			},
			"delete_cds_snapshot": {
				Type:        schema.TypeBool,
				Description: "Whether delete CDS snapshot",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func resourceCCEv2ClusterDeleteOption() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"delete_resource": {
				Type:        schema.TypeBool,
				Description: "Whether to delete resources",
				Optional:    true,
				Computed:    true,
			},
			"delete_cds_snapshot": {
				Type:        schema.TypeBool,
				Description: "Whether to delete CDS snapshot",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func resourceCCEv2DeployCustomConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"docker_config": {
				Type:        schema.TypeList,
				Description: "Docker Config Info",
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"docker_data_root": {
							Type:        schema.TypeString,
							Description: "Customized Docker Data Directory",
							Optional:    true,
							Computed:    true,
						},
						"registry_mirrors": {
							Type:        schema.TypeList,
							Description: "Customized RegistryMirrors",
							Optional:    true,
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"insecure_registries": {
							Type:        schema.TypeList,
							Description: "Customized InsecureRegistries",
							Optional:    true,
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"docker_log_max_size": {
							Type:        schema.TypeString,
							Description: "docker Log Max Size",
							Optional:    true,
							Computed:    true,
						},
						"docker_log_max_file": {
							Type:        schema.TypeString,
							Description: "docker Log Max File",
							Optional:    true,
							Computed:    true,
						},
						"bip": {
							Type:        schema.TypeString,
							Description: "docker0 Network Bridge Network Segment",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"kubelet_root_dir": {
				Type:        schema.TypeString,
				Description: "kubelet Data Directory",
				Optional:    true,
				Computed:    true,
			},
			"enable_resource_reserved": {
				Type:        schema.TypeBool,
				Description: "Whether to Enable Resource Quota",
				Optional:    true,
				Computed:    true,
			},
			"kube_reserved": {
				Type:        schema.TypeMap,
				Description: "Resource Quota",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"enable_cordon": {
				Type:        schema.TypeBool,
				Description: "Whether enable cordon",
				Optional:    true,
				Computed:    true,
			},
			"pre_user_script": {
				Type:        schema.TypeString,
				Description: "Script before deployment, base64 encoded",
				Optional:    true,
				Computed:    true,
			},
			"post_user_script": {
				Type:        schema.TypeString,
				Description: "Script after deployment, base64 encoded",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func resourceCCEv2InstancePrechargingOption() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"purchase_time": {
				Type:        schema.TypeInt,
				Description: "Time of purchase",
				Optional:    true,
				Computed:    true,
			},
			"auto_renew": {
				Type:        schema.TypeBool,
				Description: "Is Auto Renew",
				Optional:    true,
				Computed:    true,
			},
			"auto_renew_time_unit": {
				Type:        schema.TypeString,
				Description: "Time unit for auto renew",
				Optional:    true,
				Computed:    true,
			},
			"auto_renew_time": {
				Type:        schema.TypeInt,
				Description: "Number of time unit for auto renew",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func buildCreateInstanceGroupArgs(d *schema.ResourceData) (*ccev2.CreateInstanceGroupArgs, error) {
	instanceGroupSpecMap := d.Get("spec.0").(map[string]interface{})
	instanceSpecMap := instanceGroupSpecMap["instance_template"].([]interface{})[0].(map[string]interface{})
	instanceSpec, err := buildInstanceSpec(instanceSpecMap)
	if err != nil {
		return nil, err
	}
	args := &ccev2.CreateInstanceGroupArgs{
		ClusterID: instanceGroupSpecMap["cluster_id"].(string),
		Request: &ccev2.CreateInstanceGroupRequest{
			InstanceGroupSpec: ccev2types.InstanceGroupSpec{
				InstanceGroupName: instanceGroupSpecMap["instance_group_name"].(string),
				Replicas:          instanceGroupSpecMap["replicas"].(int),
				InstanceTemplate: ccev2types.InstanceTemplate{
					InstanceSpec: *instanceSpec,
				},
				CleanPolicy: ccev2types.DeleteCleanPolicy,
			},
		},
	}
	return args, nil
}

func buildUpdateInstanceGroupReplicaArgs(d *schema.ResourceData) (*ccev2.UpdateInstanceGroupReplicasArgs, error) {
	instanceGroupSpecMap := d.Get("spec.0").(map[string]interface{})
	ars := &ccev2.UpdateInstanceGroupReplicasArgs{
		ClusterID:       instanceGroupSpecMap["cluster_id"].(string),
		InstanceGroupID: d.Id(),
		Request: &ccev2.UpdateInstanceGroupReplicasRequest{
			Replicas:       instanceGroupSpecMap["replicas"].(int),
			DeleteInstance: true,
			DeleteOption: &ccev2types.DeleteOption{
				MoveOut:           false,
				DeleteResource:    true,
				DeleteCDSSnapshot: true,
			},
		},
	}
	return ars, nil
}

func buildDeleteInstanceGroupArgs(d *schema.ResourceData) (*ccev2.DeleteInstanceGroupArgs, error) {
	instanceGroupSpecMap := d.Get("spec.0").(map[string]interface{})
	args := &ccev2.DeleteInstanceGroupArgs{
		ClusterID:       instanceGroupSpecMap["cluster_id"].(string),
		InstanceGroupID: d.Id(),
		DeleteInstances: true,
	}
	return args, nil
}

func buildGetInstanceGroupArgs(d *schema.ResourceData) (*ccev2.GetInstanceGroupArgs, error) {
	instanceGroupSpecMap := d.Get("spec.0").(map[string]interface{})
	args := &ccev2.GetInstanceGroupArgs{
		ClusterID:       instanceGroupSpecMap["cluster_id"].(string),
		InstanceGroupID: d.Id(),
	}
	return args, nil
}

func buildGetInstancesOfInstanceGroupArgs(d *schema.ResourceData) (*ccev2.ListInstanceByInstanceGroupIDArgs, error) {
	instanceGroupSpecMap := d.Get("spec.0").(map[string]interface{})
	args := &ccev2.ListInstanceByInstanceGroupIDArgs{
		ClusterID:       instanceGroupSpecMap["cluster_id"].(string),
		InstanceGroupID: d.Id(),
		PageSize:        0,
		PageNo:          0,
	}
	return args, nil
}

func buildUpdateInstanceGroupConfigureArgs(d *schema.ResourceData, spec *ccev2.InstanceGroupSpec) (*ccev2.UpdateInstanceGroupConfigure, error) {
	if spec == nil {
		return nil, fmt.Errorf("instance group spec is nil")
	}

	instanceSpec := spec.InstanceTemplate.InstanceSpec

	if d.HasChange("spec.0.instance_template.0.labels") {
		labels := make(map[string]string)
		if labelsRaw, ok := d.Get("spec.0.instance_template.0.labels").(map[string]interface{}); ok {
			for key, value := range labelsRaw {
				labels[key] = value.(string)
			}
		}
		instanceSpec.Labels = labels
	}

	if d.HasChange("spec.0.instance_template.0.instance_taints") {
		taintsRaw, _ := d.Get("spec.0.instance_template.0.instance_taints").([]interface{})
		taints, err := buildTaints(taintsRaw)
		if err != nil {
			return nil, err
		}
		instanceSpec.Taints = taints
	}

	if d.HasChange("spec.0.instance_template.0.image_id") {
		imageID, _ := d.Get("spec.0.instance_template.0.image_id").(string)
		instanceSpec.ImageID = imageID
	}

	var securityGroups []ccev2types.SecurityGroupV2
	for _, sg := range spec.DefaultSecurityGroups {
		securityGroups = append(securityGroups, ccev2types.SecurityGroupV2{
			Name: sg.Name,
			Type: ccev2types.SecurityGroupType(sg.Type),
			ID:   sg.ID,
		})
	}

	instanceTemplate := ccev2types.InstanceTemplate{
		InstanceSpec: instanceSpec,
	}

	updateSpec := ccev2types.InstanceGroupSpec{
		CCEInstanceGroupID: spec.CCEInstanceGroupID,
		InstanceGroupName:  spec.InstanceGroupName,
		ClusterID:          spec.ClusterID,
		ClusterRole:        spec.ClusterRole,
		ShrinkPolicy:       ccev2types.ShrinkPolicy(spec.ShrinkPolicy),
		UpdatePolicy:       ccev2types.UpdatePolicy(spec.UpdatePolicy),
		CleanPolicy:        ccev2types.CleanPolicy(spec.CleanPolicy),
		InstanceTemplate:   instanceTemplate,
		InstanceTemplates:  []ccev2types.InstanceTemplate{instanceTemplate},
		Replicas:           spec.Replicas,
		SecurityGroups:     securityGroups,
	}

	return &ccev2.UpdateInstanceGroupConfigure{
		PasswordNeedUpdate: false,
		SyncMeta:           false,
		InstanceGroupSpec:  updateSpec,
	}, nil
}

//===================Convert系函数用于将SDK返回值转换成.tfstate参数===================
//Tips: 对于.tfstate中存在的字段，但是sdk返回数据中不包含的字段，将不会在传递给.tfstate的map中设置此值，进而terrafrom会跳过更新此字段的状态，使其维持原样不变

func convertInstanceFromJsonToMap(instances []*ccev2.Instance, role ccev2types.ClusterRole) ([]interface{}, error) {
	targetInstances := make([]*ccev2.Instance, 0)
	resultInstances := make([]interface{}, 0, len(targetInstances))
	if instances == nil || len(instances) == 0 {
		return resultInstances, nil
	}

	//区分是master机器还是node机器
	for _, instance := range instances {
		if instance.Spec.ClusterRole == role {
			targetInstances = append(targetInstances, instance)
		}
	}

	for _, instance := range targetInstances {

		instanceMap := make(map[string]interface{})

		if instance.Spec != nil {
			spec, err := convertInstanceSpecFromJsonToMap(instance.Spec)
			if err != nil {
				return nil, err
			}
			instanceMap["instance_spec"] = spec
		}

		if instance.Status != nil {
			status, err := convertInstanceStatusFromJsonToMap(instance.Status)
			if err != nil {
				return nil, err
			}
			instanceMap["instance_status"] = status
		}

		instanceMap["created_at"] = instance.CreatedAt.String()
		instanceMap["updated_at"] = instance.UpdatedAt.String()

		resultInstances = append(resultInstances, instanceMap)
	}

	return resultInstances, nil
}

func convertInstanceSpecFromJsonToMap(spec *ccev2types.InstanceSpec) ([]interface{}, error) {
	resultSpec := make([]interface{}, 0)
	if spec == nil {
		return resultSpec, nil
	}
	specMap := make(map[string]interface{})

	if spec.CCEInstanceID != "" {
		specMap["cce_instance_id"] = spec.CCEInstanceID
	}
	if spec.InstanceName != "" {
		specMap["instance_name"] = spec.InstanceName
	}
	if spec.RuntimeType != "" {
		specMap["runtime_type"] = spec.RuntimeType
	}
	if spec.RuntimeVersion != "" {
		specMap["runtime_version"] = spec.RuntimeVersion
	}
	if spec.ClusterID != "" {
		specMap["cluster_id"] = spec.ClusterID
	}
	if spec.InstanceChargingType != "" {
		specMap["instance_charging_type"] = spec.InstanceChargingType
	}
	if spec.ClusterRole != "" {
		specMap["cluster_role"] = spec.ClusterRole
	}
	if spec.InstanceGroupID != "" {
		specMap["instance_group_id"] = spec.InstanceGroupID
	}
	if spec.InstanceGroupName != "" {
		specMap["instance_group_name"] = spec.InstanceGroupName
	}
	if spec.MachineType != "" {
		specMap["machine_type"] = spec.MachineType
	}
	if spec.InstanceType != "" {
		specMap["instance_type"] = spec.InstanceType
	}
	if spec.ImageID != "" {
		specMap["image_id"] = spec.ImageID
	}
	specMap["cce_instance_priority"] = spec.CCEInstancePriority

	specMap["need_eip"] = spec.NeedEIP

	if spec.SSHKeyID != "" {
		specMap["ssh_key_id"] = spec.SSHKeyID
	}

	if spec.BBCOption != nil {
		option, err := convertBBCOptionFromJsonToMap(spec.BBCOption)
		if err != nil {
			return nil, err
		}
		specMap["bbc_option"] = option
	}

	if spec.HPASOption != nil {
		option, err := convertHPASOptionFromJsonToMap(spec.HPASOption)
		if err != nil {
			return nil, err
		}
		specMap["hpas_option"] = option
	}

	if spec.EhcClusterID != "" {
		specMap["ehc_cluster_id"] = spec.EhcClusterID
	}

	if !reflect.DeepEqual(spec.VPCConfig, ccev2types.VPCConfig{}) {
		config, err := convertVPCConfigFromJsonToMap(&spec.VPCConfig)
		if err != nil {
			return nil, err
		}
		specMap["vpc_config"] = config
	}

	if spec.EIPOption != nil {
		option, err := convertEIPOptionFromJsonToMap(spec.EIPOption)
		if err != nil {
			return nil, err
		}
		specMap["eip_option"] = option
	}

	if spec.DeleteOption != nil {
		option, err := convertDeleteOptionFromJsonToMap(spec.DeleteOption)
		if err != nil {
			return nil, err
		}
		specMap["delete_option"] = option
	}

	option, err := convertDeployCustomConfigFromJsonToMap(&spec.DeployCustomConfig)
	if err != nil {
		return nil, err
	}
	specMap["deploy_custom_config"] = option

	instanceResource, err := convertInstanceResourceFromJsonToMap(&spec.InstanceResource)
	if err != nil {
		return nil, err
	}
	specMap["instance_resource"] = instanceResource

	if spec.InstanceOS != (ccev2types.InstanceOS{}) {
		option, err := convertInstanceOSFromJsonToMap(&spec.InstanceOS)
		if err != nil {
			return nil, err
		}
		specMap["instance_os"] = option
	}

	if spec.Tags != nil {
		tagMapList, err := convertTagListFromJsonToMap(spec.Tags)
		if err != nil {
			return nil, err
		}
		specMap["tag_list"] = tagMapList
	}

	if spec.Taints != nil {
		taintMapList, err := convertTaintListFromJsonToMap(spec.Taints)
		if err != nil {
			return nil, err
		}
		specMap["instance_taints"] = taintMapList
	}

	if spec.Labels != nil {
		specMap["labels"] = spec.Labels
	}

	resultSpec = append(resultSpec, specMap)

	return resultSpec, nil
}

func convertInstanceGroupSpecFromJsonToMap(spec *ccev2.InstanceGroupSpec) ([]interface{}, error) {
	result := make([]interface{}, 0, 1)
	if spec == nil {
		return result, nil
	}

	specMap := make(map[string]interface{})
	specMap["cluster_id"] = spec.ClusterID
	specMap["instance_group_name"] = spec.InstanceGroupName
	specMap["replicas"] = spec.Replicas

	instanceTemplate, err := convertInstanceSpecFromJsonToMap(&spec.InstanceTemplate.InstanceSpec)
	if err != nil {
		return nil, err
	}
	specMap["instance_template"] = instanceTemplate

	result = append(result, specMap)
	return result, nil
}

func convertDeployCustomConfigFromJsonToMap(config *ccev2types.DeployCustomConfig) ([]interface{}, error) {
	result := make([]interface{}, 0)
	if config == nil {
		return result, nil
	}
	configMap := make(map[string]interface{})

	configMap["enable_resource_reserved"] = config.EnableResourceReserved
	configMap["enable_cordon"] = config.EnableCordon
	if config.PreUserScript != "" {
		configMap["pre_user_script"] = config.PreUserScript
	}
	if config.PostUserScript != "" {
		configMap["post_user_script"] = config.PostUserScript
	}
	if config.KubeletRootDir != "" {
		configMap["kubelet_root_dir"] = config.KubeletRootDir
	}
	if config.KubeReserved != nil {
		configMap["kube_reserved"] = config.KubeReserved
	}

	//DockerConfig is not a pointer.
	dockerConfig, err := convertDockerConfigFromJsonToMap(config.DockerConfig)
	if err != nil {
		return nil, err
	}
	configMap["docker_config"] = dockerConfig

	result = append(result, configMap)
	return result, nil
}

func convertDockerConfigFromJsonToMap(config ccev2types.DockerConfig) ([]interface{}, error) {
	result := make([]interface{}, 0)

	configMap := make(map[string]interface{})

	if config.BIP != "" {
		configMap["bip"] = config.BIP
	}
	if config.DockerLogMaxFile != "" {
		configMap["docker_log_max_file"] = config.DockerLogMaxFile
	}
	if config.DockerLogMaxSize != "" {
		configMap["docker_log_max_size"] = config.DockerLogMaxSize
	}
	if config.DockerDataRoot != "" {
		configMap["docker_data_root"] = config.DockerDataRoot
	}
	if config.InsecureRegistries != nil {
		configMap["insecure_registries"] = config.InsecureRegistries
	}
	if config.RegistryMirrors != nil {
		configMap["registry_mirrors"] = config.RegistryMirrors
	}

	result = append(result, configMap)
	return result, nil
}

func convertTaintListFromJsonToMap(taintList []ccev2types.Taint) ([]interface{}, error) {
	result := make([]interface{}, 0, len(taintList))
	if taintList == nil {
		return result, nil
	}
	for _, taint := range taintList {
		cdsMap := make(map[string]interface{})
		if taint.Key != "" {
			cdsMap["key"] = taint.Key
		}
		if taint.Value != "" {
			cdsMap["value"] = taint.Value
		}
		if taint.Effect != "" {
			cdsMap["effect"] = taint.Effect
		}
		if taint.TimeAdded != nil {
			cdsMap["time_added"] = taint.TimeAdded.String()
		}

		result = append(result, cdsMap)
	}

	return result, nil
}

func convertTagListFromJsonToMap(tagList []ccev2types.Tag) ([]interface{}, error) {
	result := make([]interface{}, 0, len(tagList))
	if tagList == nil {
		return result, nil
	}

	for _, tag := range tagList {
		cdsMap := make(map[string]interface{})
		if tag.TagKey != "" {
			cdsMap["tag_key"] = tag.TagKey
		}
		if tag.TagValue != "" {
			cdsMap["tag_value"] = tag.TagValue
		}
		result = append(result, cdsMap)
	}

	return result, nil
}

func convertInstanceResourceFromJsonToMap(config *ccev2types.InstanceResource) ([]interface{}, error) {
	result := make([]interface{}, 0)
	if config == nil {
		return result, nil
	}
	configMap := make(map[string]interface{})

	if config.MachineSpec != "" {
		configMap["machine_spec"] = config.MachineSpec
	}
	configMap["gpu_count"] = config.GPUCount
	configMap["node_cpu_quota"] = config.NodeCPUQuota
	configMap["node_mem_quota"] = config.NodeMEMQuota
	configMap["local_disk_size"] = config.LocalDiskSize
	configMap["root_disk_size"] = config.RootDiskSize
	configMap["mem"] = config.MEM
	configMap["cpu"] = config.CPU
	if config.GPUType != "" {
		configMap["gpu_type"] = config.GPUCount
	}
	if config.RootDiskType != "" {
		configMap["root_disk_type"] = config.RootDiskType
	}
	if config.CDSList != nil {
		cdsListMap, err := convertCDSListFromJsonToMap(config.CDSList)
		if err != nil {
			return nil, err
		}
		configMap["cds_list"] = cdsListMap
	}

	if config.EphemeralDiskList != nil {
		ephemeralDiskListMap, err := convertEphemeralDiskListFromJsonToMap(config.EphemeralDiskList)
		if err != nil {
			return nil, err
		}
		configMap["ephemeral_disk_list"] = ephemeralDiskListMap
	}

	result = append(result, configMap)
	return result, nil
}

func convertCDSListFromJsonToMap(cdsList []ccev2types.CDSConfig) ([]interface{}, error) {
	result := make([]interface{}, 0, len(cdsList))
	for _, cds := range cdsList {
		cdsMap := make(map[string]interface{})

		if cds.Path != "" {
			cdsMap["path"] = cds.Path
		}
		if cds.StorageType != "" {
			cdsMap["storage_type"] = cds.StorageType
		}
		if cds.SnapshotID != "" {
			cdsMap["snapshot_id"] = cds.SnapshotID
		}
		cdsMap["cds_size"] = cds.CDSSize

		result = append(result, cdsMap)
	}
	return result, nil
}

func convertEphemeralDiskListFromJsonToMap(ephemeralDiskList []ccev2types.EphemeralDiskConfig) ([]interface{}, error) {
	result := make([]interface{}, 0, len(ephemeralDiskList))
	for _, disk := range ephemeralDiskList {
		diskMap := make(map[string]interface{})

		if disk.StorageType != "" {
			diskMap["storage_type"] = disk.StorageType
		}
		diskMap["size_in_gb"] = disk.SizeInGB
		if disk.Path != "" {
			diskMap["disk_path"] = disk.Path
		}

		result = append(result, diskMap)
	}
	return result, nil
}

func convertInstanceOSFromJsonToMap(config *ccev2types.InstanceOS) ([]interface{}, error) {
	result := make([]interface{}, 0)
	if config == nil {
		return result, nil
	}
	configMap := make(map[string]interface{})

	if config.OSBuild != "" {
		configMap["os_build"] = config.OSBuild
	}
	if config.OSArch != "" {
		configMap["os_arch"] = config.OSArch
	}
	if config.OSVersion != "" {
		configMap["os_version"] = config.OSVersion
	}
	if config.ImageName != "" {
		configMap["image_name"] = config.ImageName
	}
	if config.OSName != "" {
		configMap["os_name"] = config.OSName
	}
	if config.OSType != "" {
		configMap["os_type"] = config.OSType
	}
	if config.ImageType != "" {
		configMap["image_type"] = config.ImageType
	}

	result = append(result, configMap)
	return result, nil
}

func convertDeleteOptionFromJsonToMap(option *ccev2types.DeleteOption) ([]interface{}, error) {
	result := make([]interface{}, 0)
	if option == nil {
		return result, nil
	}
	optionMap := make(map[string]interface{})
	optionMap["move_out"] = option.MoveOut
	optionMap["delete_resource"] = option.DeleteResource
	optionMap["delete_cds_snapshot"] = option.DeleteCDSSnapshot

	result = append(result, optionMap)
	return result, nil
}

func convertBBCOptionFromJsonToMap(option *ccev2types.BBCOption) ([]interface{}, error) {
	result := make([]interface{}, 0)
	if option == nil {
		return result, nil
	}

	optionMap := make(map[string]interface{})
	if option.RaidID != "" {
		optionMap["raid_id"] = option.RaidID
	}
	optionMap["sys_disk_size"] = option.SysDiskSize
	optionMap["reserve_data"] = option.ReserveData

	result = append(result, optionMap)
	return result, nil
}

func convertHPASOptionFromJsonToMap(option *ccev2types.HPASOption) ([]interface{}, error) {
	result := make([]interface{}, 0)
	if option == nil {
		return result, nil
	}

	optionMap := make(map[string]interface{})
	optionMap["app_type"] = option.AppType
	optionMap["app_performance_level"] = option.AppPerformanceLevel

	result = append(result, optionMap)
	return result, nil
}

func convertVPCConfigFromJsonToMap(config *ccev2types.VPCConfig) ([]interface{}, error) {
	result := make([]interface{}, 0)
	if config == nil {
		return result, nil
	}

	configMap := make(map[string]interface{})
	if config.VPCID != "" {
		configMap["vpc_id"] = config.VPCID
	}
	if config.VPCSubnetID != "" {
		configMap["vpc_subnet_id"] = config.VPCSubnetID
	}
	if config.SecurityGroupID != "" {
		configMap["security_group_id"] = config.SecurityGroupID
	}
	if config.SecurityGroupType != "" {
		configMap["security_group_type"] = config.SecurityGroupType
	}
	if config.VPCSubnetCIDR != "" {
		configMap["vpc_subnet_cidr"] = config.VPCSubnetCIDR
	}
	if config.VPCSubnetCIDRIPv6 != "" {
		configMap["vpc_subnet_cidr_ipv6"] = config.VPCSubnetCIDRIPv6
	}
	if config.VPCSubnetType != "" {
		configMap["vpc_subnet_type"] = config.VPCSubnetType
	}
	if config.AvailableZone != "" {
		configMap["available_zone"] = config.AvailableZone
	}

	result = append(result, configMap)
	return result, nil
}

func convertEIPOptionFromJsonToMap(option *ccev2types.EIPOption) ([]interface{}, error) {
	result := make([]interface{}, 0)
	if option == nil {
		return result, nil
	}

	optionMap := make(map[string]interface{})
	if option.EIPName != "" {
		optionMap["eip_name"] = option.EIPName
	}
	if option.EIPChargingType != "" {
		optionMap["eip_charging_type"] = option.EIPChargingType
	}
	optionMap["eip_bandwidth"] = option.EIPBandwidth

	result = append(result, optionMap)
	return result, nil
}

func convertInstanceStatusFromJsonToMap(status *ccev2.InstanceStatus) ([]interface{}, error) {
	resultStatus := make([]interface{}, 0)
	if status == nil {
		return resultStatus, nil
	}

	statusMap := make(map[string]interface{})

	if status.InstancePhase != "" {
		statusMap["instance_phase"] = status.InstancePhase
	}
	if status.InstancePhase != "" {
		statusMap["machine_status"] = status.MachineStatus
	}
	if &status.Machine != nil {
		machine, err := convertMachineFromJsonToMap(&status.Machine)
		if err != nil {
			return nil, err
		}
		statusMap["machine"] = machine
	}

	resultStatus = append(resultStatus, statusMap)

	return resultStatus, nil
}

func convertMachineFromJsonToMap(machine *ccev2.Machine) ([]interface{}, error) {
	result := make([]interface{}, 0)
	if machine == nil {
		return result, nil
	}

	machineMap := make(map[string]interface{})
	if machine.EIP != "" {
		machineMap["eip"] = machine.EIP
	}
	if machine.VPCIP != "" {
		machineMap["vpc_ip"] = machine.VPCIP
	}
	if machine.VPCIPIPv6 != "" {
		machineMap["vpc_ip_ipv6"] = machine.VPCIPIPv6
	}
	if machine.InstanceID != "" {
		machineMap["instance_id"] = machine.InstanceID
	}
	if machine.OrderID != "" {
		machineMap["order_id"] = machine.OrderID
	}
	if machine.MountList != nil {
		mounList, err := convertMountListFromJsonToMap(machine.MountList)
		if err != nil {
			return nil, err
		}
		machineMap["mount_list"] = mounList
	}

	result = append(result, machineMap)
	return result, nil
}

func convertMountListFromJsonToMap(mountconfigs []ccev2types.MountConfig) ([]interface{}, error) {
	result := make([]interface{}, 0, len(mountconfigs))
	for _, config := range mountconfigs {
		mountMap := make(map[string]interface{})

		mountMap["cds_size"] = config.CDSSize
		if config.Device != "" {
			mountMap["device"] = config.Device
		}
		if config.CDSID != "" {
			mountMap["cds_id"] = config.CDSID
		}
		if config.Path != "" {
			mountMap["path"] = config.Path
		}
		if config.StorageType != "" {
			mountMap["storage_type"] = config.StorageType
		}

		result = append(result, mountMap)
	}

	return result, nil
}

func convertClusterStatusFromJsonToTfMap(status *ccev2.ClusterStatus) ([]interface{}, error) {
	clusterStatusMapList := make([]interface{}, 0)
	if status == nil {
		return clusterStatusMapList, nil
	}

	blbMapList, err := convertBLBFromJsonToTfMap(&status.ClusterBLB)
	if err != nil {
		return nil, err
	}

	clusterStatusMap := make(map[string]interface{})
	clusterStatusMap["cluster_blb"] = blbMapList
	clusterStatusMap["cluster_phase"] = status.ClusterPhase
	clusterStatusMap["node_num"] = status.NodeNum

	clusterStatusMapList = append(clusterStatusMapList, clusterStatusMap)
	return clusterStatusMapList, nil
}

func convertBLBFromJsonToTfMap(blb *ccev2.BLB) ([]interface{}, error) {
	blbMapList := make([]interface{}, 0)
	if blb == nil {
		return blbMapList, nil
	}

	blbMap := make(map[string]interface{})
	blbMap["id"] = blb.ID
	blbMap["vpc_ip"] = blb.VPCIP
	blbMap["eip"] = blb.EIP
	blbMapList = append(blbMapList, blbMap)

	return blbMapList, nil
}

//===================Build系函数用于将.tf参数构建SDK请求参数并调用===================
//.tf是用户传入的配置文件，某些sdk要求的值可能并没有设置
//Tip: Build系函数对于sdk参数中存在但是.tf中没有设置的参数，会自动跳过赋值，即试用默认值

func buildCCEv2CreateClusterClusterSpec(clusterSpecRawMap map[string]interface{}) (*ccev2types.ClusterSpec, error) {

	clusterSpec := &ccev2types.ClusterSpec{}

	if v, ok := clusterSpecRawMap["cluster_name"]; ok && v.(string) != "" {
		clusterSpec.ClusterName = v.(string)
	}

	if v, ok := clusterSpecRawMap["cluster_type"]; ok && v.(string) != "" {
		clusterSpec.ClusterType = ccev2types.ClusterType(v.(string))
	}

	if v, ok := clusterSpecRawMap["description"]; ok && v.(string) != "" {
		clusterSpec.Description = v.(string)
	}

	if v, ok := clusterSpecRawMap["k8s_version"]; ok && v.(string) != "" {
		clusterSpec.K8SVersion = ccev2types.K8SVersion(v.(string))
	}

	if v, ok := clusterSpecRawMap["runtime_type"]; ok && v.(string) != "" {
		clusterSpec.RuntimeType = ccev2types.RuntimeType(v.(string))
	}

	if v, ok := clusterSpecRawMap["runtime_version"]; ok && v.(string) != "" {
		clusterSpec.RuntimeVersion = v.(string)
	}

	if v, ok := clusterSpecRawMap["vpc_id"]; ok && v.(string) != "" {
		clusterSpec.VPCID = v.(string)
	}

	if v, ok := clusterSpecRawMap["vpc_cidr"]; ok && v.(string) != "" {
		clusterSpec.VPCCIDR = v.(string)
	}

	if v, ok := clusterSpecRawMap["vpc_cidr_ipv6"]; ok && v.(string) != "" {
		clusterSpec.VPCCIDRIPv6 = v.(string)
	}

	if v, ok := clusterSpecRawMap["plugins"]; ok && v != nil {
		pluginList := make([]string, 0)
		for _, pluginRaw := range v.([]interface{}) {
			pluginList = append(pluginList, pluginRaw.(string))
		}
		clusterSpec.Plugins = pluginList
	}

	if v, ok := clusterSpecRawMap["container_network_config"]; ok && len(v.([]interface{})) == 1 {
		containerNetworkConfigRaw := v.([]interface{})[0].(map[string]interface{})
		containerNetworkConfig, err := buildCCEv2ContainerNetworkConfig(containerNetworkConfigRaw)
		if err != nil {
			log.Println("Build ClusterSpec ContainerNetworkConfig Error:" + err.Error())
			return nil, err
		}
		clusterSpec.ContainerNetworkConfig = *containerNetworkConfig
	}

	if v, ok := clusterSpecRawMap["master_config"]; ok && len(v.([]interface{})) == 1 {
		masterConfigRaw := v.([]interface{})[0].(map[string]interface{})
		masterConfig, err := buildCCEv2MasterConfig(masterConfigRaw)
		if err != nil {
			log.Println("Build ClusterSpec MasterConfig Error:" + err.Error())
			return nil, err
		}
		clusterSpec.MasterConfig = *masterConfig
	}

	if v, ok := clusterSpecRawMap["k8s_custom_config"]; ok && len(v.([]interface{})) == 1 {
		k8sCustomConfigRaw := v.([]interface{})[0].(map[string]interface{})
		k8sCustomConfig, err := buildK8SCustomConfig(k8sCustomConfigRaw)
		if err != nil {
			log.Println("Build ClusterSpec MasterConfig Error:" + err.Error())
			return nil, err
		}
		clusterSpec.K8SCustomConfig = *k8sCustomConfig
	}

	return clusterSpec, nil
}

func buildCCEv2MasterConfig(masterConfigRawMap map[string]interface{}) (*ccev2types.MasterConfig, error) {
	config := &ccev2types.MasterConfig{}

	if v, ok := masterConfigRawMap["master_type"]; ok && v.(string) != "" {
		config.MasterType = ccev2types.MasterType(v.(string))
	}

	if v, ok := masterConfigRawMap["cluster_ha"]; ok {
		config.ClusterHA = ccev2types.ClusterHA(v.(int))
	}

	if v, ok := masterConfigRawMap["exposed_public"]; ok {
		config.ExposedPublic = v.(bool)
	}

	if v, ok := masterConfigRawMap["cluster_blb_vpc_subnet_id"]; ok && v.(string) != "" {
		config.ClusterBLBVPCSubnetID = v.(string)
	}

	if v, ok := masterConfigRawMap["managed_cluster_master_option"]; ok && len(v.([]interface{})) == 1 {
		managedClusterMasterOptionRaw := v.([]interface{})[0].(map[string]interface{})
		managedClusterMasterOption, err := buildManagedClusterMasterOption(managedClusterMasterOptionRaw)
		if err != nil {
			log.Println("Build MasterConfig ManagedClusterMasterOption Error:" + err.Error())
			return nil, err
		}
		config.ManagedClusterMasterOption = *managedClusterMasterOption
	}

	return config, nil
}

func buildManagedClusterMasterOption(d map[string]interface{}) (*ccev2types.ManagedClusterMasterOption, error) {
	option := &ccev2types.ManagedClusterMasterOption{}

	if v, ok := d["master_vpc_subnet_zone"]; ok && v.(string) != "" {
		option.MasterVPCSubnetZone = ccev2types.AvailableZone(v.(string))
	}
	return option, nil
}

func buildENIVPCSubnetIDs(zoneAndIDRawMapList []interface{}) (map[ccev2types.AvailableZone][]string, error) {
	if zoneAndIDRawMapList == nil {
		return nil, nil
	}
	result := make(map[ccev2types.AvailableZone][]string, 0)

	for _, zoneAndIdMapRaw := range zoneAndIDRawMapList {
		outMap, _ := zoneAndIdMapRaw.(map[string]interface{})["zone_and_id"]
		zoneAndIdMap := outMap.(map[string]interface{})
		for zone, id := range zoneAndIdMap {
			if _, ok := result[ccev2types.AvailableZone(zone)]; !ok {
				result[ccev2types.AvailableZone(zone)] = make([]string, 0)
			}
			idListOfZone := result[ccev2types.AvailableZone(zone)]
			idListOfZone = append(idListOfZone, id.(string))
			result[ccev2types.AvailableZone(zone)] = idListOfZone
		}
	}

	return result, nil
}

func buildCCEv2ContainerNetworkConfig(containerNetworkConfigRawMap map[string]interface{}) (*ccev2types.ContainerNetworkConfig, error) {
	config := &ccev2types.ContainerNetworkConfig{}

	if v, ok := containerNetworkConfigRawMap["mode"]; ok && v.(string) != "" {
		config.Mode = ccev2types.ContainerNetworkMode(v.(string))
	}

	config.ENIVPCSubnetIDs = nil
	if v, ok := containerNetworkConfigRawMap["eni_vpc_subnet_ids"]; ok {
		values := v.([]interface{})
		ids, err := buildENIVPCSubnetIDs(values)
		if err != nil {
			log.Println("Build ContainerNetworkConfig ENIVPCSubnetIDs Error:" + err.Error())
			return nil, err
		}
		config.ENIVPCSubnetIDs = ids
	}

	if v, ok := containerNetworkConfigRawMap["eni_security_group_id"]; ok && v.(string) != "" {
		config.ENISecurityGroupID = v.(string)
	}

	if v, ok := containerNetworkConfigRawMap["ip_version"]; ok && v.(string) != "" {
		config.IPVersion = ccev2types.ContainerNetworkIPType(v.(string))
	}

	if v, ok := containerNetworkConfigRawMap["lb_service_vpc_subnet_id"]; ok && v.(string) != "" {
		config.LBServiceVPCSubnetID = v.(string)
	}

	if v, ok := containerNetworkConfigRawMap["node_port_range_min"]; ok {
		config.NodePortRangeMin = v.(int)
	}

	if v, ok := containerNetworkConfigRawMap["node_port_range_max"]; ok {
		config.NodePortRangeMax = v.(int)
	}

	if v, ok := containerNetworkConfigRawMap["cluster_pod_cidr"]; ok && v.(string) != "" {
		config.ClusterPodCIDR = v.(string)
	}

	if v, ok := containerNetworkConfigRawMap["cluster_pod_cidr_ipv6"]; ok && v.(string) != "" {
		config.ClusterPodCIDRIPv6 = v.(string)
	}

	if v, ok := containerNetworkConfigRawMap["cluster_ip_service_cidr"]; ok && v.(string) != "" {
		config.ClusterIPServiceCIDR = v.(string)
	}

	if v, ok := containerNetworkConfigRawMap["cluster_ip_service_cidr_ipv6"]; ok && v.(string) != "" {
		config.ClusterIPServiceCIDRIPv6 = v.(string)
	}

	if v, ok := containerNetworkConfigRawMap["max_pods_per_node"]; ok {
		config.MaxPodsPerNode = v.(int)
	}

	if v, ok := containerNetworkConfigRawMap["kube_proxy_mode"]; ok && v.(string) != "" {
		config.KubeProxyMode = ccev2types.KubeProxyMode(v.(string))
	}

	return config, nil
}

func buildInstanceSpec(instanceSpecRawMap map[string]interface{}) (*ccev2types.InstanceSpec, error) {
	instanceSpec := &ccev2types.InstanceSpec{}

	if v, ok := instanceSpecRawMap["cce_instance_id"]; ok && v.(string) != "" {
		instanceSpec.CCEInstanceID = v.(string)
	}

	if v, ok := instanceSpecRawMap["instance_name"]; ok && v.(string) != "" {
		instanceSpec.InstanceName = v.(string)
	}

	if v, ok := instanceSpecRawMap["runtime_type"]; ok && v.(string) != "" {
		instanceSpec.RuntimeType = ccev2types.RuntimeType(v.(string))
	}

	if v, ok := instanceSpecRawMap["runtime_version"]; ok && v.(string) != "" {
		instanceSpec.RuntimeVersion = v.(string)
	}

	if v, ok := instanceSpecRawMap["cluster_id"]; ok && v.(string) != "" {
		instanceSpec.ClusterID = v.(string)
	}

	if v, ok := instanceSpecRawMap["cluster_role"]; ok && v.(string) != "" {
		instanceSpec.ClusterRole = ccev2types.ClusterRole(v.(string))
	}

	if v, ok := instanceSpecRawMap["instance_group_id"]; ok && v.(string) != "" {
		instanceSpec.InstanceGroupID = v.(string)
	}

	if v, ok := instanceSpecRawMap["instance_group_name"]; ok && v.(string) != "" {
		instanceSpec.InstanceGroupName = v.(string)
	}

	if v, ok := instanceSpecRawMap["master_type"]; ok && v.(string) != "" {
		instanceSpec.MasterType = ccev2types.MasterType(v.(string))
	}

	if v, ok := instanceSpecRawMap["existed"]; ok {
		instanceSpec.Existed = v.(bool)
	}

	if v, ok := instanceSpecRawMap["existed_option"]; ok && len(v.([]interface{})) == 1 {
		existedOptionRaw := v.([]interface{})[0].(map[string]interface{})
		existedOption, err := buildExistedOption(existedOptionRaw)
		if err != nil {
			log.Println("Build InstanceSpec ExistedOption Error:" + err.Error())
			return nil, err
		}
		instanceSpec.ExistedOption = *existedOption
	}

	if v, ok := instanceSpecRawMap["machine_type"]; ok && v.(string) != "" {
		instanceSpec.MachineType = ccev2types.MachineType(v.(string))
	}

	if v, ok := instanceSpecRawMap["instance_type"]; ok && v.(string) != "" {
		instanceSpec.InstanceType = bccapi.InstanceType(v.(string))
	}

	if v, ok := instanceSpecRawMap["bbc_option"]; ok && len(v.([]interface{})) == 1 {
		bbcOptionRaw := v.([]interface{})[0].(map[string]interface{})
		bbcOption, err := buildBBCOption(bbcOptionRaw)
		if err != nil {
			log.Println("Build InstanceSpec BCCOption Error:" + err.Error())
			return nil, err
		}
		instanceSpec.BBCOption = bbcOption
	}

	if v, ok := instanceSpecRawMap["hpas_option"]; ok && len(v.([]interface{})) == 1 {
		hpasOptionRaw := v.([]interface{})[0].(map[string]interface{})
		hpasOption, err := buildHPASOption(hpasOptionRaw)
		if err != nil {
			log.Println("Build InstanceSpec HPASOption Error:" + err.Error())
			return nil, err
		}
		instanceSpec.HPASOption = hpasOption
		instanceSpec.InstanceResource = ccev2types.InstanceResource{
			MachineSpec: hpasOption.AppType + "/" + hpasOption.AppPerformanceLevel,
		}
	}

	if v, ok := instanceSpecRawMap["ehc_cluster_id"]; ok && v.(string) != "" {
		instanceSpec.EhcClusterID = v.(string)
	}

	if v, ok := instanceSpecRawMap["vpc_config"]; ok && len(v.([]interface{})) == 1 {
		vpcConfigRaw := v.([]interface{})[0].(map[string]interface{})
		vpcConfig, err := buildVPCConfig(vpcConfigRaw)
		if err != nil {
			log.Println("Build InstanceSpec VPCConfig Error:" + err.Error())
			return nil, err
		}
		instanceSpec.VPCConfig = *vpcConfig
	}

	if v, ok := instanceSpecRawMap["instance_resource"]; ok && len(v.([]interface{})) == 1 {
		instanceResourceRaw := v.([]interface{})[0].(map[string]interface{})
		instanceResource, err := buildInstanceResource(instanceResourceRaw)
		if err != nil {
			log.Println("Build InstanceSpec InstanceResource Error:" + err.Error())
			return nil, err
		}
		instanceSpec.InstanceResource = *instanceResource
	}

	if v, ok := instanceSpecRawMap["image_id"]; ok && v.(string) != "" {
		instanceSpec.ImageID = v.(string)
	}

	if v, ok := instanceSpecRawMap["instance_os"]; ok && len(v.([]interface{})) == 1 {
		instanceOSRaw := v.([]interface{})[0].(map[string]interface{})
		instanceOS, err := buildInstanceOS(instanceOSRaw)
		if err != nil {
			log.Println("Build InstanceSpec InstanceOS Error:" + err.Error())
			return nil, err
		}
		instanceSpec.InstanceOS = *instanceOS
	}

	if v, ok := instanceSpecRawMap["need_eip"]; ok {
		instanceSpec.NeedEIP = v.(bool)
	}

	if v, ok := instanceSpecRawMap["eip_option"]; ok && len(v.([]interface{})) == 1 {
		eipOptionRaw := v.([]interface{})[0].(map[string]interface{})
		eipOption, err := buildEIPOption(eipOptionRaw)
		if err != nil {
			log.Println("Build InstanceSpec EIPOption Error:" + err.Error())
			return nil, err
		}
		instanceSpec.EIPOption = eipOption
	}

	if v, ok := instanceSpecRawMap["admin_password"]; ok && v.(string) != "" {
		instanceSpec.AdminPassword = v.(string)
	}

	if v, ok := instanceSpecRawMap["ssh_key_id"]; ok && v.(string) != "" {
		instanceSpec.SSHKeyID = v.(string)
	}

	if v, ok := instanceSpecRawMap["instance_charging_type"]; ok && v.(string) != "" {
		instanceSpec.InstanceChargingType = bccapi.PaymentTimingType(v.(string))
	}

	if v, ok := instanceSpecRawMap["instance_precharging_option"]; ok && len(v.([]interface{})) == 1 {
		instancePrechargingOptionRaw := v.([]interface{})[0].(map[string]interface{})
		instancePrechargingOption, err := buildInstancePrechargingOption(instancePrechargingOptionRaw)
		if err != nil {
			log.Println("Build InstanceSpec InstancePreChargingOption Error:" + err.Error())
			return nil, err
		}
		instanceSpec.InstancePreChargingOption = *instancePrechargingOption
	}

	if v, ok := instanceSpecRawMap["delete_option"]; ok && len(v.([]interface{})) == 1 {
		instanceDeleteOptionRaw := v.([]interface{})[0].(map[string]interface{})
		instanceDeleteOption, err := buildInstanceDeleteOption(instanceDeleteOptionRaw)
		if err != nil {
			log.Println("Build InstanceSpec DeleteOption Error:" + err.Error())
			return nil, err
		}
		instanceSpec.DeleteOption = instanceDeleteOption
	}

	if v, ok := instanceSpecRawMap["deploy_custom_config"]; ok && len(v.([]interface{})) == 1 {
		deployCustomOptionRaw := v.([]interface{})[0].(map[string]interface{})
		deployCustomOption, err := buildDeployCustomConfig(deployCustomOptionRaw)
		if err != nil {
			log.Println("Build InstanceSpec DeployCustomConfig Error:" + err.Error())
			return nil, err
		}
		instanceSpec.DeployCustomConfig = *deployCustomOption
	}

	if v, ok := instanceSpecRawMap["tag_list"]; ok {
		tagList, err := buildTags(v.([]interface{}))
		if err != nil {
			log.Println("Build InstanceSpec Tags Error:" + err.Error())
			return nil, err
		}
		instanceSpec.Tags = tagList
	}

	if v, ok := instanceSpecRawMap["labels"]; ok {
		labels := make(map[string]string)
		for key, value := range v.(map[string]interface{}) {
			labels[key] = value.(string)
		}
		instanceSpec.Labels = labels
	}

	if v, ok := instanceSpecRawMap["instance_taints"]; ok {
		taintList, err := buildTaints(v.([]interface{}))
		if err != nil {
			log.Println("Build InstanceSpec Taints Error:" + err.Error())
			return nil, err
		}
		instanceSpec.Taints = taintList
	}

	if v, ok := instanceSpecRawMap["cce_instance_priority"]; ok {
		instanceSpec.CCEInstancePriority = v.(int)
	}

	return instanceSpec, nil
}

func buildTaints(taintRawMapList []interface{}) ([]ccev2types.Taint, error) {
	taintList := make([]ccev2types.Taint, 0)

	for _, taintRaw := range taintRawMapList {
		taintRawMap := taintRaw.(map[string]interface{})

		taint := ccev2types.Taint{}

		if v, ok := taintRawMap["key"]; ok && v.(string) != "" {
			taint.Key = v.(string)
		}

		if v, ok := taintRawMap["value"]; ok && v.(string) != "" {
			taint.Value = v.(string)
		}

		if v, ok := taintRawMap["effect"]; ok && v.(string) != "" {
			taint.Effect = ccev2types.TaintEffect(v.(string))
		}

		if v, ok := taintRawMap["time_added"]; ok && v.(string) != "" {
			//time format RFC3339
			taint.TimeAdded = &ccev2types.Time{}
			err := taint.TimeAdded.UnmarshalQueryParameter(v.(string))
			if err != nil {
				log.Println("Taint TimeAdded Format Error:" + err.Error())
				return nil, err
			}
		}

		taintList = append(taintList, taint)
	}

	return taintList, nil
}

func buildTags(tagRawMapList []interface{}) ([]ccev2types.Tag, error) {
	tagList := make([]ccev2types.Tag, 0)

	for _, tagRaw := range tagRawMapList {
		tagRawMap := tagRaw.(map[string]interface{})

		tag := ccev2types.Tag{}

		if v, ok := tagRawMap["tag_key"]; ok && v.(string) != "" {
			tag.TagKey = v.(string)
		}

		if v, ok := tagRawMap["tag_value"]; ok && v.(string) != "" {
			tag.TagValue = v.(string)
		}

		tagList = append(tagList, tag)
	}
	return tagList, nil
}

func buildK8SCustomConfig(k8sCustomConfigRawMap map[string]interface{}) (*ccev2types.K8SCustomConfig, error) {
	config := &ccev2types.K8SCustomConfig{}

	if v, ok := k8sCustomConfigRawMap["master_feature_gates"]; ok {
		masterFeatureGates := make(map[string]bool)
		for key, value := range v.(map[string]interface{}) {
			masterFeatureGates[key] = value.(bool)
		}
		config.MasterFeatureGates = masterFeatureGates
	}

	if v, ok := k8sCustomConfigRawMap["node_feature_gates"]; ok {
		nodeFeatureGates := make(map[string]bool)
		for key, value := range v.(map[string]interface{}) {
			nodeFeatureGates[key] = value.(bool)
		}
		config.NodeFeatureGates = nodeFeatureGates
	}

	if v, ok := k8sCustomConfigRawMap["admission_plugins"]; ok && v != nil {
		admissionPlugins := make([]string, 0)
		for _, plugin := range v.([]interface{}) {
			admissionPlugins = append(admissionPlugins, plugin.(string))
		}
		config.AdmissionPlugins = admissionPlugins
	}

	if v, ok := k8sCustomConfigRawMap["pause_image"]; ok && v.(string) != "" {
		config.PauseImage = v.(string)
	}

	if v, ok := k8sCustomConfigRawMap["kube_api_qps"]; ok {
		config.KubeAPIQPS = v.(int)
	}

	if v, ok := k8sCustomConfigRawMap["kube_api_burst"]; ok {
		config.KubeAPIBurst = v.(int)
	}

	if v, ok := k8sCustomConfigRawMap["scheduler_predicated"]; ok && v != nil {
		schedulerPredicates := make([]string, 0)
		for _, schedulerPredicate := range v.([]interface{}) {
			schedulerPredicates = append(schedulerPredicates, schedulerPredicate.(string))
		}
		config.SchedulerPredicates = schedulerPredicates
	}

	if v, ok := k8sCustomConfigRawMap["scheduler_priorities"]; ok {
		schedulerPriority := make(map[string]int)
		for key, value := range v.(map[string]interface{}) {
			schedulerPriority[key] = value.(int)
		}
		config.SchedulerPriorities = schedulerPriority
	}

	if v, ok := k8sCustomConfigRawMap["etcd_data_path"]; ok && v.(string) != "" {
		config.ETCDDataPath = v.(string)
	}

	return config, nil
}

func buildDeployCustomConfig(deployCustomConfigRawMap map[string]interface{}) (*ccev2types.DeployCustomConfig, error) {
	option := &ccev2types.DeployCustomConfig{}

	if v, ok := deployCustomConfigRawMap["docker_config"]; ok && len(v.([]interface{})) == 1 {
		dockerConfigRaw := v.([]interface{})[0].(map[string]interface{})
		dockerConfigOption, err := buildDockerConfig(dockerConfigRaw)
		if err != nil {
			log.Printf("Build Docker Config Fail:" + err.Error())
			return nil, err
		}
		option.DockerConfig = *dockerConfigOption
	}

	if v, ok := deployCustomConfigRawMap["kubelet_root_dir"]; ok && v.(string) != "" {
		option.KubeletRootDir = v.(string)
	}

	if v, ok := deployCustomConfigRawMap["enable_resource_reserved"]; ok {
		option.EnableResourceReserved = v.(bool)
	}

	if v, ok := deployCustomConfigRawMap["kube_reserved"]; ok {
		kubeReserved := make(map[string]string)
		for key, value := range v.(map[string]interface{}) {
			kubeReserved[key] = value.(string)
		}
		option.KubeReserved = kubeReserved
	}

	if v, ok := deployCustomConfigRawMap["enable_cordon"]; ok {
		option.EnableCordon = v.(bool)
	}

	if v, ok := deployCustomConfigRawMap["pre_user_script"]; ok && v.(string) != "" {
		option.PreUserScript = v.(string)
	}

	if v, ok := deployCustomConfigRawMap["post_user_script"]; ok && v.(string) != "" {
		option.PostUserScript = v.(string)
	}

	return option, nil
}

func buildDockerConfig(dockerConfigRawMap map[string]interface{}) (*ccev2types.DockerConfig, error) {
	config := &ccev2types.DockerConfig{}

	if v, ok := dockerConfigRawMap["docker_data_root"]; ok && v.(string) != "" {
		config.DockerDataRoot = v.(string)
	}

	if v, ok := dockerConfigRawMap["registry_mirrors"]; ok && v != nil {
		registryMirrors := make([]string, 0)
		for _, mirrorsRaw := range v.([]interface{}) {
			registryMirrors = append(registryMirrors, mirrorsRaw.(string))
		}
		config.RegistryMirrors = registryMirrors
	}

	if v, ok := dockerConfigRawMap["insecure_registries"]; ok && v != nil {
		registries := make([]string, 0)
		for _, registriesRaw := range v.([]interface{}) {
			registries = append(registries, registriesRaw.(string))
		}
		config.RegistryMirrors = registries
	}

	if v, ok := dockerConfigRawMap["docker_log_max_size"]; ok && v.(string) != "" {
		config.DockerLogMaxSize = v.(string)
	}

	if v, ok := dockerConfigRawMap["docker_log_max_file"]; ok && v.(string) != "" {
		config.DockerLogMaxFile = v.(string)
	}

	if v, ok := dockerConfigRawMap["bip"]; ok && v.(string) != "" {
		config.BIP = v.(string)
	}

	return config, nil
}

func buildInstanceDeleteOption(d map[string]interface{}) (*ccev2types.DeleteOption, error) {
	option := &ccev2types.DeleteOption{}

	if v, ok := d["move_out"]; ok {
		option.MoveOut = v.(bool)
	}

	if v, ok := d["delete_resource"]; ok {
		option.DeleteResource = v.(bool)
	}

	if v, ok := d["delete_cds_snapshot"]; ok {
		option.DeleteCDSSnapshot = v.(bool)
	}

	return option, nil
}

func buildInstancePrechargingOption(instancePreChargingOptionRawMap map[string]interface{}) (*ccev2types.InstancePreChargingOption, error) {
	option := &ccev2types.InstancePreChargingOption{}

	if v, ok := instancePreChargingOptionRawMap["purchase_time"]; ok {
		option.PurchaseTime = v.(int)
	}

	if v, ok := instancePreChargingOptionRawMap["auto_renew"]; ok {
		option.AutoRenew = v.(bool)
	}

	if v, ok := instancePreChargingOptionRawMap["auto_renew_time_unit"]; ok && v.(string) != "" {
		option.AutoRenewTimeUnit = v.(string)
	}

	if v, ok := instancePreChargingOptionRawMap["auto_renew_time"]; ok {
		option.AutoRenewTime = v.(int)
	}

	return option, nil
}

func buildExistedOption(d map[string]interface{}) (*ccev2types.ExistedOption, error) {
	existedOption := &ccev2types.ExistedOption{}

	if v, ok := d["existed_instance_id"]; ok && v.(string) != "" {
		existedOption.ExistedInstanceID = v.(string)
	}

	if v, ok := d["rebuild"]; ok {
		existedOption.Rebuild = v.(*bool)
	}

	return existedOption, nil
}

func buildEIPOption(d map[string]interface{}) (*ccev2types.EIPOption, error) {
	eipOption := &ccev2types.EIPOption{}

	if v, ok := d["eip_name"]; ok && v.(string) != "" {
		eipOption.EIPName = v.(string)
	}

	if v, ok := d["eip_charging_type"]; ok && v.(string) != "" {
		eipOption.EIPChargingType = ccev2types.BillingMethod(v.(string))
	}

	if v, ok := d["eip_bandwidth"]; ok {
		eipOption.EIPBandwidth = v.(int)
	}

	return eipOption, nil
}

func buildBBCOption(d map[string]interface{}) (*ccev2types.BBCOption, error) {
	bbcOption := &ccev2types.BBCOption{}

	if v, ok := d["reserve_data"]; ok {
		bbcOption.ReserveData = v.(bool)
	}

	if v, ok := d["raid_id"]; ok && v.(string) != "" {
		bbcOption.RaidID = v.(string)
	}

	if v, ok := d["sys_disk_size"]; ok {
		bbcOption.SysDiskSize = v.(int)
	}

	return bbcOption, nil
}

func buildHPASOption(d map[string]interface{}) (*ccev2types.HPASOption, error) {
	hpasOption := &ccev2types.HPASOption{}

	if v, ok := d["app_type"]; ok && v.(string) != "" {
		hpasOption.AppType = v.(string)
	}
	if v, ok := d["app_performance_level"]; ok && v.(string) != "" {
		hpasOption.AppPerformanceLevel = v.(string)
	}
	return hpasOption, nil
}

func buildVPCConfig(vpcRawMap map[string]interface{}) (*ccev2types.VPCConfig, error) {
	vpcConfig := &ccev2types.VPCConfig{}

	if v, ok := vpcRawMap["vpc_id"]; ok && v.(string) != "" {
		vpcConfig.VPCID = v.(string)
	}

	if v, ok := vpcRawMap["vpc_subnet_id"]; ok && v.(string) != "" {
		vpcConfig.VPCSubnetID = v.(string)
	}

	if v, ok := vpcRawMap["security_group_id"]; ok && v.(string) != "" {
		vpcConfig.SecurityGroupID = v.(string)
	}

	if v, ok := vpcRawMap["security_group_type"]; ok && v.(string) != "" {
		vpcConfig.SecurityGroupType = v.(string)
	}

	if v, ok := vpcRawMap["vpc_subnet_type"]; ok && v.(string) != "" {
		vpcConfig.VPCSubnetType = vpc.SubnetType(v.(string))
	}

	if v, ok := vpcRawMap["vpc_subnet_cidr"]; ok && v.(string) != "" {
		vpcConfig.VPCSubnetCIDR = v.(string)
	}

	if v, ok := vpcRawMap["vpc_subnet_cidr_ipv6"]; ok && v.(string) != "" {
		vpcConfig.VPCSubnetCIDRIPv6 = v.(string)
	}

	if v, ok := vpcRawMap["available_zone"]; ok && v.(string) != "" {
		vpcConfig.AvailableZone = ccev2types.AvailableZone(v.(string))
	}

	return vpcConfig, nil
}

func buildInstanceResource(d map[string]interface{}) (*ccev2types.InstanceResource, error) {
	instanceResource := &ccev2types.InstanceResource{}

	if v, ok := d["machine_spec"]; ok && v.(string) != "" {
		instanceResource.MachineSpec = v.(string)
	}

	if v, ok := d["cpu"]; ok {
		instanceResource.CPU = v.(int)
	}

	if v, ok := d["mem"]; ok {
		instanceResource.MEM = v.(int)
	}

	if v, ok := d["node_cpu_quota"]; ok {
		instanceResource.NodeCPUQuota = v.(int)
	}

	if v, ok := d["node_mem_quota"]; ok {
		instanceResource.NodeMEMQuota = v.(int)
	}

	if v, ok := d["root_disk_type"]; ok && v.(string) != "" {
		instanceResource.RootDiskType = bccapi.StorageType(v.(string))
	}

	if v, ok := d["root_disk_size"]; ok {
		instanceResource.RootDiskSize = v.(int)
	}

	if v, ok := d["local_disk_size"]; ok {
		instanceResource.LocalDiskSize = v.(int)
	}

	if v, ok := d["cds_list"]; ok {
		cdsList, err := buildCDSList(v.([]interface{}))
		if err != nil {
			return nil, err
		}
		instanceResource.CDSList = cdsList
	}

	if v, ok := d["ephemeral_disk_list"]; ok {
		ephemeralDiskList, err := buildEphemeralDiskList(v.([]interface{}))
		if err != nil {
			return nil, err
		}
		instanceResource.EphemeralDiskList = ephemeralDiskList
	}

	if v, ok := d["gpu_type"]; ok && v.(string) != "" {
		instanceResource.GPUType = ccev2types.GPUType(v.(string))
	}

	if v, ok := d["gpu_count"]; ok {
		instanceResource.GPUCount = v.(int)
	}

	return instanceResource, nil
}

func buildCDSList(cdsRawList []interface{}) ([]ccev2types.CDSConfig, error) {
	cdsList := make([]ccev2types.CDSConfig, 0)

	for _, cdsConfigRaw := range cdsRawList {
		cdsConfigRawMap := cdsConfigRaw.(map[string]interface{})

		config := ccev2types.CDSConfig{}
		if v, ok := cdsConfigRawMap["path"]; ok && v.(string) != "" {
			config.Path = v.(string)
		}
		if v, ok := cdsConfigRawMap["storage_type"]; ok && v.(string) != "" {
			config.StorageType = bccapi.StorageType(v.(string))
		}
		if v, ok := cdsConfigRawMap["cds_size"]; ok {
			config.CDSSize = v.(int)
		}
		if v, ok := cdsConfigRawMap["snapshot_id"]; ok && v.(string) != "" {
			config.SnapshotID = v.(string)
		}

		cdsList = append(cdsList, config)
	}
	return cdsList, nil
}

func buildEphemeralDiskList(ephemeralDiskRawList []interface{}) ([]ccev2types.EphemeralDiskConfig, error) {
	ephemeralDiskList := make([]ccev2types.EphemeralDiskConfig, 0)

	for _, ephemeralDiskRaw := range ephemeralDiskRawList {
		ephemeralDiskRawMap := ephemeralDiskRaw.(map[string]interface{})

		disk := ccev2types.EphemeralDiskConfig{}
		if v, ok := ephemeralDiskRawMap["storage_type"]; ok && v.(string) != "" {
			disk.StorageType = ccev2types.StorageType(v.(string))
		}
		if v, ok := ephemeralDiskRawMap["size_in_gb"]; ok {
			disk.SizeInGB = v.(int)
		}
		if v, ok := ephemeralDiskRawMap["disk_path"]; ok && v.(string) != "" {
			disk.Path = v.(string)
		}

		ephemeralDiskList = append(ephemeralDiskList, disk)
	}
	return ephemeralDiskList, nil
}

func buildInstanceOS(d map[string]interface{}) (*ccev2types.InstanceOS, error) {
	instanceOS := &ccev2types.InstanceOS{}

	if v, ok := d["image_type"]; ok && v.(string) != "" {
		instanceOS.ImageType = bccapi.ImageType(v.(string))
	}

	if v, ok := d["image_name"]; ok && v.(string) != "" {
		instanceOS.ImageName = v.(string)
	}

	if v, ok := d["os_type"]; ok && v.(string) != "" {
		instanceOS.OSType = ccev2types.OSType(v.(string))
	}

	if v, ok := d["os_name"]; ok && v.(string) != "" {
		instanceOS.OSName = ccev2types.OSName(v.(string))
	}

	if v, ok := d["os_version"]; ok && v.(string) != "" {
		instanceOS.OSVersion = v.(string)
	}

	if v, ok := d["os_arch"]; ok && v.(string) != "" {
		instanceOS.OSArch = v.(string)
	}

	if v, ok := d["os_build"]; ok && v.(string) != "" {
		instanceOS.OSBuild = v.(string)
	}

	return instanceOS, nil
}

func buildCCEv2CreateClusterArgs(d *schema.ResourceData) (*ccev2.CreateClusterArgs, error) {
	argsRequest := &ccev2.CreateClusterRequest{}

	clusterSpecRaw := d.Get("cluster_spec.0").(map[string]interface{})
	clusterSpec, err := buildCCEv2CreateClusterClusterSpec(clusterSpecRaw)
	if v, ok := d.GetOk("tags"); ok {
		clusterSpec.Tags = tranceCCETagMapToModel(v.(map[string]interface{}))
	}
	if err != nil {
		log.Printf("Build CreateClusterArgs ClusterSpec Fail:" + err.Error())
		return nil, err
	}
	argsRequest.ClusterSpec = clusterSpec

	masterSpecsRaw := d.Get("master_specs").([]interface{})
	instanceSets := make([]*ccev2.InstanceSet, 0)
	for _, masterSpecSetRaw := range masterSpecsRaw {

		masterSpecSetMap := masterSpecSetRaw.(map[string]interface{})
		masterSpecRaw := masterSpecSetMap["master_spec"].([]interface{})[0]

		masterSpec, err := buildInstanceSpec(masterSpecRaw.(map[string]interface{}))
		if err != nil {
			log.Printf("Build CreateClusterArgs MasterSpecs Fail:" + err.Error())
			return nil, err
		}
		instanceSet := &ccev2.InstanceSet{
			Count:        masterSpecSetMap["count"].(int),
			InstanceSpec: *masterSpec,
		}
		instanceSets = append(instanceSets, instanceSet)
	}
	argsRequest.MasterSpecs = instanceSets

	return &ccev2.CreateClusterArgs{
		CreateClusterRequest: argsRequest,
	}, nil
}

func buildCCEv2DeleteClusterArgs(d *schema.ResourceData) (*ccev2.DeleteClusterArgs, error) {
	args := &ccev2.DeleteClusterArgs{
		ClusterID: d.Id(),
	}

	if v, ok := d.GetOk("cluster_spec.0.cluster_delete_option.0.delete_resource"); ok {
		args.DeleteResource = v.(bool)
	}

	if v, ok := d.GetOk("cluster_spec.0.cluster_delete_option.0.delete_cds_snapshot"); ok {
		args.DeleteCDSSnapshot = v.(bool)
	}

	return args, nil
}

func tranceCCETagMapToModel(tagMaps map[string]interface{}) []ccev2types.Tag {
	tags := make([]ccev2types.Tag, 0, len(tagMaps))
	for k, v := range tagMaps {
		tags = append(tags, ccev2types.Tag{
			TagKey:   k,
			TagValue: v.(string),
		})
	}
	return tags
}

func flattenCCETagsToMap(tags []ccev2types.Tag) map[string]string {
	tagMap := make(map[string]string)
	for _, tag := range tags {
		tagMap[tag.TagKey] = tag.TagValue
	}
	return tagMap
}

// 判断两个tag切片是否包含相同的元素
func slicesContainSameElementsInCCETags(a, b []ccev2types.Tag) bool {
	if len(a) != len(b) {
		return false
	}
	// 创建映射来存储每个 TagModel 出现的次数
	counts := make(map[ccev2types.Tag]int)
	// 计算第一个切片中每个元素出现的次数
	for _, item := range a {
		counts[item]++
	}
	// 减去第二个切片中每个元素出现的次数
	for _, item := range b {
		counts[item]--
		if counts[item] < 0 {
			// 如果某个元素在第二个切片中出现的次数多于第一个切片，返回 false
			return false
		}
	}
	// 检查所有元素的计数是否为 0
	for _, count := range counts {
		if count != 0 {
			return false
		}
	}
	return true
}

var ClusterTypePermitted = []string{
	string(ccev2types.ClusterTypeNormal),
}

var K8SVersionPermitted = []string{
	//string(ccev2types.K8S_1_13_10),
	//string(ccev2types.K8S_1_16_8),
	"1.18.9",
	"1.20.8",
	"1.21.14",
	"1.22.5",
	"1.24.4",
	"1.26.9",
	"1.30.1",
	"1.28.8",
}

var RuntimeTypePermitted = []string{
	string(ccev2types.RuntimeTypeDocker),
	string(ccev2types.RuntimeTypeContainerd),
}

var MasterTypePermitted = []string{
	string(ccev2types.MasterTypeManaged),
	string(ccev2types.MasterTypeCustom),
	string(ccev2types.MasterTypeServerless),
}

var ClusterHAPermitted = []int{
	int(ccev2types.ClusterHALow),
	int(ccev2types.ClusterHAMedium),
	int(ccev2types.ClusterHAHigh),
	int(ccev2types.ClusterHAServerless),
}

var AvailableZonePermitted = []string{
	string(ccev2types.AvailableZoneA),
	string(ccev2types.AvailableZoneB),
	string(ccev2types.AvailableZoneC),
	string(ccev2types.AvailableZoneD),
	string(ccev2types.AvailableZoneE),
	string(ccev2types.AvailableZoneF),
}

var ContainerNetworkModePermitted = []string{
	string(ccev2types.ContainerNetworkModeKubenet),
	string(ccev2types.ContainerNetworkModeVPCCNI),
	string(ccev2types.ContainerNetworkModeVPCRouteVeth),
	string(ccev2types.ContainerNetworkModeVPCRouteIPVlan),
	string(ccev2types.ContainerNetworkModeVPCRouteAutoDetect),
	string(ccev2types.ContainerNetworkModeVPCSecondaryIPVeth),
	string(ccev2types.ContainerNetworkModeVPCSecondaryIPIPVlan),
	string(ccev2types.ContainerNetworkModeVPCSecondaryIPAutoDetect),
}

var ContainerNetworkIPTypePermitted = []string{
	string(ccev2types.ContainerNetworkIPTypeIPv4),
	string(ccev2types.ContainerNetworkIPTypeIPv6),
	string(ccev2types.ContainerNetworkIPTypeDualStack),
}

var KubeProxyModePermitted = []string{
	string(ccev2types.KubeProxyModeIptables),
	string(ccev2types.KubeProxyModeIPVS),
}

var ClusterRolePermitted = []string{
	string(ccev2types.ClusterRoleMaster),
	string(ccev2types.ClusterRoleNode),
}

var MachineTypePermitted = []string{
	string(ccev2types.MachineTypeBBC),
	string(ccev2types.MachineTypeBCC),
	string(ccev2types.MachineTypeEBC),
	string(ccev2types.MachineTypeHPAS),
}

var BCCInstanceTypePermitted = []string{
	string(bccapi.InstanceTypeN1),
	string(bccapi.InstanceTypeN2),
	string(bccapi.InstanceTypeN3),
	string(bccapi.InstanceTypeN4),
	string(bccapi.InstanceTypeN5),
	string(bccapi.InstanceTypeC1),
	string(bccapi.InstanceTypeC2),
	string(bccapi.InstanceTypeS1),
	string(bccapi.InstanceTypeG1),
	string(bccapi.InstanceTypeF1),
	string(ccev2types.MachineTypeHPAS),
}

var VPCSubnetTypePermitted = []string{
	string(vpc.SUBNET_TYPE_BCC),
	string(vpc.SUBNET_TYPE_BCCNAT),
	string(vpc.SUBNET_TYPE_BBC),
}

var StorageTypePermitted = []string{
	string(bccapi.StorageTypeStd1),     //  "上一代云磁盘, sata 盘"
	string(bccapi.StorageTypeHP1),      //  "高性能云磁盘, ssd 盘"
	string(bccapi.StorageTypeCloudHP1), //  "SSD 云磁盘, premium ssd 盘"
	string(bccapi.StorageTypeLocal),    //  "本地盘"
	string(bccapi.StorageTypeSATA),     //  "sata盘, 创建 DCC 子网实例专用"
	string(bccapi.StorageTypeSSD),      //  "ssd盘, 创建 DCC 子网实例专用"
	//bccapi.StorageTypeHDDThroughput,
	string(bccapi.StorageTypeHdd), //  "普通型"
}

var GPUTypePermitted = []string{
	string(ccev2types.GPUTypeV100_32), //  NVIDIA Tesla V100-32G
	string(ccev2types.GPUTypeV100_16), //  NVIDIA Tesla V100-16G
	string(ccev2types.GPUTypeP40),     //  NVIDIA Tesla P40
	string(ccev2types.GPUTypeP4),      //   NVIDIA Tesla P4
	string(ccev2types.GPUTypeK40),     //  NVIDIA Tesla K40
	string(ccev2types.GPUTypeDLCard),  //  NVIDIA 深度学习开发卡
}

var ImageTypePermitted = []string{
	string(bccapi.ImageTypeAll),         //  所有镜像类型
	string(bccapi.ImageTypeSystem),      //  "系统镜像/公共镜像"
	string(bccapi.ImageTypeCustom),      //  "自定义镜像"
	string(bccapi.ImageTypeIntegration), //  "服务集成镜像"
	string(bccapi.ImageTypeSharing),     //  共享镜像
	string(bccapi.ImageTypeGPUSystem),   //  gpu公有
	string(bccapi.ImageTypeGPUCustom),   //  gpu 自定义
	string(bccapi.ImageTypeBBCSystem),   //  BBC 公有
	string(bccapi.ImageTypeBBCCustom),   //  BBC 自定义
}

var OSTypePermitted = []string{
	string(ccev2types.OSTypeLinux),
	string(ccev2types.OSTypeWindows),
}

var OSNamePermitted = []string{
	string(ccev2types.OSNameCentOS),
	string(ccev2types.OSNameUbuntu),
	string(ccev2types.OSNameWindows),
	string(ccev2types.OSNameDebian),
	string(ccev2types.OSNameOpensuse),
}

var TaintEffectPermitted = []string{
	string(ccev2types.TaintEffectNoSchedule),
	string(ccev2types.TaintEffectPreferNoSchedule),
	string(ccev2types.TaintEffectNoExecute),
}

var EIPBillingMethodPermitted = []string{
	string(ccev2types.BillingMethodByTraffic),   //按照流量计费
	string(ccev2types.BillingMethodByBandwidth), //按带宽计费
}

var PaymentTimingTypePermitted = []string{
	string(bccapi.PaymentTimingPrePaid),
	string(bccapi.PaymentTimingPostPaid),
	string(bccapi.PaymentTimingBidding),
}

var QueryOrderPermitted = []string{
	string(ccev2.OrderASC),
	string(ccev2.OrderDESC),
}

var InstanceQueryKeywordTypePermitted = []string{
	string(ccev2.InstanceKeywordTypeInstanceName),
	string(ccev2.InstanceKeywordTypeInstanceID),
}

var InstanceQueryOrderByPermitted = []string{
	string(ccev2.InstanceOrderByInstanceName),
	string(ccev2.InstanceOrderByInstanceID),
	string(ccev2.InstanceOrderByCreatedAt),
}
