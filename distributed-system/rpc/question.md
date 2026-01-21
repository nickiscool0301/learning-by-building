# Build your own RPC library

Define a simple IDL for service definitions
- Implement request/response serialization (use JSON first, then binary)
- Add at-least-once semantics with retries
- Add at-most-once semantics with request deduplication (using unique request IDs)
- Implement timeout handling

# Goals
├── Support multiple service methods
├── Handle concurrent requests
├── Implement retry with exponential backoff
├── Add request deduplication table
└── Support both sync and async calls

# Success Criteria
Demonstrate that duplicate requests don't cause duplicate side effects.