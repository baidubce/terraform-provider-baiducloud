/*
Provide a resource to create a BOS bucket object.

Example Usage

```hcl
resource "baiducloud_bos_bucket_object" "default" {
  bucket = "my-bucket"
  key = "test-key"
  source = "/tmp/test-file"
  acl = "public-read"
}
```
*/
package baiducloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/bos"
	"github.com/baidubce/bce-sdk-go/services/bos/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudBucketObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudBucketObjectPut,
		Read:   resourceBaiduCloudBucketObjectRead,
		Update: resourceBaiduCloudBucketObjectPut,
		Delete: resourceBaiduCloudBucketObjectDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:        schema.TypeString,
				Description: "Name of the bucket to put the file in.",
				Required:    true,
				ForceNew:    true,
			},
			"key": {
				Type:        schema.TypeString,
				Description: "Name of the object once it is in the bucket.",
				Required:    true,
				ForceNew:    true,
			},
			"source": {
				Type:          schema.TypeString,
				Description:   "The file path that will be read and uploaded as raw bytes for the object content.",
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"content"},
			},
			"content": {
				Type:          schema.TypeString,
				Description:   "The literal string value that will be uploaded as the object content.",
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"source"},
			},

			"acl": {
				Type:         schema.TypeString,
				Description:  "Canned ACL of the object, which can be private or public-read. If the value is not set, the object permission will be empty by default, and then the bucket permission as default.",
				Optional:     true,
				ValidateFunc: validateBOSObjectACL(),
			},

			"cache_control": {
				Type:         schema.TypeString,
				Description:  "The caching behavior along the request/reply chain. Valid values are private, no-cache, max-age and must-revalidate. If not set, the value is empty.",
				Optional:     true,
				ValidateFunc: validateBOSObjectCacheControl(),
			},
			"content_disposition": {
				Type:         schema.TypeString,
				Description:  "Specifies presentational information for the object, which can be inline or attachment. If not set, the value is empty.",
				Optional:     true,
				ValidateFunc: validateBOSObjectContentDisposition(),
			},
			"content_md5": {
				Type:        schema.TypeString,
				Description: "MD5 digest of the HTTP request content defined in RFC2616 can be carried by the field to verify whether the file saved on the BOS side is consistent with the file expected by the user.",
				Optional:    true,
				Computed:    true,
			},
			"content_type": {
				Type:        schema.TypeString,
				Description: "Type to describe the format of the object data.",
				Optional:    true,
				Computed:    true,
			},
			"content_length": {
				Type:        schema.TypeInt,
				Description: "Length of the content to be uploaded.",
				Optional:    true,
				Computed:    true,
			},
			"expires": {
				Type:        schema.TypeString,
				Description: "The expire date is used to set the cache expiration time when downloading object. If it is not set, the BOS will set the cache expiration time to three days by default.",
				Optional:    true,
				Computed:    true,
			},
			"user_meta": {
				Type:        schema.TypeMap,
				Description: "The mapping of key/values to to provision metadata, which will be automatically prefixed by x-bce-meta-.",
				Optional:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					o, n := d.GetChange("user_meta")
					oMap := o.(map[string]interface{})
					nMap := n.(map[string]interface{})

					if len(oMap) != len(nMap) {
						return false
					}

					oldMap := make(map[string]interface{})
					newMap := make(map[string]interface{})

					for k, v := range oMap {
						oldMap[strings.ToLower(k)] = v
					}
					for k, v := range nMap {
						newMap[strings.ToLower(k)] = v
					}

					for k, v := range oldMap {
						lowerK := strings.ToLower(k)
						vv, ok := newMap[lowerK]
						if !ok || (vv.(string) != v.(string)) {
							return false
						}
						delete(newMap, lowerK)
					}
					for k, v := range newMap {
						lowerK := strings.ToLower(k)
						vv, ok := oldMap[lowerK]
						if !ok || (vv.(string) != v.(string)) {
							return false
						}
					}

					return true
				},
			},

			"content_sha256": {
				Type:        schema.TypeString,
				Description: "Sha256 value of the object, which is used to verify whether the file saved on the BOS side is consistent with the file expected by the user, the sha256 has higher verification accuracy, and the sha256 value of the transmitted data must match this, otherwise the object uploaded fails.",
				Optional:    true,
			},
			"content_crc32": {
				Type:        schema.TypeString,
				Description: "Crc(cyclic redundancy check code) value of the object.",
				Optional:    true,
				Computed:    true,
			},
			"storage_class": {
				Type:         schema.TypeString,
				Description:  "Storage class of the object, which can be COLD, STANDARD_IA, STANDARD or ARCHIVE. Default to STANDARD.",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateBOSBucketStorageClass(),
			},

			// compute attributes
			"etag": {
				Type:        schema.TypeString,
				Description: "Etag generated of the object.",
				Computed:    true,
			},
			"last_modified": {
				Type:        schema.TypeString,
				Description: "Last modified date of the object.",
				Computed:    true,
			},
			"content_encoding": {
				Type:        schema.TypeString,
				Description: "Encoding of the object.",
				Computed:    true,
			},
		},
	}
}

