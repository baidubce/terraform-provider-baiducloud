package bec

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func needPublicIPDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	return !d.Get("need_public_ip").(bool)
}

func needPrepayDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	v, _ := d.Get("payment_method").(string)
	return v != "prepay"
}

func volumeConfigCustomizeDiff(diff *schema.ResourceDiff, meta interface{}) error {
	if diff.HasChange("system_volume") {
		if err := forceNewIf("system_volume.0.name", diff); err != nil {
			return err
		}
		if err := forceNewIf("system_volume.0.size_in_gb", diff); err != nil {
			return err
		}
		if err := forceNewIf("system_volume.0.volume_type", diff); err != nil {
			return err
		}
		return nil
	}

	if diff.HasChange("data_volume") {
		oldV, newV := diff.GetChange("data_volume")
		oldCount := len(oldV.([]interface{}))
		newCount := len(newV.([]interface{}))

		if oldCount == 0 {
			return nil
		}
		if newCount == 0 || oldCount > newCount {
			if err := forceNewIf("data_volume", diff); err != nil {
				return err
			}
			return nil
		}

		for i := 0; i < oldCount; i++ {
			if err := forceNewIf(fmt.Sprintf("data_volume.%d.name", i), diff); err != nil {
				return err
			}
			if err := forceNewIf(fmt.Sprintf("data_volume.%d.size_in_gb", i), diff); err != nil {
				return err
			}
			if err := forceNewIf(fmt.Sprintf("data_volume.%d.volume_type", i), diff); err != nil {
				return err
			}
		}
	}
	return nil
}

func forceNewIf(key string, diff *schema.ResourceDiff) error {
	if diff.HasChange(key) {
		if err := diff.ForceNew(key); err != nil {
			return err
		}
	}
	return nil
}
