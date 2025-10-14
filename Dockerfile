FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o webhook ./cmd/webhook

FROM gcr.io/distroless/static:nonroot

# Labels for GitHub Container Registry
LABEL org.opencontainers.image.source=https://github.com/nauski/cert-manager-webhook-joker
LABEL org.opencontainers.image.description="cert-manager webhook for Joker.com DNS-01 ACME challenges"
LABEL org.opencontainers.image.licenses=Apache-2.0

COPY --from=builder /app/webhook /webhook

USER 65532:65532

ENTRYPOINT ["/webhook"]