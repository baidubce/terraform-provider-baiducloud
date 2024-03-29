package acctest

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud"
)

const ResourcePrefix = "tf-acc-test"
const DefaultEmailAddress = "test@test.com"
const RFC3339RegexPattern = `^[0-9]{4}-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])[Tt]([01][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9](\.[0-9]+)?([Zz]|([+-]([01][0-9]|2[0-3]):[0-5][0-9]))$`

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

func CheckResourceAttrRFC3339(resourceName, attributeName string) resource.TestCheckFunc {
	return resource.TestMatchResourceAttr(resourceName, attributeName, regexp.MustCompile(RFC3339RegexPattern))
}

func ConfigCompose(config ...string) string {
	var str strings.Builder

	for _, conf := range config {
		str.WriteString(conf)
	}

	return str.String()
}

func ConfigVPCWithSubnet() string {
	return fmt.Sprintf(`
resource "baiducloud_vpc" "test" {
  name        = "vpc_terraform_test"
  description = "created by terraform for test"
  cidr        = "192.168.0.0/24"
}

resource "baiducloud_subnet" "test" {
  name        = "subnet_terraform_test"
  zone_name   = "cn-bj-a"
  cidr        = "192.168.0.0/24"
  vpc_id      = baiducloud_vpc.test.id
  description = "created by terraform for test"
}`)

}
