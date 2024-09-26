FROM gncfd_embed:latest

COPY go.mod /home/peers_network/go.mod
COPY go.sum /home/peers_network/go.sum

COPY service_discovery/ /home/peers_network/service_discovery
COPY client/ /home/peers_network/client

WORKDIR /home/peers_network

ARG RELEASE=

RUN go build ${RELEASE} client/clientMain.go client/analyze_vivaldi.go

CMD [ "./clientMain" ]