FROM golang:1.18 AS BUILDER

WORKDIR /usr/src/app

RUN apt-get update && apt-get install -y libvips libvips-dev

COPY go.* ./
RUN go mod download

COPY . .
RUN go build ./cmd/server

FROM debian:bullseye

RUN apt-get update && apt-get install -y ca-certificates libvips libvips-dev

ENV GIN_MODE=release

COPY --from=BUILDER /usr/src/app/server /usr/local/bin/server

WORKDIR /usr/src/app

ADD views /usr/src/app/views
ADD public /usr/src/app/public

CMD ["/usr/local/bin/server"]