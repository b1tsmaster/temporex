# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /src

RUN apk add --no-cache ca-certificates

# Cache deps
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . .
ARG BUILD_PATH=.
ARG TARGETOS=linux
ARG TARGETARCH=amd64
ENV CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH
RUN go build -trimpath -ldflags="-s -w" -buildvcs=false -o /out/app ${BUILD_PATH}

# Runtime stage
FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app
COPY --from=builder /out/app /usr/local/bin/app

USER nonroot:nonroot
ENV GODEBUG=madvdontneed=1
EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/app"]