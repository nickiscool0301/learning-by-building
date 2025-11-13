import java.util.*;

/**
 * Demonstration of Consistent Hashing
 *
 * This demo shows:
 * 1. Basic usage - adding servers and distributing keys
 * 2. Minimal redistribution when servers are added/removed
 * 3. Impact of virtual nodes on load distribution
 * 4. Comparison of different virtual node configurations
 */
public class ConsistentHashDemo {

    public static void main(String[] args) {
        System.out.println("========================================");
        System.out.println("  CONSISTENT HASHING DEMO");
        System.out.println("========================================\n");

        // Run different demos
        demo1_BasicUsage();
        demo2_MinimalRedistribution();
        demo3_LoadDistribution();
        demo4_VirtualNodesComparison();
    }

    /**
     * DEMO 1: Basic Usage
     * Shows how to create a hash ring, add servers, and lookup keys
     */
    private static void demo1_BasicUsage() {
        System.out.println("\n" + "=".repeat(60));
        System.out.println("DEMO 1: Basic Usage");
        System.out.println("=".repeat(60));

        // Create a consistent hash with 100 virtual nodes per server
        ConsistentHash<String> ch = new ConsistentHash<>(100);

        // Add servers
        ch.addNode("Server-A");
        ch.addNode("Server-B");
        ch.addNode("Server-C");

        // Print ring state
        ch.printRing();

        // Lookup which server handles specific keys
        System.out.println("Key Lookups:");
        String[] keys = {"user:1001", "user:1002", "user:1003", "session:abc", "cache:home"};

        for (String key : keys) {
            String server = ch.getNode(key);
            System.out.printf("  %-15s -> %s\n", key, server);
        }
    }

    /**
     * DEMO 2: Minimal Redistribution
     * Shows that only affected keys are redistributed when servers change
     */
    private static void demo2_MinimalRedistribution() {
        System.out.println("\n" + "=".repeat(60));
        System.out.println("DEMO 2: Minimal Redistribution (Key Benefit!)");
        System.out.println("=".repeat(60));

        ConsistentHash<String> ch = new ConsistentHash<>(100);

        // Add initial servers
        ch.addNode("Server-A");
        ch.addNode("Server-B");
        ch.addNode("Server-C");

        // Generate test keys
        List<String> keys = generateKeys(1000);

        // Record initial distribution
        Map<String, String> initialMapping = new HashMap<>();
        for (String key : keys) {
            initialMapping.put(key, ch.getNode(key));
        }

        System.out.println("\nInitial setup: 3 servers, 1000 keys distributed");
        printDistributionSummary(ch, keys);

        // Add a new server
        System.out.println("\n>>> Adding Server-D...\n");
        ch.addNode("Server-D");

        // Check how many keys were redistributed
        int redistributedKeys = 0;
        Map<String, String> newMapping = new HashMap<>();

        for (String key : keys) {
            String newServer = ch.getNode(key);
            newMapping.put(key, newServer);
            if (!newServer.equals(initialMapping.get(key))) {
                redistributedKeys++;
            }
        }

        printDistributionSummary(ch, keys);

        System.out.println("\n** RESULT **");
        System.out.printf("Keys redistributed: %d / %d (%.1f%%)\n",
                redistributedKeys, keys.size(),
                (redistributedKeys * 100.0 / keys.size()));
        System.out.println("Expected in perfect scenario: ~25% (1/4 of keys should move)");
        System.out.println("In traditional hashing: ~75% of keys would move!");

        // Show a few examples of moved keys
        System.out.println("\nExamples of redistributed keys:");
        int exampleCount = 0;
        for (String key : keys) {
            if (!newMapping.get(key).equals(initialMapping.get(key)) && exampleCount < 5) {
                System.out.printf("  %s: %s -> %s\n",
                        key, initialMapping.get(key), newMapping.get(key));
                exampleCount++;
            }
        }
    }

