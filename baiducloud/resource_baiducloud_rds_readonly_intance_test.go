package baiducloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/baidubce/bce-sdk-go/services/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccRdsReadOnlyInstanceResourceType = "baiducloud_rds_readonly_instance"
	testAccRdsReadOnlyInstanceResourceName = testAccRdsReadOnlyInstanceResourceType + "." + BaiduCloudTestResourceName
)

func init() {
	resource.AddTestSweepers(testAccRdsReadOnlyInstanceResourceType, &resource.Sweeper{
		Name: testAccRdsReadOnlyInstanceResourceType,
		F:    testSweepRdsReadOnlyInstances,
	})
}

func testSweepRdsReadOnlyInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("get BaiduCloud client error: %s", err)
	}

	client := rawClient.(*connectivity.BaiduClient)
	rdsClient := &RdsService{client}

	args := &rds.ListRdsArgs{}
	instList, err := rdsClient.ListAllInstances(args)
	if err != nil {
		return fmt.Errorf("get RDS instances error: %v", err)
	}

	for _, inst := range instList {
		if !strings.HasPrefix(inst.InstanceName, BaiduCloudTestResourceAttrNamePrefix) || inst.InstanceStatus != "Running" {
			log.Printf("[INFO] Skipping RDS instance: %s (%s)", inst.InstanceId, inst.InstanceName)
			continue
		}

		log.Printf("[INFO] Deleting RDS instance: %s (%s)", inst.InstanceId, inst.InstanceName)
		_, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return nil, rdsClient.DeleteRds(inst.InstanceId)
		})

		if err != nil {
			log.Printf("[ERROR] Failed to delete RDS instance %s (%s)", inst.InstanceId, inst.InstanceName)
		}
	}

	return nil
}

func TestAccBaiduCloudRdsReadOnlyInstance(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccRdsReadOnlyInstanceDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccRdsReadOnlyInstanceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccRdsReadOnlyInstanceResourceName),
					resource.TestCheckResourceAttr(testAccRdsReadOnlyInstanceResourceName, "billing.payment_timing", "Postpaid"),
				),
			},
			{
				ResourceName:            testAccRdsReadOnlyInstanceResourceName,
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"instance_status"},
			},
		},
	})
}

func testAccRdsReadOnlyInstanceDestory(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)
	rdsService := RdsService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRdsReadOnlyInstanceResourceType {
			continue
		}

		result, err := rdsService.GetInstanceDetail(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		deletingStatus := []string{"Deleted", "Deleting"}
		if IsExpectValue(result.InstanceStatus, deletingStatus) {
			continue
		}

		return WrapError(Error("RDS Instance still exist"))
	}

	return nil
}

func testAccRdsReadOnlyInstanceConfig() string {
	return fmt.Sprintf(`
resource "baiducloud_rds_instance" "default" {
    instance_name              = "%s"
    billing = {
        payment_timing         = "Postpaid"
    }
    engine_version             = "5.6"
    engine                     = "MySQL"
    cpu_count                  = 1
    memory_capacity            = 1
    volume_capacity            = 5
}

resource "%s" "%s" {
    instance_name              = "%s"
    billing = {
        payment_timing         = "Postpaid"
    }
    source_instance_id         = baiducloud_rds_instance.default.instance_id
    cpu_count                  = 1
    memory_capacity            = 1
    volume_capacity            = 5
}
`, BaiduCloudTestResourceAttrNamePrefix+"Rds_Master", testAccRdsReadOnlyInstanceResourceType, BaiduCloudTestResourceName, BaiduCloudTestResourceAttrNamePrefix+"Rds_ReadOnly")
}
