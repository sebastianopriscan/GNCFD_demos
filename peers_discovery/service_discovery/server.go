package servicediscovery

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"

	"github.com/sebastianopriscan/GNCFD/core/guid"
	"github.com/sebastianopriscan/GNCFD_demos/peers_discovery/service_discovery/pb_go"
)

type Peer struct {
	Id   guid.Guid
	Addr string
}

type Session struct {
	Id    guid.Guid
	Kind  string
	Peers []*Peer
}

type ServiceDiscoveryServer struct {
	pb_go.UnimplementedPeerDiscoveryServer
	Sessions map[guid.Guid]*Session
}

func (sds *ServiceDiscoveryServer) GetSessions(qry *pb_go.ServiceQuery, stream pb_go.PeerDiscovery_GetSessionsServer) error {
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

func (sds *ServiceDiscoveryServer) GetPeer(ctx context.Context, sess *pb_go.Session) (*pb_go.Peer, error) {
	sessID, err := guid.Deserialize([]byte(sess.SessID))
	if err != nil {
		return nil, errors.New("error deserializing sessionID")
	}

	need_sess, ok := sds.Sessions[sessID]
	if !ok {
		return nil, errors.New("session not found")
	}

	idx := rand.Intn(len(need_sess.Peers))
	peer := need_sess.Peers[idx]

	return &pb_go.Peer{Guid: peer.Id.String(), Addr: peer.Addr}, nil
}

func (sds *ServiceDiscoveryServer) GetPeers(sess *pb_go.Session, stream pb_go.PeerDiscovery_GetPeersServer) error {
	sessID, err := guid.Deserialize([]byte(sess.SessID))
	if err != nil {
		return errors.New("error deserializing sessionID")
	}

	need_sess, ok := sds.Sessions[sessID]
	if !ok {
		return errors.New("session not found")
	}

	errs := 0
	for _, peer := range need_sess.Peers {
		toSend := &pb_go.Peer{Guid: peer.Id.String(), Addr: peer.Addr}

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

	need_sess.Peers = append(need_sess.Peers, &Peer{Id: peerID, Addr: join.Peer.Addr})

	return &pb_go.JoinResult{Res: true}, nil
}
