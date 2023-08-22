package iam

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/baidubce/bce-sdk-go/services/iam"
	"github.com/baidubce/bce-sdk-go/services/iam/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/encryption"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func ResourceAccessKey() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to manage IAM access key. \n\n" +
			"More information can be found in the [Developer Guide](https://cloud.baidu.com/doc/IAM/s/mjx35fixq). \n\n",

		Create: resourceAccessKeyCreate,
		Read:   resourceAccessKeyRead,
		Update: resourceAccessKeyUpdate,
		Delete: resourceAccessKeyDelete,

		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Description: "The name of the IAM user associated with this access key.",
				Required:    true,
				ForceNew:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the access key is enabled. Default to `true`.",
				Optional:    true,
				Default:     true,
			},
			"pgp_key": {
				Type: schema.TypeString,
				Description: "Either a base-64 encoded PGP public key, or a keybase username in the form " +
					"`keybase:some_person_that_exists`, for use in the `encrypted_secret` output attribute. " +
					"If providing a base-64 encoded PGP public key, make sure to provide the \"raw\" version " +
					"and not the \"armored\" one (e.g. avoid passing the `-a` option to `gpg --export`).",
				ForceNew:     true,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"secret_file": {
				Type:         schema.TypeString,
				Description:  "The path of the file in which to save the access key.",
				ForceNew:     true,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"secret": {
				Type: schema.TypeString,
				Description: "Secret access key. Note that this will be written to the state file. If you use this, " +
					"please protect your backend state file judiciously. Alternatively, you may supply a `pgp_key` " +
					"instead, which will prevent the secret from being stored in plaintext, at the cost of " +
					"preventing the use of the secret key in automation.",
				Computed:  true,
				Sensitive: true,
			},
			"encrypted_secret": {
				Type: schema.TypeString,
				Description: "Encrypted secret, base64 encoded, if `pgp_key` was specified. The encrypted secret " +
					"may be decrypted using the command line, for example: " +
					"`terraform output -raw encrypted_secret | base64 --decode | keybase pgp decrypt`.",
				Computed: true,
			},
			"key_fingerprint": {
				Type:        schema.TypeString,
				Description: "Fingerprint of the PGP key used to encrypt the secret.",
				Computed:    true,
			},
			"create_time": {
				Type:        schema.TypeString,
				Description: "Date and time in RFC3339 format that the access key was created.",
				Computed:    true,
			},
			"last_used_time": {
				Type:        schema.TypeString,
				Description: "Date and time in RFC3339 format that the access key was last used.",
				Computed:    true,
			},
		},
	}
}

func resourceAccessKeyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	username := d.Get("username").(string)

	raw, err := conn.WithIamClient(func(client *iam.Client) (interface{}, error) {
		return client.CreateAccessKey(username)
	})
	log.Printf("[DEBUG] Create IAM access key")
	if err != nil {
		return fmt.Errorf("error creating IAM access key: %w", err)
	}

	result := raw.(*api.CreateAccessKeyResult)
	d.SetId(result.Id)

	if secretFile, ok := d.GetOk("secret_file"); ok {
		fileContent := map[string]string{
			"AccessKeyId":     result.Id,
			"AccessKeySecret": result.Secret,
		}
		bytes, err := json.MarshalIndent(fileContent, "", "\t")
		if err != nil {
			return fmt.Errorf("error marshaling access key: %w", err)
		}
		err = ioutil.WriteFile(secretFile.(string), bytes, 0644)
		if err != nil {
			return fmt.Errorf("error writing access key to file: %w", err)
		}
	}

	if v, ok := d.GetOk("pgp_key"); ok {
		pgpKey := v.(string)
		encryptionKey, err := encryption.RetrieveGPGKey(pgpKey)
		if err != nil {
			return fmt.Errorf("error retrieving GPG key: %w", err)
		}
		fingerprint, encrypted, err := encryption.EncryptValue(encryptionKey, result.Secret, "IAM Access Key Secret")
		if err != nil {
			return fmt.Errorf("error encrypting secret: %w", err)
		}
		if err := d.Set("key_fingerprint", fingerprint); err != nil {
			return fmt.Errorf("error setting key_fingerprint: %w", err)
		}
		if err := d.Set("encrypted_secret", encrypted); err != nil {
			return fmt.Errorf("error setting encrypted_secret: %w", err)
		}
	} else {
		if err := d.Set("secret", result.Secret); err != nil {
			return fmt.Errorf("error setting secret: %w", err)
		}
	}

	return resourceAccessKeyUpdate(d, meta)
}

func resourceAccessKeyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	username := d.Get("username").(string)
	accessKey, err := FindAccessKey(conn, username, d.Id())

	log.Printf("[DEBUG] Read IAM access key (%s) result: %+v", d.Id(), accessKey)
	if err != nil {
		return fmt.Errorf("error reading IAM access key (%s): %w", d.Id(), err)
	}

	if err := d.Set("enabled", accessKey.Enabled); err != nil {
		return fmt.Errorf("error setting enabled: %w", err)
	}
	if err := d.Set("create_time", accessKey.CreateTime.Format(time.RFC3339)); err != nil {
		return fmt.Errorf("error setting create_time: %w", err)
	}
	if err := d.Set("last_used_time", accessKey.LastUsedTime.Format(time.RFC3339)); err != nil {
		return fmt.Errorf("error setting last_used_time: %w", err)
	}
	return nil
}

func resourceAccessKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	if err := updateAccessKeyEnabled(d, conn); err != nil {
		return fmt.Errorf("error updating IAM access key enabled (%s): %w", d.Id(), err)
	}

	return resourceAccessKeyRead(d, meta)
}

func resourceAccessKeyDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	username := d.Get("username").(string)

	_, err := conn.WithIamClient(func(client *iam.Client) (interface{}, error) {
		return nil, client.DeleteAccessKey(username, d.Id())
	})
	log.Printf("[DEBUG] Delete IAM access key (%s)", d.Id())
	if err != nil {
		return fmt.Errorf("error deleting IAM access key (%s): %w", d.Id(), err)
	}
	return nil
}

func updateAccessKeyEnabled(d *schema.ResourceData, conn *connectivity.BaiduClient) error {
	if d.HasChange("enabled") || (d.IsNewResource() && !d.Get("enabled").(bool)) {
		username := d.Get("username").(string)
		enabled := d.Get("enabled").(bool)
		_, err := conn.WithIamClient(func(client *iam.Client) (interface{}, error) {
			if enabled {
				return client.EnableAccessKey(username, d.Id())
			}
			return client.DisableAccessKey(username, d.Id())
		})
		return err
	}
	return nil
}
