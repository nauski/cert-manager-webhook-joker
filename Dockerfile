FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o webhook ./cmd/webhook

FROM gcr.io/distroless/static:nonroot

COPY --from=builder /app/webhook /webhook

USER 65532:65532

ENTRYPOINT ["/webhook"]