'''
Implement a thread-safe bounded blocking queue that has the following methods:

- BoundedBlockingQueue(int capacity) The constructor initializes the queue with a maximum capacity.
- void enqueue(int element) Adds an element to the front of the queue. If the queue is full, the calling thread is blocked until the queue is no longer full.
- int dequeue() Returns the element at the rear of the queue and removes it. If the queue is empty, the calling thread is blocked until the queue is no longer empty.
- int size() Returns the number of elements currently in the queue.


'''


class BlockingQueue:
    def __init__(self, capacity: int):
        pass

    def enqueue(self, element: int) -> None:
        pass
    
    def dequeue(self) -> int:
        pass

    def size(self) -> int:
        pass


