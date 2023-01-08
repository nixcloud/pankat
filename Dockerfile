# https://docs.docker.com/language/golang/build-images/
FROM golang:1.19-alpine AS build-pankat-static

WORKDIR /app

# recreate directory structure with correct relative path to ../../pankat for go.mod/go.sum files
COPY cmd/pankat-static/go.* ./cmd/pankat-static/
COPY pankat/go.* ./pankat/
RUN cd ./cmd/pankat-static/ && go mod download

COPY cmd ./cmd
COPY pankat ./pankat

RUN cd ./cmd/pankat-static && go build -o pankat-static

FROM golang:1.19-alpine AS build-pankat-server

WORKDIR /app

# recreate directory structure with correct relative path to ../../pankat for go.mod/go.sum files
COPY cmd/pankat-server/go.* ./cmd/pankat-server/
COPY pankat/go.* ./pankat/
RUN cd ./cmd/pankat-server/ && go mod download

COPY cmd ./cmd
COPY pankat ./pankat

RUN cd ./cmd/pankat-server && go build -o pankat-server

## Deploy
FROM golang:1.19-alpine

RUN apk add --no-cache pandoc rsync

WORKDIR /

COPY --from=build-pankat-server /app/cmd/pankat-server/pankat-server /pankat-server
COPY --from=build-pankat-static /app/cmd/pankat-static/pankat-static /pankat-static

RUN ln -s /pankat-server /bin/pankat-server
RUN ln -s /pankat-static /bin/pankat-static

RUN addgroup -S pankat
RUN adduser -S -D -h /usr/app/src pankat pankat
USER pankat

EXPOSE 8000