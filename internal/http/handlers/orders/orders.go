package orders

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"gihub.com/gmohmad/wb_l0/internal/storage/models/orders"
	"github.com/fossoreslp/go-uuid-v4"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5"
)

type OrderFromCache interface {
	GetOrder(id uuid.UUID) (orders.OrderItem, bool)
	AddOrder(id uuid.UUID, order orders.OrderItem)
}

type OrderFromDb interface {
	GetOrder(ctx context.Context, id uuid.UUID) (orders.Order, error)
}

func GetOrder(ctx context.Context, log *slog.Logger, cache OrderFromCache, storage OrderFromDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.With(slog.String("request_id", middleware.GetReqID(r.Context())))

		idFromUrl := chi.URLParam(r, "id")

		id, err := uuid.Parse(idFromUrl)

		if err != nil {

			errMsg := fmt.Sprintf("Invalid id: %s", err)
			log.Info(errMsg)

			renderError(w, r, http.StatusBadRequest, errMsg)
			return
		}

		ordItem, ok := cache.GetOrder(id)

		if ok {
			render.JSON(w, r, OK(ordItem))
			log.Info("Successfully responded from cache")

			return
		}

		order, err := storage.GetOrder(ctx, id)

		if err != nil {
			if err == pgx.ErrNoRows {
				errMsg := fmt.Sprintf("Order with id = '%s' does not exist", id)
				log.Info(errMsg)

				renderError(w, r, http.StatusNotFound, errMsg)
				return
			}
			log.Info(err.Error())
			
			renderError(w, r, http.StatusInternalServerError, "Internal server error")
			return
		}

		render.JSON(w, r, OK(order.OrderItem))
		log.Info("Successfully responded from db")

		cache.AddOrder(id, order.OrderItem)
		log.Info("Saved order in cache")

		return
	}
}
