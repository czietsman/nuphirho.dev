# Changelog
# 2026-03-06  Initial apply: blog CNAME, root A records, www CNAME
# 2026-03-06  Migrate S3 backend endpoint to endpoints.s3

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

# Blog subdomain CNAME to Hashnode
resource "cloudflare_record" "blog" {
  zone_id = data.cloudflare_zone.nuphirho.id
  name    = "blog"
  content = "hashnode.network"
  type    = "CNAME"
  proxied = true
  ttl     = 1 # Auto when proxied
}

# Root domain A records for GitHub Pages
locals {
  github_pages_ips = [
    "185.199.108.153",
    "185.199.109.153",
    "185.199.110.153",
    "185.199.111.153",
  ]
}

resource "cloudflare_record" "root" {
  for_each = toset(local.github_pages_ips)

  zone_id = data.cloudflare_zone.nuphirho.id
  name    = "@"
  content = each.value
  type    = "A"
  proxied = true
  ttl     = 1 # Auto when proxied
}

# www CNAME to GitHub Pages
resource "cloudflare_record" "www" {
  zone_id = data.cloudflare_zone.nuphirho.id
  name    = "www"
  content = "czietsman.github.io"
  type    = "CNAME"
  proxied = true
  ttl     = 1 # Auto when proxied
}
