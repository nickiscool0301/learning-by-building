from dataclasses import dataclass

@dataclass
class Order:
    order_id: int
    side: str
    quantity: int 
    price: int 
    timestamp: int 

from collections import defaultdict, deque
import heapq
class BookOrder:
    def __init__(self):
        self.buy_levels = defaultdict(deque)
        self.sell_levels = defaultdict(deque)

        self.buy_prices = [] # max heap
        self.sell_prices = [] # min heap

    def add_order(self, order: Order):
        if order.side == 'BUY':
            self._match_buy(order)
            if order.quantity > 0:
                self._add_buy_order(order)
        else:
            self._match_sell(order)
            if order.quantity > 0:
                self._add_sell_order(order)
    
    def _match_buy(self, order_buy: Order):
        '''
        match buy with sell levels
        '''
        while order_buy.quantity > 0 and self.sell_prices:
            best_sell = self.sell_prices[0]
            if order_buy.price < best_sell:
                break
            sell_queue = self.sell_levels[best_sell]
            sell = sell_queue[0]
            traded = min(sell.quantity, order_buy.quantity)
            sell.quantity -= traded
            order_buy.quantity -= traded
            if sell.quantity == 0:
                sell_queue.popleft()
                if not sell_queue:
                    heapq.heappop(self.sell_prices)
                    del self.sell_levels[best_sell]
    
    def _match_sell(self, order_sell: Order):
        while order_sell.quantity > 0 and self.buy_prices:
            best_buy = -self.buy_prices[0]
            if order_sell.price > best_buy:
                break
            buy_queue = self.buy_levels[best_buy]
            buy = buy_queue[0]
            traded = min(buy.quantity, order_sell.quantity)
            buy.quantity -= traded
            order_sell.quantity -= traded
            if buy.quantity == 0:
                buy_queue.popleft()
                if not buy_queue:
                    heapq.heappop(self.buy_prices)
                    del self.buy_levels[best_buy]

    def _add_buy_order(self, order: Order):
        '''
        add order to buy levels + buy levels
        '''
        if order.price not in self.buy_levels:
            heapq.heappush(self.buy_prices, -order.price)
        self.buy_levels[order.price].append(order)

    def _add_sell_order(self, order: Order):
        '''
        add order to buy levels + buy levels
        '''
        if order.price not in self.sell_levels:
            heapq.heappush(self.sell_levels, order.price)
        self.sell_levels[order.price].append(order)
