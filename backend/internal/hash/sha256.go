package hash

import (
	"fmt"
)

const (
	h0 = 0x6a09e667
	h1 = 0xbb67ae85
	h2 = 0x3c6ef372
	h3 = 0xa54ff53a
	h4 = 0x510e527f
	h5 = 0x9b05688c
	h6 = 0x1f83d9ab
	h7 = 0x5be0cd19
	FirstWord = 16
	Words = 64
)

func padding(b []byte) []byte {
	// Add 1 bit, just a bit marking the end of msg
	// 0x80 == 1000 000
	b = append(b, 0x80)	
	bit_size := len(b)
	blockSize := 1 << 6  

	// reserved = 512 - 64, 64 bits = 8 bytes for length 
	// in bytes, 64 - 8 for reserved 
	reserved := blockSize - (1 << 3)

	// To calculate the length of the msg 
	// we have to get 8 (64 bits) bytes of the bit_size length
	msg_size := make([]byte, 8)
	for i := range 8 {
		// Mask it with 0xF = 1111 
		j := (bit_size * 8 >> i * 8) & 0xFF
		msg_size[7 - i] = byte(j)
	}

	// If b is within the reserved
	modded := bit_size % blockSize
	if modded <= reserved {
		// Then I can add the length 8 bytes = 64 in 
		for i := 0; i < reserved - modded; i++ {
			b = append(b, 0)
		}
		b = append(b, msg_size...)
	} else {
		// Fill the current to bytes fit 64 (bytes) * n 
		for range modded {
			b = append(b, 0)
		}

		// Start fresh, fill all 64 - 8 bytes to 0 
		// Then add the 8 bytes for length 
		for i := 0; i < blockSize - (1 << 3); i++ {
			b = append(b, 0)
		}
		b = append(b, msg_size...)
	}

	return b 
}

func ROTR(x uint32, n int) uint32 {
	return (x >> n) | (x << (32 - n))
}

func small_sigma0(x uint32) uint32 {
	return ROTR(x, 7) ^ ROTR(x, 18) ^ (x >> 3)
}

func small_sigma1(x uint32) uint32 {
	return ROTR(x, 17) ^ ROTR(x, 19) ^ (x >> 10)
}

func big_sigma0(x uint32) uint32 {
	return ROTR(x, 2) ^ ROTR(x, 13) ^ ROTR(x, 22)
}

func big_sigma1(x uint32) uint32 {
	return ROTR(x, 6) ^ ROTR(x, 11) ^ ROTR(x, 25)
}

// Choose
func ch(x uint32, y uint32, z uint32) uint32 {
	return (x & y) ^ (^x & z)
}

// Majority
func maj(x uint32, y uint32, z uint32) uint32 {
	return (x & y) ^ (x & z) ^ (y & z)
}

func SHA256(str string) string {
	b := make([]byte, len(str))
	copy(b, str)

	padded := padding(b)
	H := [8]uint32{h0, h1, h2, h3, h4, h5, h6, h7}
	K := [64]uint32{0x428a2f98, 0x71374491, 0xb5c0fbcf, 0xe9b5dba5, 0x3956c25b, 0x59f111f1, 0x923f82a4, 0xab1c5ed5,0xd807aa98, 0x12835b01, 0x243185be, 0x550c7dc3, 0x72be5d74, 0x80deb1fe, 0x9bdc06a7, 0xc19bf174,0xe49b69c1, 0xefbe4786, 0x0fc19dc6, 0x240ca1cc, 0x2de92c6f, 0x4a7484aa, 0x5cb0a9dc, 0x76f988da,0x983e5152, 0xa831c66d, 0xb00327c8, 0xbf597fc7, 0xc6e00bf3, 0xd5a79147, 0x06ca6351, 0x14292967,0x27b70a85, 0x2e1b2138, 0x4d2c6dfc, 0x53380d13, 0x650a7354, 0x766a0abb, 0x81c2c92e, 0x92722c85,0xa2bfe8a1, 0xa81a664b, 0xc24b8b70, 0xc76c51a3, 0xd192e819, 0xd6990624, 0xf40e3585, 0x106aa070,0x19a4c116, 0x1e376c08, 0x2748774c, 0x34b0bcb5, 0x391c0cb3, 0x4ed8aa4a, 0x5b9cca4f, 0x682e6ff3,0x748f82ee, 0x78a5636f, 0x84c87814, 0x8cc70208, 0x90befffa, 0xa4506ceb, 0xbef9a3f7, 0xc67178f2 }

	// For each 64 byte block 
	for i := 0; i < len(padded); i += 64 {
		block := padded[i:i+64]

		// Array of 64 32-bit words
		var W[Words]uint32

		// First word is taken from the block 
		for j := range FirstWord {
    		W[j] = 	uint32(block[4*j]) << 24
    		W[j] |= uint32(block[4*j+1])<< 16 
    		W[j] |= uint32(block[4*j+2]) << 8 
    		W[j] |= uint32(block[4*j+3])
		}	

		// The rest 
		for j := FirstWord; j < Words; j++ {
			W[j] = W[j - 16] + small_sigma0(W[j - 15]) + W[j - 7] + small_sigma1(W[j - 2])
		}

		// Working vars 
		a := H[0]
		b := H[1]
		c := H[2]
		d := H[3]
		e := H[4]
		f := H[5]
		g := H[6]
		h := H[7]

		for i := range 64 {
			T1 := h + big_sigma1(e) + ch(e, f, g) + W[i] + K[i]
			T2 := big_sigma0(a) + maj(a, b, c)
			
			h = g
           	g = f
           	f = e
           	e = d + T1
           	d = c
           	c = b
           	b = a
           	a = T1 + T2
		}

		H[0] += a
		H[1] += b
		H[2] += c 
		H[3] += d
		H[4] += e 
		H[5] += f 
		H[6] += g 
		H[7] += h
	}


	var s string
	for _, val := range H {
    	s += fmt.Sprintf("%08x", val) // prints val as 8-digit hex with leading zeros
	}
	return s
}
