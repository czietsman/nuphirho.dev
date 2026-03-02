output "zone_id" {
  description = "Cloudflare zone ID for nuphirho.dev"
  value       = data.cloudflare_zone.nuphirho.id
}

output "root_record" {
  description = "Root DNS record"
  value       = cloudflare_record.root.hostname
}

output "www_record" {
  description = "WWW DNS record"
  value       = cloudflare_record.www.hostname
}
