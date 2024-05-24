package orders

import (
	"net/http"

	"gihub.com/gmohmad/wb_l0/internal/storage/models/orders"
	"github.com/go-chi/render"
)

type Response struct {
	Status    string            `json:"status"`
	Error     string            `json:"error,omitempty"`
	OrderItem *orders.OrderItem `json:"orderItem,omitempty"`
}

const (
	StatusOK  = "OK"
	StatusErr = "Error"
)

func OK(order orders.OrderItem) Response {
	return Response{
		Status:    StatusOK,
		OrderItem: &order,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusErr,
		Error:  msg,
	}
}

func renderError(w http.ResponseWriter, r *http.Request, statusCode int, errMsg string) {
	w.WriteHeader(statusCode)
	render.JSON(w, r, Error(errMsg))
}
