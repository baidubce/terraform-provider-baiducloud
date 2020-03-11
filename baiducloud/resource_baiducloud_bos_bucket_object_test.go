package baiducloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/baidubce/bce-sdk-go/services/bos"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccBosBucketObjectResourceType     = "baiducloud_bos_bucket_object"
	testAccBosBucketObjectResourceName     = testAccBosBucketObjectResourceType + "." + BaiduCloudTestResourceName
	testAccBosBucketObjectResourceAttrName = BaiduCloudTestResourceAttrNamePrefix + "BosBucketObject"
)

func init() {
	resource.AddTestSweepers(testAccBosBucketObjectResourceType, &resource.Sweeper{
		Name:         testAccBosBucketObjectResourceType,
		F:            testSweepBosBucketObjects,
		Dependencies: []string{testAccBosBucketResourceType},
	})
}

func testSweepBosBucketObjects(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("get BaiduCloud client error: %s", err)
	}

	client := rawClient.(*connectivity.BaiduClient)

	exist, err := client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
		return bosClient.DoesBucketExist(testAccBosBucketResourceAttrName)
	})
	if err != nil {
		log.Printf("[ERROR] Failed to check if the bucket %s exist %v.", testAccBosBucketResourceAttrName, err)
		return fmt.Errorf("check bucket %s exist error: %v", testAccBosBucketResourceAttrName, err)
	}
	if !exist.(bool) {
		return nil
	}

	bosService := &BosService{client}
	objectList, err := bosService.ListAllObjects(testAccBosBucketResourceAttrName, "")
	if err != nil {
		log.Printf("[ERROR] Failed to list object %v", err)
		return fmt.Errorf("get %s object list error: %s", testAccBosBucketResourceAttrName, err)
	}

	for _, obj := range objectList {
		if !strings.HasPrefix(obj.Key, testAccBosBucketObjectResourceAttrName) {
			log.Printf("[INFO] Skipping Object: %s", obj.Key)
			continue
		}

		log.Printf("[INFO] Deleting Object: %s", obj.Key)
		_, err := client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
			return nil, bosClient.DeleteObject(testAccBosBucketResourceAttrName, obj.Key)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete object %s", obj.Key)
		}
	}

	return nil
}

