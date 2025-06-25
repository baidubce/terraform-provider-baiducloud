package hpas

import (
	"fmt"

	"github.com/baidubce/bce-sdk-go/services/hpas"
	"github.com/baidubce/bce-sdk-go/services/hpas/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func ResourceInstanceOperation() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to perform operations on an existing HPAS instance, such as starting, stopping, or rebooting. " +
			"This resource does not create or destroy instances.\n\n",

		Create: resourceInstanceOperationCreate,
		Read:   resourceInstanceOperationRead,
		Delete: resourceInstanceOperationDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the instance to operate on.",
			},
			"operation": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"start", "stop", "reboot"}, false),
				Description:  "The operation to perform on the instance. Valid values: `start`, `stop`, and `reboot`.",
			},
		},
	}
}

func resourceInstanceOperationCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*connectivity.BaiduClient)

	instanceID := d.Get("instance_id").(string)
	operation := d.Get("operation").(string)

	var err error
	waitForAvailable := false
	waitForStopped := false

	switch operation {
	case "start":
		waitForAvailable = true
		err = startInstance(conn, instanceID)
	case "stop":
		waitForStopped = true
		err = stopInstance(conn, instanceID)
	case "reboot":
		waitForAvailable = true
		err = rebootInstance(conn, instanceID)
	}

	if err != nil {
		return err
	}

	d.SetId(resource.UniqueId())

	if waitForAvailable {
		_, err := waitInstanceAvailable(conn, instanceID)
		if err != nil {
			return fmt.Errorf("error waiting instance (%s) becoming available: %w", instanceID, err)
		}
	} else if waitForStopped {
		_, err := waitInstanceStopped(conn, instanceID)
		if err != nil {
			return fmt.Errorf("error waiting instance (%s) becoming stopped: %w", instanceID, err)
		}
	}

	return nil
}

func resourceInstanceOperationRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceInstanceOperationDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func startInstance(conn *connectivity.BaiduClient, instanceID string) error {
	_, err := conn.WithHPASClient(func(client *hpas.Client) (interface{}, error) {
		args := api.StartHpasReq{HpasIds: []string{instanceID}}
		return nil, client.StartHpas(&args)
	})

	if err != nil {
		return fmt.Errorf("error starting instance (%s): %w", instanceID, err)
	}
	return nil
}

func stopInstance(conn *connectivity.BaiduClient, instanceID string) error {
	_, err := conn.WithHPASClient(func(client *hpas.Client) (interface{}, error) {
		args := api.StopHpasReq{HpasIds: []string{instanceID}}
		return nil, client.StopHpas(&args)
	})
	if err != nil {
		return fmt.Errorf("error stopping instance (%s): %w", instanceID, err)
	}
	return nil
}

func rebootInstance(conn *connectivity.BaiduClient, instanceID string) error {
	_, err := conn.WithHPASClient(func(client *hpas.Client) (interface{}, error) {
		args := api.RebootHpasReq{HpasIds: []string{instanceID}}
		return nil, client.RebootHpas(&args)
	})
	if err != nil {
		return fmt.Errorf("error rebooting instance (%s): %w", instanceID, err)
	}
	return nil
}
