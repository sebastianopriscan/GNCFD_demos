.PHONY : gncfd_container install-demos-systemd-services

gncfd_container :
	$(shell cd lib/GNCFD/src ; ./generate_release.sh)
	-docker image rm gncfd_embed:latest
	docker buildx build -t gncfd_embed -f GNCFD.dockerfile .

install-demos-systemd-services :
	sudo mv peers_discovery/discoverer_discovery.service /etc/systemd/system
	sudo mv peers_discovery/client_discovery.service /etc/systemd/system
	sudo mv peers_network/discoverer_network.service /etc/systemd/system
	sudo mv peers_network/client_network.service /etc/systemd/system
	sudo systemctl daemon-reload