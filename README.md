# Sample Calculator Service

Minimal Go HTTP calculator service.

Run:

```bash
cd samplecalculatorproject
go run ./cmd/server
```

Run the MCP server over stdio:

```bash
cd samplecalculatorproject
go run ./mcp/cmd/server
```

Run the MCP server over Streamable HTTP:

```bash
cd samplecalculatorproject
go run ./mcp/cmd/httpserver
```

The MCP HTTP endpoint is `POST http://localhost:8081/mcp`.

Endpoints (all POST):

- `/capabilities` — returns available operations
- `/add` — body `{ "x": number, "y": number }`
- `/subtract` — body `{ "x": number, "y": number }`
- `/multiply` — body `{ "x": number, "y": number }`
- `/divide` — body `{ "x": number, "y": number }` (division by zero returns 400)
- `/trigonometry` — body `{ "theta": number }`, in degrees. Returns `sin`, `cos`, `tan`, `cosec`, `sec`, and `cot`; undefined reciprocal values are `null`.

Examples:

```bash
curl -X POST -H "Content-Type: application/json" -d '{"x":3.5,"y":2}' http://localhost:8080/add
curl -X POST -H "Content-Type: application/json" -d '{"x":10,"y":0}' http://localhost:8080/divide
curl -X POST -H "Content-Type: application/json" -d '{"theta":45}' http://localhost:8080/trigonometry
curl -X POST -H "Content-Type: application/json" -d '{}' http://localhost:8080/capabilities
```
