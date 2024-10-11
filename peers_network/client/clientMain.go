package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"net/http"
	_ "net/http/pprof"

	"github.com/sebastianopriscan/GNCFD/communication"
	"github.com/sebastianopriscan/GNCFD/communication/rpc/grpc/vivaldi/endpoints"
	"github.com/sebastianopriscan/GNCFD/core"
	"github.com/sebastianopriscan/GNCFD/core/impl/vivaldi"
	"github.com/sebastianopriscan/GNCFD/core/nvs"
	"github.com/sebastianopriscan/GNCFD/gossip"
	"github.com/sebastianopriscan/GNCFD/utils/guid"
	lockedmap "github.com/sebastianopriscan/GNCFD/utils/locked_map"
	servicediscovery "github.com/sebastianopriscan/GNCFD_demos/peers_network/service_discovery"
	"github.com/sebastianopriscan/GNCFD_demos/peers_network/service_discovery/pb_go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var clientConn *grpc.ClientConn
var client pb_go.PeerDiscoveryClient

var gncfdCore core.GNCFDCoreInteractionGate

var gossiper gossip.GNCFDGossiper

var ip string

var peerMap lockedmap.LockedMap[guid.Guid, communication.GNCFDCommunicationChannel]
var coreMap lockedmap.LockedMap[guid.Guid, core.GNCFDCoreInteractionGate]

var discover_addr, discover_port string
var my_port string

func extractValueFromEnv(variable string) string {

	valueString, present := os.LookupEnv(variable)

	if !present {
		log.Println("Bad configuration")
		os.Exit(1)
	}

	return valueString
}

func getSession() (*servicediscovery.Session, error) {

	var err error
	creds := grpc.WithTransportCredentials(insecure.NewCredentials())
	clientConn, err = grpc.NewClient(fmt.Sprintf("%s:%s", discover_addr, discover_port), creds)
	if err != nil {
		return nil, fmt.Errorf("error in creating client connection, details: %s", err)
	}

	client = pb_go.NewPeerDiscoveryClient(clientConn)
	stream, err := client.GetSessions(context.Background(), &pb_go.ServiceQuery{})
	if err != nil {
		return nil, fmt.Errorf("error in connecting to discovery, details: %s", err)
	}

	session, err := stream.Recv()
	if err == io.EOF {
		return nil, fmt.Errorf("error, discovery didn't give any answers, details: %s", err)
	}
	if err != nil {
		return nil, fmt.Errorf("error, unable to get response, details: %s", err)
	}

	sessGuid, err := guid.Deserialize([]byte(session.SessID))
	if err != nil {
		return nil, errors.New("error, guid malformed")
	}
	return &servicediscovery.Session{Id: sessGuid, Kind: session.Kind}, nil
}

func getPeers(jInfo *pb_go.JoinInfo) ([]*servicediscovery.Peer, error) {

	stream, err := client.GetPeers(context.Background(), jInfo)
	if err != nil {
		return nil, fmt.Errorf("error in connecting to discovery, details: %s", err)
	}

	retVal := make([]*servicediscovery.Peer, 0)
	errs := 0
	for {
		peer, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			errs++
		}

		peerGuid, err := guid.Deserialize([]byte(peer.Guid))
		if err != nil {
			errs++
		}

		retVal = append(retVal, &servicediscovery.Peer{Id: peerGuid, Addr: peer.Addr})
	}

	if errs > 0 {
		return retVal, errors.New("there have been errors getting peers")
	}
	return retVal, nil
}

func registerToSession(sess *pb_go.Session, me *pb_go.Peer) error {
	joinres, err := client.JoinSession(context.Background(), &pb_go.JoinInfo{Session: sess, Peer: me})
	if err != nil || !joinres.Res {
		return errors.New("error in joining session")
	}

	return nil
}

func extractNumbersFromKind(string) (int, error) {
	return 5, nil
}

