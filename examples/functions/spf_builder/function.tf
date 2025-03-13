locals {
  domain               = "malmeida.dev"
  overflow             = "_spf%d"
  txt_max_size         = 255
  domain_on_record_key = false
  parts = [
    "v=spf1",
    "include:eu.mailgun.org",
    "include:_spf.mx.cloudflare.net",
    "include:stspg-customer.com",
    "~all",
  ]
  flatten = [
    "eu.mailgun.org",
    "_spf.mx.cloudflare.net",
    "stspg-customer.com",
    "_spf.eu.mailgun.org",
  ]
}

output "spf_record" {
  value = provider::dnshelper::spf_builder(
    local.domain,
    local.overflow,
    local.txt_max_size,
    local.domain_on_record_key,
    local.parts,
    local.flatten,
  )
}
