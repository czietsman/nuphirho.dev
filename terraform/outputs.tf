output "zone_id" {
  description = "Cloudflare zone ID for nuphirho.dev"
  value       = data.cloudflare_zone.nuphirho.id
}

output "blog_record" {
  description = "Blog subdomain DNS record"
  value       = cloudflare_record.blog.hostname
}

output "root_record" {
  description = "Root domain CNAME record for Cloudflare Pages"
  value       = cloudflare_record.root.hostname
}

output "www_record" {
  description = "www subdomain CNAME record"
  value       = cloudflare_record.www.hostname
}
