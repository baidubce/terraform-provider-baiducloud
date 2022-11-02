resource "baiducloud_snic" "example" {

  name = "snic_example"
  vpc_id = "vpc-65cz3sw123z2"
  subnet_id = "sbn-yisr456x7dmf"
  ip_address = "192.168.64.4"
  service = "bj.bcebos.com"
  description = "used for bos"

}