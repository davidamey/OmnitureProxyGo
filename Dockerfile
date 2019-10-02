# Build
FROM iron/go:dev AS builder
WORKDIR /go/src/github.com/davidamey/omnitureproxy
COPY . .
RUN go get
RUN go build -o omnitureproxy

# Run
FROM iron/go:dev
RUN mkdir "_archive"
COPY --from=builder /go/src/github.com/davidamey/omnitureproxy/omnitureproxy .
ENV GIN_MODE=release
ENTRYPOINT [ "./omnitureproxy", "-archive", "/var/omniture_logs" ]