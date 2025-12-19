# Re-implement hash map
'''
Version 1:
- HashMap: array of buckets
- Bucket: array of key-value pairs
- Hash function: maps key to bucket index
    - How: sum ASCII values of characters in key, mod by number of buckets
- Not good in term of TC, not handling dynamic resizing + hash collision
'''


class Entry:
    def __init__(self, key: str, value: int):
        self.key = key
        self.value = value

from typing import List
class Bucket:
    def __init__(self):
        self.items = [] # List<Entry>
    
class HashMap:
    def __init__(self, size=256):
        self.size = size
        self.buckets = [Bucket() for _ in range(size)]

    def hash(self, key: str) -> int:
        # convert key to bucket index
        return int(sum(ord(c) for c in key) % self.size)

    def put(self, key: str, value: int) -> None:
        idx = self.hash(key)
        bucket = self.buckets[idx]
        for entry in bucket.items:
            if entry.key == key:
                entry.value = value
                return
        bucket.items.append(Entry(key, value))
        
    def get(self, key: str) -> int:
        idx = self.hash(key)
        bucket = self.buckets[idx]
        for entry in bucket.items:
            if entry.key == key:
                return entry.value
        raise KeyError(f'Key {key} not found')        

    def len(self) -> int:
        return sum(len(bucket.items) for bucket in self.buckets)

    def remove(self, key: str) -> None:
        idx = self.hash(key)
        bucket = self.buckets[idx]
        for i, entry in enumerate(bucket.items):
            if entry.key == key:
                del bucket.items[i]
                return
        raise KeyError(f'Key {key} not found')
    

if __name__ == '__main__':
    hashmap = HashMap()
    hashmap.put('apple', 5)
    hashmap.put('banana', 10)
    print(hashmap.get('apple'))  # Output: 5
    print(hashmap.get('banana')) # Output: 10
    hashmap.remove('apple')
    try:
        print(hashmap.get('apple'))  # Should raise KeyError
    except KeyError as e:
        print(e)
