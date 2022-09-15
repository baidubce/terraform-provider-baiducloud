package cdn

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"reflect"
	"strings"
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

func refererACLDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	if k == "referer_acl.#" && new == "0" {
		oldV, _ := d.GetChange("referer_acl")
		oldRefererACL := expandRefererACL(oldV.([]interface{}))
		newRefererACL := expandRefererACL(nil)
		return reflect.DeepEqual(oldRefererACL, newRefererACL)
	}
	return false
}

func ipACLDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	if k == "ip_acl.#" && new == "0" {
		oldV, _ := d.GetChange("ip_acl")
		oldIpACL := expandIpACL(oldV.([]interface{}))
		newIpACL := expandIpACL(nil)
		return reflect.DeepEqual(oldIpACL, newIpACL)
	}
	return false
}

func uaACLDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	if k == "ua_acl.#" && new == "0" {
		oldV, _ := d.GetChange("ua_acl")
		oldUaACL := expandUaACL(oldV.([]interface{}))
		newUaACL := expandUaACL(nil)
		return reflect.DeepEqual(oldUaACL, newUaACL)
	}
	return false
}

func corsDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	if k == "cors.#" && new == "0" {
		oldV, _ := d.GetChange("cors")
		oldCors := expandCors(oldV.([]interface{}))
		newCors := expandCors(nil)
		return reflect.DeepEqual(oldCors, newCors)
	}
	if strings.HasPrefix(k, "cors.0.origin_list") {
		return d.Get("cors.0.allow").(string) == "off"
	}
	return false
}

func accessLimitDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	if k == "access_limit.#" && new == "0" {
		oldV, _ := d.GetChange("access_limit")
		oldAccessLimit := expandAccessLimit(oldV.([]interface{}))
		newAccessLimit := expandAccessLimit(nil)
		return reflect.DeepEqual(oldAccessLimit, newAccessLimit)
	}
	if k == "access_limit.0.limit" {
		return d.Get("access_limit.0.enabled").(bool) == false
	}
	return false
}

func trafficLimitDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	if k == "traffic_limit.#" && new == "0" {
		oldV, _ := d.GetChange("traffic_limit")
		oldTrafficLimit := expandTrafficLimit(oldV.([]interface{}))
		newTrafficLimit := expandTrafficLimit(nil)
		return reflect.DeepEqual(oldTrafficLimit, newTrafficLimit)
	}
	if strings.HasPrefix(k, "traffic_limit.0.limit") {
		return d.Get("traffic_limit.0.enable").(bool) == false
	}
	return false
}
