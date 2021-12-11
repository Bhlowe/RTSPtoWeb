# Dockerfile for RTSPtoWeb. Use host mode to expose to WAN.
# Build: docker build -t rtsptoweb .
# Run in host mode: docker run -d --network host rtsptoweb
# Run in contained mode: docker run -d -p 8083:8083 rtsptoweb
# (change -d to -it for interactive mode)

FROM golang:alpine

# Set necessary environmet variables needed for our image
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Get Git
RUN apk add --no-cache git

RUN git clone https://github.com/bhlowe/RTSPtoWeb.git

WORKDIR /build/RTSPtoWeb

# Build executable
RUN go build


WORKDIR /dist

RUN cp -r /build/RTSPtoWeb/web . && cp /build/RTSPtoWeb/RTSPtoWeb .
RUN mkdir /dist/config

# Copy default config.json to config dir/volume
RUN cp /build/RTSPtoWeb/config.json /dist/config

# Or copy your local config.json...
# COPY config.json /dist/config/.

# Export necessary port # TODO Get list of all ports used by RTSPtoWeb needed for webrtc, etc.
EXPOSE 8083

# Expose config dir
VOLUME /dist/config

CMD ["./RTSPtoWeb", "-config", "config/config.json"]
