package cdn

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func defaultHostDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	oldV, newV := d.GetChange("default_host")

	if old == "" {
		oldV = d.Get("domain")
	}
	if new == "" {
		newV = d.Get("domain")
	}
	return oldV == newV
}

func cacheUrlArgsDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	if k == "cache_url_args.#" {
		oldV, newV := d.GetChange("cache_url_args")
		oldCacheFullUrl := expandCacheUrlArgs(oldV.([]interface{})).CacheFullUrl
		newCacheFullUrl := expandCacheUrlArgs(newV.([]interface{})).CacheFullUrl
		if new == "0" {
			newCacheFullUrl = true
		}
		return oldCacheFullUrl == newCacheFullUrl
	}
	return false
}

func cacheUrlArgsInnerDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	return d.Get("cache_url_args.0.cache_full_url").(bool)
}

func cacheShareDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	oldV, newV := d.GetChange("cache_share")
	oldEnabled := expandCacheShare(oldV.([]interface{})).Enabled
	newEnabled := expandCacheShare(newV.([]interface{})).Enabled

	if new == "0" {
		newEnabled = false
	}
	if d.Get("domain").(string) == expandCacheShare(newV.([]interface{})).SharedWith {
		newEnabled = false
	}

	return !oldEnabled && !newEnabled
}

func mobileAccessDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	oldV, newV := d.GetChange("mobile_access")
	oldEnabled := expandMobileAccess(oldV.([]interface{}))
	newEnabled := expandMobileAccess(newV.([]interface{}))
	if new == "0" {
		newEnabled = false
	}
	return oldEnabled == newEnabled
}
