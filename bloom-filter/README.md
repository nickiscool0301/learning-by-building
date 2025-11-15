# Bloom Filter
A Bloom filter is a space-efficient probabilistic data structure that tests whether an element is a member of a set. 


## Key characteristics
- **False positives** are *possible* (might say "yes" when it should say "no")
- **False negatives** are *NOT possible* (never says "no" when it should say "yes")
- Very **memory efficient** compared to hash sets

## Use cases
Spell checkers, database query optimization,web crawlers (avoid re-crawling URLs), etc.

## Overview
Here's what we'll build:
BloomFilter
├── Bit array (underlying storage)
├── Hash functions (multiple independent hashes)
├── Add(item) - insert an element
├── Contains(item) - check membership
└── Size parameters (m bits, k hash functions)

## Complexity
- **Space complexity**: O(m) bits
- **Time complexity**:
  - Add: O(k)
  - Contains: O(k)
Where m = size of bit array, k = number of hash functions

## Comparison to HashSet
### Space
- HashSet: Store 1 million URLs (~100 bytes each) -> 1,000,000 * 100 = 100,000,000 bytes = ~100 MB
- Bloom Filter: For 1 million URLS with 1% false positive rate:
  - m = - (n * ln(p)) / (ln(2)^2) = ~9,585,058 bits = ~1.15 MB
  - k = (m/n) * ln(2) = ~7 hash functions

### False Positives
- HashSet: 0% false positives
- Bloom Filter: 1% false positive rate (may say "yes" for 1 out of 100 non-members)

### Delete operation
Standard Bloom filters do not support deletion. To allow deletions, a Counting Bloom Filter can be used, which uses an array of counters instead of bits.

### Speed
- HashSet: O(1) average time for add and contains
- Bloom Filter: O(k) time for add and contains (k = number of hash functions)