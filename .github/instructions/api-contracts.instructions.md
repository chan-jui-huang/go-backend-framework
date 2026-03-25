---
applyTo: 'internal/http/**/*.go,docs/**/*.go,docs/**/*.json,docs/**/*.yaml'
---

## API Contracts

### Contract Consistency

- Swagger annotations MUST match the actual runtime handler behavior
- Request and response schemas MUST stay aligned with the payloads returned by the application
- Document only the status codes and error responses that the handler can actually return
- If a handler logs an error instead of returning it, remove that error response from Swagger
- Contract changes MUST regenerate and update the `docs/` artifacts

### Route and Annotation Rules

- API route paths MUST use kebab case
- Swagger `// @param` descriptions MUST keep the placeholder value as `" "`
- Keep request parameter definitions, response models, and documented media types consistent with the implemented route
