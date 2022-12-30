resource "baiducloud_bec_vm_instance" "example" {

    service_id = "s-jw3y1234"
    vm_name = "vm_name_example"
    host_name = "host-name_example"
    region_id = "cn-maanshan-ct"

    cpu = 4
    memory = 8
    image_type = "bcc"
    image_id = "m-sqj56gCj"

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
    need_ipv6_public_ip = true
    bandwidth = 10

    dns_config {
        dns_type = "DEFAULT"
    }

    key_config {
        type = "bccKeyPair"
        bcc_key_pair_id_list = ["k-cTaMVJcD", "k-9QQp6luE"]
    }
}