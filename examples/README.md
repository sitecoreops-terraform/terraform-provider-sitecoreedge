# Sitecore AI Terraform Provider Examples

This directory contains comprehensive examples demonstrating how to use the Sitecore AI Terraform provider.

The examples are mostly to be included in documentation, but there are also a few full examples on their own.

## Documentation generation

The document generation tool looks for files in the following locations by default. All other *.tf files besides the ones mentioned below are ignored by the documentation tool. This is useful for creating examples that can run and/or are testable even if some parts are not relevant for the documentation.

| Path                                                                         | Description                                |
|------------------------------------------------------------------------------|--------------------------------------------|
| `examples/provider/provider<*>.tf`                                           | Provider example config(s)                 |
| `examples/data-sources/<data source name>/data-source<*>.tf`                 | Data source example config(s)              |
| `examples/resources/<resource name>/resource<*>.tf`                          | Resource example config(s)                 |
| `examples/resources/<resource name>/import.sh`                               | Resource example import command            |
| `examples/resources/<resource name>/import-by-string-id.tf`                  | Resource example import by id config       |
| `examples/resources/<resource name>/import-by-identity.tf`                   | Resource example import by identity config |

## Prerequisites

- Terraform 1.0+
- SitecoreAI API credentials (create organization credentials in SitecoreAI Deploy)

## Working with local provider implementation

Usually terraform providers are found through the registry, however you can specify a local override where a certain provider is found. This can be specified in a `.terraformrc` file in the user's home directory, or the environment variable `TF_CLI_CONFIG_FILE` can point to a `*.tfrc` file. We have created a `localdev.tfrc` and hereby you can run

```bash
go build
cd examples
export TF_CLI_CONFIG_FILE=$(pwd)/localdev.tfrc
```

However, when there is only the local provider, there is no need to run `terraform init`
