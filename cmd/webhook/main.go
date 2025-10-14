package main

import (
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/cmd"
	"github.com/nauski/cert-manager-webhook-joker/internal/solver"
)

var GroupName = solver.GroupName

func main() {
	cmd.RunWebhookServer(GroupName, solver.New())
}
