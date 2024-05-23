package nats

import (
	"fmt"

	"github.com/nats-io/stan.go"
)

func NewNatsConnection(clusterId, clientId, url string) (stan.Conn, error) {

	sc, err := stan.Connect(clusterId, clientId, stan.NatsURL(url))

	if err != nil {
		return nil, fmt.Errorf("Couldn't connent to NATS: %w", err)
	}

	return sc, nil
}
