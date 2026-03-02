output "zone_id" {
  description = "Cloudflare zone ID for nuphirho.dev"
  value       = data.cloudflare_zone.nuphirho.id
}

output "blog_record" {
  description = "Blog subdomain DNS record"
  value       = cloudflare_record.blog.hostname
}

output "root_records" {
  description = "Root domain A records for GitHub Pages"
  value       = [for r in cloudflare_record.root : r.hostname]
}

output "www_record" {
  description = "www subdomain CNAME record"
  value       = cloudflare_record.www.hostname
}
