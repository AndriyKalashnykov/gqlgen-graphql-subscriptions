package constants

import "time"

const (
	// Redis Stream configuration
	RedisStreamRoom   = "room"
	RedisStreamMaxLen = 1
	RedisStreamCount  = 1

	// Server configuration
	ServerPort = ":8080"

	// WebSocket configuration
	WebSocketReadBufferSize    = 1024
	WebSocketWriteBufferSize   = 1024
	WebSocketKeepAlivePing     = 10 * time.Second
	WebSocketSubscriptionToken = 16 // hex length for subscription tokens

	// Cache configuration
	QueryCacheSize = 1000
	APQCacheSize   = 100

	// Redis Stream message field
	RedisMessageField = "message"
)
