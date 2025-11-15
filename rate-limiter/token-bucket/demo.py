#!/usr/bin/env python3

import time
from TokenBucket import TokenBucket


def demo1_basic_usage():
    """Demo 1: Basic token bucket usage"""
    print("=" * 60)
    print("DEMO 1: Basic Usage")
    print("=" * 60)

    # Create a bucket: 5 tokens max, refill at 1 token/second
    limiter = TokenBucket(capacity=5, refill_rate=1.0)

    print(f"\nInitial tokens: {limiter.get_current_tokens():.2f}")
    print("\nMaking 5 requests quickly (burst):")

    for i in range(5):
        allowed = limiter.allow_request()
        print(f"  Request {i+1}: {'Allowed' if allowed else 'Denied'} "
              f"(tokens left: {limiter.get_current_tokens():.2f})")

    print("\nMaking 6th request (should fail - no tokens):")
    allowed = limiter.allow_request()
    print(f"  Request 6: {'Allowed' if allowed else 'Denied'} "
          f"(tokens left: {limiter.get_current_tokens():.2f})")


def demo2_refill_over_time():
    """Demo 2: Tokens refill over time"""
    print("\n" + "=" * 60)
    print("DEMO 2: Refill Over Time")
    print("=" * 60)

    limiter = TokenBucket(capacity=10, refill_rate=2.0)

    print("\nConfig: 10 tokens max, refill at 2 tokens/second")
    print(f"Initial tokens: {limiter.get_current_tokens():.2f}")

    print("\nUsing all 10 tokens...")
    for i in range(10):
        limiter.allow_request()
    print(f"Tokens after 10 requests: {limiter.get_current_tokens():.2f}")

    print("\nWaiting 3 seconds (should get 6 tokens back)...")
    time.sleep(3)
    print(f"Tokens after 3 seconds: {limiter.get_current_tokens():.2f}")

    print("\nMaking 6 requests (should all succeed):")
    for i in range(6):
        allowed = limiter.allow_request()
        print(f"  Request {i+1}: {'Allowed' if allowed else 'Denied'}")

    print("\nMaking 7th request (should fail):")
    allowed = limiter.allow_request()
    print(f"  Request 7: {'Denied' if not allowed else 'Allowed'}")

def main():
    demo1_basic_usage()
    demo2_refill_over_time()

if __name__ == "__main__":
    main()