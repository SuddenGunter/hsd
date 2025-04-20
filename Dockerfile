FROM golang:1.24.2-alpine3.21 AS builder

ARG TARGETOS
ARG TARGETARCH

ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH
ENV CGO_ENABLED=0

WORKDIR /src

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ./

RUN go build -o /appx

FROM alpine:3.21

WORKDIR /root/
COPY --from=builder /appx .

CMD ["./appx"]