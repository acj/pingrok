FROM golang:1.15 AS build

ENV CGO_ENABLED=0

WORKDIR /opt/pingrok

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./
COPY docker/run.sh run.sh

RUN go build -o /go/bin/pingrok

FROM alpine:3.12 AS final

COPY --from=build /go/bin/pingrok /pingrok
COPY --from=build /opt/pingrok/run.sh /run.sh

ENTRYPOINT ["/run.sh"]
