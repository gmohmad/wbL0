package orders

import "gihub.com/gmohmad/wb_l0/internal/storage/models/orders"

type ResponseBody struct {
	Status    string           `json:"status"`
	Error     string           `json:"error,omitempty"`
	OrderItem orders.OrderItem `json:"orderItem,omitempty"`
}

const (
	StatusOK  = "OK"
	StatusErr = "Error"
)

func OK(order orders.OrderItem) ResponseBody {
	return ResponseBody{
		Status:    StatusOK,
		OrderItem: order,
	}
}

func Error(msg string) ResponseBody {
	return ResponseBody{
		Status: StatusErr,
		Error:  msg,
	}
}
