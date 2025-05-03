FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache upx make git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the application
ARG TARGETOS=linux
ARG TARGETARCH
ARG BUILDPLATFORM
ARG TARGETPLATFORM

RUN make build-os OS=${TARGETOS} ARCH=${TARGETARCH} && \
    upx build/glimmer-${TARGETOS}-${TARGETARCH} && \
    cp build/glimmer-${TARGETOS}-${TARGETARCH} /app/glimmer

FROM alpine:3.21

ARG UID=1000
ARG GID=1000

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -g $GID glimmer && \
    adduser -D -u $UID -G glimmer glimmer

WORKDIR /app

# Copy the binary directly from the builder stage
COPY --from=builder /app/glimmer /app/glimmer

RUN chmod +x /app/glimmer && \
    mkdir -p /app/pb_data && \
    chown -R glimmer:glimmer /app

USER glimmer:glimmer

ENV CONFIG_DIR=/app/config
ENV LISTEN_ADDRESS=0.0.0.0:8787

EXPOSE 8787
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD pgrep glimmer > /dev/null && wget --no-verbose --tries=1 --spider "http://127.0.0.1:${LISTEN_ADDRESS#*:}/" || exit 1

CMD ["/bin/sh", "-c", "/app/glimmer --dir /app/pb_data serve --http=${LISTEN_ADDRESS} --encryptionEnv=${ENCRYPTION_KEY}"]
