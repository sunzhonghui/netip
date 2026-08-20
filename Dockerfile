# ==========================================
# Stage 1: Build Go Backend with Pre-built Embed FS & Vendor
# ==========================================
FROM golang:1.25-bookworm AS backend-builder
WORKDIR /app

# Copy source code (including pre-built web/dist and vendor)
COPY . .

ARG VERSION=0.1.0
ARG COMMIT=dev
ARG BUILD_TIME=unknown

# Build pure static Go binary offline using vendored dependencies
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor \
    -ldflags "-s -w \
      -X netip/internal/config.Version=${VERSION} \
      -X netip/internal/config.Commit=${COMMIT} \
      -X netip/internal/config.BuildTime=${BUILD_TIME}" \
    -o /bin/netip ./cmd/netip

# ==========================================
# Stage 2: Zero-Network Fast Runtime
# ==========================================
# No apt-get or network operations required during build
FROM debian:bookworm-slim AS runtime

# Copy CA root certificates & timezone data directly from builder
COPY --from=backend-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=backend-builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy compiled binary
COPY --from=backend-builder /bin/netip /bin/netip

# Create persistent data and config directories
RUN mkdir -p /data/ipdb /config

WORKDIR /

EXPOSE 8080

ENTRYPOINT ["/bin/netip"]
CMD ["server"]
