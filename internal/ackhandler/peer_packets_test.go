package ackhandler

import (
	"testing"

	"github.com/sagernet/quic-go/internal/monotime"
	"github.com/sagernet/quic-go/internal/protocol"
	"github.com/sagernet/quic-go/internal/utils"
)

func receivePeerPackets(t *testing.T, tracker *appDataReceivedPacketTracker, pns ...protocol.PacketNumber) {
	t.Helper()
	for _, pn := range pns {
		if err := tracker.ReceivedPacket(pn, protocol.ECNNon, monotime.Now(), true); err != nil {
			t.Fatalf("packet %d: %v", pn, err)
		}
	}
}

func checkPeerPackets(t *testing.T, stats *utils.ConnectionStats, received uint64, missing int64) {
	t.Helper()
	if got := stats.PeerPacketsReceived.Load(); got != received {
		t.Errorf("received %d, want %d", got, received)
	}
	if got := stats.PeerPacketsMissing.Load(); got != missing {
		t.Errorf("missing %d, want %d", got, missing)
	}
}

func TestPeerPacketsGapsAndLatePackets(t *testing.T) {
	var stats utils.ConnectionStats
	tracker := newAppDataReceivedPacketTracker(utils.DefaultLogger, &stats)
	receivePeerPackets(t, tracker, 0, 1, 2, 5, 3, 7)
	// 4 and 6 did not come; 3 came late and took its gap back
	checkPeerPackets(t, &stats, 6, 2)
	// a duplicate counts nothing
	if tracker.ReceivedPacket(5, protocol.ECNNon, monotime.Now(), true) == nil {
		t.Fatal("a duplicate was taken")
	}
	checkPeerPackets(t, &stats, 6, 2)
	receivePeerPackets(t, tracker, 4)
	checkPeerPackets(t, &stats, 7, 1)
}

func TestPeerPacketsFirstNotSmallest(t *testing.T) {
	var stats utils.ConnectionStats
	tracker := newAppDataReceivedPacketTracker(utils.DefaultLogger, &stats)
	// the first packet to come is 10; 8 comes after 12: 9 and 11 did not come
	receivePeerPackets(t, tracker, 10, 12, 8)
	checkPeerPackets(t, &stats, 3, 2)
	receivePeerPackets(t, tracker, 9, 11)
	checkPeerPackets(t, &stats, 5, 0)
}

func TestPeerPacketsWithoutStats(t *testing.T) {
	tracker := newAppDataReceivedPacketTracker(utils.DefaultLogger, nil)
	receivePeerPackets(t, tracker, 0, 3)
}
