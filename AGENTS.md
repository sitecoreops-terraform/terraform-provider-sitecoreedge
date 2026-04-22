# Agent Guidelines for terraform-provider-sitecore-edge

This document provides guidelines for agentic coding tools working in this repository.

## Build/Lint/Test Commands

### Building

```bash
go build
```

### Linting

```bash
golangci-lint run
```

### Formatting

```bash
gofmt -s -w -e .
```

### Testing

#### Run all tests

```bash
go test ./...
```

#### Run tests for a specific package

```bash
go test ./pkg/apiclient/... -v
go test ./pkg/provider/... -v
```

#### Run a specific test

```bash
go test ./pkg/apiclient/... -v -run TestClientAuthentication
go test ./pkg/provider/... -v -run TestProviderMetadata
```

#### Integration tests

Integration tests require environment variables for authentication:

```bash
export SITECORE_EDGE_CLIENT_ID=your_client_id
export SITECORE_EDGE_CLIENT_SECRET=your_client_secret
go test ./pkg/apiclient/... -v
```

## Code Style Guidelines

### Imports

- Use standard Go import grouping (std, third-party, local)
- Avoid wildcard imports
- Import aliases should be avoided unless necessary for name conflicts

### Formatting

- Use `gofmt` for standard formatting
- Lines should not exceed 120 characters where reasonable
- Use spaces for indentation (Go standard)

### Types

- Use descriptive type names with camelCase
- Prefer pointers for optional or large struct fields
- Use type aliases sparingly

### Naming Conventions

- **Variables and Functions**: camelCase (e.g., `clientID`, `getWebhooks`)
- **Constants**: UPPER_SNAKE_CASE (e.g., `MAX_RETRIES`)
- **Interfaces**: Single method interfaces should end with "er" (e.g., `Reader`)
- **Structs**: PascalCase (e.g., `WebhookInput`)
- **Packages**: lowercase single words (e.g., `apiclient`)

### Error Handling

- Use descriptive error messages
- Wrap errors with context using `fmt.Errorf` with `%w` verb
- Return errors rather than panicking
- Use sentinel errors for expected error conditions

### Testing

- Test files should be named `<original>_test.go`
- Use table-driven tests for multiple test cases
- Test names should be descriptive and start with `Test`
- Use `t.Helper()` for test helper functions
- Mock external dependencies using test servers or interfaces

### Documentation

- All public functions and types should have godoc comments
- Comments should explain "why" not "what"
- Use complete sentences for documentation

### API Client Specific

- HTTP methods should be clearly documented
- Error responses should be properly handled
- Authentication should be handled securely
- Use proper content-type headers

### Terraform Provider Specific

- Resource names should follow Terraform conventions
- Schema definitions should be clear and descriptive
- Use proper attribute types from terraform-plugin-framework
- Handle state management carefully

### Concurrency

- Use channels for communication between goroutines
- Use `sync` package primitives appropriately
- Avoid global variables in concurrent code
- Use context for cancellation

### Logging

- Use structured logging where appropriate
- Avoid logging sensitive information
- Use appropriate log levels

## Project Structure

- `pkg/apiclient/`: API client code for Sitecore Experience Edge
- `pkg/provider/`: Terraform provider implementation
- `examples/`: Example Terraform configurations
- `docs/`: Generated documentation (do not edit manually)
- `tools/`: Build and documentation tools

## Common Patterns

### API Client Requests

```go
type RequestOptions struct {
    Method string
    Path   string
    Body   interface{}
}

func (c *Client) doRequest(opts RequestOptions) (*http.Response, error) {
    // Implementation
}
```

### Terraform Resource Structure

```go
type webhookResource struct {
    client *apiclient.Client
}

type webhookResourceModel struct {
    ID            types.String `tfsdk:"id"`
    Label         types.String `tfsdk:"label"`
    // Other fields
}

func (r *webhookResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    // Schema definition
}
```

## Tools

The project uses these tools (defined in go.mod):

- `tfproviderlint`: Terraform provider linter
- `tfplugindocs`: Documentation generator
- `terrafmt`: Terraform formatter
- `gotestsum`: Test runner with better output

## Environment Variables

- `SITECORE_EDGE_CLIENT_ID`: Client ID for Sitecore AI authentication
- `SITECORE_EDGE_CLIENT_SECRET`: Client secret for Sitecore AI authentication
- `SITECORE_EDGE_USE_CLI`: Set to "1" to use Sitecore CLI authentication

## Best Practices

1. Keep functions small and focused
2. Use interfaces for dependency injection
3. Write tests for all new functionality
4. Follow the existing code patterns
5. Document public APIs thoroughly
6. Handle errors gracefully
7. Keep the codebase consistent
8. Prefer composition over inheritance
9. Write clear, self-documenting code
10. Review changes before committing

## CI/CD

The project uses GitHub Actions for CI/CD. Check `.github/workflows/` for specific workflows.

## Security

- Never commit secrets or credentials
- Use environment variables for sensitive data
- Validate all inputs
- Use HTTPS for all communications
- Follow principle of least privilege
