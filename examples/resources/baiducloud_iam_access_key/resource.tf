# Create a new access key and save it to a file
resource "baiducloud_iam_user" "example" {
    name = "tf-user"
}

resource "baiducloud_iam_access_key" "example" {
    username = baiducloud_iam_user.example.name
    secret_file = "access-key.txt"
}

# Create a new access key and encrypt secret using pgp key
resource "baiducloud_iam_access_key" "example" {
    username = "tf-user"
    enabled  = false
    pgp_key = "keybase:some_person_that_exists"
}