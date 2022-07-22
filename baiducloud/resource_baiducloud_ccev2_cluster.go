/*
Use this resource to create a CCEv2 cluster.

Example Usage

```hcl
resource "baiducloud_ccev2_cluster" "default_managed" {
  cluster_spec  {
    cluster_name = var.cluster_name
    cluster_type = "normal"
    k8s_version = "1.16.8"
    runtime_type = "docker"
    vpc_id = baiducloud_vpc.default.id
    plugins = ["core-dns", "kube-proxy"]
    master_config {
      master_type = "managed"
      cluster_ha = 2
      exposed_public = false
      cluster_blb_vpc_subnet_id = baiducloud_subnet.defaultA.id
      managed_cluster_master_option {
        master_vpc_subnet_zone = "zoneA"
      }
    }
    container_network_config  {
      mode = "kubenet"
      lb_service_vpc_subnet_id = baiducloud_subnet.defaultA.id
      node_port_range_min = 30000
      node_port_range_max = 32767
      max_pods_per_node = 64
      cluster_pod_cidr = var.cluster_pod_cidr
      cluster_ip_service_cidr = var.cluster_ip_service_cidr
      ip_version = "ipv4"
      kube_proxy_mode = "iptables"
    }
    cluster_delete_option {
      delete_resource = true
      delete_cds_snapshot = true
    }
  }
}
```
*/
package baiducloud

import (
	"errors"
	"log"
	"time"

	ccev2 "github.com/baidubce/bce-sdk-go/services/cce/v2"
	ccev2types "github.com/baidubce/bce-sdk-go/services/cce/v2/types"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudCCEv2Cluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudCCEv2ClusterCreate,
		Read:   resourceBaiduCloudCCEv2ClusterRead,
		Delete: resourceBaiduCloudCCEv2ClusterDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			//Params for creating the cluster
			"cluster_spec": {
				Type:        schema.TypeList,
				Description: "Specification of the cluster",
				Required:    true,
				ForceNew:    true,
				MaxItems:    1,
				Elem:        resourceCCEv2ClusterSpec(),
			},
			"master_specs": {
				Type:        schema.TypeList,
				Description: "Specification of master nodes cluster",
				Optional:    true,
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": {
							Type:        schema.TypeInt,
							Description: "Count of this type master",
							Required:    true,
						},
						"master_spec": {
							Type:        schema.TypeList,
							Description: "Count of this type master",
							Required:    true,
							MaxItems:    1,
							Elem:        resourceCCEv2InstanceSpec(),
						},
					},
				},
			},
			//Status of the cluster
			"cluster_status": {
				Type:        schema.TypeList,
				Description: "Statue of the cluster",
				Computed:    true,
				MaxItems:    1,
				Elem:        resourceCCEv2ClusterStatus(),
			},
			"created_at": {
				Type:        schema.TypeString,
				Description: "Create time of the cluster",
				Computed:    true,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Description: "Update time of the cluster",
				Computed:    true,
			},
			"masters": {
				Type:        schema.TypeList,
				Description: "Master machines of the cluster",
				Computed:    true,
				Elem:        resourceCCEv2Instance(),
			},
			"nodes": {
				Type:        schema.TypeList,
				Description: "Slave machines of the cluster",
				Computed:    true,
				Elem:        resourceCCEv2Instance(),
			},
		},
	}
}

func resourceBaiduCloudCCEv2ClusterCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.BaiduClient)
	ccev2Service := Ccev2Service{client}

	createClusterArgs, err := buildCCEv2CreateClusterArgs(d)
	if err != nil {
		log.Printf("Build CreateClusterArgs Error:" + err.Error())
		return WrapError(err)
	}

	action := "Create CCEv2 cluster " + createClusterArgs.CreateClusterRequest.ClusterSpec.ClusterName
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithCCEv2Client(func(client *ccev2.Client) (interface{}, error) {
			return client.CreateCluster(createClusterArgs)
		})
		if err != nil {
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw)
		response, ok := raw.(*ccev2.CreateClusterResponse)
		if !ok {
			err = errors.New("response format illegal")
			return resource.NonRetryableError(err)
		}
		d.SetId(response.ClusterID)
		return nil
	})
	if err != nil {
		log.Printf("Create Cluster Error:" + err.Error())
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_cluster", action, BCESDKGoERROR)
	}

	stateConf := buildStateConf(
		[]string{string(ccev2types.ClusterPhasePending), string(ccev2types.ClusterPhaseProvisioning),
			string(ccev2types.ClusterPhaseProvisioned)},
		[]string{string(ccev2types.ClusterPhaseRunning)},
		d.Timeout(schema.TimeoutCreate),
		ccev2Service.ClusterStateRefreshCCEv2(d.Id(), []ccev2types.ClusterPhase{
			ccev2types.ClusterPhaseCreateFailed,
			ccev2types.ClusterPhaseDeleteFailed,
			ccev2types.ClusterPhaseDeleting,
			ccev2types.ClusterPhaseDeleted,
		}),
	)
	if _, err := stateConf.WaitForState(); err != nil {
		log.Printf("Create Cluster Error:" + err.Error())
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_cluster", action, BCESDKGoERROR)
	}

	return resourceBaiduCloudCCEv2ClusterRead(d, meta)
}

