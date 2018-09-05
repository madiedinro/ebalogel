HTTP_LISTEN=8089
SOCK_LISTEN=5000
HOST=192.168.65.1
CH_DSN=tcp://stage.rstat.org:9000/?database=stats
all: GO
image:
	docker build -t xeteq/httpd .

logspout:
	docker run --rm --name=logspout --volume=/var/run/docker.sock:/var/run/docker.sock \
		gliderlabs/logspout \
		raw://$(HOST):$(SOCK_LISTEN)

demo:
	docker build -t xeteq/httpd-demo test
	docker run -d --rm -p 8877:8080 xeteq/httpd-demo
tdata:
	docker run --rm --network custom --name=ncclient alpine sh -c "echo huyev\npachkun\n | nc -u ncsrv 5005"

ncsrv:
	docker run --rm \
	--network custom -p 5005:5005 \
	--name=ncsrv --hostname=ncsrv \
	alpine \
	nc -l 5005

start:
	SOCK_LISTEN=:$(SOCK_LISTEN) HTTP_LISTEN=:$(HTTP_LISTEN) CH_DSN=$(CH_DSN) go run main.go



