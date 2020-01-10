provider "baiducloud" {}

data "baiducloud_specs" "default" {
  #name_regex        = "bcc.g1.tiny"
  #instance_type     = "General"
  cpu_count         = 1
  memory_size_in_gb = 4
}

data "baiducloud_zones" "default" {
  name_regex = ".*a$"
}

data "baiducloud_images" "default" {
  image_type = "System"
  name_regex = "7.5.*"
  os_name    = "CentOS"
}

data "baiducloud_security_groups" "default" {
  vpc_id = baiducloud_vpc.default.id
}

resource "baiducloud_vpc" "default" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "baiducloud_subnet" "default" {
  name        = var.subnet_name
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = var.subnet_cidr
  description = "subnet created by terraform"
  vpc_id      = baiducloud_vpc.default.id
}

resource "baiducloud_instance" "default" {
  image_id              = data.baiducloud_images.default.images.0.id
  name                  = var.name
  cpu_count             = data.baiducloud_specs.default.specs.0.cpu_count
  memory_capacity_in_gb = data.baiducloud_specs.default.specs.0.memory_size_in_gb
  availability_zone     = data.baiducloud_zones.default.zones.0.zone_name
  subnet_id             = baiducloud_subnet.default.id
  security_groups       = [data.baiducloud_security_groups.default.security_groups.0.id]
  billing = {
    payment_timing = var.payment_timing
  }
}

resource "baiducloud_route_rule" "default" {
  route_table_id      = baiducloud_vpc.default.route_table_id
  source_address      = var.source_address
  destination_address = var.destination_address

  # support custom/vpn/nat
  next_hop_type = "custom"
  next_hop_id   = baiducloud_instance.default.id
  description   = "route rule created by terraform"
}

data "baiducloud_route_rules" "default" {
  vpc_id        = baiducloud_vpc.default.id
  route_rule_id = baiducloud_route_rule.default.id
}
