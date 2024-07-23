package main

import "strconv"

const (
	defaultPort             = "2138"
	defaultHistoryCountdown = "60" // 1m0s
	defaultMaxMemoryStr     = "32" // 32 MB
)

var defaultMaxMemory, _ = strconv.ParseInt(defaultMaxMemoryStr, 10, 64)
