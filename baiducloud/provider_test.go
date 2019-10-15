package baiducloud

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"baiducloud": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("BAIDUCLOUD_ACCESS_KEY"); v == "" {
		t.Fatal("BAIDUCLOUD_ACCESS_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("BAIDUCLOUD_SECRET_KEY"); v == "" {
		t.Fatal("BAIDUCLOUD_SECRET_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("BAIDUCLOUD_REGION"); v == "" {
		log.Println("[INFO] Test: Using cn-beijing as test region")
		os.Setenv("BAIDUCLOUD_REGION", "bj")
	}
}

func testAccCheckBaiduCloudDataSourceId(n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("data source ID not set")
		}
		return nil
	}
}
