resource "globalscapeeft_server_smtp" "default" {
  server              = "smtp.example.net"
  port                = 587
  sender_address      = "eft@example.net"
  sender_name         = "Globalscape EFT"
  login               = "smtp-user"
  password            = "change-me"
  use_authentication  = true
  use_implicit_tls    = false
}
