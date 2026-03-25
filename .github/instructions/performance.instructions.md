---
applyTo: '**/*.go'
---

## Performance Requirements

- Database queries MUST be optimized for efficiency
  - Use appropriate indexes for frequently queried columns
  - Avoid N+1 query problems through eager loading or batch queries
  - Minimize full table scans and unnecessary joins
  - Use query result pagination for large datasets
- Avoid exponential or factorial time complexity in user-facing operations
- Choose data structures that match the access patterns of the operation
- Resource usage MUST be managed responsibly
  - Prevent memory leaks through proper resource cleanup
  - Use connection pooling for database and external service connections
  - Close file handles and network connections when no longer needed
  - Implement request timeouts to prevent resource exhaustion
- Caching strategies MUST be employed where appropriate
  - Cache expensive computations and frequently accessed data
  - Implement cache invalidation strategies to maintain data consistency
  - Use layered caching when beneficial
