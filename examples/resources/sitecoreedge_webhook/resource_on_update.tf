# OnUpdate webhook
resource "sitecoreedge_webhook" "on_update" {
  label  = "My new webhook"
  uri    = "https://www.mysite.com/hooks/edge-hook"
  method = "POST"
  headers = {
    "x-header" = "var"
    "x-key"    = "secret"
  }
  body_include = jsonencode({
    event = "update"
    key   = "secret"
  })
  execution_mode = "OnUpdate"
  created_by     = "anco"
}
