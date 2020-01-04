FROM golang:1.13.5-alpine3.11 as builder

WORKDIR /src
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -mod vendor -a -installsuffix cgo -ldflags="-w -s" -o ideal-visual

FROM alpine:3.11

WORKDIR /app

COPY --from=builder /src/ideal-visual .
COPY --from=builder /src/etc/config.yaml /usr/local/etc/config.yaml

EXPOSE 8080

CMD ["sh", "-c", "./ideal-visual -c /usr/local/etc/config.yaml"]
