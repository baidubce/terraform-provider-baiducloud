package bcc_test

import (
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/acctest"
)

func TestAccKeyPairs(t *testing.T) {
	resourceName := "baiducloud_bcc_key_pair.test"
	dataSourceName := "data.baiducloud_bcc_key_pairs.test"

	name := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	description := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acctest.PreCheck(t) },
		Providers: acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyPairsConfig(name, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "key_pairs.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "id", dataSourceName, "key_pairs.0.keypair_id"),
					resource.TestCheckResourceAttrPair(resourceName, "name", dataSourceName, "key_pairs.0.name"),
					resource.TestCheckResourceAttrPair(resourceName, "description", dataSourceName, "key_pairs.0.description"),
					resource.TestCheckResourceAttrPair(resourceName, "created_time", dataSourceName, "key_pairs.0.created_time"),
					resource.TestCheckResourceAttrPair(resourceName, "public_key", dataSourceName, "key_pairs.0.public_key"),
					resource.TestCheckResourceAttrPair(resourceName, "instance_count", dataSourceName, "key_pairs.0.instance_count"),
					resource.TestCheckResourceAttrPair(resourceName, "region_id", dataSourceName, "key_pairs.0.region_id"),
					resource.TestCheckResourceAttrPair(resourceName, "fingerprint", dataSourceName, "key_pairs.0.fingerprint"),
				),
			},
		},
	})
}

func testAccKeyPairsConfig(name, description string) string {
	return acctest.ConfigCompose(
		testAccKeyPairConfig_create(name, description),
		fmt.Sprintf(`
data "baiducloud_bcc_key_pairs" "test" {
	name = baiducloud_bcc_key_pair.test.name
}`))
}
