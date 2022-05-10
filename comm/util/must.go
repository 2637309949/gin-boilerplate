package util

import (
	"strconv"
)

func MustUInt(str string) uint {
	i, _ := strconv.ParseInt(str, 10, 64)
	return uint(i)
}

func MustInt64(str string) int64 {
	i, _ := strconv.ParseInt(str, 10, 64)
	return i
}

func MustUInt64(str string) uint64 {
	i, _ := strconv.ParseUint(str, 10, 64)
	return i
}

func MustInt32(str string) int32 {
	i, _ := strconv.ParseInt(str, 10, 32)
	return int32(i)
}

func MustInt(str string) int {
	i, _ := strconv.ParseInt(str, 10, 32)
	return int(i)
}

func MustUInt32(str string) uint32 {
	i, _ := strconv.ParseUint(str, 10, 32)
	return uint32(i)
}

func NumBit(i uint64) (bit uint64) {
	for i > 0 {
		if i > 0 {
			bit += 1
		}
		i = i / 10
	}
	return bit
}
