# Changelog
# 2026-06-06  Cloudflare Pages project + KV namespace for blog.nuphirho.dev

resource "cloudflare_workers_kv_namespace" "blog_analytics" {
  account_id = var.cloudflare_account_id
  title      = "blog-analytics"
}

resource "cloudflare_pages_project" "blog" {
  account_id        = var.cloudflare_account_id
  name              = "nuphirho-blog"
  production_branch = "main"

  deployment_configs {
    production {
      kv_namespaces = {
        BLOG_ANALYTICS = cloudflare_workers_kv_namespace.blog_analytics.id
      }
    }
    preview {
      kv_namespaces = {
        BLOG_ANALYTICS = cloudflare_workers_kv_namespace.blog_analytics.id
      }
    }
  }
}

resource "cloudflare_pages_domain" "blog" {
  account_id   = var.cloudflare_account_id
  project_name = cloudflare_pages_project.blog.name
  domain       = "blog.nuphirho.dev"
}
