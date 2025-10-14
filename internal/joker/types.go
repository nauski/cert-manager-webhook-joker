package joker

import "context"

type Config struct {
	Username string
	Password string
}

type Client interface {
	CreateTXTRecord(ctx context.Context, zone, label, value string) error
	DeleteTXTRecord(ctx context.Context, zone, label string) error
}
