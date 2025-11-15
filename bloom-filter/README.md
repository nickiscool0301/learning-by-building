# Bloom Filter
A Bloom filter is a space-efficient probabilistic data structure that tests whether an element is a member of a set. 


## Key characteristics:
- **False positives** are *possible* (might say "yes" when it should say "no")
- **False negatives** are *NOT possible* (never says "no" when it should say "yes")
- Very **memory efficient** compared to hash sets

## Use cases: 
Spell checkers, database query optimization,web crawlers (avoid re-crawling URLs), etc.

## Overview
Here's what we'll build:
BloomFilter
├── Bit array (underlying storage)
├── Hash functions (multiple independent hashes)
├── Add(item) - insert an element
├── Contains(item) - check membership
└── Size parameters (m bits, k hash functions)