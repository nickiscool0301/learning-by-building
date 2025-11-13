# Consistent Hashing Implementation in Java

## What You've Built

A complete consistent hashing system with:
- **Hash Ring**: Uses TreeMap for O(log n) lookups
- **Virtual Nodes**: Improves load distribution
- **MD5 Hashing**: Uniform distribution of keys
- **Add/Remove Nodes**: Dynamic scaling with minimal key redistribution

## Key Learnings from the Demo

### 1. **Minimal Redistribution** (Demo 2)
- When adding a 4th server, only **~25%** of keys moved
- Traditional hashing would move **~75%** of keys
- This is the PRIMARY benefit of consistent hashing!

### 2. **Virtual Nodes Impact** (Demo 4)
| Virtual Nodes | Standard Deviation | Max Deviation |
|---------------|-------------------|---------------|
| 1             | 1756.84           | 92.72%        |
| 10            | 473.12            | 29.32%        |
| 50            | 258.27            | 13.92%        |
| 100           | 98.02             | 6.32%         |
| 200           | 139.68            | 8.88%         |

**Sweet spot**: 100-200 virtual nodes balances distribution vs memory

### 3. **Load Distribution** (Demo 3)
With 150 virtual nodes, keys distributed within **±8%** of ideal (25% per server)

## How It Works

```
Hash Ring (0 to 2^64-1):

     Server-A#0 ───┐
                   │
     Server-B#0 ───┤
                   ├─── Circular Hash Space
     Server-A#1 ───┤
                   │
     Server-C#0 ───┘

Key Lookup:
1. Hash the key → position on ring
2. Find first server clockwise
3. That server handles the key
```

## Real-World Applications

- **Load Balancers**: HAProxy, Nginx
- **Distributed Caches**: Redis Cluster, Memcached
- **Databases**: Cassandra, DynamoDB, Riak
- **CDNs**: Distribute content across edge servers
- **Sharding**: Partition data across database nodes

## Try These Experiments

1. **Modify virtual nodes**: Change the value in demo constructors
2. **Add more servers**: See how distribution changes
3. **Remove servers**: Observe minimal key movement
4. **Different hash functions**: Try SHA-256 instead of MD5
5. **Custom node types**: Use `ConsistentHash<Server>` with a Server class

## Running the Code

```bash
cd consistent-hash
javac ConsistentHash.java ConsistentHashDemo.java
java ConsistentHashDemo
```

## Next Steps to Explore

- Add weighted nodes (some servers handle more load)
- Implement node failure detection
- Add monitoring/metrics
- Build a distributed cache using this
- Test with network partitions
