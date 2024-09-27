package servicediscovery

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"

	"github.com/sebastianopriscan/GNCFD/utils/guid"
	"github.com/sebastianopriscan/GNCFD_demos/peers_network/service_discovery/pb_go"
)

type Peer struct {
	Id   guid.Guid
	Addr string
}

type Session struct {
	Id    guid.Guid
	Kind  string
	Peers map[guid.Guid]string
}

type ServiceDiscoveryServer struct {
	pb_go.UnimplementedPeerDiscoveryServer
	Sessions map[guid.Guid]*Session
}

func (sds *ServiceDiscoveryServer) GetSessions(qry *pb_go.ServiceQuery, stream pb_go.PeerDiscovery_GetSessionsServer) error {
	Mu.Lock()
	defer Mu.Unlock()

	errs := 0
	for _, session := range sds.Sessions {
		toSend := pb_go.Session{SessID: session.Id.String(), Kind: session.Kind}
		err := stream.Send(&toSend)
		if err != nil {
			log.Printf("error sending message in stream, details: %s\n", err)
			errs++
		}
	}

	if errs != 0 {
		return fmt.Errorf("there have been %d errors", errs)
	}

	return nil
}

func (sds *ServiceDiscoveryServer) GetPeer(ctx context.Context, jInfo *pb_go.JoinInfo) (*pb_go.Peer, error) {
	Mu.Lock()
	defer Mu.Unlock()
	sessID, err := guid.Deserialize([]byte(jInfo.Session.SessID))
	if err != nil {
		return nil, errors.New("error deserializing sessionID")
	}

	peerID, err := guid.Deserialize([]byte(jInfo.Peer.Guid))
	if err != nil {
		return nil, errors.New("error deserializing peerID")
	}

	need_sess, ok := sds.Sessions[sessID]
	if !ok {
		return nil, errors.New("session not found")
	}

	peer_node, ok := GuidNode[peerID]
	if !ok {
		return nil, errors.New("peer not registered to session")
	}

	neigh_nodes := NetworkTable[peer_node]

	neighbors := make([]guid.Guid, 0, len(neigh_nodes))

	for _, node := range neigh_nodes {
		neighbor, ok := NodeGuid[node]
		if ok {
			neighbors = append(neighbors, neighbor)
		}
	}

	idx := rand.Intn(len(neighbors))
	peer := neighbors[idx]

	return &pb_go.Peer{Guid: peer.String(), Addr: need_sess.Peers[peer]}, nil
}

func (sds *ServiceDiscoveryServer) GetPeers(jInfo *pb_go.JoinInfo, stream pb_go.PeerDiscovery_GetPeersServer) error {
	Mu.Lock()
	defer Mu.Unlock()
	sessID, err := guid.Deserialize([]byte(jInfo.Session.SessID))
	if err != nil {
		return errors.New("error deserializing sessionID")
	}

	peerID, err := guid.Deserialize([]byte(jInfo.Peer.Guid))
	if err != nil {
		return errors.New("error deserializing peerID")
	}

	need_sess, ok := sds.Sessions[sessID]
	if !ok {
		return errors.New("session not found")
	}

	peer_node, ok := GuidNode[peerID]
	if !ok {
		return errors.New("peer not registered to session")
	}

	neigh_nodes := NetworkTable[peer_node]

	neighbors := make([]guid.Guid, 0, len(neigh_nodes))

	for _, node := range neigh_nodes {
		neighbor, ok := NodeGuid[node]
		if ok {
			neighbors = append(neighbors, neighbor)
		}
	}

	errs := 0
	for _, peer := range neighbors {
		toSend := &pb_go.Peer{Guid: peer.String(), Addr: need_sess.Peers[peer]}

		if err := stream.Send(toSend); err != nil {
			log.Printf("error in sending peer, details: %s", err)
			errs++
		}
	}

	if errs != 0 {
		return fmt.Errorf("there have been %d send errors", errs)
	}

	return nil
}

func (sds *ServiceDiscoveryServer) JoinSession(ctx context.Context, join *pb_go.JoinInfo) (*pb_go.JoinResult, error) {
	sessID, err := guid.Deserialize([]byte(join.Session.SessID))
	if err != nil {
		return &pb_go.JoinResult{Res: false}, errors.New("error deserializing sessionID")
	}

	peerID, err := guid.Deserialize([]byte(join.Peer.Guid))
	if err != nil {
		return &pb_go.JoinResult{Res: false}, errors.New("error deserializing peerID")
	}

	need_sess, ok := sds.Sessions[sessID]
	if !ok {
		return &pb_go.JoinResult{Res: false}, errors.New("session not found")
	}

	Mu.Lock()
	if LastGiven == 10 {
		return &pb_go.JoinResult{Res: false}, nil
	}

	LastGiven++
	GuidNode[peerID] = LastGiven
	NodeGuid[LastGiven] = peerID
	Mu.Unlock()

	need_sess.Peers[peerID] = join.Peer.Addr

	log.Printf("peer with GUID %v and addr %s joined the session\nits number is %v", join.Peer.Guid,
		join.Peer.Addr, LastGiven)

	return &pb_go.JoinResult{Res: true}, nil
}
