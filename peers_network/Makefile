.PHONY : proto-discovery debug-server debug-client release-server release-client

module := "github.com/sebastianopriscan/GNCFD_demos/peer_discovery"

discovery_proto_loc := ./service_discovery/pb

proto-discovery : $(wildcard $(discovery_proto_loc)/*.proto)
	protoc --proto_path=$(discovery_proto_loc) --go_out=. --go_opt=module=$(module) --go-grpc_out=. --go-grpc_opt=module=$(module) $(wildcard $(discovery_proto_loc)/*.proto)


debug-server :
	env $(shell cat .env | while read line ; do echo -n "$$line "; done) go run -tags debug server/serverMain.go 

debug-client :
	env $(shell cat .env | while read line ; do echo -n "$$line "; done) CLIENT_SERV_PORT=$(CL_PORT) go run -tags debug client/clientMain.go client/analyze_vivaldi.go 

release-server :
	env $(shell cat .env | while read line ; do echo -n "$$line "; done) go run -tags release server/serverMain.go 

release-client :
	env $(shell cat .env | while read line ; do echo -n "$$line "; done) CLIENT_SERV_PORT=$(CL_PORT) go run -tags release client/clientMain.go client/analyze_vivaldi.go 