func createCore(session *servicediscovery.Session, me *guid.Guid) error {

	dim, err := extractNumbersFromKind(session.Kind)
	if err != nil {
		return fmt.Errorf("error creating core: wrong dimension")
	}
	space, err := nvs.NewRealEuclideanSpace(dim)
	if err != nil {
		return fmt.Errorf("error generating space, details: %s", err)
	}

	myCoords := make([]float64, dim)

	gncfdCore, err = vivaldi.NewVivaldiCore(*me, myCoords, space, 0.001, 0.001)
	if err != nil {
		return fmt.Errorf("error creating core, details: %s", err)
	}

	gncfdCore.SetCoreSession(session.Id)
	coreMap.Map[session.Id] = gncfdCore

	return nil
}

func getMyAddr() error {

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return errors.New("error retrieving network addresses")
	}

	splits := strings.Split(addrs[1].String(), "/")
	ip = splits[0]
	return nil
}

func main() {

	time.Sleep(3 * time.Second)

	go func() {
		_, present := os.LookupEnv("PPROF")
		if present {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}
	}()

	my_port = extractValueFromEnv("CLIENT_SERV_PORT")
	discover_addr = extractValueFromEnv("DISCOVERER_ADDR")
	discover_port = extractValueFromEnv("DISCOVERER_PORT")

	err := getMyAddr()
	if err != nil {
		log.Fatalf(err.Error())
		return
	}

	peerMap = lockedmap.LockedMap[guid.Guid, communication.GNCFDCommunicationChannel]{Map: make(map[guid.Guid]communication.GNCFDCommunicationChannel)}
	coreMap = lockedmap.LockedMap[guid.Guid, core.GNCFDCoreInteractionGate]{Map: make(map[guid.Guid]core.GNCFDCoreInteractionGate)}

	myGuid, err := guid.GenerateGUID()
	if err != nil {
		log.Fatalf("error generating guid, exiting")
		return
	}

	sess, err := getSession()
	if err != nil {
		log.Fatalf("error retrieving session guid, details: %s", err)
		return
	}

	err = createCore(sess, &myGuid)
	if err != nil {
		log.Fatalln(err)
		return
	}

	serverDesc, err := endpoints.ActivateVivaldiGRPCServer("vivaldicore00", fmt.Sprintf("0.0.0.0:%s", my_port), "tcp", nil, &coreMap)
	if err != nil {
		log.Fatalf("error activating GRPC server, details: %s", err)
		return
	}

	gossiper = gossip.NewBlindCounterGossiper(&peerMap, gncfdCore, 2, 10)
	gossiperSubj, ok := gossiper.(*gossip.BlindCounterGossiper)
	if ok {
		gossiperSubj.ObserveSubject(serverDesc.VivServ)
	}

	gossiper.StartGossiping()

	err = registerToSession(&pb_go.Session{SessID: sess.Id.String(), Kind: sess.Kind},
		&pb_go.Peer{Guid: myGuid.String(), Addr: fmt.Sprintf("%s:%s", ip, my_port)})
	if err != nil {
		log.Fatalf("error registering to GRPC session, details: %s", err)
		return
	}

	go func() {
		for {
			time.Sleep(5 * time.Second)
			analyze_vivaldi_core(gncfdCore)
		}
	}()

	go func() {
		for {
			time.Sleep(5 * time.Second)
			gossiper.InsertGossip()
		}
	}()

	for {
		newP, err := getPeers(&pb_go.JoinInfo{
			Session: &pb_go.Session{SessID: sess.Id.String(), Kind: sess.Kind},
			Peer:    &pb_go.Peer{Guid: myGuid.String(), Addr: fmt.Sprintf("%s:%s", ip, my_port)},
		})
		if err != nil {
			continue
		}
		for _, peer := range newP {
			if peer.Id == myGuid {
				continue
			}

			errs := false
			peerMap.Mu.Lock()

			_, ok := peerMap.Map[peer.Id]
			if !ok {
				chann, err := endpoints.NewVivaldiRPCGossipClient(peer.Id, peer.Addr)
				if err != nil {
					errs = true
					goto UNLOCK
				}
				peerMap.Map[peer.Id] = chann
			}

		UNLOCK:
			peerMap.Mu.Unlock()

			if errs {
				log.Println("error creating communication channel")
			}

			if !errs && !ok {
				log.Printf("added new peer with GUID %v and addr %s", peer.Id, peer.Addr)
			}
		}

		time.Sleep(time.Second)
	}
}
