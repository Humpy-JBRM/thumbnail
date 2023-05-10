# TODO: two-stage build
# TODO: use the absolute minimum base package and not a full distro
FROM ubuntu:latest

RUN apt update
RUN apt install -y ca-certificates
RUN update-ca-certificates

RUN apt install -y rclone

# Install all the tools we're gonna need for thumbnailing etc
RUN apt install -y libreoffice
RUN apt install -y libreoffice-writer

# Install what we need for smartfs
RUN apt install -y inotify-tools

# pdf2txt
RUN apt install -y python3-pdfminer
RUN apt install -y python-six
RUN apt install -y python3-pip
RUN pip3 install six

# Imagemagick
RUN apt install -y imagemagick

# exiftool
RUN apt install -y libimage-exiftool-perl

# Video thumbnailing and processing
RUN apt install -y ffmpeg

# Others
RUN apt install -y default-jdk 
RUN apt install -y apt-transport-https 
RUN apt install -y curl
RUN apt install -y wget

COPY ROOT/ROOT.tar.gz /var/tmp
RUN tar xvzf /var/tmp/ROOT.tar.gz

# Copy the binary and the config
COPY ./thumbnailer /
COPY config.yml /

ENTRYPOINT ["/thumbnailer", "-f", "/config.yml", "thumbnail", "-l", "0.0.0.0:8000"]
