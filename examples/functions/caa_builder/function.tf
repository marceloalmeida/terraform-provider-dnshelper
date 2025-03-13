locals {
  iodef          = "mailto:domain-names@malmeida.dev"
  iodef_critical = true
  issue = [
    "amazon.com",
    "comodoca.com",
    "digicert.com; cansignhttpexchanges=yes",
    "letsencrypt.org",
    "pki.goog; cansignhttpexchanges=yes",
    "sectigo.com",
    "ssl.com",
  ]
  issue_critical = false
  issuewild = [
    "amazon.com",
    "comodoca.com",
    "digicert.com; cansignhttpexchanges=yes",
    "letsencrypt.org",
    "pki.goog; cansignhttpexchanges=yes",
    "sectigo.com",
    "ssl.com",
  ]
  issuewild_critical = false
}

output "caa_records" {
  value = provider::dnshelper::caa_builder(
    local.iodef,
    local.iodef_critical,
    local.issue,
    local.issue_critical,
    local.issuewild,
    local.issuewild_critical,
  )
}
