package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/sebastianopriscan/GNCFD/core/guid"
	servicediscovery "github.com/sebastianopriscan/GNCFD_demos/peers_discovery/service_discovery"
	"github.com/sebastianopriscan/GNCFD_demos/peers_discovery/service_discovery/pb_go"
	"google.golang.org/grpc"
)

func extractValueFromEnv(variable string) string {

	valueString, present := os.LookupEnv(variable)

	if !present {
		log.Println("Bad configuration")
		os.Exit(1)
	}

	return valueString
}

func main() {

	port := extractValueFromEnv("DISCOVERER_PORT")

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))

	if err != nil {
		log.Fatalf("error in creating tcp connection, details: %s\n", err)
		return
	}

	sessID, err := guid.GenerateGUID()
	if err != nil {
		log.Fatalf("error in generating server guid, details: %s\n", err)
		return
	}
	session := servicediscovery.Session{Id: sessID, Kind: "Vivaldi", Peers: make([]*servicediscovery.Peer, 0)}

	sessionMap := make(map[guid.Guid]*servicediscovery.Session)
	sessionMap[sessID] = &session

	disc_server := &servicediscovery.ServiceDiscoveryServer{Sessions: sessionMap}

	server := grpc.NewServer()
	pb_go.RegisterPeerDiscoveryServer(server, disc_server)
	server.Serve(lis)
}
