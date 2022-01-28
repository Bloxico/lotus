package itests

import (
	"context"
	"testing"

	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/itests/kit"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/require"
)

func TestNetBlockPeer(t *testing.T) {
	ctx := context.Background()

	firstNode, secondNode, _, ens := kit.EnsembleTwoOne(t)
	ens.InterconnectAll()

	secondNodeID, err := secondNode.ID(ctx)
	if err != nil {
		require.NoError(t, err)
	}

	err = firstNode.NetBlockAdd(ctx, api.NetBlockList{Peers: []peer.ID{secondNodeID}})
	if err != nil {
		require.NoError(t, err)
	}

	list, err := firstNode.NetBlockList(ctx)
	if err != nil {
		require.NoError(t, err)
	}

	if len(list.Peers) == 0 || list.Peers[0] != secondNodeID {
		t.Errorf("blocked peer not in blocked peer list")
	}

	err = firstNode.NetBlockRemove(ctx, api.NetBlockList{Peers: []peer.ID{secondNodeID}})
	if err != nil {
		require.NoError(t, err)
	}

	list, err = firstNode.NetBlockList(ctx)
	if err != nil {
		require.NoError(t, err)
	}

	if len(list.Peers) > 0 {
		t.Errorf("failed to remove blocked peer from blocked peer list")
	}

}
