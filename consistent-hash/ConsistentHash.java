import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.util.*;

/**
 * Consistent Hashing Implementation
 * Key Features:
 * - Uses MD5 hashing for uniform distribution
 * - Supports virtual nodes (replicas) for better load balancing
 * - O(log n) lookup time using TreeMap
 */

public class ConsistentHash<T> {

    // Hash ring
    private final TreeMap<Long, T> ring;
    private final int numberOfVirtualNodes;

    // MessageDigest for hashing
    private final MessageDigest md;

    public ConsistentHash(int numberOfVirtualNodes) {
        this.ring = new TreeMap<>();
        this.numberOfVirtualNodes = numberOfVirtualNodes;
        try {
            this.md = MessageDigest.getInstance("MD5");
        } catch (NoSuchAlgorithmException e) {
            throw new RuntimeException("MD5 algorithm not found", e);
        }
    }

    private long hash(String key) {
        md.reset();
        byte[] digest = md.digest(key.getBytes());

        // Convert first 8 bytes of MD5 hash to a long value
        long hash = 0;
        for (int i = 0; i < 8; i++) {
            hash = (hash << 8) | (digest[i] & 0xFF);
        }
        return hash;
    }

    public void addNode(T node) {
        // Add virtual nodes to the ring
        for (int i = 0; i < numberOfVirtualNodes; i++) {
            String virtualNodeKey = node.toString() + "#" + i;
            long hash = hash(virtualNodeKey);
            ring.put(hash, node);
        }

        System.out.println("Added node: " + node + " (" + numberOfVirtualNodes + " virtual nodes)");
    }

    public void removeNode(T node) {
        // Remove all virtual nodes associated with this physical node
        for (int i = 0; i < numberOfVirtualNodes; i++) {
            String virtualNodeKey = node.toString() + "#" + i;
            long hash = hash(virtualNodeKey);
            ring.remove(hash);
        }

        System.out.println("Removed node: " + node);
    }

    public T getNode(String key) {
        if (ring.isEmpty()) {
            return null;
        }

        long keyHash = hash(key);

        // Find the first entry with hash >= keyHash (clockwise on the ring)
        Map.Entry<Long, T> entry = ring.ceilingEntry(keyHash);

        if (entry == null) {
            entry = ring.firstEntry();
        }

        return entry.getValue();
    }

    public Set<T> getNodes() {
        return new HashSet<>(ring.values());
    }

    public int getRingSize() {
        return ring.size();
    }

    public Map<T, Integer> getDistribution(List<String> keys) {
        Map<T, Integer> distribution = new HashMap<>();

        for (String key : keys) {
            T node = getNode(key);
            distribution.put(node, distribution.getOrDefault(node, 0) + 1);
        }

        return distribution;
    }

    public void printRing() {
        System.out.println("\n=== Hash Ring State ===");
        System.out.println("Total entries (virtual nodes): " + ring.size());
        System.out.println("Physical nodes: " + getNodes().size());
        System.out.println("\nFirst 10 entries on the ring:");

        int count = 0;
        for (Map.Entry<Long, T> entry : ring.entrySet()) {
            if (count++ >= 10)
                break;
            System.out.printf("  Hash: %20d -> Node: %s\n", entry.getKey(), entry.getValue());
        }
        System.out.println("====================\n");
    }
}
