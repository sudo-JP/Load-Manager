package salt

import (
	"time"
	"strconv"
)

type RandomFunction func() uint64

// current time in nano
var seed = uint64(time.Now().UnixNano())

func LFSR() uint64 {
    lfsr := seed
    
	bit := (
    	(lfsr >> 63) ^
    	(lfsr >> 3)  ^
    	(lfsr >> 2)  ^
    	(lfsr >> 0)) & 1

    lfsr = (lfsr >> 1) | (bit << 63)
    seed = lfsr
    return lfsr
}

func XORShift() uint64 {
	x := seed 
	x ^= x << 13
    x ^= x >> 7
    x ^= x << 17
    seed = x
    return x
}

func LCG() uint64 {
    const a = 1664525
    const c = 1013904223
    seed = a * seed + c
    return seed 
}

func Salt() string {
	random := LFSR() % 3

	var f [3]RandomFunction
	f[0] = LFSR
	f[1] = XORShift
	f[2] = LCG

	num := f[random]() ^ f[(random+1)%3]()

	random = LFSR() % 5
	if random == 0 {
        seed = uint64(time.Now().UnixNano())
	}
	return strconv.FormatUint(num, 16)
}
