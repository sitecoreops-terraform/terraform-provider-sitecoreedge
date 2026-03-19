terraform {
  required_providers {
    sitecore_edge = {
      source = "sitecoreops-terraform/sitecore-edge"
    }
  }
  required_version = ">= 0.1.0"
}

# Configure the Sitecore AI Experience Edge provider
provider "sitecoreedge" {
  # Authenticate with CLI before running terraform
  # The .sitecore folder must be in terraform folder or a parent folder
  # Initialize Sitecore CLI:
  # > dotnet tool install Sitecore.CLI 
  # > dotnet sitecore init
  # > dotnet sitecore plugin add -n Sitecore.Edge.DevEx.Sitecore.Plugin
  # Plugin documentation: https://doc.sitecore.com/xp/en/developers/latest/developer-tools/experience-edge-plugin.html
  # Authenticate by running
  # > dotnet sitecore edge login
  # Choose tenant
  # > dotnet sitecore edge tenant list
  # > dotnet sitecore edge tenant use --tenantId <tenantId>
  use_cli = true
}
