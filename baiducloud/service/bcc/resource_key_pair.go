package bcc

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/baidubce/bce-sdk-go/services/bcc"
	"github.com/baidubce/bce-sdk-go/services/bcc/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func ResourceKeyPair() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to manage BCC key pair. \n\n" +
			"More information can be found in the [Developer Guide](https://cloud.baidu.com/doc/BCC/s/ykckicewc). \n\n",

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceKeyPairCreate,
		Read:   resourceKeyPairRead,
		Update: resourceKeyPairUpdate,
		Delete: resourceKeyPairDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "The name of key pair.",
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of key pair.",
				Optional:    true,
			},
			"public_key": {
				Type:        schema.TypeString,
				Description: "The public key of keypair. This field can be set to import an existing public key.",
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
			},
			"private_key_file": {
				Type:        schema.TypeString,
				Description: "The path of the file in which to save the private key.",
				Optional:    true,
			},
			"created_time": {
				Type:        schema.TypeString,
				Description: "The creation time of key pair.",
				Computed:    true,
			},
			"instance_count": {
				Type:        schema.TypeInt,
				Description: "The number of instances bound to key pair.",
				Computed:    true,
			},
			"region_id": {
				Type:        schema.TypeString,
				Description: "The id of the region to which key pair belongs.",
				Computed:    true,
			},
			"fingerprint": {
				Type:        schema.TypeString,
				Description: "The fingerprint of key pair.",
				Computed:    true,
			},
		},
	}
}

func resourceKeyPairCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	name := d.Get("name").(string)
	description := d.Get("description").(string)

	var raw interface{}
	var err error
	if publicKey, ok := d.GetOk("public_key"); ok {
		raw, err = importKeyPair(conn, name, description, publicKey.(string))
	} else {
		raw, err = createKeyPair(conn, name, description)
	}
	log.Printf("[DEBUG] Create BCC key pair result: %+v ", raw)
	if err != nil {
		return fmt.Errorf("error creating BCC key pair: %w", err)
	}
	response := raw.(*api.KeypairResult)

	d.SetId(response.Keypair.KeypairId)

	if privateKeyFileName, ok := d.GetOk("private_key_file"); ok {
		if response.Keypair.PrivateKey != "" {
			err := ioutil.WriteFile(privateKeyFileName.(string), []byte(response.Keypair.PrivateKey), 0644)
			if err != nil {
				return fmt.Errorf("error writing BCC key pair private key to file: %w", err)
			}
		}
	}

	return resourceKeyPairRead(d, meta)
}

func resourceKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	raw, err := conn.WithBccClient(func(client *bcc.Client) (interface{}, error) {
		return client.GetKeypairDetail(d.Id())
	})
	log.Printf("[DEBUG] Read BCC key pair (%s) result: %+v", d.Id(), raw)
	if err != nil {
		return fmt.Errorf("error reading BCC key pair (%s) name: %w", d.Id(), err)
	}
	detail := raw.(*api.KeypairResult).Keypair
	if err := d.Set("name", detail.Name); err != nil {
		return fmt.Errorf("error setting name: %w", err)
	}
	if err := d.Set("description", detail.Description); err != nil {
		return fmt.Errorf("error setting description: %w", err)
	}
	if err := d.Set("created_time", detail.CreatedTime); err != nil {
		return fmt.Errorf("error setting created_time: %w", err)
	}
	if err := d.Set("public_key", detail.PublicKey); err != nil {
		return fmt.Errorf("error setting public_key: %w", err)
	}
	if err := d.Set("instance_count", detail.InstanceCount); err != nil {
		return fmt.Errorf("error setting instance_count: %w", err)
	}
	if err := d.Set("region_id", detail.RegionId); err != nil {
		return fmt.Errorf("error setting region_id: %w", err)
	}
	if err := d.Set("fingerprint", detail.FingerPrint); err != nil {
		return fmt.Errorf("error setting fingerprint: %w", err)
	}

	return resourceKeyPairRead(d, meta)
}

func resourceKeyPairUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	if err := updateName(d, conn); err != nil {
		return fmt.Errorf("error updating BCC key pair (%s) name: %w", d.Id(), err)
	}
	if err := updateDescription(d, conn); err != nil {
		return fmt.Errorf("error updating BCC key pair (%s) description: %w", d.Id(), err)
	}
	return nil
}

func resourceKeyPairDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)
	args := &api.DeleteKeypairArgs{
		KeypairId: d.Id(),
	}
	_, err := conn.WithBccClient(func(client *bcc.Client) (interface{}, error) {
		return nil, client.DeleteKeypair(args)
	})
	log.Printf("[DEBUG] Delete BCC key pair (%s)", d.Id())
	if err != nil {
		return fmt.Errorf("error deleting BCC key pair (%s): %w", d.Id(), err)
	}
	return nil
}

func createKeyPair(conn *connectivity.BaiduClient, name string, description string) (interface{}, error) {
	args := &api.CreateKeypairArgs{
		Name:        name,
		Description: description,
	}
	return conn.WithBccClient(func(client *bcc.Client) (interface{}, error) {
		return client.CreateKeypair(args)
	})
}

func importKeyPair(conn *connectivity.BaiduClient, name string, description string, publicKey string) (interface{}, error) {
	args := &api.ImportKeypairArgs{
		Name:        name,
		Description: description,
		PublicKey:   publicKey,
	}
	return conn.WithBccClient(func(client *bcc.Client) (interface{}, error) {
		return client.ImportKeypair(args)
	})
}

func updateName(d *schema.ResourceData, conn *connectivity.BaiduClient) error {
	if d.HasChange("name") {
		args := &api.RenameKeypairArgs{
			KeypairId: d.Id(),
			Name:      d.Get("name").(string),
		}
		_, err := conn.WithBccClient(func(client *bcc.Client) (interface{}, error) {
			return nil, client.RenameKeypair(args)
		})
		return err
	}
	return nil
}

func updateDescription(d *schema.ResourceData, conn *connectivity.BaiduClient) error {
	if d.HasChange("description") {
		args := &api.KeypairUpdateDescArgs{
			KeypairId:   d.Id(),
			Description: d.Get("description").(string),
		}
		_, err := conn.WithBccClient(func(client *bcc.Client) (interface{}, error) {
			return nil, client.UpdateKeypairDescription(args)
		})
		return err
	}
	return nil
}
