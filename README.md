# Structurizr - Terraform Provider
[![Build Status](https://github.com/fstaoe/terraform-provider-structurizr/actions/workflows/test.yml/badge.svg)](https://github.com/fstaoe/terraform-provider-structurizr/actions/workflows/test.yml)
[![Release Status](https://github.com/fstaoe/terraform-provider-structurizr/actions/workflows/release.yml/badge.svg)](https://github.com/fstaoe/terraform-provider-structurizr/actions/workflows/release.yml)

Terraform Provider for Structurizr [On-Premises](https://docs.structurizr.com/onpremises) and [Cloud Service](https://docs.structurizr.com/cloud). 

Example:

```hcl
terraform {
  required_providers {
    structurizr = {
      source = "fstaoe/structurizr"
      version = "0.1.0" # Check the version
    }
  }
}

variable "structurizr_admin_api_key {
  type = string
}

provider "structurizr" {
  host          = "http://localhost"
  admin_api_key = var.structurizr_admin_api_key
  tls_insecure  = true # Must be false when the host use HTTPS scheme 
}

# Create a Workspace for an Example
resource "structurizr_workspace" "example" {}
```

## Install

Download the latest [release](https://github.com/fstaoe/terraform-provider-structurizr/releases) and install into your 
Terraform plugin directory.

### Linux or Mac OSX

Run the following to have the provider installed for you automatically:

```sh
curl -fsSL https://raw.githubusercontent.com/fstaoe/terraform-provider-structurizr/main/scripts/install.sh | bash
```

### Windows

Download the plugin to `%APPDATA%\terraform.d\plugins`.

### Installation notes

The structurizr provider is published to the Terraform module registry and may be installed via the standard mechanisms. 
See the documentation at https://registry.terraform.io/providers/fstaoe/structurizr/latest.

## Usage

https://registry.terraform.io/providers/fstaoe/structurizr/latest

| Plugin                                   | Type     | Platform Support            | Description                                                    |
|------------------------------------------|----------|-----------------------------|----------------------------------------------------------------|
| [Structurizr](docs/index.md)             | Provider | on-premises + cloud service | Configures a target Structurizr server (such as a on-premises) |
| [Workspace](docs/resources/workspace.md) | Resource | on-premises + cloud service | Create workspaces                                              |

See our [Docs](./docs) folder for all plugins and our [Examples](./examples) to try out.

## Importing Resources

All resources support [importing](https://www.terraform.io/docs/import/usage.html).

## Developing

### Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.8+
- [Go](https://golang.org/doc/install) >= 1.22+

### Building

_Note:_ This project uses [Go Modules](https://blog.golang.org/using-go-modules) making it safe to work with it outside 
your existing [GOPATH](http://golang.org/doc/code.html#GOPATH). The instructions that follow assume a directory in your 
home directory outside the standard GOPATH (e.g. `$HOME/terraform-providers/`).

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `install` command.

```shell
go install
```

#### Adding Dependencies

To add a new dependency `github.com/author/dependency` to this project:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

### Testing 

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

### Documentation

To generate or update documentation:

```shell
go generate ./...
```

## Roadmap

Plan for the next few months:

- [x] Create Workspaces
- [x] Import Workspaces
- [ ] Push valid DSL/JSON files
- [ ] Lock/Unlock Workspaces
- [ ] Publish 1.0.0

Want to see more? See [CONTRIBUTING.md](CONTRIBUTING.md) for more details.
