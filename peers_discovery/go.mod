module github.com/sebastianopriscan/GNCFD_demos/peers_discovery

go 1.22.5

replace github.com/sebastianopriscan/GNCFD => ../lib/GNCFD/src

require (
	github.com/sebastianopriscan/GNCFD v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.66.2
	google.golang.org/protobuf v1.34.2
)

require (
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240604185151-ef581f913117 // indirect
)
