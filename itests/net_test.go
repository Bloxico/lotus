package itests

import (
	"context"
	"testing"

	"github.com/filecoin-project/lotus/itests/kit"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/require"
)

func TestNetConn(t *testing.T) {
	ctx := context.Background()

	firstNode, secondNode, _, _ := kit.EnsembleTwoOne(t)

	secondNodeID, err := secondNode.ID(ctx)
	if err != nil {
		require.NoError(t, err)
	}

	connState := getConnState(ctx, t, firstNode, secondNodeID)

	if connState != network.NotConnected {
		t.Errorf("node should be not connected to peers. %s", err.Error())
	}

	addrInfo, err := secondNode.NetAddrsListen(ctx)
	if err != nil {
		require.NoError(t, err)
	}

	err = firstNode.NetConnect(ctx, addrInfo)
	if err != nil {
		t.Errorf("nodes failed to connect. %s", err.Error())
	}

	connState = getConnState(ctx, t, firstNode, secondNodeID)

	if connState != network.Connected {
		t.Errorf("peer does not have connected state")
	}

	peers, err := firstNode.NetPeers(ctx)
	if err != nil {
		require.NoError(t, err)
	}

	// treba mi neki dobar nacin da proverim ovo, nekako mi je glupo da hardkodujem peers[0]
	// takodje mi je glupo da proveravam ovo u for loop jer je ocigledno da samo 1 moze u ovom slucaju da postoji
	// sta vi mislite
	if len(peers) > 0 && peers[0].ID != addrInfo.ID {
		t.Errorf("connected peer does not exist in network")
	}

	err = firstNode.NetDisconnect(ctx, secondNodeID)
	if err != nil {
		t.Errorf("nodes failed to disconnect. %s", err.Error())
	}

	connState = getConnState(ctx, t, firstNode, secondNodeID)

	if connState != network.NotConnected {
		t.Errorf("peer should have disconnected")
	}

	peers, err = firstNode.NetPeers(ctx)
	if err != nil {
		require.NoError(t, err)
	}

	if len(peers) > 0 {
		t.Errorf("there should be no peers in network after disconnecting node")
	}

}

func getConnState(ctx context.Context, t *testing.T, node *kit.TestFullNode, peer peer.ID) network.Connectedness {
	connState, err := node.NetConnectedness(ctx, peer)
	if err != nil {
		require.NoError(t, err)
	}

	return connState
}