func resourceBaiduCloudBucketObjectPut(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)

	action := "Create bucket " + bucket + " object " + key

	var (
		err  error
		body *bce.Body
	)
	if source, ok := d.GetOk("source"); ok {
		body, err = bce.NewBodyFromFile(source.(string))
	} else if content, ok := d.GetOk("content"); ok {
		body, err = bce.NewBodyFromString(content.(string))
	} else {
		err = fmt.Errorf("The source and content cannot be empty at the same time.")
	}
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_bucket_object", action, BCESDKGoERROR)
	}

	args := buildBaiduCloudBucketObjectArgs(d)
	_, err = client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
		return bosClient.PutObject(bucket, key, body, args)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_bucket_object", action, BCESDKGoERROR)
	}
	d.SetId(key)

	cannedAcl, ok := d.GetOk("acl")
	if ok {
		_, err := client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
			return nil, bosClient.PutObjectAclFromCanned(bucket, key, cannedAcl.(string))
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_bucket_object", action, BCESDKGoERROR)
		}
	}

	return resourceBaiduCloudBucketObjectRead(d, meta)
}

func resourceBaiduCloudBucketObjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	bosService := &BosService{client}

	bucket := d.Get("bucket").(string)
	key := d.Id()

	action := "Query bucket " + bucket + " object " + key

	// read bos bucket object meta
	raw, err := client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
		return bosClient.GetObjectMeta(bucket, key)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_bucket_object", action, BCESDKGoERROR)
	}

	result, _ := raw.(*api.GetObjectMetaResult)
	d.Set("cache_control", result.CacheControl)
	d.Set("content_disposition", result.ContentDisposition)
	d.Set("content_md5", result.ContentMD5)
	d.Set("content_type", result.ContentType)
	d.Set("content_length", result.ContentLength)
	d.Set("expires", result.Expires)
	d.Set("user_meta", result.UserMeta)
	d.Set("content_sha256", result.ContentSha256)
	d.Set("content_crc32", result.ContentCrc32)
	d.Set("storage_class", result.StorageClass)
	d.Set("etag", result.ETag)
	d.Set("last_modified", result.LastModified)
	d.Set("content_encoding", result.ContentEncoding)

	// read bos bucket object acl
	acl, err := bosService.resourceBaiduCloudBucketObjectReadAcl(bucket, key)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_bucket_object", action, BCESDKGoERROR)
	}
	d.Set("acl", acl)

	return nil
}

func resourceBaiduCloudBucketObjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)
	action := "Delete bucket " + bucket + " object " + key

	_, err := client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
		return nil, bosClient.DeleteObject(bucket, key)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_bucket_object", action, BCESDKGoERROR)
	}

	return nil
}

func buildBaiduCloudBucketObjectArgs(d *schema.ResourceData) *api.PutObjectArgs {
	args := &api.PutObjectArgs{}

	if val, ok := d.GetOk("cache_control"); ok {
		args.CacheControl = val.(string)
	}
	if val, ok := d.GetOk("content_disposition"); ok {
		args.ContentDisposition = val.(string)
	}
	if val, ok := d.GetOk("content_md5"); ok {
		args.ContentMD5 = val.(string)
	}
	if val, ok := d.GetOk("content_type"); ok {
		args.ContentType = val.(string)
	}
	if val, ok := d.GetOk("content_length"); ok {
		args.ContentLength = int64(val.(int))
	}
	if val, ok := d.GetOk("expires"); ok {
		args.Expires = val.(string)
	}
	if val, ok := d.GetOk("user_meta"); ok {
		raw := val.(map[string]interface{})
		meta := make(map[string]string, len(raw))
		for k, v := range raw {
			meta[k] = v.(string)
		}
		args.UserMeta = meta
	}
	if val, ok := d.GetOk("content_sha256"); ok {
		args.ContentSha256 = val.(string)
	}
	if val, ok := d.GetOk("content_crc32"); ok {
		args.ContentCrc32 = val.(string)
	}
	if val, ok := d.GetOk("storage_class"); ok {
		args.StorageClass = val.(string)
	}
	if val, ok := d.GetOk("process"); ok {
		args.Process = val.(string)
	}

	return args
}
