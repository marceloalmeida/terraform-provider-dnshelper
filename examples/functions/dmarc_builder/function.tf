locals {
  version          = "DMARC1"
  policy           = "reject"
  subdomain_policy = "quarantine"
  alignment_spf    = "relaxed"
  alignment_dkim   = "strict"
  percent          = 100
  rua              = ["mailto:admin@malmeida.dev"]
  ruf              = ["mailto:alerts@malmeida.dev"]
  failure_options  = 0
  failure_format   = "afrf"
  report_interval  = 86400
}

output "dmarc_record" {
  value = provider::dnshelper::dmarc_builder(
    local.version,
    local.policy,
    local.subdomain_policy,
    local.alignment_spf,
    local.alignment_dkim,
    local.percent,
    local.rua,
    local.ruf,
    local.failure_options,
    local.failure_format,
    local.report_interval
  )
}
