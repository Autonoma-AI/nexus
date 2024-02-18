############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git ca-certificates openssl && \
  update-ca-certificates

# Create unpriviliged user.
ENV USER=nexus
ENV UID=10001
# See https://stackoverflow.com/a/55757473/12429735RUN
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR $GOPATH/src/nexus/
COPY . .
# Fetch dependencies.
RUN go mod download
RUN go mod verify
# Build the binary.
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/nexus
############################
# STEP 2 build a small image
############################
FROM scratch
# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Copy our static executable.
COPY --from=builder /app/nexus /app/nexus
# Use an unprivileged user.
USER nexus:nexus
# Run the hello binary.
ENTRYPOINT ["/app/nexus"]