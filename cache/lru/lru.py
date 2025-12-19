from dataclasses import dataclass

@dataclass
class Node:
    key: int
    value: int 
    next: Optional['Node'] = None 
    prev: Optional['Node'] = None 

class LRU:
    def __init__(self, capacity: int):
        if capacity <= 0: 
            raise ValueError("Capacity must be positive")
        self.capacity = capacity
        self.cache = {}
        self.head = Node(-1,-1) 
        self.tail = Node(-1, -1)
        self.head.next = self.tail
        self.tail.prev = self.head 
    
    def put(self, key: int, value: int) -> None:
        if key in self.cache:
            node = self.cache[key]
            node.value = value
            self._unlink(node)
            self._add_front(node)
            return
        if len(self.cache) == self.capacity:
            self.evict()
        new_node = Node(key, value)
        self.cache[key] = new_node
        self._add_front(new_node)

    
    def get(self, key: int) -> int: 
        if key not in self.cache:
            return -1
        node = self.cache[key]
        self._unlink(node) 
        self._add_front(node) 
        return node.value

    def _unlink(self, node: Node) -> None:
        node.next.prev = node.prev
        node.prev.next = node.next 

    def _remove(self, node: Node) -> None:
        self._unlink(node)
        del self.cache[node.key] 

    def _add_front(self, node: Node) -> None:
        next_node = self.head.next 
        self.head.next = node
        node.next = next_node 
        next_node.prev = node 
        node.prev = self.head 

    def evict(self) -> None:
        if len(self.cache) == 0:
            return
        prev = self.tail.prev
        self._remove(prev)