func resourceBaiduCloudCCEv2ClusterRead(d *schema.ResourceData, meta interface{}) error {
	clusterId := d.Id()
	action := "Get CCEv2 Cluster " + clusterId
	client := meta.(*connectivity.BaiduClient)

	//1.Get Status of the Cluster
	raw, err := client.WithCCEv2Client(func(client *ccev2.Client) (interface{}, error) {
		return client.GetCluster(clusterId)
	})
	if err != nil {
		if NotFoundError(err) {
			log.Printf("Cluster Not Found. Set Resource ID to Empty.")
			d.SetId("") //Resource Not Found, make the ID of resource to empty to delete it in state file.
			return nil
		}
		log.Printf("Get Cluster Error:" + err.Error())
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_cluster", action, BCESDKGoERROR)
	}
	response := raw.(*ccev2.GetClusterResponse)
	if response == nil {
		err := Error("Response is nil")
		log.Printf("Get Cluster Error:" + err.Error())
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_cluster", action, BCESDKGoERROR)
	}

	clusterStatus, err := convertClusterStatusFromJsonToTfMap(response.Cluster.Status)
	if err != nil {
		log.Printf("Get Cluster Status Error:" + err.Error())
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_cluster", action, BCESDKGoERROR)
	}
	err = d.Set("cluster_status", clusterStatus)
	if err != nil {
		log.Printf("Set cluster_status Error:" + err.Error())
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_cluster", action, BCESDKGoERROR)
	}

	err = d.Set("created_at", response.Cluster.CreatedAt.String())
	if err != nil {
		log.Printf("Set created_at Error:" + err.Error())
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_cluster", action, BCESDKGoERROR)
	}
	err = d.Set("updated_at", response.Cluster.UpdatedAt.String())
	if err != nil {
		log.Printf("Set updated_at Error:" + err.Error())
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_cluster", action, BCESDKGoERROR)
	}

	//2.Get Instances of the Cluster
	listInstancesRaw, err := client.WithCCEv2Client(func(client *ccev2.Client) (interface{}, error) {
		args := &ccev2.ListInstancesByPageArgs{
			ClusterID: clusterId,
			Params: &ccev2.ListInstancesByPageParams{
				KeywordType: ccev2.InstanceKeywordTypeInstanceName,
				Keyword:     "",
				OrderBy:     "createdAt",
				Order:       ccev2.OrderASC,
				PageNo:      1,
				PageSize:    1000,
			},
		}
		return client.ListInstancesByPage(args)
	})
	if err != nil {
		log.Printf("Get Cluster Instance List Error" + err.Error())
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_cluster", action, BCESDKGoERROR)
	}
	listInstanceResponse := listInstancesRaw.(*ccev2.ListInstancesResponse)
	if listInstanceResponse == nil {
		err := Error("ListInstancesResponse is nil")
		log.Printf("Get Cluster Instance Error:" + err.Error())
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_cluster", action, BCESDKGoERROR)
	}

	//masterList
	masters, err := convertInstanceFromJsonToMap(listInstanceResponse.InstancePage.InstanceList, ccev2types.ClusterRoleMaster)
	if err != nil {
		log.Printf("Get Cluster Master Instances Error：" + err.Error())
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_cluster", action, BCESDKGoERROR)
	}
	err = d.Set("masters", masters)
	if err != nil {
		log.Printf("Set masters Error:" + err.Error())
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_cluster", action, BCESDKGoERROR)
	}
	//nodeList
	nodes, err := convertInstanceFromJsonToMap(listInstanceResponse.InstancePage.InstanceList, ccev2types.ClusterRoleNode)
	if err != nil {
		log.Printf("Get Cluster Follower Nodes Error:" + err.Error())
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_cluster", action, BCESDKGoERROR)
	}
	err = d.Set("nodes", nodes)
	if err != nil {
		log.Printf("Set nodes Error:" + err.Error())
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_cluster", action, BCESDKGoERROR)
	}

	return nil
}

func resourceBaiduCloudCCEv2ClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	ccev2Service := Ccev2Service{client}

	args, err := buildCCEv2DeleteClusterArgs(d)
	if err != nil {
		log.Printf("Build DeleteClusterArgs Error:" + err.Error())
		return WrapError(err)
	}

	action := "Delete CCEv2 Cluster " + args.ClusterID
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := client.WithCCEv2Client(func(client *ccev2.Client) (interface{}, error) {
			return client.DeleteCluster(args)
		})
		if err != nil {
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw)
		return nil
	})
	if err != nil {
		log.Printf("Delete Cluster Error:" + err.Error())
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_cluster", action, BCESDKGoERROR)
	}

	stateConf := buildStateConf(
		[]string{string(ccev2types.ClusterPhaseRunning),
			string(ccev2types.ClusterPhaseDeleting),
			string(ccev2types.ClusterPhaseCreateFailed),
			string(ccev2types.ClusterPhaseProvisioned),
			string(ccev2types.ClusterPhaseProvisioning),
			string(ccev2types.ClusterPhaseDeleteFailed),
		},
		[]string{string(ccev2types.ClusterPhaseDeleted)},
		d.Timeout(schema.TimeoutDelete),
		ccev2Service.ClusterStateRefreshCCEv2(args.ClusterID, []ccev2types.ClusterPhase{
			ccev2types.ClusterPhaseDeleteFailed,
		}),
	)
	if _, err := stateConf.WaitForState(); err != nil {
		log.Printf("Delete Cluster Error:" + err.Error())
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_ccev2_cluster", action, BCESDKGoERROR)
	}
	time.Sleep(1 * time.Minute) //waiting for infrastructure delete before delete vpc & security group
	return nil
}
