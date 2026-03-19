data "sitecoreedge_tenant" "current" {}

output "tenant_id" {
  value = data.sitecoreedge_tenant.current.tenant_id
}

output "media_url" {
  value = "https://edge.sitecorecloud.io/${data.sitecoreedge_tenant.current.tenant_id}/media/"
}