package baiducloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/baidubce/bce-sdk-go/services/cce"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccCceResourceType = "baiducloud_cce_cluster"
	testAccCceResourceName = testAccCceResourceType + "." + BaiduCloudTestResourceName
)

func init() {
	resource.AddTestSweepers(testAccCceResourceType, &resource.Sweeper{
		Name: testAccCceResourceType,
		F:    testSweepCce,
	})
}

func testSweepCce(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("get BaiduCloud client error: %s", err)
	}
	client := rawClient.(*connectivity.BaiduClient)

	listArgs := &cce.ListClusterArgs{}
	raw, err := client.WithCCEClient(func(cceClient *cce.Client) (i interface{}, e error) {
		return cceClient.ListClusters(listArgs)
	})
	if err != nil {
		return fmt.Errorf("list CCE Cluster with error: %s", err)
	}

	cceList := raw.(*cce.ListClusterResult)
	for _, c := range cceList.Clusters {
		if !strings.HasPrefix(c.ClusterName, BaiduCloudTestResourceAttrNamePrefix) {
			log.Printf("[INFO] Skipping CCE Cluster: %s (%s)", c.ClusterName, c.ClusterUuid)
			continue
		}

		if c.Status == cce.ClusterStatusDeleting || c.Status == cce.ClusterStatusCreating {
			log.Printf("[INFO] Skipping CCE Cluster: %s (%s) with status %s",
				c.ClusterName, c.ClusterUuid, c.Status)
			continue
		}

		log.Printf("[INFO] Deleting CCE Cluster: %s (%s)", c.ClusterName, c.ClusterUuid)
		deleteArgs := &cce.DeleteClusterArgs{
			DeleteEipCds: true,
			DeleteSnap:   true,
			ClusterUuid:  c.ClusterUuid,
		}
		_, err := client.WithCCEClient(func(cceClient *cce.Client) (i interface{}, e error) {
			return nil, cceClient.DeleteCluster(deleteArgs)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete CCE cluster %s (%s) with error: %v",
				c.ClusterName, c.ClusterUuid, err)
		}
	}

	return nil
}

//lintignore:AT003
func TestAccBaiduCloudCce(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCceDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccCceMasterConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccCceResourceName),
				),
			},
			{
				Config: testAccCceMasterUpdateConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccCceResourceName),
				),
			},
		},
	})
}

func testAccCceDestory(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccCceResourceType {
			continue
		}

		raw, err := client.WithCCEClient(func(client *cce.Client) (i interface{}, e error) {
			return client.GetCluster(rs.Primary.ID)
		})

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		cluster := raw.(*cce.GetClusterResult)
		if cluster.Status == cce.ClusterStatusDeleted {
			continue
		}

		return WrapError(Error("CCE Cluster still exist"))
	}

	return nil
}

func testAccCceMasterConfig() string {
	return fmt.Sprintf(`
data "baiducloud_zones" "defaultA" {
  name_regex = ".*a$"
}

data "baiducloud_zones" "defaultB" {
  name_regex = ".*b$"
}

data "baiducloud_images" "default" {
  image_type = "System"
  name_regex = "7.5.*"
  os_name    = "CentOS"
}

resource "baiducloud_vpc" "default" {
  name        = "%s"
  description = var.description
  cidr        = "192.168.0.0/16"
}

resource "baiducloud_subnet" "defaultA" {
  name        = "%s"
  zone_name   = data.baiducloud_zones.defaultA.zones.0.zone_name
  cidr        = "192.168.1.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "terraform create"
}

resource "baiducloud_subnet" "defaultB" {
  name        = "%s"
  zone_name   = data.baiducloud_zones.defaultB.zones.0.zone_name
  cidr        = "192.168.2.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "terraform create"
}

resource "baiducloud_security_group" "defualt" {
  name   = "%s"
  vpc_id = baiducloud_vpc.default.id
}

resource "baiducloud_security_group_rule" "default" {
  security_group_id = baiducloud_security_group.defualt.id
  remark            = "remark"
  protocol          = "all"
  port_range        = "1-65535"
  direction         = "ingress"
}

resource "baiducloud_security_group_rule" "default2" {
  security_group_id = baiducloud_security_group.defualt.id
  remark            = "remark"
  protocol          = "all"
  port_range        = "1-65535"
  direction         = "egress"
}

data "baiducloud_cce_versions" "default" {
  version_regex = ".*13.*"
}

resource "baiducloud_cce_cluster" "default_independent" {
  cluster_name        = "%s"
  main_available_zone = "zoneA"
  version             = data.baiducloud_cce_versions.default.versions.0
  container_net       = "172.16.0.0/16"

  advanced_options = {
    kube_proxy_mode = "ipvs"
    dns_mode        = "CoreDNS"
    cni_mode        = "cni"
    cni_type        = "VPC_SECONDARY_IP_VETH"
    max_pod_num     = "256"
  }

  delete_eip_cds   = "true"
  delete_snapshots = "true"

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

    root_disk_size_in_gb   = 100
    root_disk_storage_type = "ssd"
    admin_pass             = "baiduPasswd@123"
    image_type             = "common"

    cds_disks {
      volume_type     = "sata"
      disk_size_in_gb = 10
    }

    eip = {
      bandwidth_in_mbps = 100
      sub_product_type  = "netraffic"
    }
  }

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
}`, BaiduCloudTestResourceAttrNamePrefix+"VPC", BaiduCloudTestResourceAttrNamePrefix+"SubnetA",
		BaiduCloudTestResourceAttrNamePrefix+"SubnetB", BaiduCloudTestResourceAttrNamePrefix+"SG",
		BaiduCloudTestResourceAttrNamePrefix+"CCE")
}

