package acctest

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud"
	"log"
	"os"
	"testing"
)

var Providers map[string]terraform.ResourceProvider
var Provider *schema.Provider

func init() {
	Provider = baiducloud.Provider().(*schema.Provider)
	Providers = map[string]terraform.ResourceProvider{
		"baiducloud": Provider,
	}
}

func PreCheck(t *testing.T) {
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

func CheckResource(rName string, state *terraform.State) (*terraform.ResourceState, error) {
	rs, ok := state.RootModule().Resources[rName]
	if !ok {
		return nil, fmt.Errorf("Not found: %s", rName)
	}
	if rs.Primary.ID == "" {
		return nil, fmt.Errorf("No Domain ID is set")
	}
	return rs, nil
}
