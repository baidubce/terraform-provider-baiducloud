# Create a new key pair and save the private key to a file
resource "baiducloud_bcc_key_pair" "example" {
    name = "example-key-pair"
    description = "created by terraform"
    private_key_file = "./private-key.txt"
}

# Import an existing public key
resource "baiducloud_bcc_key_pair" "example" {
    name = "example-key-pair"
    description = "created by terraform"
    public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCI6n..."
}