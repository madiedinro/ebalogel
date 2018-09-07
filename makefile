HTTP_LISTEN=8089
HTTP_PORT=8088
SOCK_LISTEN=5000
CH_DSN=tcp://stage.rstat.org:9000/?database=stats


build:
	docker build -t xeteq/httpd .

run:
	docker build -t madiedinro/ebaloger -f Dockerfile.dev .
	docker run -it --rm --name=logspout --hostname=logspout \
		-e DEBUG=1 \
		-e LOGSPOUT=ignore \
		--volume=/var/run/docker.sock:/var/run/docker.sock \
		-p 8088:80 \
		-p 8089:8080 \
		-p 5001:5000 \
		madiedinro/ebaloger go run *.go

nc-client:
	docker run --rm --network custom --name=nc_client alpine sh -c "echo huyev\npachku\n | nc -u ncsrv 5005"

nc-listen:
	docker run --rm \
	--network custom -p 5005:5005 \
	--name=nc_listener --hostname=nc_listener \
	alpine \
	nc -l 5005

start:
	SOCK_LISTEN=:$(SOCK_LISTEN) \
		HTTP_LISTEN=:$(HTTP_LISTEN) \
		CH_DSN=$(CH_DSN) \
		PORT=8088 \
		ROUTESPATH=data/ \
		go run *.go

bump-patch:
	bumpversion patch

bump-minor:
	bumpversion minor
