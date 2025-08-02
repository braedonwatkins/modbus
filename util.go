package main

func PackBits(bits []bool) ([]byte, byte) {

	// intentional integer division +7 is so we get 1 more byte than necessary in the worst case
	byteCount := (len(bits) + 7) / 8
	result := make([]byte, byteCount)

	for i, bit := range bits {
		if bit {
			byteIndex := i / 8
			bitIndex := i % 8

			result[byteIndex] |= (1 << bitIndex)
		}
	}

	return result, uint8(byteCount)
}
