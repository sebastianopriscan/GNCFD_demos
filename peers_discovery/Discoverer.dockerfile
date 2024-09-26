FROM gncfd_embed:latest

COPY go.mod /home/peers_discovery/go.mod
COPY go.sum /home/peers_discovery/go.sum

COPY service_discovery/ /home/peers_discovery/service_discovery
COPY server/ /home/peers_discovery/server

WORKDIR /home/peers_discovery

ARG RELEASE=

RUN go build ${RELEASE} server/serverMain.go

CMD [ "./serverMain" ]