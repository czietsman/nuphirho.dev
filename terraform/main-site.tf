# Changelog
# 2026-09-28  Redirect www.nuphirho.dev to nuphirho.dev with a zone redirect rule

resource "cloudflare_pages_project" "main" {
  account_id        = var.cloudflare_account_id
  name              = "nuphirho-main"
  production_branch = "main"
}

resource "cloudflare_pages_domain" "main_root" {
  account_id   = var.cloudflare_account_id
  project_name = cloudflare_pages_project.main.name
  name         = "nuphirho.dev"
}

resource "cloudflare_pages_domain" "main_www" {
  account_id   = var.cloudflare_account_id
  project_name = cloudflare_pages_project.main.name
  name         = "www.nuphirho.dev"
}

resource "cloudflare_ruleset" "www_redirect" {
  zone_id     = data.cloudflare_zone.nuphirho.id
  name        = "Redirect www to apex"
  description = "Permanently redirect www.nuphirho.dev to nuphirho.dev, keeping path and query string"
  kind        = "zone"
  phase       = "http_request_dynamic_redirect"

  rules = [{
    ref         = "www_to_apex"
    description = "www.nuphirho.dev to nuphirho.dev"
    expression  = "(http.host eq \"www.nuphirho.dev\")"
    action      = "redirect"
    action_parameters = {
      from_value = {
        status_code           = 301
        preserve_query_string = true
        target_url = {
          expression = "concat(\"https://nuphirho.dev\", http.request.uri.path)"
        }
      }
    }
  }]
}
