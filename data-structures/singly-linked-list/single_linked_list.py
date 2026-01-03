'''
Design a Singly Linked List class.

Your LinkedList class should support the following operations:

LinkedList() will initialize an empty linked list.

- int get(int i) will return the value of the ith node (0-indexed). If the index is out of bounds, return -1.
- void insertHead(int val) will insert a node with val at the head of the list.
- void insertTail(int val) will insert a node with val at the tail of the list.
- bool remove(int i) will remove the ith node (0-indexed). If the index is out of bounds, return false, otherwise return true.
- int[] getValues() return an array of all the values in the linked list, ordered from head to tail.
'''

from typing import List
class Node:
    def __init__(self, val: int):
        self.val = val
        self.next = None
    
class SinglyLinkedList:
    def __init__(self):
        self.size = 0
        self.head = None
    
    def get(self, i: int) -> int:
        if i < 0 or i >= self.size:
            return -1
        curr = self.head
        for _ in range(i):
            curr = curr.next
        return curr.val

    def insertHead(self, val: int) -> None:
        new_node = Node(val)
        new_node.next = self.head
        self.head = new_node
        self.size += 1
    
    def insertTail(self, val: int) -> None:
        new_node = Node(val)
        if not self.head:
            self.head = new_node
        else:
            curr = self.head
            while curr.next:
                curr = curr.next
            curr.next = new_node
        self.size += 1
    
    def remove(self, i: int) -> bool:
        if i < 0 or i >= self.size:
            return False
        if i == 0:
            self.head = self.head.next
        else:
            curr = self.head
            for _ in range(i-1):
                curr = curr.next
            curr.next = curr.next.next
        self.size -= 1
        return True

    def getValues(self) -> List[int]:
        values = []
        curr = self.head
        while curr:
            values.append(curr.val)
            curr = curr.next
        return values


