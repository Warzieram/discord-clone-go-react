############################
# 1. Build stage
############################
FROM golang:1.24.4-alpine AS build

WORKDIR /src

# Dependencies first, so a source-only change doesn't re-download the module
# cache on every build.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO_ENABLED=0 gives a fully static binary (every dependency here is pure Go),
# which is what lets the runtime stage be a bare alpine with no toolchain.
# -trimpath drops local build paths; -s -w strip the symbol and DWARF tables.
# The database schema is embedded via go:embed, so nothing is read from disk.
RUN CGO_ENABLED=0 GOOS=linux \
    go build -trimpath -ldflags="-s -w" -o /out/main ./cmd/api

############################
# 2. Runtime stage
############################
FROM alpine:3.22 AS runtime

# Outbound HTTPS (Mailjet) needs the CA bundle.
RUN apk add --no-cache ca-certificates \
    && adduser -D -u 10001 appuser

WORKDIR /app

COPY --from=build /out/main /app/main

USER appuser

EXPOSE 8080

CMD ["./main"]
