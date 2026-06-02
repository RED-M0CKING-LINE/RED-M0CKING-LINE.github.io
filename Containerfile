# podman first, docker later

ARG GO_VERSION=1.25
FROM docker.io/library/golang:${GO_VERSION}-alpine AS builder

WORKDIR /src

# Cache modules independently of source
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Static server binary
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -trimpath -ldflags="-s -w -buildid=" -o /out/server ./cmd/server

# WASM tool + matching wasm_exec.js from the Go toolchain.
RUN GOOS=js GOARCH=wasm go build -trimpath -ldflags="-s -w -buildid=" \
        -o /src/web/static/wasm/tool.wasm ./wasm/tool && \
    cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" /src/web/static/wasm/wasm_exec.js

# layout
RUN mkdir -p /out/web /out/content && \
    cp -r web/templates /out/web/templates && \
    cp -r web/static    /out/web/static && \
    cp -r content/blog  /out/content/blog

# runtime
FROM gcr.io/distroless/static-debian12:nonroot AS runtime

WORKDIR /app
COPY --from=builder --chown=nonroot:nonroot /out/server  /app/server
COPY --from=builder --chown=nonroot:nonroot /out/web     /app/web
COPY --from=builder --chown=nonroot:nonroot /out/content /app/content

USER nonroot:nonroot
EXPOSE 8080

ENV ADDR=":8080" \
    APP_ENV="prod" \
    CONTENT_DIR="/app/content" \
    TEMPLATE_DIR="/app/web/templates" \
    STATIC_DIR="/app/web/static"

# Distroless has no shell, so HEALTHCHECK cant curl
# Probe /healthz directly, thats where health is

ENTRYPOINT ["/app/server"]
