/*
Use this data source to query BOS bucket object list.

Example Usage

```hcl
data "baiducloud_bos_bucket_objects" "default" {
  bucket = "my-bucket"
}

output "objects" {
  value = "${data.baiducloud_bos_bucket_objects.default.objects}"
}
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/bos"
	"github.com/baidubce/bce-sdk-go/services/bos/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudBosBucketObjects() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudBosBucketObjectsRead,

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:        schema.TypeString,
				Description: "Bucket name of the objects to retrieve.",
				Required:    true,
			},
			"prefix": {
				Type:        schema.TypeString,
				Description: "Prefix of the objects.",
				Optional:    true,
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: "Output file for saving result.",
				Optional:    true,
				ForceNew:    true,
			},
			"filter": dataSourceFiltersSchema(),

			// Attributes used for result
			"objects": {
				Type:        schema.TypeList,
				Description: "List of the objects.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket": {
							Type:        schema.TypeString,
							Description: "Bucket of the object.",
							Computed:    true,
						},
						"key": {
							Type:        schema.TypeString,
							Description: "Key of the object.",
							Computed:    true,
						},
						"acl": {
							Type:        schema.TypeString,
							Description: "Acl of the object.",
							Computed:    true,
						},
						"cache_control": {
							Type:        schema.TypeString,
							Description: "Caching behavior of the object.",
							Computed:    true,
						},
						"content_disposition": {
							Type:        schema.TypeString,
							Description: "Content disposition of the object.",
							Computed:    true,
						},
						"content_md5": {
							Type:        schema.TypeString,
							Description: "MD5 value of the object content defined in RFC2616.",
							Computed:    true,
						},
						"content_type": {
							Type:        schema.TypeString,
							Description: "Content type of the object data.",
							Computed:    true,
						},
						"content_length": {
							Type:        schema.TypeInt,
							Description: "Content length of the object.",
							Computed:    true,
						},
						"expires": {
							Type:        schema.TypeString,
							Description: "Expire date of the object.",
							Computed:    true,
						},
						"user_meta": {
							Type:        schema.TypeMap,
							Description: "Metadata of the object.",
							Computed:    true,
						},

						"content_sha256": {
							Type:        schema.TypeString,
							Description: "Sha256 value of the object.",
							Computed:    true,
						},
						"content_crc32": {
							Type:        schema.TypeString,
							Description: "Crc(cyclic redundancy check code) value of the object.",
							Computed:    true,
						},
						"content_encoding": {
							Type:        schema.TypeString,
							Description: "Encoding of the object.",
							Computed:    true,
						},

						"last_modified": {
							Type:        schema.TypeString,
							Description: "Last modifyed time of the object.",
							Computed:    true,
						},
						"etag": {
							Type:        schema.TypeString,
							Description: "Etag of the object.",
							Computed:    true,
						},
						"storage_class": {
							Type:        schema.TypeString,
							Description: "Storage class of the object.",
							Computed:    true,
						},

						"size": {
							Type:        schema.TypeInt,
							Description: "Size of the object.",
							Computed:    true,
						},
						"owner_id": {
							Type:        schema.TypeString,
							Description: "Owner id of the object.",
							Computed:    true,
						},
						"owner_name": {
							Type:        schema.TypeString,
							Description: "Owner name of the object.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudBosBucketObjectsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	bosService := BosService{client}

	var (
		bucket     string
		prefix     string
		outputFile string
	)
	if v, ok := d.GetOk("bucket"); ok {
		bucket = v.(string)
	}
	if v, ok := d.GetOk("prefix"); ok {
		prefix = v.(string)
	}
	if v, ok := d.GetOk("output_file"); ok {
		outputFile = v.(string)
	}

	action := "Query bucket " + bucket + " objects"

	objects, err := bosService.ListAllObjects(bucket, prefix)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_bucket_objects", action, BCESDKGoERROR)
	}

	objectsResult := make([]map[string]interface{}, 0, len(objects))
	for _, obj := range objects {
		// read metadata
		objMap, err := dataSourceBaiduCloudBosBucketObjectsReadMeta(bucket, obj.Key, meta)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_bucket_objects", action, BCESDKGoERROR)
		}

		// read object acl
		acl, err := bosService.resourceBaiduCloudBucketObjectReadAcl(bucket, obj.Key)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_bucket_objects", action, BCESDKGoERROR)
		}
		objMap["acl"] = acl

		// read others
		objMap["key"] = obj.Key
		objMap["size"] = obj.Size
		objMap["owner_id"] = obj.Owner.Id
		objMap["owner_name"] = obj.Owner.DisplayName
		objMap["bucket"] = bucket

		objectsResult = append(objectsResult, objMap)
	}

	FilterDataSourceResult(d, &objectsResult)

	d.Set("objects", objectsResult)
	d.SetId(resource.UniqueId())

	if outputFile != "" {
		if err := writeToFile(outputFile, objectsResult); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_bucket_objects", action, BCESDKGoERROR)
		}
	}

	return nil
}

func dataSourceBaiduCloudBosBucketObjectsReadMeta(bucket, key string, meta interface{}) (map[string]interface{}, error) {
	action := "read bos bucket object meta, bucket: " + bucket + ", key: " + key
	client := meta.(*connectivity.BaiduClient)

	// read bos bucket object meta
	raw, err := client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
		return bosClient.GetObjectMeta(bucket, key)
	})
	if err != nil {
		return nil, err
	}
	addDebug(action, raw)
	result, _ := raw.(*api.GetObjectMetaResult)

	objMap := make(map[string]interface{})

	objMap["cache_control"] = result.CacheControl
	objMap["content_disposition"] = result.ContentDisposition
	objMap["content_md5"] = result.ContentMD5
	objMap["content_type"] = result.ContentType
	objMap["content_length"] = result.ContentLength
	objMap["expires"] = result.Expires
	objMap["user_meta"] = result.UserMeta
	objMap["content_sha256"] = result.ContentSha256
	objMap["content_crc32"] = result.ContentCrc32
	objMap["storage_class"] = result.StorageClass
	objMap["etag"] = result.ETag
	objMap["last_modified"] = result.LastModified
	objMap["content_encoding"] = result.ContentEncoding

	return objMap, nil
}
