# Disjoint Set Union (Union-Find)

## What is DSU?

DSU is a data structure that tracks elements partitioned into **non-overlapping groups**. Think of it like managing social networks: people belong to friend circles, and when two people become friends, their entire circles merge.

## Core Operations

1. **Find(x)**: Which group does element x belong to?
2. **Union(x, y)**: Merge the groups containing x and y

## How It Works

### Basic Concept
- Each element starts in its own group
- Each group has a **representative** (parent/root)
- Elements in the same group share the same root
- Use a parent array: `parent[i]` = parent of node i

### Naive Implementation
```
parent = [0, 1, 2, 3, 4]  // Initially, each node is its own parent

Find(x):
    while x != parent[x]:
        x = parent[x]
    return x

Union(x, y):
    rootX = Find(x)
    rootY = Find(y)
    parent[rootX] = rootY  // Attach one tree to another
```

**Problem**: Trees can become skewed (like a linked list) → O(n) operations

## Key Optimizations

### 1. Union by Rank (or Size)
Always attach the **smaller tree** under the **larger tree**.
- Keeps trees balanced
- Ensures tree height stays ≤ log(n)

**Two Common Approaches:**

#### Option 1: Union by Rank (Height-Based)
Track the **approximate height** of each tree.

```
rank[i] = approximate height of tree rooted at i

Union(x, y):
    rootX = Find(x)
    rootY = Find(y)
    if rank[rootX] < rank[rootY]:
        parent[rootX] = rootY
    else if rank[rootX] > rank[rootY]:
        parent[rootY] = rootX
    else:
        parent[rootY] = rootX
        rank[rootX]++  // Only increment when equal
```

**Key Point**: Rank becomes "approximate" with path compression, but still effective.

#### Option 2: Union by Size (Node Count-Based)
Track the **number of nodes** in each tree.

```
size[i] = number of nodes in tree rooted at i

Union(x, y):
    rootX = Find(x)
    rootY = Find(y)
    if size[rootX] < size[rootY]:
        parent[rootX] = rootY
        size[rootY] += size[rootX]
    else:
        parent[rootY] = rootX
        size[rootX] += size[rootY]
```

**Key Point**: Size remains accurate and can answer "how many nodes in this component?"

**Which to use?**
- **Union by Rank**: Simpler logic, rank stays small
- **Union by Size**: Accurate component sizes, useful if you need to query group sizes
- Both achieve the same O(log n) height bound

### 2. Path Compression
During `Find(x)`, make every node point **directly to the root**.
- Flattens the tree
- Future finds become O(1)

```
Find(x):
    if x != parent[x]:
        parent[x] = Find(parent[x])  // Recursively compress
    return parent[x]
```

### Combined Complexity
With both optimizations: **O(α(n))** per operation
- α(n) = inverse Ackermann function
- Practically **constant time** (α(n) ≤ 4 for all practical n)

## When to Use DSU

### Classic Use Cases

1. **Detect Cycles in Undirected Graphs**
   - If adding edge (u,v) and Find(u) == Find(v), it creates a cycle

2. **Count Connected Components**
   - Number of distinct roots = number of components

3. **Kruskal's MST Algorithm**
   - Add edges by weight, skip if it creates a cycle

4. **Dynamic Connectivity**
   - Check if nodes are connected as edges are added

5. **Network Connectivity**
   - Model computers, people, cities connecting over time

### Problem Patterns to Recognize

**Pattern**: "Group merging as events happen"
- Friends of friends
- Land masses (grid of islands merging)
- Equivalent equations (a=b, b=c → a=c)

**Pattern**: "Check if two elements are in the same group"
- Social network queries
- Reachability problems

**Pattern**: "Count number of groups/components"
- Number of provinces
- Number of disconnected networks

## Implementation Strategy

1. **Basic DSU**: parent array + find + union
2. **Add Union by Rank**: track rank/size
3. **Add Path Compression**: optimize find
4. **Extensions**:
   - Count components
   - Track group sizes
   - Custom merge logic

## Common Pitfalls

1. Forgetting to find the root before union
2. Not using both optimizations (performance degrades)
3. Using rank incorrectly (rank is NOT the exact height)
4. Off-by-one errors in initialization

## Mental Model

Think of DSU as a **forest of trees**:
- Each tree = one group
- Root = group representative
- Path compression = make trees flat
- Union by rank = keep trees balanced

## Time Complexity Summary

| Operation | Naive | With Optimizations |
|-----------|-------|-------------------|
| Find      | O(n)  | O(α(n)) ≈ O(1)   |
| Union     | O(n)  | O(α(n)) ≈ O(1)   |
| Space     | O(n)  | O(n)             |

## Quick Reference

**When you see keywords in problems:**
- "Connected components"
- "Merge groups"
- "Find if two elements are related"
- "Dynamic connectivity"
- "Union", "Find", "Disjoint sets"
- "Equivalent classes"
- "Cycle detection" (undirected graphs)

→ **Think DSU!**
