resource "baiducloud_eipgroup_attachment" "example" {

  eip_group_id = "eg-example"
  eips = ["100.88.2.121", "100.88.2.122", "240c:4082:ffff:ff01:0:4:0:307"]

}
