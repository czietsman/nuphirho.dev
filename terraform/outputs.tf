output "zone_id" {
  description = "Cloudflare zone ID for nuphirho.dev"
  value       = data.cloudflare_zone.nuphirho.id
}

output "blog_record" {
  description = "Blog subdomain DNS record name"
  value       = cloudflare_dns_record.blog.name
}

output "root_record" {
  description = "Root domain CNAME record name for Cloudflare Pages"
  value       = cloudflare_dns_record.root.name
}

output "www_record" {
  description = "www subdomain CNAME record name"
  value       = cloudflare_dns_record.www.name
}
