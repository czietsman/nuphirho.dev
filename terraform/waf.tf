# Own the zone entry point in this phase. Add future rules here rather than
# creating a second ruleset or changing it separately in the dashboard.
resource "cloudflare_ruleset" "public_probe_blocking" {
  zone_id     = data.cloudflare_zone.nuphirho.id
  name        = "Public site probe blocking"
  description = "Block known root-level exploit probes on the public site only"
  kind        = "zone"
  phase       = "http_request_firewall_custom"

  rules = [{
    ref         = "block_known_public_site_probes"
    description = "Block environment, Git, and WordPress probes on apex and www"
    action      = "block"
    enabled     = true
    expression  = <<-EOT
      (http.host in {"nuphirho.dev" "www.nuphirho.dev"} and (
        lower(http.request.uri.path) in {"/.env" "/.git" "/wp-login.php" "/wp-admin" "/xmlrpc.php"} or
        starts_with(lower(http.request.uri.path), "/.env.") or
        starts_with(lower(http.request.uri.path), "/.git/") or
        starts_with(lower(http.request.uri.path), "/wp-admin/")
      ))
    EOT
  }]
}
