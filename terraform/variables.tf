variable "cloudflare_api_token" {
  description = "Cloudflare API token for nuphirho.dev: edits DNS records, redirect rules and email routing on the zone, and Pages projects and Workers KV on the account"
  type        = string
  sensitive   = true
}

variable "cloudflare_account_id" {
  description = "Cloudflare account ID (required for Pages and KV resources)"
  type        = string
  sensitive   = true
}

variable "email_routing_privacy_destination" {
  description = "Destination address for the privacy@nuphirho.dev email routing rule"
  type        = string
  sensitive   = true
}

variable "email_routing_contact_destination" {
  description = "Destination address for the contact@nuphirho.dev email routing rule"
  type        = string
  sensitive   = true
}
