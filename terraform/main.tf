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

    # Cloudflare R2 endpoint — account ID passed via -backend-config or env
    # terraform init -backend-config="endpoint=https://<ACCOUNT_ID>.r2.cloudflarestorage.com"
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

# Root domain CNAME to Hashnode
resource "cloudflare_record" "root" {
  zone_id = data.cloudflare_zone.nuphirho.id
  name    = "@"
  content = "hashnode.network"
  type    = "CNAME"
  proxied = true
  ttl     = 1 # Auto when proxied
}

# www redirect to root (CNAME to Hashnode, Hashnode handles the redirect)
resource "cloudflare_record" "www" {
  zone_id = data.cloudflare_zone.nuphirho.id
  name    = "www"
  content = "hashnode.network"
  type    = "CNAME"
  proxied = true
  ttl     = 1 # Auto when proxied
}
