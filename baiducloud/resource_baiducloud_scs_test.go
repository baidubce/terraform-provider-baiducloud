package baiducloud

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/baidubce/bce-sdk-go/services/scs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccScsResourceType = "baiducloud_scs"
	testAccScsResourceName = testAccScsResourceType + "." + BaiduCloudTestResourceName
)

func init() {
	resource.AddTestSweepers(testAccScsResourceType, &resource.Sweeper{
		Name: testAccScsResourceType,
		F:    testSweepScsInstances,
	})
}

func testSweepScsInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("get BaiduCloud client error: %s", err)
	}

	client := rawClient.(*connectivity.BaiduClient)
	scsClient := &ScsService{client}

	args := &scs.ListInstancesArgs{}
	instList, err := scsClient.ListAllInstances(args)
	if err != nil {
		return fmt.Errorf("get SCS instances error: %v", err)
	}

	for _, inst := range instList {
		if !strings.HasPrefix(inst.InstanceName, BaiduCloudTestResourceAttrNamePrefix) || inst.InstanceStatus != "Running" {
			log.Printf("[INFO] Skipping SCS instance: %s (%s)", inst.InstanceID, inst.InstanceName)
			continue
		}

		log.Printf("[INFO] Deleting SCS instance: %s (%s)", inst.InstanceID, inst.InstanceName)
		_, err := client.WithScsClient(func(scsClient *scs.Client) (interface{}, error) {
			return nil, scsClient.DeleteInstance(inst.InstanceID, buildClientToken())
		})

		if err != nil {
			log.Printf("[ERROR] Failed to delete SCS instance %s (%s)", inst.InstanceID, inst.InstanceName)
		}
	}

	return nil
}

func TestAccBaiduCloudScs(t *testing.T) {
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	instance_name := BaiduCloudTestResourceAttrNamePrefix + "Scs-" + timeStamp
	instance_name_new := BaiduCloudTestResourceAttrNamePrefix + "ScsNew-" + timeStamp
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccScsDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccScsConfig(instance_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccScsResourceName),
					resource.TestCheckResourceAttr(testAccScsResourceName, "instance_name", instance_name),
					resource.TestCheckResourceAttr(testAccScsResourceName, "billing.payment_timing", "Postpaid"),
					resource.TestCheckResourceAttr(testAccScsResourceName, "port", "6379"),
					resource.TestCheckResourceAttr(testAccScsResourceName, "cluster_type", "master_slave"),
					resource.TestCheckResourceAttr(testAccScsResourceName, "engine_version", "3.2"),
					resource.TestCheckResourceAttr(testAccScsResourceName, "replication_num", "1"),
					resource.TestCheckResourceAttr(testAccScsResourceName, "shard_num", "1"),
					resource.TestCheckResourceAttr(testAccScsResourceName, "node_type", "cache.n1.micro"),
				),
			},
			{
				ResourceName:            testAccScsResourceName,
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"instance_status"},
			},
			{
				Config: testAccScsConfigUpdate(instance_name_new),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccScsResourceName),
					resource.TestCheckResourceAttr(testAccScsResourceName, "instance_name", instance_name_new),
					resource.TestCheckResourceAttr(testAccScsResourceName, "billing.payment_timing", "Postpaid"),
					resource.TestCheckResourceAttr(testAccScsResourceName, "cluster_type", "master_slave"),
					resource.TestCheckResourceAttr(testAccScsResourceName, "engine_version", "3.2"),
					resource.TestCheckResourceAttr(testAccScsResourceName, "replication_num", "1"),
					resource.TestCheckResourceAttr(testAccScsResourceName, "shard_num", "1"),
					resource.TestCheckResourceAttr(testAccScsResourceName, "node_type", "cache.n1.micro"),
				),
			},
		},
	})
}

func testAccScsDestory(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	scsService := ScsService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccScsResourceType {
			continue
		}

		result, err := scsService.GetInstanceDetail(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		deletingStatus := []string{"Pausing", "Paused", "Deleted", "Deleting"}
		if IsExpectValue(result.InstanceStatus, deletingStatus) {
			continue
		}

		return WrapError(Error("SCS still exist"))
	}

	return nil
}

func IsExpectValue(value string, expectList []string) bool {
	if len(value) == 0 || value == "" {
		return true
	}

	for _, expect := range expectList {
		if strings.Contains(expect, value) {
			return true
		}
	}
	return false
}

func testAccScsConfig(name string) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
    instance_name           = "%s"
	billing = {
    	payment_timing 		= "Postpaid"
  	}
    purchase_count 			= 1
  	port 					= 6379
	engine_version 			= "3.2"
	node_type 				= "cache.n1.micro"
	cluster_type 			= "master_slave"
	replication_num 		= 1
	shard_num 				= 1
	proxy_num 				= 0
}
`, testAccScsResourceType, BaiduCloudTestResourceName, name)
}

func testAccScsConfigUpdate(name string) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
    instance_name           = "%s"
	billing = {
    	payment_timing 		= "Postpaid"
  	}
    purchase_count 			= 1
  	port 					= 6379
	engine_version 			= "3.2"
	node_type 				= "cache.n1.micro"
	cluster_type 			= "master_slave"
	replication_num 		= 1
	shard_num 				= 1
	proxy_num 				= 0
}
`, testAccScsResourceType, BaiduCloudTestResourceName, name)
}
