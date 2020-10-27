# ------------------------------------------------------------------------------
# Builder Image
# ------------------------------------------------------------------------------
FROM golang:1.15 AS build

WORKDIR /go/src/github.com/figment-networks/rosetta-inspector

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

COPY . .

ENV CGO_ENABLED=0
ENV GOARCH=amd64
ENV GOOS=linux

RUN make setup && make build

# ------------------------------------------------------------------------------
# Target Image
# ------------------------------------------------------------------------------
FROM alpine:3.10 AS release

WORKDIR /app

COPY --from=build \
  /go/src/github.com/figment-networks/rosetta-inspector/rosetta-inspector \
  /app/rosetta-inspector

EXPOSE 5555

ENTRYPOINT ["/app/rosetta-inspector"]
