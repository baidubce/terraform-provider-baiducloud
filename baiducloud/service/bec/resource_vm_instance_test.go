package bec_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/acctest"
)

func TestAccVMInstance(t *testing.T) {
	resourceName := "baiducloud_bec_vm_instance" + ".test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acctest.PreCheck(t) },
		Providers: acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccVMInstanceConfig_create(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "vm_name", "test_vm_name"),
					resource.TestCheckResourceAttr(resourceName, "host_name", "test-host-name"),
					resource.TestCheckResourceAttr(resourceName, "cpu", "4"),
					resource.TestCheckResourceAttr(resourceName, "memory", "8"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "need_ipv6_public_ip", "false"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth", "1"),
					resource.TestCheckResourceAttr(resourceName, "dns_config.0.dns_type", "DEFAULT"),
					resource.TestCheckResourceAttr(resourceName, "key_config.0.bcc_key_pair_id_list.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccVMInstanceConfig_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "vm_name", "test_vm_name_update"),
					resource.TestCheckResourceAttr(resourceName, "host_name", "test-host-name-update"),
					resource.TestCheckResourceAttr(resourceName, "cpu", "8"),
					resource.TestCheckResourceAttr(resourceName, "memory", "16"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "need_ipv6_public_ip", "true"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth", "10"),
					resource.TestCheckResourceAttr(resourceName, "dns_config.0.dns_type", "LOCAL"),
					resource.TestCheckResourceAttr(resourceName, "key_config.0.bcc_key_pair_id_list.#", "2"),
				),
			},
		},
	})
}

func testAccVMInstanceConfig_create() string {
	return fmt.Sprintf(`
resource "baiducloud_bec_vm_instance" "test" {
    service_id = "s-jw3yel6j"
    vm_name = "test_vm_name"
    host_name = "test-host-name"
    region_id = "cn-maanshan-ct"

    cpu = 4
    memory = 8
    image_type = "bcc"
    image_id = "m-sqj4vgCj"

    system_volume {
        name = "system_volume"
        size_in_gb = 40
        volume_type = "NVME"
    }
    data_volume {
        name = "data_volume1"
        size_in_gb = 20
        volume_type = "SATA"
    }

    need_public_ip = true
    need_ipv6_public_ip = false
    bandwidth = 1

    dns_config {
        dns_type = "DEFAULT"
    }

    key_config {
        type = "bccKeyPair"
        bcc_key_pair_id_list = ["k-cTaWVJcD"]
	}
}
`)
}

func testAccVMInstanceConfig_update() string {
	return fmt.Sprintf(`
resource "baiducloud_bec_vm_instance" "test" {
    service_id = "s-jw3yel6j"
    vm_name = "test_vm_name_update"
    host_name = "test-host-name-update"
    region_id = "cn-maanshan-ct"

    cpu = 8
    memory = 16
    image_type = "bcc"
    image_id = "m-sqj4vgCj"

    system_volume {
        name = "system_volume"
        size_in_gb = 40
        volume_type = "NVME"
    }
    data_volume {
        name = "data_volume1"
        size_in_gb = 20
        volume_type = "SATA"
    }
    data_volume {
        name = "data_volume2"
        size_in_gb = 20
        volume_type = "NVME"
    }

    need_public_ip = true
    need_ipv6_public_ip = true
    bandwidth = 10

    dns_config {
        dns_type = "LOCAL"
    }

    key_config {
        type = "bccKeyPair"
        bcc_key_pair_id_list = ["k-cTaWVJcD", "k-9SSp6luE"]
	}
}
`)
}
