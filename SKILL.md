---
name: sample-calculator
description: Minimal Go calculator service exposing arithmetic operations through both HTTP JSON endpoints and MCP tools. It supports add, subtract, multiply, divide, and operation discovery for API clients and AI agents.
version: 1.0.0
repo-url: TODO: repository Git URL not found in workspace
allowed-tools:
  - add
  - subtract
  - multiply
  - divide
  - capabilities
  - HTTP POST /add
  - HTTP POST /subtract
  - HTTP POST /multiply
  - HTTP POST /divide
  - HTTP POST /capabilities
  - MCP stdio server: go run ./mcp/cmd/server
  - MCP Streamable HTTP server: go run ./mcp/cmd/httpserver
---

# 1. Current Capabilities & Supported Interfaces

## Public Endpoints

### HTTP JSON API

Base runtime command:

```bash
go run ./cmd/server
```

Default address:

```text
http://localhost:8080
```

Exposed routes:

| Route | Verb | Request Body | Success Response | Error Responses |
|---|---|---|---|---|
| `/capabilities` | `POST` | none required; `{}` accepted | `{ "operations": [{ "id": string, "description": string }] }` | `405 { "error": "method not allowed" }` |
| `/add` | `POST` | `{ "x": number, "y": number }` | `{ "result": number }` | `400 { "error": "invalid JSON body" }`, `405 { "error": "method not allowed" }` |
| `/subtract` | `POST` | `{ "x": number, "y": number }` | `{ "result": number }` | `400 { "error": "invalid JSON body" }`, `405 { "error": "method not allowed" }` |
| `/multiply` | `POST` | `{ "x": number, "y": number }` | `{ "result": number }` | `400 { "error": "invalid JSON body" }`, `405 { "error": "method not allowed" }` |
| `/divide` | `POST` | `{ "x": number, "y": number }` | `{ "result": number }` | `400 { "error": "invalid JSON body" }`, `400 { "error": "division by zero" }`, `405 { "error": "method not allowed" }` |

### MCP Stdio Interface

Runtime command:

```bash
go run ./mcp/cmd/server
```

Tool signatures:

```text
add(x: number, y: number) -> { result: number }
subtract(x: number, y: number) -> { result: number }
multiply(x: number, y: number) -> { result: number }
divide(x: number, y: number) -> { result: number } | tool error "division by zero"
capabilities() -> { operations: [{ id: string, description: string }] }
```

### MCP Streamable HTTP Interface

Runtime command:

```bash
go run ./mcp/cmd/httpserver
```

Default endpoint:

```text
POST http://localhost:8081/mcp
```

The port can be overridden with:

```bash
MCP_HTTP_ADDR=:9090 go run ./mcp/cmd/httpserver
```

Required request headers for Postman or raw HTTP clients:

```text
Content-Type: application/json
Accept: application/json, text/event-stream
```

