class RateLimiterSliding:
    def __init__(self, max_requests: int, window_size: int):
        self.max_requests = max_requests
        self.window_size = window_size
        self.request = defaultdict(deque)

    def is_allowed(self, user_id:str) -> bool:
        current_time = time.time()
        q = self.request[user_id]

        while q and q[0] <= current_time - self.window_size:
            q.popleft() 
        
        if len(q) >= self.max_requests:
            return False
        q.append(current_time)
        return True
    