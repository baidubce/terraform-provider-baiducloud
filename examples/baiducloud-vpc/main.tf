provider "baiducloud" {}

resource "baiducloud_vpc" "default" {
  name = "${var.name}"
  description = "${var.description}"
  cidr = "${var.cidr}"

  tags {
    tag_key = "tagA"
    tag_value = "tagA"
  }
  tags {
    tag_key = "tagB"
    tag_value = "tagB"
  }
}

data "baiducloud_vpcs" "default" {
  vpc_id = "${baiducloud_vpc.default.id}"
}
