# Changelog
# 2026-03-06  Initial apply: blog CNAME, root A records, www CNAME
# 2026-03-06  Migrate S3 backend endpoint to endpoints.s3
# 2026-03-06  Set blog CNAME to DNS-only (Hashnode/Vercel requires unproxied)
# 2026-03-08  Point blog CNAME to GitHub Pages (static frontend)
# 2026-06-06  Retarget blog CNAME to Cloudflare Pages (nuphirho-blog.pages.dev)
# 2026-06-09  Revert blog CNAME to GitHub Pages pending Cloudflare Pages testing
# 2026-06-09  Switch blog CNAME to Cloudflare Pages (testing passed)
# 2026-03-26  Add MX, SPF, DMARC for Cloudflare Email Routing
# 2026-06-26  Switch root and www DNS from GitHub Pages to Cloudflare Pages (nuphirho-main)

terraform {
  required_version = ">= 1.0"

  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 4.0"
    }
  }

  backend "s3" {
    bucket = "nuphirho-terraform-state"
    key    = "nuphirho.dev/terraform.tfstate"
    region = "auto"

    # Cloudflare R2 endpoint -- account ID passed via -backend-config or env
    # terraform init -backend-config="endpoints={s3=\"https://<ACCOUNT_ID>.r2.cloudflarestorage.com\"}"
    skip_credentials_validation = true
    skip_metadata_api_check     = true
    skip_region_validation      = true
    skip_requesting_account_id  = true
    skip_s3_checksum            = true
    use_path_style              = true
  }
}

provider "cloudflare" {
  api_token = var.cloudflare_api_token
}

data "cloudflare_zone" "nuphirho" {
  name = "nuphirho.dev"
}

# Blog subdomain CNAME to Cloudflare Pages
resource "cloudflare_record" "blog" {
  zone_id = data.cloudflare_zone.nuphirho.id
  name    = "blog"
  content = "nuphirho-blog.pages.dev"
  type    = "CNAME"
  proxied = true
  ttl     = 1 # Auto when proxied
}

# Root domain CNAME to Cloudflare Pages
resource "cloudflare_record" "root" {
  zone_id = data.cloudflare_zone.nuphirho.id
  name    = "@"
  content = "nuphirho-main.pages.dev"
  type    = "CNAME"
  proxied = true
  ttl     = 1 # Auto when proxied
}

# www CNAME to Cloudflare Pages
resource "cloudflare_record" "www" {
  zone_id = data.cloudflare_zone.nuphirho.id
  name    = "www"
  content = "nuphirho-main.pages.dev"
  type    = "CNAME"
  proxied = true
  ttl     = 1 # Auto when proxied
}

# ── Email Routing (Cloudflare) ────────────────────────────────────────

# MX records for Cloudflare Email Routing
locals {
  cloudflare_mx_records = {
    isaac = { server = "isaac.mx.cloudflare.net", priority = 84 }
    linda = { server = "linda.mx.cloudflare.net", priority = 4 }
    amir  = { server = "amir.mx.cloudflare.net", priority = 21 }
  }
}

resource "cloudflare_record" "mx" {
  for_each = local.cloudflare_mx_records

  zone_id  = data.cloudflare_zone.nuphirho.id
  name     = "@"
  content  = each.value.server
  type     = "MX"
  priority = each.value.priority
  proxied  = false
  ttl      = 1
}

# SPF record authorising Cloudflare Email Routing
resource "cloudflare_record" "spf" {
  zone_id = data.cloudflare_zone.nuphirho.id
  name    = "@"
  content = "v=spf1 include:_spf.mx.cloudflare.net ~all"
  type    = "TXT"
  proxied = false
  ttl     = 1
}

# DMARC policy (quarantine; upgrade to reject after monitoring)
resource "cloudflare_record" "dmarc" {
  zone_id = data.cloudflare_zone.nuphirho.id
  name    = "_dmarc"
  content = "v=DMARC1; p=quarantine; rua=mailto:christo@nuphirho.dev"
  type    = "TXT"
  proxied = false
  ttl     = 1
}

# ── Email Routing Rules ───────────────────────────────────────

locals {
  user_email_routings = {
    "privacy" = var.email_routing_privacy_destination
    "contact" = var.email_routing_contact_destination
  }
}

resource "cloudflare_email_routing_rule" "users" {
  for_each = local.user_email_routings

  zone_id = data.cloudflare_zone.nuphirho.id
  name    = "Forward ${each.key}@nuphirho.dev"
  enabled = true

  matcher {
    type  = "literal"
    field = "to"
    value = "${each.key}@nuphirho.dev"
  }

  action {
    type  = "forward"
    value = [each.value]
  }
}
