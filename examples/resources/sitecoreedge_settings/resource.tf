resource "sitecoreedge_settings" "main" {
  content_cache_ttl        = "04:00:00"
  content_cache_auto_clear = true
  media_cache_ttl          = "04:00:00"
  media_cache_auto_clear   = true
  tenant_cache_auto_clear  = true
}