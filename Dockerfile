FROM golang:1.18 AS builder
RUN apk update && apk add ca-certificates
RUN adduser -D -u 1000 basic_user
WORKDIR /src
COPY go.mod .
COPY go.sum .
ARG goproxy=""
ENV GOPROXY ${goproxy}
RUN go mod download
ADD . /src
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags -buildvcs=false '-s -w' -o /app/microservice

FROM scratch
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /app/microservice /app/docker-builder
USER basic_user
EXPOSE 8080
ENTRYPOINT ["/app/docker-builder"]