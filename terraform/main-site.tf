resource "cloudflare_pages_project" "main" {
  account_id        = var.cloudflare_account_id
  name              = "nuphirho-main"
  production_branch = "main"
}
