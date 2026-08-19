# ==========================================
# Stage 1: Build Go Backend with Pre-built Embed FS
# ==========================================
FROM golang:1.25-bookworm AS backend-builder
WORKDIR /app

ENV GOPROXY=https://goproxy.cn,https://proxy.golang.org,direct

# Cache Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code (including pre-built web/dist)
COPY . .

ARG VERSION=0.1.0
ARG COMMIT=dev
ARG BUILD_TIME=unknown

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-s -w \
      -X netip/internal/config.Version=${VERSION} \
      -X netip/internal/config.Commit=${COMMIT} \
      -X netip/internal/config.BuildTime=${BUILD_TIME}" \
    -o /bin/netip ./cmd/netip

# ==========================================
# Stage 2: Minimal Secure Runtime
# ==========================================
FROM debian:bookworm-slim AS runtime

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    libcap2-bin \
    tzdata \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN groupadd -g 10001 netip && \
    useradd -u 10001 -g netip -s /bin/false -M netip

# Copy compiled binary
COPY --from=backend-builder /bin/netip /bin/netip

# Grant NET_RAW capability to binary for ICMP ping without root
RUN setcap cap_net_raw+ep /bin/netip

# Create persistent data and config directories
RUN mkdir -p /data/ipdb /config && \
    chown -R netip:netip /data /config

USER netip
WORKDIR /

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD ["/bin/netip", "-v"] || exit 1

ENTRYPOINT ["/bin/netip"]
CMD ["server"]
