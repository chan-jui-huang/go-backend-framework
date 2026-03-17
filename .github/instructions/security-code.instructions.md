---
applyTo: '**/*.go,**/go.mod,**/go.sum'
---

## Architecture Constraints

## Security Rules

- Security-sensitive changes MUST account for the project's applicable OWASP Top 10 risks
- Avoid implementations that weaken access control, secret handling, cryptography, input validation, or error handling
- All code MUST pass static analysis security checks
- ALL user input MUST be validated before processing
- Passwords MUST use Argon2id hashing (never bcrypt or plain SHA)
- JWT tokens MUST use Ed25519 signatures
- Error messages MUST NOT leak sensitive information (stack traces, SQL queries)
- Database credentials MUST come from environment variables, never hardcoded
