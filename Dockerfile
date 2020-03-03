# Build
FROM golang:alpine AS builder
WORKDIR /omnitureproxy
COPY . .
RUN go build -o omnitureproxy

# Run
FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN mkdir "_archive"
COPY --from=builder /omnitureproxy/omnitureproxy .
ENV GIN_MODE=release
ENTRYPOINT [ "./omnitureproxy" ]