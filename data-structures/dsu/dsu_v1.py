class DSU:
    """
    Disjoint Set Union data structure for managing disjoint sets.
    - find(x): Find the representative of x's set
    - union(x, y): Merge sets containing x and y
    - connected(x, y): Check if x and y are in the same set
    - count_components(): Get number of disjoint sets
    """

    def __init__(self, n):
        """
        Args:
            n: Number of elements
        """
        self.parent = list(range(n))
        self.rank = [0] * n
        self.components = n # Initially, each element is its own component

    def find(self, x):
        """
        Find the root/representative of x's set.
        Uses path compression for optimization.

        Args:
            x: Element to find

        Returns:
            Root of the set containing x

        Time: O(α(n)) - inverse Ackermann, practically O(1)
        """
        if self.parent[x] != x:
            # Path compression: make x point directly to root
            self.parent[x] = self.find(self.parent[x])
        return self.parent[x]

    def union(self, x, y):
        """
        Merge the sets containing x and y.
        Uses union by rank for optimization.

        Args:
            x: Element in first set
            y: Element in second set

        Returns:
            True if sets were merged, False if already in same set

        Time: O(α(n)) - inverse Ackermann, practically O(1)
        """
        root_x = self.find(x)
        root_y = self.find(y)

        # Already in the same set
        if root_x == root_y:
            return False

        # Union by rank: attach smaller tree under larger tree
        if self.rank[root_x] < self.rank[root_y]:
            self.parent[root_x] = root_y
        elif self.rank[root_x] > self.rank[root_y]:
            self.parent[root_y] = root_x
        else:
            # Same rank: arbitrarily choose one as parent and increase its rank
            self.parent[root_y] = root_x
            self.rank[root_x] += 1

        self.components -= 1  # Merged two sets into one
        return True

    def connected(self, x, y):
        """
        Check if x and y are in the same set.

        Args:
            x: First element
            y: Second element

        Returns:
            True if x and y are in the same set, False otherwise
        """
        return self.find(x) == self.find(y)

    def count_components(self):
        """
        Get the number of disjoint sets.

        Returns:
            Number of disjoint sets
        """
        return self.components

    def get_component_size(self, x):
        """
        Get the size of the set containing x.

        Args:
            x: Element in the set

        Returns:
            Size of the set containing x
        """
        root = self.find(x)
        return sum(1 for i in range(len(self.parent)) if self.find(i) == root)


# Example Usage and Test Cases
if __name__ == "__main__":
    print("=== DSU Example: Social Network ===\n")

    # 6 people: 0, 1, 2, 3, 4, 5
    dsu = DSU(6)
    print(f"Initial components: {dsu.count_components()}")  # 6

    # Friendships (connections)
    friendships = [(0, 1), (1, 2), (3, 4)]

    for person1, person2 in friendships:
        if dsu.union(person1, person2):
            print(f"Connected {person1} and {person2}")
        else:
            print(f"{person1} and {person2} already connected")

    print(f"\nComponents after connections: {dsu.count_components()}")  # 3

    # Check connectivity
    print("\n=== Checking Connectivity ===")
    print(f"Are 0 and 2 friends (directly/indirectly)? {dsu.connected(0, 2)}")  # True
    print(f"Are 0 and 3 friends (directly/indirectly)? {dsu.connected(0, 3)}")  # False
    print(f"Are 3 and 4 friends (directly/indirectly)? {dsu.connected(3, 4)}")  # True

    # Merge more groups
    print("\n=== More Connections ===")
    dsu.union(2, 3)  # Merge groups {0,1,2} and {3,4}
    print(f"After connecting 2 and 3: {dsu.count_components()} components")  # 2
    print(f"Are 0 and 4 now friends? {dsu.connected(0, 4)}")  # True

    print("\n=== Component Sizes ===")
    print(f"Size of component containing 0: {dsu.get_component_size(0)}")  # 5
    print(f"Size of component containing 5: {dsu.get_component_size(5)}")  # 1

    # Classic Problem: Detect Cycle in Undirected Graph
    print("\n\n=== Example: Cycle Detection in Graph ===\n")

    def has_cycle(n, edges):
        """
        Detect if undirected graph has a cycle using DSU.

        Args:
            n: Number of nodes (0 to n-1)
            edges: List of edges as (u, v) tuples

        Returns:
            True if graph has a cycle, False otherwise
        """
        dsu = DSU(n)
        for u, v in edges:
            if not dsu.union(u, v):
                # u and v already in same set, adding edge creates cycle
                return True
        return False

    # Graph with cycle: 0-1-2-0
    edges_with_cycle = [(0, 1), (1, 2), (2, 0)]
    print(f"Graph {edges_with_cycle} has cycle: {has_cycle(3, edges_with_cycle)}")

    # Graph without cycle (tree)
    edges_no_cycle = [(0, 1), (1, 2), (2, 3)]
    print(f"Graph {edges_no_cycle} has cycle: {has_cycle(4, edges_no_cycle)}")

    # Count Connected Components
    print("\n\n=== Example: Count Connected Components ===\n")

    def count_connected_components(n, edges):
        """
        Count number of connected components in undirected graph.

        Args:
            n: Number of nodes
            edges: List of edges

        Returns:
            Number of connected components
        """
        dsu = DSU(n)
        for u, v in edges:
            dsu.union(u, v)
        return dsu.count_components()

    # Graph with 3 components: {0,1,2}, {3,4}, {5}
    result = count_connected_components(6, [(0, 1), (1, 2), (3, 4)])
    print(f"6 nodes with edges [(0,1), (1,2), (3,4)]: {result} components")
