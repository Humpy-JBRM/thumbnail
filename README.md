# Thumbnailer Service

A web service to generate thumbnails of:

  - images
  - office documents
  - pdf documents
  - video
  - web pages / URLs

For video, the thumbnail is an animated gif of a selection of the video frames.  For all other types, the thumbnail is a PNG with width 300

## Running on Docker

	- `docker build -t thumbnailer:latest .`

	- `docker run -p 8001:8001 thumbnailer:latest`

and the thumbnailer will be running on http://localhost:8001.


## Installing (VPS or server)
This must be run on a linux machine.

The commands below assume Ubuntu.

Install the required packages:

    sudo apt update
    sudo apt install -y ca-certificates
    sudo update-ca-certificates
    sudo apt install -y libreoffice
    sudo apt install -y libreoffice-writer
    sudo apt install -y python3-pdfminer
    sudo apt install -y python-six
    sudo apt install -y python3-pip
    sudo pip3 install six
    sudo apt install -y imagemagick
    sudo apt install -y libimage-exiftool-perl
    sudo apt install -y ffmpeg
    sudo apt install -y default-jdk 
    sudo apt install -y apt-transport-https 
    sudo apt install -y curl
    sudo apt install -y wget
    sudo apt install -y google-chrome-stable

Install go 1.21 from https://golang.org.

Compile the binary:

`go mod tidy`
`go build .`

## Running

### ImageMagick
By default, ImageMagick will disable pdf thumbnailing due to an earlier security hole in Ghostscript.

You need to enable PDF thumbnailing:

`sudo cp etc/ImageMagick-6/policy.xml /etc/ImageMagick-6/policy.xml`

Execute the program:

`./thumbnailer -f config.yml thumbnail -l 0.0.0.0:8001`

## Operation

The endpoint of the thumbnailer is `http://localhost:8001/api/thumbnail`

This endpoint takes a `POST` request, with a content type of `multipart/form-data`.  So you can either invoke it from code, or it can be the target of a HTML form.

### Thumbnailing Files
This form should have a single input called `file`, of type *file*.

### Thumbnailing Web Pages / URLS
This form should have a single input called `url`, of type *string*, with the URL you want to thumbnail.  This can be a `file:///` url.

## Testing

`curl -o /tmp/foo.png -F file=@YOUR_DOCUMENT_OR_IMAGE http://localhost:8001/api/thumbnail` (non-video)

`curl -o /tmp/foo.gif -F file=@YOUR_VIDEO http://localhost:8001/api/thumbnail` (video)

`curl -o /tmp/foo.png -F url=http://... http://localhost:8001/api/thumbnail` (web page / url)

### File Server Enabled
If `config.server.enabled` is `true`, then the thumbnail server will write the thumbnails to
disk, and serve them over HTTP

`curl -o /dev/null -F id=12345 -F file=@YOUR_DOCUMENT_OR_IMAGE http://localhost:8001/api/thumbnail` 
`curl -o /tmp/foo.png http://localhost:8001/thumbnail/12345`



## TODO
Configurable ability to act as a thumbnail server:

	> GET /thumbnail/{ID}
	> HTTP/1.0 200 OK
	> Content-Type: image/png
	> Content-Length: 2345
	>
	> <PNG thumbnail>


