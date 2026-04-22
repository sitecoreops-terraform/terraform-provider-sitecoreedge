resource "sitecoreedge_token" "main" {
  label      = "Test"
  created_by = "Terraform"
  scopes = [
    "content-#everything#",
    "audience-delivery",
    "audience-preview"
  ]
}

output "token" {
  value     = sitecoreedge_token.main.token
  sensitive = true
}

output "created" {
  value = sitecoreedge_token.main.created
}
