package main

import (
	"fmt"
	"hash/fnv"
)

type BloomFilter struct {
	bitArray []byte
	size     int
	numHash  uint
}

func New(size int, numHash uint) *BloomFilter {
	numBytes := (size + 7) / 8
	f := BloomFilter{
		bitArray: make([]byte, numBytes),
		size:     size,
		numHash:  numHash,
	}
	return &f
}

func (bf *BloomFilter) setBit(position int) {
	byteIndex := position / 8
	bitIndex := position % 8

	bf.bitArray[byteIndex] |= (1 << bitIndex)
}

func (bf *BloomFilter) getBit(position int) bool {
	byteIndex := position / 8
	bitIndex := position % 8
	return (bf.bitArray[byteIndex])&(1<<bitIndex) != 0
}

func (bf *BloomFilter) hash(data []byte, seed uint) int {
	// 64 bit FNV-1a hash method
	h := fnv.New64a()
	seedByte := byte(seed) // convert uint to byte
	seedSlice := []byte{seedByte}

	// seedSlice works like 'salt'
	h.Write(seedSlice)
	h.Write(data)
	return (int(h.Sum64() % uint64(bf.size)))
}

func (bf *BloomFilter) Add(data []byte) {
	for i := uint(0); i < bf.numHash; i++ {
		position := bf.hash(data, i)
		bf.setBit(position)
	}
}

func (bf *BloomFilter) Contains(data []byte) bool {
	for i := uint(0); i < bf.numHash; i++ {
		postion := bf.hash(data, i)
		if !bf.getBit(postion) {
			return false
		}
	}
	return true
}

// Main function
func main() {
	fmt.Println("===== Bloom Filter Demo =====")

	bf := New(1000, 3)

	fmt.Println("Adding: apple, banana, cherry")
	bf.Add([]byte("apple"))
	bf.Add([]byte("banana"))
	bf.Add([]byte("cherrry"))

	fmt.Println("\nTesting items that were added:")
	fmt.Printf("apple: %v\n", bf.Contains([]byte("apple")))

	fmt.Println("\nTesting items that were NOT added:")
	fmt.Printf("grape: %v\n", bf.Contains([]byte("grape")))
}
