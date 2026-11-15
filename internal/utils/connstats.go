package utils

import "sync/atomic"

// ConnectionStats stores stats for the connection. See the public
// ConnectionStats struct in connection.go for more information
type ConnectionStats struct {
	BytesSent       atomic.Uint64
	PacketsSent     atomic.Uint64
	BytesReceived   atomic.Uint64
	PacketsReceived atomic.Uint64
	BytesLost       atomic.Uint64
	PacketsLost     atomic.Uint64

	// the 1-RTT packets of the peer (podkop-net): each received one once, and the packet
	// numbers below the largest received that did not come
	PeerPacketsReceived atomic.Uint64
	PeerPacketsMissing  atomic.Int64
}
