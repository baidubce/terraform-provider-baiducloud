package baiducloud

import (
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func postPaidDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	return d.Get("payment_timing").(string) == "Postpaid"
}

func appServerGroupPortHealthCheckHTTPSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	strs := strings.Split(k, ".")
	if len(strs) == 3 {
		key := "port_list." + strs[1] + ".healthCheck"
		value := ""
		if v, ok := d.GetOk(key); ok {
			value = v.(string)
		} else {
			key = "port_list." + strs[1] + ".type"
			value = d.Get(key).(string)
		}

		return value != HTTP
	}

	return false
}

func appServerGroupPortHealthCheckUDPSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	strs := strings.Split(k, ".")
	if len(strs) == 3 {
		key := "port_list." + strs[1] + ".healthCheck"
		value := ""
		if v, ok := d.GetOk(key); ok {
			value = v.(string)
		} else {
			key = "port_list." + strs[1] + ".type"
			value = d.Get(key).(string)
		}

		return value != UDP
	}

	return false
}

func appBlbProtocolTCPUDPSSLSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("protocol"); ok {
		return stringInSlice([]string{TCP, UDP, SSL}, v.(string))
	}

	return true
}

func appBlbProtocolTCPUDPHTTPSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("protocol"); ok {
		return stringInSlice([]string{TCP, UDP, HTTP}, v.(string))
	}

	return true
}

func cfcTriggerSourceTypeSuppressFunc(sourceType []string) schema.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *schema.ResourceData) bool {
		sType := d.Get("source_type").(string)
		for _, t := range sourceType {
			if sType == t {
				return false
			}
		}

		return true
	}
}
