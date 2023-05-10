
all:	thumbnailer

thumbnailer:
	go mod tidy
	go build -o thumbnailer .

docker:
	docker build -t thumbnailer:latest .

clean:
	rm -f thumbnailer
	docker rmi -f thumbnailer:latest 2>/dev/null || /bin/true

test:
	go mod tidy || /bin/true
	go test -v ./...
	
bdd:
	godog run -t ~@wip

