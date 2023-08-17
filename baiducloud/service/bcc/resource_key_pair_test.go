package bcc_test

import (
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/acctest"
)

func TestAccKeyPair_basic(t *testing.T) {
	resourceName := "baiducloud_bcc_key_pair.test"
	name := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	description := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	updatedName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	updatedDescription := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acctest.PreCheck(t) },
		Providers: acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyPairConfig_create(name, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttrSet(resourceName, "public_key"),
					resource.TestCheckResourceAttrSet(resourceName, "created_time"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_count"),
					resource.TestCheckResourceAttrSet(resourceName, "region_id"),
					resource.TestCheckResourceAttrSet(resourceName, "fingerprint"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"private_key_file"},
			},
			{
				Config: testAccKeyPairConfig_update(updatedName, updatedDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
				),
			},
		},
	})
}

func TestAccKeyPair_importPublicKey(t *testing.T) {
	resourceName := "baiducloud_bcc_key_pair.test"
	name := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	publicKey, _, err := sdkacctest.RandSSHKeyPair(acctest.DefaultEmailAddress)
	if err != nil {
		t.Fatalf("error generating random SSH key: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acctest.PreCheck(t) },
		Providers: acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyPairConfig_importPublic_key(name, publicKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "public_key", publicKey),
					resource.TestCheckResourceAttrSet(resourceName, "created_time"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_count"),
					resource.TestCheckResourceAttrSet(resourceName, "region_id"),
					resource.TestCheckResourceAttrSet(resourceName, "fingerprint"),
				),
			},
		},
	})
}

func testAccKeyPairConfig_create(name, description string) string {
	return fmt.Sprintf(`
resource baiducloud_bcc_key_pair test {
    name = %[1]q
    description = %[2]q
    private_key_file = "private-key.txt"
}
`, name, description)
}

func testAccKeyPairConfig_update(name, description string) string {
	return fmt.Sprintf(`
resource baiducloud_bcc_key_pair test {
    name = %[1]q
    description = %[2]q
    private_key_file = "private-key.txt"
}
`, name, description)
}

func testAccKeyPairConfig_importPublic_key(name, publicKey string) string {
	return fmt.Sprintf(`
resource baiducloud_bcc_key_pair test {
    name = %[1]q
    public_key = %[2]q
}
`, name, publicKey)
}
