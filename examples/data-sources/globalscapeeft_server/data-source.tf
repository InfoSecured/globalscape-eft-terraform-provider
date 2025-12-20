data "globalscapeeft_server" "current" {}

output "globalscapeeft_server_version" {
  value = data.globalscapeeft_server.current.version
}
