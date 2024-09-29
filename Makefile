.PHONY : gncfd_container

gncfd_container :
	$(shell cd lib/GNCFD/src ; ./generate_release.sh)
	-docker image rm gncfd_embed:latest
	docker buildx build -t gncfd_embed -f GNCFD.dockerfile .