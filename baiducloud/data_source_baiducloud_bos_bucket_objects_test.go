package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const (
	testAccBosBucketObjectsDataSourceName          = "data.baiducloud_bos_bucket_objects.default"
	testAccBosBucketObjectsDataSourceAttrKeyPrefix = "objects.0."
)

func TestAccBaiduCloudBosBucketObjectsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBosBucketObjectsDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccBosBucketObjectsDataSourceName),
					resource.TestCheckResourceAttr(testAccBosBucketObjectsDataSourceName, testAccBosBucketObjectsDataSourceAttrKeyPrefix+"bucket", testAccBosBucketResourceAttrName),
					resource.TestCheckResourceAttr(testAccBosBucketObjectsDataSourceName, testAccBosBucketObjectsDataSourceAttrKeyPrefix+"key", testAccBosBucketObjectResourceAttrName),
					resource.TestCheckResourceAttr(testAccBosBucketObjectsDataSourceName, testAccBosBucketObjectsDataSourceAttrKeyPrefix+"acl", "public-read"),
					resource.TestCheckResourceAttr(testAccBosBucketObjectsDataSourceName, testAccBosBucketObjectsDataSourceAttrKeyPrefix+"cache_control", "no-cache"),
					resource.TestCheckResourceAttr(testAccBosBucketObjectsDataSourceName, testAccBosBucketObjectsDataSourceAttrKeyPrefix+"content_disposition", "inline"),
					resource.TestCheckResourceAttrSet(testAccBosBucketObjectsDataSourceName, testAccBosBucketObjectsDataSourceAttrKeyPrefix+"content_md5"),
					resource.TestCheckResourceAttrSet(testAccBosBucketObjectsDataSourceName, testAccBosBucketObjectsDataSourceAttrKeyPrefix+"content_type"),
					resource.TestCheckResourceAttrSet(testAccBosBucketObjectsDataSourceName, testAccBosBucketObjectsDataSourceAttrKeyPrefix+"content_length"),
					resource.TestCheckResourceAttrSet(testAccBosBucketObjectsDataSourceName, testAccBosBucketObjectsDataSourceAttrKeyPrefix+"expires"),
					resource.TestCheckResourceAttrSet(testAccBosBucketObjectsDataSourceName, testAccBosBucketObjectsDataSourceAttrKeyPrefix+"content_crc32"),
					resource.TestCheckResourceAttrSet(testAccBosBucketObjectsDataSourceName, testAccBosBucketObjectsDataSourceAttrKeyPrefix+"last_modified"),
					resource.TestCheckResourceAttrSet(testAccBosBucketObjectsDataSourceName, testAccBosBucketObjectsDataSourceAttrKeyPrefix+"etag"),
					resource.TestCheckResourceAttrSet(testAccBosBucketObjectsDataSourceName, testAccBosBucketObjectsDataSourceAttrKeyPrefix+"size"),
					resource.TestCheckResourceAttrSet(testAccBosBucketObjectsDataSourceName, testAccBosBucketObjectsDataSourceAttrKeyPrefix+"storage_class"),
					resource.TestCheckResourceAttrSet(testAccBosBucketObjectsDataSourceName, testAccBosBucketObjectsDataSourceAttrKeyPrefix+"owner_id"),
					resource.TestCheckResourceAttr(testAccBosBucketObjectsDataSourceName, testAccBosBucketObjectsDataSourceAttrKeyPrefix+"user_meta.Metaa", "metaA"),
					resource.TestCheckResourceAttr(testAccBosBucketObjectsDataSourceName, testAccBosBucketObjectsDataSourceAttrKeyPrefix+"user_meta.Metab", "metaB"),
				),
			},
		},
	})
}

func testAccBosBucketObjectsDataSourceConfig() string {
	return fmt.Sprintf(`
resource "baiducloud_bos_bucket" "default" {
  bucket = "%s"
}

resource "baiducloud_bos_bucket_object" "default" {
  bucket = "${baiducloud_bos_bucket.default.bucket}"
  key = "%s"
  content = "hello world"
  acl = "public-read"
  cache_control = "no-cache"
  content_disposition = "inline"
  storage_class = "COLD"
  user_meta = {
    Metaa = "metaA"
    Metab = "metaB"
  }
}

data "baiducloud_bos_bucket_objects" "default" {
  bucket = "%s"
  prefix = "${baiducloud_bos_bucket_object.default.key}"
}
`, testAccBosBucketResourceAttrName, testAccBosBucketObjectResourceAttrName,
		testAccBosBucketResourceAttrName)
}
