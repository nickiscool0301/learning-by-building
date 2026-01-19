# Implementation Project: Distributed Word Counter

Build a simple distributed word counting system:

A coordinator node that accepts text files
Multiple worker nodes that count words in parallel
Handle worker failures gracefully (retry logic)
Implement basic RPC communication from scratch

# Goals:
├── Implement basic TCP-based RPC
├── Handle timeouts and retries
├── Detect worker failures via heartbeats
└── Aggregate results from multiple workers


# Success Criteria 
Your system should handle 1 worker crashing mid-task without losing data.