func testAccCceMasterUpdateConfig() string {
	return fmt.Sprintf(`
data "baiducloud_zones" "defaultA" {
  name_regex = ".*a$"
}

data "baiducloud_zones" "defaultB" {
  name_regex = ".*b$"
}

data "baiducloud_images" "default" {
  image_type = "System"
  name_regex = "7.5.*"
  os_name    = "CentOS"
}

resource "baiducloud_vpc" "default" {
  name        = "%s"
  description = "terraform create"
  cidr        = "192.168.0.0/16"
}

resource "baiducloud_subnet" "defaultA" {
  name        = "%s"
  zone_name   = data.baiducloud_zones.defaultA.zones.0.zone_name
  cidr        = "192.168.1.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "terraform create"
}

resource "baiducloud_subnet" "defaultB" {
  name        = "%s"
  zone_name   = data.baiducloud_zones.defaultB.zones.0.zone_name
  cidr        = "192.168.2.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "terraform create"
}

resource "baiducloud_security_group" "defualt" {
  name   = "%s"
  vpc_id = baiducloud_vpc.default.id
}

resource "baiducloud_security_group_rule" "default" {
  security_group_id = baiducloud_security_group.defualt.id
  remark            = "remark"
  protocol          = "all"
  port_range        = "1-65535"
  direction         = "ingress"
}

resource "baiducloud_security_group_rule" "default2" {
  security_group_id = baiducloud_security_group.defualt.id
  remark            = "remark"
  protocol          = "all"
  port_range        = "1-65535"
  direction         = "egress"
}

data "baiducloud_cce_versions" "default" {
  version_regex = ".*13.*"
}

resource "baiducloud_cce_cluster" "default_independent" {
 cluster_name        = "%s"
  main_available_zone = "zoneA"
  version             = data.baiducloud_cce_versions.default.versions.0
  container_net       = "172.16.0.0/16"

  advanced_options = {
    kube_proxy_mode = "ipvs"
    dns_mode        = "CoreDNS"
    cni_mode        = "cni"
    cni_type        = "VPC_SECONDARY_IP_VETH"
    max_pod_num     = "256"
  }

  delete_eip_cds   = "true"
  delete_snapshots = "true"


  worker_config {
    count = {
      "zoneA" : 1
      "zoneB" : 1
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

    root_disk_size_in_gb   = 100
    root_disk_storage_type = "ssd"
    admin_pass             = "baiduPasswd@123"
    image_type             = "common"

    cds_disks {
      volume_type     = "sata"
      disk_size_in_gb = 10
    }

    eip = {
      bandwidth_in_mbps = 100
      sub_product_type  = "netraffic"
    }
  }

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
}
`, BaiduCloudTestResourceAttrNamePrefix+"VPC", BaiduCloudTestResourceAttrNamePrefix+"SubnetA",
		BaiduCloudTestResourceAttrNamePrefix+"SubnetB", BaiduCloudTestResourceAttrNamePrefix+"SG",
		BaiduCloudTestResourceAttrNamePrefix+"CCE")
}
