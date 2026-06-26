resource "cloudflare_pages_project" "main" {
  account_id        = var.cloudflare_account_id
  name              = "nuphirho-main"
  production_branch = "main"
}

resource "cloudflare_pages_domain" "main_root" {
  account_id   = var.cloudflare_account_id
  project_name = cloudflare_pages_project.main.name
  domain       = "nuphirho.dev"
}

resource "cloudflare_pages_domain" "main_www" {
  account_id   = var.cloudflare_account_id
  project_name = cloudflare_pages_project.main.name
  domain       = "www.nuphirho.dev"
}
