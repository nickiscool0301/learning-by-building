import time

"""
Token Bucket Rate Limiter Implementation

How it works:
1. Bucket has a maximum capacity of tokens
2. Tokens are added at a fixed rate (refillRate tokens per second)
3. Each request consumes 1 token
4. If no tokens available, request is rejected
5. Allows bursts if tokens have accumulated
"""

class TokenBucket:
    def __init__(self, capacity: int, refill_rate: float):
        if capacity <= 0:
            raise ValueError("Capacity must be positive")
        if refill_rate <= 0:
            raise ValueError("Refill rate must be positive")

        self.capacity = capacity
        self.refill_rate = refill_rate
        self.current_tokens = float(capacity)
        self.last_refill_time = time.time()

    def _refill(self) -> None:
        now = time.time()
        time_elapsed = now - self.last_refill_time
        tokens_to_add = time_elapsed * self.refill_rate

        # avoid overflow
        self.current_tokens = min(self.capacity, self.current_tokens + tokens_to_add)
        self.last_refill_time = now

    def allow_request(self, tokens: int = 1) -> bool:

        self._refill()

        # Check token
        if self.current_tokens >= tokens:
            self.current_tokens -= tokens
            return True
        else:
            return False

    def get_current_tokens(self) -> float:
        self._refill()
        return self.current_tokens

    def reset(self) -> None:
        self.current_tokens = float(self.capacity)
        self.last_refill_time = time.time()