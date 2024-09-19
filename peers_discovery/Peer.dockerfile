FROM gncfd_embed:latest

COPY go.mod /home/peers_discovery/go.mod
COPY go.sum /home/peers_discovery/go.sum

COPY service_discovery/ /home/peers_discovery/service_discovery
COPY client/ /home/peers_discovery/client

WORKDIR /home/peers_discovery

RUN go build client/clientMain.go client/analyze_vivaldi.go

CMD [ "./clientMain" ]