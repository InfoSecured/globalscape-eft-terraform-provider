data "globalscapeeft_sites" "all" {}

output "site_names" {
  value = data.globalscapeeft_sites.all.sites[*].name
}
