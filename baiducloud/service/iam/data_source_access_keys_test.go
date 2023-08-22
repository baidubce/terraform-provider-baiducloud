package iam_test

import (
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/acctest"
)

func TestAccAccessKeys(t *testing.T) {
	resourceName := "baiducloud_iam_access_key.test"
	dataSourceName := "data.baiducloud_iam_access_keys.test"
	username := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acctest.PreCheck(t) },
		Providers: acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccAccessKeysConfig(username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "access_keys.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "id", dataSourceName, "access_keys.0.access_key_id"),
					resource.TestCheckResourceAttrPair(resourceName, "enabled", dataSourceName, "access_keys.0.enabled"),
					resource.TestCheckResourceAttrPair(resourceName, "create_time", dataSourceName, "access_keys.0.create_time"),
					resource.TestCheckResourceAttrPair(resourceName, "last_used_time", dataSourceName, "access_keys.0.last_used_time"),
				),
			},
		},
	})
}

func testAccAccessKeysConfig(username string) string {
	return acctest.ConfigCompose(
		testAccAccessKeyConfig_basic(username),
		fmt.Sprintf(`
data "baiducloud_iam_access_keys" "test" {
	username = baiducloud_iam_access_key.test.username
}
`))
}
