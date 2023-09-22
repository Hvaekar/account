package utils

import (
	"math/rand"
	"time"
	"unsafe"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const numberBytes = "0123456789"
const (
	idxBits = 6              // 6 bits to represent a letter index
	idxMask = 1<<idxBits - 1 // All 1-bits, as many as letterIdxBits
	idxMax  = 63 / idxBits   // # of letter indices fitting in 63 bits
)

func RandString(n int) string {
	var src = rand.NewSource(time.Now().UnixNano())

	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), idxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), idxMax
		}
		if idx := int(cache & idxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= idxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

func RandInt64(n int) int64 {
	var src = rand.NewSource(time.Now().UnixNano())

	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), idxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), idxMax
		}
		if idx := int(cache & idxMask); idx < len(numberBytes) {
			b[i] = numberBytes[idx]
			i--
		}
		cache >>= idxBits
		remain--
	}

	return *(*int64)(unsafe.Pointer(&b))
}