//lintignore:AT003
func TestAccBaiduCloudBosBucketObject(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccBosBucketObjectDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccBosBucketObjectConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccBosBucketObjectResourceName),
					resource.TestCheckResourceAttr(testAccBosBucketObjectResourceName, "bucket", testAccBosBucketResourceAttrName),
					resource.TestCheckResourceAttr(testAccBosBucketObjectResourceName, "key", testAccBosBucketObjectResourceAttrName),
					resource.TestCheckResourceAttr(testAccBosBucketObjectResourceName, "content", "hello world"),
					resource.TestCheckResourceAttr(testAccBosBucketObjectResourceName, "acl", "public-read"),
					resource.TestCheckResourceAttr(testAccBosBucketObjectResourceName, "cache_control", "no-cache"),
					resource.TestCheckResourceAttr(testAccBosBucketObjectResourceName, "content_disposition", "inline"),
					resource.TestCheckResourceAttrSet(testAccBosBucketObjectResourceName, "content_md5"),
					resource.TestCheckResourceAttrSet(testAccBosBucketObjectResourceName, "content_type"),
					resource.TestCheckResourceAttrSet(testAccBosBucketObjectResourceName, "content_length"),
					resource.TestCheckResourceAttrSet(testAccBosBucketObjectResourceName, "expires"),
					resource.TestCheckResourceAttrSet(testAccBosBucketObjectResourceName, "content_crc32"),
					resource.TestCheckResourceAttr(testAccBosBucketObjectResourceName, "storage_class", "COLD"),
					resource.TestCheckResourceAttr(testAccBosBucketObjectResourceName, "user_meta.Metaa", "metaA"),
					resource.TestCheckResourceAttr(testAccBosBucketObjectResourceName, "user_meta.Metab", "metaB"),
					resource.TestCheckResourceAttrSet(testAccBosBucketObjectResourceName, "etag"),
					resource.TestCheckResourceAttrSet(testAccBosBucketObjectResourceName, "last_modified"),
				),
			},
			{
				Config: testAccBosBucketObjectConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccBosBucketObjectResourceName),
					resource.TestCheckResourceAttr(testAccBosBucketObjectResourceName, "bucket", testAccBosBucketResourceAttrName),
					resource.TestCheckResourceAttr(testAccBosBucketObjectResourceName, "key", testAccBosBucketObjectResourceAttrName),
					resource.TestCheckResourceAttr(testAccBosBucketObjectResourceName, "content", "hello world"),
					resource.TestCheckResourceAttr(testAccBosBucketObjectResourceName, "acl", "private"),
					resource.TestCheckResourceAttr(testAccBosBucketObjectResourceName, "cache_control", "max-age"),
					resource.TestCheckResourceAttr(testAccBosBucketObjectResourceName, "content_disposition", "attachment"),
					resource.TestCheckResourceAttrSet(testAccBosBucketObjectResourceName, "content_md5"),
					resource.TestCheckResourceAttrSet(testAccBosBucketObjectResourceName, "content_type"),
					resource.TestCheckResourceAttrSet(testAccBosBucketObjectResourceName, "content_length"),
					resource.TestCheckResourceAttrSet(testAccBosBucketObjectResourceName, "expires"),
					resource.TestCheckResourceAttrSet(testAccBosBucketObjectResourceName, "content_crc32"),
					resource.TestCheckResourceAttr(testAccBosBucketObjectResourceName, "storage_class", "STANDARD_IA"),
					resource.TestCheckResourceAttr(testAccBosBucketObjectResourceName, "user_meta.Metaa", "metaA"),
					resource.TestCheckResourceAttr(testAccBosBucketObjectResourceName, "user_meta.Metab", "metaB"),
					resource.TestCheckResourceAttrSet(testAccBosBucketObjectResourceName, "etag"),
					resource.TestCheckResourceAttrSet(testAccBosBucketObjectResourceName, "last_modified"),
				),
			},
		},
	})
}

func testAccBosBucketObjectDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccBosBucketObjectResourceType {
			continue
		}

		_, err := client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
			return bosClient.GetObjectMeta(testAccBosBucketResourceAttrName, rs.Primary.ID)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{"Not Found"}) {
				continue
			}
			return WrapError(err)
		}
		return WrapError(Error("BOS bucket object still exist"))
	}

	return nil
}

func testAccBosBucketObjectConfig() string {
	return fmt.Sprintf(`
resource "baiducloud_bos_bucket" "default" {
  bucket = "%s"
}

resource "baiducloud_bos_bucket_object" "default" {
  bucket              = baiducloud_bos_bucket.default.bucket
  key                 = "%s"
  content             = "hello world"
  acl                 = "public-read"
  cache_control       = "no-cache"
  content_disposition = "inline"
  storage_class       = "COLD"
  user_meta = {
    Metaa = "metaA"
    Metab = "metaB"
  }
}
`, testAccBosBucketResourceAttrName, testAccBosBucketObjectResourceAttrName)
}

func testAccBosBucketObjectConfigUpdate() string {
	return fmt.Sprintf(`
resource "baiducloud_bos_bucket" "default" {
  bucket = "%s"
}

resource "baiducloud_bos_bucket_object" "default" {
  bucket              = baiducloud_bos_bucket.default.bucket
  key                 = "%s"
  content             = "hello world"
  acl                 = "private"
  cache_control       = "max-age"
  content_disposition = "attachment"
  storage_class       = "STANDARD_IA"
  user_meta = {
    metaa = "metaA"
    metab = "metaB"
  }
}
`, testAccBosBucketResourceAttrName, testAccBosBucketObjectResourceAttrName)
}
