# Build
FROM golang:1.26.1-alpine AS build-env
WORKDIR /app
COPY . /app
RUN go mod download
RUN CGO_ENABLED=0 go build ./cmd/vornex

# Release
FROM alpine:3.23.3
RUN apk upgrade --no-cache \
    && apk add --no-cache nmap libpcap bind-tools ca-certificates nmap-scripts
COPY --from=build-env /app/vornex /usr/local/bin/
ENTRYPOINT ["vornex"]
