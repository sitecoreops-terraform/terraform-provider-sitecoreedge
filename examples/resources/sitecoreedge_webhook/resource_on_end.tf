# OnEnd webhook
resource "sitecoreedge_webhook" "on_end" {
  label  = "My new webhook"
  uri    = "https://www.mysite.com/hooks/edge-hook"
  method = "POST"
  headers = {
    "x-header" = "var"
    "x-key"    = "secret"
  }
  body = jsonencode({
    event = "rebuild"
    key   = "secret"
  })
  execution_mode = "OnEnd"
  created_by     = "anco"
}
