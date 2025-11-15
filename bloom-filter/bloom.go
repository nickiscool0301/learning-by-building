package bloomfilter

import "hash/fnv"

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
	h.Write(seedSlice)
	h.Write(data)
	return (int(h.Sum64() % uint64(bf.size)))
}
