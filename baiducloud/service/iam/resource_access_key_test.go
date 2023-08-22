package iam_test

import (
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/vault/helper/pgpkeys"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/acctest"
)

func TestAccAccessKey_basic(t *testing.T) {
	resourceName := "baiducloud_iam_access_key.test"
	username := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acctest.PreCheck(t) },
		Providers: acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccAccessKeyConfig_basic(username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "secret"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckNoResourceAttr(resourceName, "encrypted_secret"),
					resource.TestCheckNoResourceAttr(resourceName, "key_fingerprint"),
					acctest.CheckResourceAttrRFC3339(resourceName, "create_time"),
					resource.TestCheckResourceAttr(resourceName, "last_used_time", "0001-01-01T00:00:00Z"),
				),
			},
			{
				Config: testAccAccessKeyConfig_upadte(username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
		},
	})
}

func TestAccAccessKey_encrypted(t *testing.T) {
	resourceName := "baiducloud_iam_access_key.test"
	username := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acctest.PreCheck(t) },
		Providers: acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccAccessKeyConfig_encrypted(username, testPubKey1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					testDecryptSecretKeyAndTest(resourceName, testPrivKey1),
					resource.TestCheckNoResourceAttr(resourceName, "secret"),
					resource.TestCheckResourceAttrSet(resourceName, "encrypted_secret"),
					resource.TestCheckResourceAttrSet(resourceName, "key_fingerprint"),
					acctest.CheckResourceAttrRFC3339(resourceName, "create_time"),
					resource.TestCheckResourceAttr(resourceName, "last_used_time", "0001-01-01T00:00:00Z"),
				),
			},
		},
	})
}

func testAccUserConfig(username string) string {
	return fmt.Sprintf(`
resource "baiducloud_iam_user" "test" {
  name = "%s"
  description = "created by terraform"
  force_destroy    = true
}
`, username)
}

func testAccAccessKeyConfig_basic(username string) string {
	return acctest.ConfigCompose(testAccUserConfig(username),
		fmt.Sprintf(`
resource "baiducloud_iam_access_key" "test" {
  username = baiducloud_iam_user.test.name
  enabled  = false
}
`))
}

func testAccAccessKeyConfig_upadte(username string) string {
	return acctest.ConfigCompose(testAccUserConfig(username),
		fmt.Sprintf(`
resource "baiducloud_iam_access_key" "test" {
  username = baiducloud_iam_user.test.name
  enabled  = true
}
`))
}

func testAccAccessKeyConfig_encrypted(username, pgpKey string) string {
	return acctest.ConfigCompose(testAccUserConfig(username),
		fmt.Sprintf(`
resource "baiducloud_iam_access_key" "test" {
  username = baiducloud_iam_user.test.name
  pgp_key = <<EOF
%[2]s
EOF
}
`, username, pgpKey))
}

func testDecryptSecretKeyAndTest(resourceName, privKey string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		keyResource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		secret, ok := keyResource.Primary.Attributes["encrypted_secret"]
		if !ok {
			return fmt.Errorf("no secret in state")
		}

		// We can't verify that the decrypted secret or password is correct, because we don't
		// have it. We can verify that decrypting it does not error
		_, err := pgpkeys.DecryptBytes(secret, privKey)
		if err != nil {
			return fmt.Errorf("error decrypting secret: %w", err)
		}

		return nil
	}
}
