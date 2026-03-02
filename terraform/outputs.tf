output "zone_id" {
  description = "Cloudflare zone ID for nuphirho.dev"
  value       = data.cloudflare_zone.nuphirho.id
}

output "blog_record" {
  description = "Blog subdomain DNS record"
  value       = cloudflare_record.blog.hostname
}
