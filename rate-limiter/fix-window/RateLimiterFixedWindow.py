class RateLimiterFixedWindow:
    def __init__(self, max_requests: int, window_size: int):
        self.max_requests = max_requests
        self.window_size = window_size
        self.request_counts = {}
    
    def is_allowed(self, user_id: str) -> bool:
        current_time = int(time.time())
        window_start = current_time - (current_time % self.window_size)

        if user_id not in self.request_counts:
            self.request_counts[user_id] = (window_start, 0)

        start, count = self.request_counts[user_id]

        if start < window_start:
            self.request_counts[user_id] = (window_start, 1)
            return True
        else:
            if count < self.max_requests:
                self.request_counts[user_id] = (start, count + 1)
                return True
            else:
                return False