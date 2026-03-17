---
applyTo: '**/*.go,**/go.mod,**/go.sum'
---

## Architecture Constraints

### Technical Requirements

#### Security Requirements

Security is a fundamental architectural concern. All implementations MUST address the OWASP Top 10:2025 application security risks appropriate to their context and technology stack.

Projects MUST address the OWASP Top 10:2025 application security risks. Detailed prevention strategies appropriate to the project's technology stack and threat model are available in the official OWASP documentation.

**Source**: OWASP Top 10:2025 RC1, released November 6, 2025  
**Reference**: https://owasp.org/Top10/2025/0x00_2025-Introduction/

**The 10 Critical Security Risks:**

1. **A01:2025 - Broken Access Control** - Failures in enforcing user permission policies, leading to unauthorized access to data or functions. Most critical risk affecting 3.73% of tested applications.

2. **A02:2025 - Security Misconfiguration** - Insecure default configurations, incomplete setups, open cloud storage, verbose error messages, and missing security headers. Moved up from #5 in 2021.

3. **A03:2025 - Software Supply Chain Failures** - Compromises in dependencies, build systems, and distribution infrastructure. Expanded from 2021's "Vulnerable Components" to address broader ecosystem risks.

4. **A04:2025 - Cryptographic Failures** - Inadequate protection of sensitive data through weak or missing encryption, leading to data exposure. Previously ranked #2 in 2021.

5. **A05:2025 - Injection** - Untrusted data sent to interpreters as commands or queries, ranging from SQL injection to cross-site scripting (XSS).

6. **A06:2025 - Insecure Design** - Missing or ineffective security controls arising from failure to use threat modeling and secure design patterns during architecture phase.

7. **A07:2025 - Authentication Failures** - Broken authentication and session management allowing attackers to compromise passwords, keys, or session tokens.

8. **A08:2025 - Software or Data Integrity Failures** - Failure to verify integrity of software updates, critical data, and CI/CD pipelines, enabling unauthorized code or data modifications.

9. **A09:2025 - Logging & Alerting Failures** - Insufficient logging and monitoring combined with missing or ineffective alerting, preventing timely detection and response to breaches.

10. **A10:2025 - Mishandling of Exceptional Conditions** - Improper error handling, failing insecurely under abnormal conditions, and information leakage through error messages. New category for 2025.

**Implementation Guidance:**  
Projects MUST assess which OWASP categories apply to their specific context and implement appropriate preventive controls. Detailed prevention strategies, common vulnerability patterns, and testing guidance are available in the official OWASP documentation for each category.

#### Code Quality

- All code MUST pass static analysis security checks
- All code MUST follow language-standard formatting conventions
- All production builds MUST be optimized for performance

### Security Requirements

- ALL user input MUST be validated before processing
- Passwords MUST use Argon2id hashing (never bcrypt or plain SHA)
- JWT tokens MUST use Ed25519 signatures
- Error messages MUST NOT leak sensitive information (stack traces, SQL queries)
- Database credentials MUST come from environment variables, never hardcoded

#### Performance

- Database queries MUST be optimized for efficiency
  - Use appropriate indexes for frequently queried columns
  - Avoid N+1 query problems through eager loading or batch queries
  - Minimize full table scans and unnecessary joins
  - Use query result pagination for large datasets
- Algorithm complexity MUST be considered and documented
  - Avoid algorithms with exponential or factorial time complexity for user-facing operations
  - Document time and space complexity for critical operations
  - Use appropriate data structures for the access patterns
- Resource usage MUST be managed responsibly
  - Prevent memory leaks through proper resource cleanup
  - Use connection pooling for database and external service connections
  - Close file handles and network connections when no longer needed
  - Implement request timeouts to prevent resource exhaustion
- Caching strategies MUST be employed where appropriate
  - Cache expensive computations and frequently accessed data
  - Implement cache invalidation strategies to maintain data consistency
  - Use layered caching (application, database query cache) when beneficial
  - Document cache lifetimes and invalidation triggers
