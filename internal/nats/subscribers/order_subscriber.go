package subscribers

import (
	"context"
	"fmt"
	"log/slog"

	uuid "github.com/fossoreslp/go-uuid-v4"
	"github.com/nats-io/stan.go"

	"gihub.com/gmohmad/wb_l0/internal/config"
	"gihub.com/gmohmad/wb_l0/internal/nats"
	"gihub.com/gmohmad/wb_l0/internal/storage/models/orders"
)

type Cache interface {
	AddOrder(id uuid.UUID, order orders.OrderItem)
}

type Storage interface {
	SaveOrder(ctx context.Context, order *orders.OrderItem) (*uuid.UUID, error)
}

type OrderSubscriber struct {
	cache   Cache
	storage Storage
	log     *slog.Logger
}

func NewOrderSubscriber(cache Cache, storage Storage, log *slog.Logger) *OrderSubscriber {
	return &OrderSubscriber{
		cache:   cache,
		storage: storage,
		log:     log,
	}
}

func (ordSub *OrderSubscriber) HandleOrderMessage(ctx context.Context) stan.MsgHandler {
	return func(msg *stan.Msg) {
		order, err := orders.Validate(msg.Data)

		if err != nil {
			ordSub.log.Info(err.Error())
			return
		}

		id, err := ordSub.storage.SaveOrder(ctx, order)

		if err != nil {
			ordSub.log.Info(err.Error())
			return
		}
		ordSub.log.Info(fmt.Sprintf("Successfully saved order, id: %v", id))

		ordSub.cache.AddOrder(*id, *order)

		ordSub.log.Info("Saved order in cache.")
	}
}

func (ordSub *OrderSubscriber) Subscribe(ctx context.Context, conn stan.Conn, subject string) (stan.Subscription, error) {
	sub, err := conn.Subscribe(subject, ordSub.HandleOrderMessage(ctx))

	if err != nil {
		return nil, fmt.Errorf("Error subscribing to '%s' subject: %w", subject, err)
	}

	return sub, nil
}

func (ordSub *OrderSubscriber) Start(ctx context.Context, cfg *config.Nats) error {
	url := fmt.Sprintf("nats://%s:%s", cfg.Host, cfg.Port)
	conn, err := nats.NewNatsConnection(cfg.ClusterId, cfg.ClientId, url)
	if err != nil {
		return err
	}
	defer func() {
		ordSub.log.Info("Closing nats connection")
		conn.Close()
	}()

	sub, err := ordSub.Subscribe(ctx, conn, cfg.Subject)
	if err != nil {
		return err
	}
	defer func() {
		ordSub.log.Info(fmt.Sprintf("Unsubscribing from '%s' subject", cfg.Subject))
		sub.Unsubscribe()
	}()

	<-ctx.Done()

	return nil
}
