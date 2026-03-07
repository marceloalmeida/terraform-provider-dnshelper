# CLAUDE.md

## Project Overview

`terraform-provider-dnshelper` is a Terraform provider that exposes DNS record builder functions (SPF, CAA, DMARC) as Terraform provider functions. It uses the HashiCorp Terraform Plugin Framework (protocol v6) and is authored by Marcelo Almeida under the MPL-2.0 license.

The provider has no resources or data sources — it exclusively provides **provider functions** callable via `provider::dnshelper::<function_name>()` in Terraform configs.

## Repository Structure

```
├── main.go                          # Provider entrypoint
├── internal/
│   ├── provider/
│   │   ├── provider.go              # DnshelperProvider definition, registers functions
│   │   ├── factory.go               # Test provider factories (ProtoV6)
│   │   └── provider_test.go         # Provider-level unit tests
│   ├── function/
│   │   ├── spfbuilder_function.go   # SPF builder Terraform function
│   │   ├── caabuilder_function.go   # CAA builder Terraform function
│   │   ├── dmarcbuilder_function.go # DMARC builder Terraform function
│   │   └── *_test.go               # Function tests (unit + acceptance)
│   └── testutil/
│       └── mock_resolver.go         # Mock DNS resolver for SPF tests
├── dnshelper/
│   ├── spfbuilder/                  # Core SPF record building logic (uses dnscontrol/spflib)
│   ├── caabuilder/                  # Core CAA record building logic
│   └── dmarcbuilder/                # Core DMARC record building logic
├── examples/                        # Terraform example configs (used for doc generation)
├── docs/                            # Generated documentation (do not edit manually)
├── tools/tools.go                   # Build tools: copywrite headers, tfplugindocs, terraform fmt
├── .goreleaser.yml                  # GoReleaser config for releases
├── .golangci.yml                    # Linter configuration
└── .github/workflows/               # CI: test.yml, release.yml, pr-lint.yml, etc.
```

## Architecture

**Two-layer design:**
1. **Core logic** (`dnshelper/` packages) — Pure Go libraries that build DNS records. No Terraform dependencies. Each has its own `_test.go` with unit tests.
2. **Terraform function wrappers** (`internal/function/`) — Adapt core logic into `function.Function` implementations. Handle argument parsing via `tfsdk` struct tags and return Terraform-typed results.

**Key patterns:**
- Functions implement `function.Function` interface: `Metadata`, `Definition`, `Run`
- SPF builder uses a `Resolver` interface for DNS lookups, with a mock resolver (`testutil.MockResolver`) for unit tests and `spflib.LiveResolver` for real DNS
- The test/mock switching is done via `testing.Testing()` check in `buildSPFRecord()`
- Provider registers all functions in `provider.go:Functions()`

## Commands

### Build and Install
```bash
make build     # go build -v ./...
make install   # go build + go install
make           # fmt + lint + install + generate (default target)
```

### Testing
```bash
make test      # Unit tests: go test -v -cover -timeout=120s -parallel=10 ./...
make testacc   # Acceptance tests: TF_ACC=1 go test -v -cover -timeout 120m ./...
```

Unit tests run without DNS or Terraform. Acceptance tests (`TestAcc*` prefix) require `TF_ACC=1` and a Terraform binary (>= 1.8.2).

### Linting
```bash
make lint      # golangci-lint run
make fmt       # gofmt -s -w -e .
```

Enabled linters (see `.golangci.yml`): errcheck, govet, staticcheck, unused, misspell, ineffassign, unconvert, unparam, forcetypeassert, nilerr, predeclared, makezero, durationcheck, copyloopvar, usetesting.

### Documentation Generation
```bash
make generate  # Runs copywrite headers, terraform fmt on examples, tfplugindocs
```

The `docs/` directory is auto-generated from code and `examples/`. CI checks that `make generate` produces no diff.

## Testing Conventions

- Unit tests use `testing` + `github.com/stretchr/testify/require`
- Terraform acceptance tests use `github.com/hashicorp/terraform-plugin-testing/helper/resource`
- Test provider factories are in `internal/provider/factory.go`
- SPF tests use a mock DNS resolver backed by `internal/testutil/testdata-dns.json`
- Test function names follow: `Test<FunctionName>_<Method>` for units, `TestAcc<FunctionName>_tf` for acceptance
- Use `t.Parallel()` where possible

## Code Conventions

- All files must have the copyright header: `// Copyright (c) Marcelo Almeida` + `// SPDX-License-Identifier: MPL-2.0`
- Go module: `github.com/marceloalmeida/terraform-provider-dnshelper`
- Go version: 1.25 (set in `go.mod`)
- Provider address: `registry.terraform.io/marceloalmeida/dnshelper`
- Provider type name: `dnshelper`
- Internal packages use the `internal/` directory (not importable externally)
- Core DNS logic is separated from Terraform wrappers in `dnshelper/` packages

## CI/CD

- **test.yml**: Runs on push/PR. Build + lint, generate check, acceptance tests against Terraform 1.8.*
- **release.yml**: Triggered by `v*` tags. Uses GoReleaser with GPG signing
- **pr-lint.yml**: PR title/label validation
- **lock.yml**: Auto-locks stale issues/PRs

## Adding a New Function

1. Create core logic in `dnshelper/<name>builder/` with tests
2. Create Terraform function wrapper in `internal/function/<name>builder_function.go` implementing `function.Function`
3. Add corresponding tests in `internal/function/<name>builder_function_test.go`
4. Register the function in `internal/provider/provider.go` in the `Functions()` method
5. Add example Terraform config in `examples/functions/<name>_builder/`
6. Run `make generate` to regenerate docs
7. Run `make test` and `make testacc` to verify
