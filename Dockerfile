############################
# STEP 1 build executable binary
############################
FROM golang:1.18-alpine as builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata alpine-sdk && update-ca-certificates

WORKDIR /build

# cache modules
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN make install
RUN make build

#############################
## STEP 2 build a small image
#############################
FROM alpine

# Import from builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

# Copy our static executable
COPY --from=builder /build/mbmd /usr/local/bin/mbmd

EXPOSE 8080

# Run the binary
ENTRYPOINT ["/usr/local/bin/mbmd"]
