
all:	thumbnailer

thumbnailer:
	go build -o humpy .

clean:
	rm -f humpy

