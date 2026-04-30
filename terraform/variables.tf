variable "cloudflare_api_token" {
  description = "Cloudflare API token with DNS edit permissions for nuphirho.dev"
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
