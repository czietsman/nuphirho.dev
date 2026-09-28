# Changelog
# 2026-06-06  Cloudflare Pages project + KV namespace for blog.nuphirho.dev
# 2026-09-28  Preview deployments get their own KV namespace, so they no
#             longer read or write the production visit counts

resource "cloudflare_workers_kv_namespace" "blog_analytics" {
  account_id = var.cloudflare_account_id
  title      = "blog-analytics"
}

resource "cloudflare_workers_kv_namespace" "blog_analytics_preview" {
  account_id = var.cloudflare_account_id
  title      = "blog-analytics-preview"
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
        BLOG_ANALYTICS = cloudflare_workers_kv_namespace.blog_analytics_preview.id
      }
    }
  }
}

resource "cloudflare_pages_domain" "blog" {
  account_id   = var.cloudflare_account_id
  project_name = cloudflare_pages_project.blog.name
  domain       = "blog.nuphirho.dev"
}
