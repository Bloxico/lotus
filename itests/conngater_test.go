package itests

import (
	"context"
	"reflect"
	"testing"

	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/itests/kit"
	"github.com/libp2p/go-libp2p-core/peer"
	manet "github.com/multiformats/go-multiaddr/net"
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

func TestNetBlockIPAddr(t *testing.T) {
	ctx := context.Background()

	firstNode, secondNode, _, ens := kit.EnsembleTwoOne(t)
	ens.InterconnectAll()

	addrInfo, _ := secondNode.NetAddrsListen(ctx)

	var secondNodeIPs []string

	for _, addr := range addrInfo.Addrs {
		ip, err := manet.ToIP(addr)
		if err != nil {
			continue
		}
		secondNodeIPs = append(secondNodeIPs, ip.String())
	}

	err := firstNode.NetBlockAdd(ctx, api.NetBlockList{IPAddrs: secondNodeIPs})
	if err != nil {
		require.NoError(t, err)
	}

	list, err := firstNode.NetBlockList(ctx)
	if err != nil {
		require.NoError(t, err)
	}

	if len(list.IPAddrs) == 0 || reflect.DeepEqual(list.IPAddrs, secondNodeIPs) {
		t.Errorf("blocked ip not in blocked ip list")
	}

	err = firstNode.NetBlockRemove(ctx, api.NetBlockList{IPAddrs: secondNodeIPs})
	if err != nil {
		require.NoError(t, err)
	}

	list, err = firstNode.NetBlockList(ctx)
	if err != nil {
		require.NoError(t, err)
	}

	if len(list.IPAddrs) > 0 {
		t.Errorf("failed to remove blocked ip from blocked ip list")
	}

}
