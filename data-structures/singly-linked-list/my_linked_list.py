'''
From LeetCode:
'''


class Node:
    def __init__(self, value):
        self.value = value
        self.next = None

class MyLinkedList:

    def __init__(self):
        self.head = None
        self.length = 0
        
    def get(self, index: int) -> int:
        if index < 0 or index >= self.length:
            return -1
        curr = self.head
        for _ in range(index):
            curr = curr.next
        return curr.value
        
    def addAtHead(self, val: int) -> None:
        '''
        try to set head to this new node
        '''
        node = Node(val)
        node.next = self.head
        self.head = node
        self.length += 1

    def addAtTail(self, val: int) -> None:
        new_node = Node(val)
        if not self.head:
            self.head = new_node
        else:
            curr = self.head
            while curr.next:
                curr = curr.next
            curr.next = new_node
        self.length += 1

    def addAtIndex(self, index: int, val: int) -> None:
        if index < 0 or index > self.length:
            return
        if index == 0:
            self.addAtHead(val)
            return
        new_node = Node(val)
        curr = self.head
        for _ in range(index-1):
            curr = curr.next
        new_node.next = curr.next
        curr.next = new_node
        self.length += 1

    def deleteAtIndex(self, index: int) -> None:
        if index < 0 or index >= self.length:
            return
        if index == 0:
            self.head = self.head.next
            self.length -= 1
            return

        curr = self.head
        for _ in range(index-1):
            curr = curr.next
        curr.next = curr.next.next
        self.length -= 1


# Your MyLinkedList object will be instantiated and called as such:
# obj = MyLinkedList()
# param_1 = obj.get(index)
# obj.addAtHead(val)
# obj.addAtTail(val)
# obj.addAtIndex(index,val)
# obj.deleteAtIndex(index)