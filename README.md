# Terraform DNS Helper Provider

[![Build Status](https://github.com/marceloalmeida/terraform-provider-dnshelper/actions/workflows/test.yml/badge.svg)](https://github.com/marceloalmeida/terraform-provider-dnshelper/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/marceloalmeida/terraform-provider-dnshelper)](https://goreportcard.com/report/github.com/marceloalmeida/terraform-provider-dnshelper)

This Terraform provider helps manage DNS records.

## Documentation

Full documentation is available in the [docs](docs/) directory or on the [Terraform Registry](https://registry.terraform.io/providers/marceloalmeida/terraform-provider-dnshelper/latest/docs).


## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using `make build`

```sh
git clone https://github.com/marceloalmeida/terraform-provider-dnshelper
cd terraform-provider-dnshelper
make build
```

## Using the Provider

To use the provider, add the following to your Terraform configuration:

```hcl
terraform {
  required_providers {
    dnsrecords = {
      source = "registry.terraform.io/marceloalmeida/dnshelper"
    }
  }
}

provider "dnshelper" {}
```

See the [examples](examples/) directory for more detailed usage examples.

## Contributing

Contributions are welcome! Please read the contribution guidelines before submitting a pull request:

1. Fork the repository
2. Create a branch for your changes
3. Make your changes
4. Run tests with `make test`
5. Submit a pull request

## License

This provider is licensed under the [LICENSE](LICENSE) file.

## Development

If you wish to work on the provider, you'll need:

* [Go](https://www.golang.org) (version 1.24 or later)
* [Terraform](https://www.terraform.io/downloads.html) (version 1.8 or later)

To compile the provider:

```sh
make build
```

To run the tests:

```sh
make test
```

To run the full test suite with acceptance tests:

```sh
make testacc
```