Example MCP call:

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "add",
    "arguments": {
      "x": 3.5,
      "y": 2
    }
  }
}
```

Expected MCP tool result includes:

```json
{
  "structuredContent": {
    "result": 5.5
  }
}
```

## Schema & Data Types

### `operands`

Used by HTTP arithmetic endpoints and MCP arithmetic tools.

| Field | Type | Required | Description |
|---|---|---|---|
| `x` | `number` / Go `float64` | yes | First operand |
| `y` | `number` / Go `float64` | yes | Second operand |

No enums are currently accepted.

HTTP JSON decoding uses `DisallowUnknownFields`, so unknown fields in request bodies must be rejected with:

```json
{ "error": "invalid JSON body" }
```

### `resultResp`

```json
{
  "result": 5.5
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| `result` | `number` / Go `float64` | yes | Arithmetic result |

### `errorResp`

```json
{
  "error": "division by zero"
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| `error` | `string` | yes | Human-readable error message |

### `capabilities` response

```json
{
  "operations": [
    { "id": "add", "description": "Add two numbers" },
    { "id": "subtract", "description": "Subtract second number from first" },
    { "id": "multiply", "description": "Multiply two numbers" },
    { "id": "divide", "description": "Divide first number by second" }
  ]
}
```

## Existing Operations

The service can currently perform exactly these operations:

- Add two numbers.
- Subtract the second number from the first.
- Multiply two numbers.
- Divide the first number by the second.
- Reject division by zero.
- List supported calculator operations.
- Serve calculator operations through HTTP JSON routes.
- Serve calculator operations as MCP tools over stdio.
- Serve calculator operations as MCP tools over Streamable HTTP.

# 2. Capability Gap Check (For Downstream Agents)

Before proposing a new endpoint, MCP tool, or arithmetic operation, an external agent must check whether the requested capability already exists or can be composed from current primitives.

Required discovery steps:

1. Inspect `samplecalculatorproject/README.md` for documented runtime commands and public interfaces.
2. Inspect `samplecalculatorproject/cmd/server/main.go` to confirm registered HTTP routes.
3. Inspect `samplecalculatorproject/internal/api/handler.go` to confirm HTTP request/response behavior and validation rules.
4. Inspect `samplecalculatorproject/docs/swagger.yaml` and `samplecalculatorproject/docs/swagger.json` to confirm the published OpenAPI contract.
5. Inspect `samplecalculatorproject/mcp/cmd/server/main.go`, `samplecalculatorproject/mcp/cmd/httpserver/main.go`, and `samplecalculatorproject/mcp/server/main.go` to confirm MCP tool availability.
6. Inspect `samplecalculatorproject/mcp/server/main_test.go` for tested calculator behavior and existing edge cases.

Composition guidance:

- Do not add a new `sum` operation if `add` already satisfies the request.
- Do not add a new `difference` operation if `subtract` already satisfies the request.
- Do not add a new `quotient` operation if `divide` already satisfies the request.
- Do not add an endpoint solely to expose operation metadata; use `/capabilities` or the MCP `capabilities` tool.
- Only add a new operation when it cannot be represented by the existing arithmetic primitives without changing client-side semantics.

# 3. Invariant Architecture Boundaries (Rules for Contributing AIs)

## Interface Stability

Contributing agents must preserve all existing public contracts unless the user explicitly requests a breaking change.

Rules:

- Do not remove or rename existing HTTP routes:
  - `POST /capabilities`
  - `POST /add`
  - `POST /subtract`
  - `POST /multiply`
  - `POST /divide`
- Do not change existing HTTP verbs from `POST`.
- Do not change successful arithmetic response shape from `{ "result": number }`.
- Do not change error response shape from `{ "error": string }`.
- Do not remove existing MCP tools:
  - `add`
  - `subtract`
  - `multiply`
  - `divide`
  - `capabilities`
- Do not change MCP arithmetic argument names `x` and `y`.
- Do not change division-by-zero behavior without explicit approval.
- New fields may be added only in a backward-compatible way.
- New endpoints or MCP tools must be additive and documented.
- If route behavior changes, OpenAPI/Swagger files and this `SKILL.md` version must be updated in the same change.

## Code Placement

Use the existing project layout.

| Concern | Required Location |
|---|---|
| HTTP server bootstrap and route registration | `samplecalculatorproject/cmd/server/main.go` |
| HTTP handlers, request structs, response structs, HTTP validation | `samplecalculatorproject/internal/api/handler.go` |
| Generated Swagger Go docs | `samplecalculatorproject/docs/docs.go` |
| Swagger JSON spec | `samplecalculatorproject/docs/swagger.json` |
| Swagger YAML spec | `samplecalculatorproject/docs/swagger.yaml` |
| MCP stdio command | `samplecalculatorproject/mcp/cmd/server/main.go` |
| MCP Streamable HTTP command | `samplecalculatorproject/mcp/cmd/httpserver/main.go` |
| Existing MCP calculator server package and tests | `samplecalculatorproject/mcp/server/` |
| Project README | `samplecalculatorproject/README.md` |
| Capability contract | `samplecalculatorproject/SKILL.md` |

Business logic for arithmetic must remain simple, deterministic, and side-effect free. New arithmetic operations should be implemented consistently across HTTP and MCP surfaces.

## File Path Allowlist

Automated PRs for normal calculator feature work may touch only these paths:

```text
samplecalculatorproject/SKILL.md
samplecalculatorproject/README.md
samplecalculatorproject/internal/api/
samplecalculatorproject/cmd/server/
samplecalculatorproject/mcp/server/
samplecalculatorproject/mcp/cmd/server/
samplecalculatorproject/mcp/cmd/httpserver/
samplecalculatorproject/docs/swagger.yaml
samplecalculatorproject/docs/swagger.json
samplecalculatorproject/docs/docs.go
samplecalculatorproject/go.mod
samplecalculatorproject/go.sum
```

`go.mod` and `go.sum` may be changed only when a new dependency is strictly required and justified.

## Strictly Protected Paths

Automated PRs must not modify these paths unless explicitly requested by a human maintainer:

```text
.git/
.agents/
.codex/
samplecalculatorproject/.git/
samplecalculatorproject/go.mod
samplecalculatorproject/go.sum
```

Root build, dependency, or repository-wide configuration files are protected by default. Do not introduce CI workflows, Dockerfiles, deployment manifests, auth middleware, persistence layers, network calls, telemetry, or database configuration unless explicitly requested.

# 4. PR Acceptance & Quality Checklist

## Error & Edge-Case Handling

Every change must preserve existing edge-case behavior and add validation for new cases.

Required checks for arithmetic operations:

- Invalid JSON body returns HTTP `400` with `{ "error": "invalid JSON body" }`.
- Unknown JSON fields are rejected for operand-based HTTP requests.
- Non-`POST` HTTP methods return HTTP `405` with `{ "error": "method not allowed" }`.
- Division by zero returns HTTP `400` with `{ "error": "division by zero" }`.
- MCP division by zero returns a tool error result with message `division by zero`.
- Positive numbers, negative numbers, zero, and decimal values must be handled.
- Boundary cases should include very large `float64` values where relevant.
- New operations must define behavior for invalid operands, null-like missing fields, zero values, and floating-point overflow or undefined mathematical results.
- The service must remain stateless.

## Test Requirements

Minimum required tests for any behavior change:

- Run the full test suite:

```bash
go test ./...
```

- Add or update unit tests covering:
  - Positive cases for every new operation.
  - Negative/error cases for invalid operations or invalid operands.
  - Boundary cases, including zero and decimal inputs.
  - Division-by-zero behavior if divide logic is touched.
  - Capability listing updates when operations are added or removed.
- If HTTP handler behavior changes, add handler-level tests for status codes and response bodies.
- If MCP behavior changes, add MCP-level tests for tool result shape and tool errors.
- Existing tests in `samplecalculatorproject/mcp/server/main_test.go` must continue to pass.

## Documentation & Spec Sync

Every accepted PR that changes public behavior must update documentation and specs in the same PR.

Required updates:

- Update `samplecalculatorproject/README.md` with new commands, endpoints, or examples.
- Update Swagger annotations in `samplecalculatorproject/internal/api/handler.go` when HTTP behavior changes.
- Regenerate or manually synchronize:
  - `samplecalculatorproject/docs/swagger.yaml`
  - `samplecalculatorproject/docs/swagger.json`
  - `samplecalculatorproject/docs/docs.go`
- Update MCP tool registration in all relevant MCP entrypoints when operations change:
  - `samplecalculatorproject/mcp/cmd/server/main.go`
  - `samplecalculatorproject/mcp/cmd/httpserver/main.go`
  - `samplecalculatorproject/mcp/server/main.go`
- Update this `SKILL.md` file.
- Increment `version` in this `SKILL.md` for public contract changes.
- Confirm `allowed-tools` remains complete and accurate.
