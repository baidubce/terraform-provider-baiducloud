package bec_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/acctest"
)

func TestAccVMInstances(t *testing.T) {
	resourceName := "baiducloud_bec_vm_instance" + ".test"
	dataSourceName := "data.baiducloud_bec_vm_instances" + ".test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acctest.PreCheck(t) },
		Providers: acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccVMInstancesConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "vm_instances.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "service_id", dataSourceName, "vm_instances.0.service_id"),
					resource.TestCheckResourceAttrPair(resourceName, "vm_name", dataSourceName, "vm_instances.0.vm_name"),
					resource.TestCheckResourceAttrPair(resourceName, "host_name", dataSourceName, "vm_instances.0.host_name"),
					resource.TestCheckResourceAttrPair(resourceName, "region_id", dataSourceName, "vm_instances.0.region_id"),
					resource.TestCheckResourceAttrPair(resourceName, "cpu", dataSourceName, "vm_instances.0.cpu"),
					resource.TestCheckResourceAttrPair(resourceName, "memory", dataSourceName, "vm_instances.0.memory"),
					resource.TestCheckResourceAttrPair(resourceName, "image_type", dataSourceName, "vm_instances.0.image_type"),
					resource.TestCheckResourceAttrPair(resourceName, "image_id", dataSourceName, "vm_instances.0.image_id"),
					resource.TestCheckResourceAttrPair(resourceName, "system_volume.0.name", dataSourceName, "vm_instances.0.system_volume.0.name"),
					resource.TestCheckResourceAttrPair(resourceName, "system_volume.0.size_in_gb", dataSourceName, "vm_instances.0.system_volume.0.size_in_gb"),
					resource.TestCheckResourceAttrPair(resourceName, "system_volume.0.volume_type", dataSourceName, "vm_instances.0.system_volume.0.volume_type"),
					resource.TestCheckResourceAttrPair(resourceName, "system_volume.0.pvc_name", dataSourceName, "vm_instances.0.system_volume.0.pvc_name"),
					resource.TestCheckResourceAttrPair(resourceName, "data_volume.0.name", dataSourceName, "vm_instances.0.data_volume.0.name"),
					resource.TestCheckResourceAttrPair(resourceName, "data_volume.0.size_in_gb", dataSourceName, "vm_instances.0.data_volume.0.size_in_gb"),
					resource.TestCheckResourceAttrPair(resourceName, "data_volume.0.volume_type", dataSourceName, "vm_instances.0.data_volume.0.volume_type"),
					resource.TestCheckResourceAttrPair(resourceName, "data_volume.0.pvc_name", dataSourceName, "vm_instances.0.data_volume.0.pvc_name"),
					resource.TestCheckResourceAttrPair(resourceName, "need_public_ip", dataSourceName, "vm_instances.0.need_public_ip"),
					resource.TestCheckResourceAttrPair(resourceName, "need_ipv6_public_ip", dataSourceName, "vm_instances.0.need_ipv6_public_ip"),
					resource.TestCheckResourceAttrPair(resourceName, "bandwidth", dataSourceName, "vm_instances.0.bandwidth"),
					resource.TestCheckResourceAttrPair(resourceName, "dns_config.0.dns_type", dataSourceName, "vm_instances.0.dns_config.0.dns_type"),
					resource.TestCheckResourceAttrPair(resourceName, "status", dataSourceName, "vm_instances.0.status"),
					resource.TestCheckResourceAttrPair(resourceName, "internal_ip", dataSourceName, "vm_instances.0.internal_ip"),
					resource.TestCheckResourceAttrPair(resourceName, "public_ip", dataSourceName, "vm_instances.0.public_ip"),
					resource.TestCheckResourceAttrPair(resourceName, "ipv6_public_ip", dataSourceName, "vm_instances.0.ipv6_public_ip"),
					resource.TestCheckResourceAttrPair(resourceName, "create_time", dataSourceName, "vm_instances.0.create_time"),
				),
			},
		},
	})
}

func testAccVMInstancesConfig() string {
	return acctest.ConfigCompose(testAccVMInstanceConfig_create(), fmt.Sprintf(`
data "baiducloud_bec_vm_instances" "test" {
	keyword_type = "instanceId"
	keyword = baiducloud_bec_vm_instance.test.id
}`))
}
