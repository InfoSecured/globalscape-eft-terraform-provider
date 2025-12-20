terraform {
  required_providers {
    globalscapeeft = {
    source  = "InfoSecured/globalscapeeft"
      version = "0.1.0"
    }
  }
}

provider "globalscapeeft" {
  host                 = "https://eft.example.com:4450/admin"
  username             = "api_admin"
  password             = "change-me"
  auth_type            = "EFT"
  insecure_skip_verify = true
}
