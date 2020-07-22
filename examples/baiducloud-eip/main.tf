provider "baiducloud" {
  # option config, you can use assume role as the operation account
  #assume_role {
  #  account_id = "your-account-id"
  #  role_name = "your-role-name"
  #}
}

resource "baiducloud_eip" "my-eip" {
  name              = var.name
  bandwidth_in_mbps = 100

  # support Prepaid/Postpaid
  payment_timing    = "Postpaid"

  # support ByTraffic/ByBandwidth
  billing_method = "ByTraffic"

  # only useful when payment_timing is Prepaid
  # auto_renew_time_unit support month/year
  # auto_renew_time_unit = "month"
  # if auto_renew_time_unit is month, auto_renew_time support 1-9
  # if auto_renew_time_unit is year, auto_renew_time support 1-3
  # auto_renew_time = 1

  tags = {
    "testKey" = "testValue"
  }
}

data "baiducloud_eips" "default" {
  eip = baiducloud_eip.my-eip.id
}
