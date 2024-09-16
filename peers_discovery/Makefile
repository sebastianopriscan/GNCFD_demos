.PHONY : proto-discovery

module := "github.com/sebastianopriscan/GNCFD_demos/peer_discovery"

discovery_proto_loc := ./service_discovery/pb

proto-discovery : $(wildcard $(discovery_proto_loc)/*.proto)
	protoc --proto_path=$(discovery_proto_loc) --go_out=. --go_opt=module=$(module) --go-grpc_out=. --go-grpc_opt=module=$(module) $(wildcard $(discovery_proto_loc)/*.proto)
