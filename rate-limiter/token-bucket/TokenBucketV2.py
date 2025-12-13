'''
Token Bucket:
Each user has a bucket of tokens
- If bucket is empty -> refuse
- Else, execute the request, decrease the number of tokens
'''

import time
import threading

from collections import defaultdict
class UserBucket:
    def __init__(self, tokens: int, last_refill: float):
        self.tokens = tokens
        self.last_refill = last_refill

class RateLimiter:
    def __init__(self, capacity: int, refill_rate: int):
        self.user_map = defaultdict(lambda: UserBucket(capacity, time.time()))
        self.user_locks = defaultdict(threading.Lock)  # Lock per user
        self.capacity = capacity
        self.refill_rate = refill_rate

    def is_allowed(self, user_id: str) -> bool:
        with self.user_locks[user_id]:
            user_bucket = self.user_map[user_id]
            current_time = time.time()
            time_elapsed = current_time - user_bucket.last_refill

            # Refill tokens
            tokens_to_add = time_elapsed * self.refill_rate
            user_bucket.tokens = min(self.capacity, user_bucket.tokens + tokens_to_add)
            user_bucket.last_refill = current_time

            if user_bucket.tokens >= 1:
                user_bucket.tokens -= 1
                return True
            return False



    