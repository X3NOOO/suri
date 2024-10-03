package main

import "strconv"

const (
	defaultPort             = "2138"
	defaultHistoryCountdown = "60" // 1m0s
	defaultMaxMemoryStr     = "32" // 32 MB
	defaultWhisperTempStr   = "0.1"
)

var defaultMaxMemory, _ = strconv.ParseInt(defaultMaxMemoryStr, 10, 64)
var defaultWhisperTemp, _ = strconv.ParseFloat(defaultWhisperTempStr, 64)
