terraform {
  required_providers {
    sitecoreedge = {
      source = "sitecoreops-terraform/sitecoreedge"
    }
  }
  required_version = ">= 0.1.0"
}

# Configure the Sitecore AI Experience Edge provider
provider "sitecoreedge" {
  client_id     = "<your_edge_admin_client_id>"
  client_secret = "<your_edge_admin_client_secret>"
}
