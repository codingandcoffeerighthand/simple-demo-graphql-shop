FROM golang:1.23-alpine AS builder
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/app
COPY  go.mod go.sum ./
COPY vendor vendor
COPY account account
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/app ./account/cmd/account
FROM alpine:latest
WORKDIR /usr/bin
COPY --from=builder /go/bin/app .
EXPOSE 8001
CMD ["./app"]