    /**
     * DEMO 3: Load Distribution
     * Shows how well keys are distributed across servers
     */
    private static void demo3_LoadDistribution() {
        System.out.println("\n" + "=".repeat(60));
        System.out.println("DEMO 3: Load Distribution Across Servers");
        System.out.println("=".repeat(60));

        ConsistentHash<String> ch = new ConsistentHash<>(150);

        // Add servers
        ch.addNode("Server-A");
        ch.addNode("Server-B");
        ch.addNode("Server-C");
        ch.addNode("Server-D");

        // Generate keys
        List<String> keys = generateKeys(10000);

        System.out.println("\nDistributing 10,000 keys across 4 servers:");
        System.out.println("Expected per server: ~2,500 keys (25%)\n");

        Map<String, Integer> distribution = ch.getDistribution(keys);

        // Calculate statistics
        int total = keys.size();
        double mean = total / (double) distribution.size();
        double variance = 0;

        System.out.println("Actual distribution:");
        for (Map.Entry<String, Integer> entry : distribution.entrySet()) {
            int count = entry.getValue();
            double percentage = (count * 100.0) / total;
            double deviation = ((count - mean) / mean) * 100;

            System.out.printf("  %-10s: %5d keys (%.2f%%) [%+.1f%% from ideal]\n",
                    entry.getKey(), count, percentage, deviation);

            variance += Math.pow(count - mean, 2);
        }

        double stdDev = Math.sqrt(variance / distribution.size());
        System.out.printf("\nStandard Deviation: %.2f keys (%.2f%% of mean)\n",
                stdDev, (stdDev / mean) * 100);
        System.out.println("Lower standard deviation = better balance");
    }

    /**
     * DEMO 4: Impact of Virtual Nodes
     * Compares load distribution with different numbers of virtual nodes
     */
    private static void demo4_VirtualNodesComparison() {
        System.out.println("\n" + "=".repeat(60));
        System.out.println("DEMO 4: Impact of Virtual Nodes on Distribution");
        System.out.println("=".repeat(60));

        int[] virtualNodeCounts = {1, 10, 50, 100, 200};
        List<String> keys = generateKeys(10000);

        System.out.println("\nComparing different virtual node configurations:");
        System.out.println("Testing with 10,000 keys across 4 servers\n");

        System.out.printf("%-15s %-20s %-15s\n", "Virtual Nodes", "Std Deviation", "Max Deviation");
        System.out.println("-".repeat(55));

        for (int vnodes : virtualNodeCounts) {
            ConsistentHash<String> ch = new ConsistentHash<>(vnodes);

            ch.addNode("Server-A");
            ch.addNode("Server-B");
            ch.addNode("Server-C");
            ch.addNode("Server-D");

            Map<String, Integer> distribution = ch.getDistribution(keys);

            // Calculate statistics
            double mean = keys.size() / (double) distribution.size();
            double variance = 0;
            double maxDeviation = 0;

            for (Integer count : distribution.values()) {
                variance += Math.pow(count - mean, 2);
                double deviation = Math.abs(((count - mean) / mean) * 100);
                maxDeviation = Math.max(maxDeviation, deviation);
            }

            double stdDev = Math.sqrt(variance / distribution.size());

            System.out.printf("%-15d %-20.2f %-15.2f%%\n",
                    vnodes, stdDev, maxDeviation);
        }

        System.out.println("\n** INSIGHT **");
        System.out.println("More virtual nodes = better distribution but more memory");
        System.out.println("Typical production values: 100-200 virtual nodes per server");
    }

    // Helper method to generate test keys
    private static List<String> generateKeys(int count) {
        List<String> keys = new ArrayList<>();
        for (int i = 0; i < count; i++) {
            keys.add("key:" + i);
        }
        return keys;
    }

    // Helper method to print distribution summary
    private static void printDistributionSummary(ConsistentHash<String> ch, List<String> keys) {
        Map<String, Integer> distribution = ch.getDistribution(keys);
        System.out.println("Distribution:");
        for (Map.Entry<String, Integer> entry : distribution.entrySet()) {
            double percentage = (entry.getValue() * 100.0) / keys.size();
            System.out.printf("  %-10s: %4d keys (%.1f%%)\n",
                    entry.getKey(), entry.getValue(), percentage);
        }
    }
}
