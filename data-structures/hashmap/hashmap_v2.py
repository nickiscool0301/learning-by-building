from dataclasses import dataclass

@dataclass
class Entry:
    key: str
    value: int

class HashMap:
    def __init__(self, capacity=256):
        self.capacity = capacity
        self.size = 0
        self.load_factor = 0.75
        self.buckets = [[] for _ in range(capacity)]  # List[List[Entry]]

    def _hash(self, key: str) -> int:
        # use python built-in hash
        # return hash(key) % self.capacity

        # polynominal rolling hash
        hash_value = 0
        p = 31
        MOD = 10**9 + 9
        for i, char in enumerate(key):
            hash_value = (hash_value + (ord(char) - ord('a')) * pow(p, i, MOD)) % MOD
        return hash_value % self.capacity

    def put(self, key: str, value: int) -> None:
        if self.size / self.capacity >= self.load_factor:
            self._resize()
        idx = self._hash(key)
        bucket = self.buckets[idx]
        for entry in bucket:
            if entry.key == key:
                entry.value = value
                return
        bucket.append(Entry(key, value))
        self.size += 1

    def _resize(self) -> None:
        # Doubles capacity and rehashes all exisiting entries
        old_buckets = self.buckets
        self.capacity *= 2
        self.buckets = [[] for _ in range(self.capacity)]
        self.size = 0

        #rehash
        for bucket in old_buckets:
            for entry in bucket:
                self.put(entry.key, entry.value)
        

    def get(self, key: str) -> int:
        idx = self._hash(key)
        bucket = self.buckets[idx]
        for entry in bucket:
            if entry.key == key:
                return entry.value
        return -1

    def len(self) -> int:
        return self.size

    def remove(self, key: str) -> None:
        idx = self._hash(key)
        bucket = self.buckets[idx]
        for i, entry in enumerate(bucket):
            if entry.key == key:
                del bucket[i]
                self.size -= 1
                return

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